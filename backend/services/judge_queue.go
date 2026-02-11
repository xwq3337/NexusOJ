package services

import (
	"context"
	"fmt"
	"log"
	"pkg/config"
	"pkg/models"
	"pkg/utils/chttp"
	"sync"
	"time"
)

type JudgeTask struct {
	SubmissionID int64
	ProblemID    string
	UserID       string
	Code         string
	Language     string
	TestCases    []models.JudgeTestCase
	JudgeConfig  models.JudgeConfig
}

type JudgeQueue struct {
	taskChan   chan *JudgeTask
	workerNum  int
	wg         sync.WaitGroup
	judgeQueue *JudgeServerQueue
}

type JudgeServerQueue struct {
	servers []string
	index   int
	mu      sync.Mutex
}

var GlobalJudgeQueue *JudgeQueue

func InitJudgeQueue(workerNum, queueSize int) {
	GlobalJudgeQueue = &JudgeQueue{
		taskChan:   make(chan *JudgeTask, queueSize),
		workerNum:  workerNum,
		judgeQueue: NewJudgeServerQueue(),
	}

	for i := 0; i < workerNum; i++ {
		GlobalJudgeQueue.wg.Add(1)
		go GlobalJudgeQueue.worker(i)
	}

	log.Printf("Judge queue initialized with %d workers and queue size %d", workerNum, queueSize)
}

func NewJudgeServerQueue() *JudgeServerQueue {
	// 使用配置文件中的判题服务器地址
	judgeServerAddr := fmt.Sprintf("%s:%s", config.JudgeServer, config.JudgeServerPort)
	return &JudgeServerQueue{
		servers: []string{judgeServerAddr},
		index:   0,
	}
}

func (q *JudgeServerQueue) GetNextServer() string {
	q.mu.Lock()
	defer q.mu.Unlock()

	server := q.servers[q.index]
	q.index = (q.index + 1) % len(q.servers)
	return server
}

func (q *JudgeQueue) SubmitTask(task *JudgeTask) error {
	select {
	case q.taskChan <- task:
		return nil
	default:
		return fmt.Errorf("judge queue is full")
	}
}

func (q *JudgeQueue) worker(id int) {
	defer q.wg.Done()

	log.Printf("Judge worker %d started", id)

	for task := range q.taskChan {
		q.processTask(task, id)
	}

	log.Printf("Judge worker %d stopped", id)
}

func (q *JudgeQueue) processTask(task *JudgeTask, workerID int) {
	log.Printf("Worker %d processing submission %d", workerID, task.SubmissionID)

	// 生成判题配置
	config := q.buildJudgeConfig(task)

	// 获取判题服务器
	serverAddr := q.judgeQueue.GetNextServer()

	// 执行判题
	result, err := q.evaluateCode(serverAddr, config)
	if err != nil {
		log.Printf("Worker %d: submission %d failed - %v", workerID, task.SubmissionID, err)
		q.handleJudgeError(task, err)
		return
	}

	// 保存判题结果
	if err := q.saveResult(task, result); err != nil {
		log.Printf("Worker %d: failed to save result for submission %d - %v", workerID, task.SubmissionID, err)
		return
	}

	// 更新题目统计
	q.updateProblemStats(task.ProblemID, result.Verdict)

	log.Printf("Worker %d: submission %d completed with verdict %s", workerID, task.SubmissionID, result.Verdict)
}

func (q *JudgeQueue) buildJudgeConfig(task *JudgeTask) models.JudgeInputStruct {
	// 根据题目配置设置资源限制
	memoryLimit := int(task.JudgeConfig.MemoryLimit)
	if memoryLimit <= 0 {
		memoryLimit = 64 // 默认64MB
	}

	return models.JudgeInputStruct{
		SubmissionID: task.SubmissionID,
		Language:     task.Language,
		Code:         task.Code,
		TestCases:    task.TestCases,
		ResourcesLimits: models.JudgeResourcesLimits{
			CpuTime:     100000,                    // CPU配额份额(固定值,非时间限制)
			MemoryBytes: memoryLimit * 1024 * 1024, // 内存限制(默认64MB)
			StackBytes:  64 * 1024 * 1024,          // 栈空间限制10MB
			OutputBytes: 10485760,                  // 输出限制10MB
		},
		Message:        "",
		SeccompProfile: q.getSeccompProfile(task.Language),
	}
}

func (q *JudgeQueue) getSeccompProfile(language string) string {
	profiles := map[string]string{
		"cpp":        "cpp",
		"c":          "c",
		"python":     "python",
		"python3":    "python",
		"java":       "java",
		"go":         "go",
		"rust":       "rust",
		"javascript": "javascript",
		"typescript": "typescript",
	}

	if profile, ok := profiles[language]; ok {
		return profile
	}
	return "general"
}

func (q *JudgeQueue) evaluateCode(serverAddr string, config models.JudgeInputStruct) (*models.JudgeOutputResult, error) {
	// 设置更长的超时时间以容纳编译和执行时间
	// 包括: 编译时间 + 实际运行时间 + 网络延迟
	timeout := 60 * time.Second

	client := chttp.New(
		chttp.WithBaseURL(serverAddr),
		chttp.WithTimeout(timeout),
		chttp.WithHeader("Accept", "application/json"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var result models.JudgeOutputResult
	err := client.Post(ctx, "/submit", &config, &result)
	if err != nil {
		return nil, fmt.Errorf("judge request failed: %w", err)
	}

	return &result, nil
}

func (q *JudgeQueue) handleJudgeError(task *JudgeTask, err error) {
	record := &models.Record{
		ID:        task.SubmissionID,
		UserId:    task.UserID,
		ProblemId: task.ProblemID,
		Code:      task.Code,
		Language:  task.Language,
		Verdict:   models.SystemError,
		MaxTime:   0,
		MaxMemory: 0,
	}

	if createErr := models.CreateRecord(record); createErr != nil {
		log.Printf("Failed to create error record for submission %d: %v", task.SubmissionID, createErr)
	}
}

func (q *JudgeQueue) saveResult(task *JudgeTask, result *models.JudgeOutputResult) error {
	record := &models.Record{
		ID:          task.SubmissionID,
		UserId:      task.UserID,
		ProblemId:   task.ProblemID,
		Code:        task.Code,
		Language:    task.Language,
		MaxTime:     result.MaxTime,
		MaxMemory:   result.MaxMemory,
		Verdict:     result.Verdict,
		JudgeResult: result.Result,
	}

	return models.CreateRecord(record)
}

func (q *JudgeQueue) updateProblemStats(problemID string, verdict models.JudgeVerdict) {
	problem, err := models.Problem{}.GetProblemInfoWithoutUsername(problemID)
	if err != nil {
		log.Printf("Failed to get problem %s for stats update: %v", problemID, err)
		return
	}

	problem.Submission++
	if verdict == models.Accepted {
		problem.Accept++
	}

	problemModel := models.Problem{}
	if updateErr := problemModel.UpdateProblem(&problem); updateErr != nil {
		log.Printf("Failed to update stats for problem %s: %v", problemID, updateErr)
	}
}

func (q *JudgeQueue) Stop() {
	close(q.taskChan)
	q.wg.Wait()
	log.Println("Judge queue stopped")
}

func GetQueueStatus() map[string]interface{} {
	if GlobalJudgeQueue == nil {
		return map[string]interface{}{
			"status":       "not_initialized",
			"queue_length": 0,
			"worker_num":   0,
		}
	}

	return map[string]interface{}{
		"status":       "running",
		"queue_length": len(GlobalJudgeQueue.taskChan),
		"worker_num":   GlobalJudgeQueue.workerNum,
		"capacity":     cap(GlobalJudgeQueue.taskChan),
	}
}
