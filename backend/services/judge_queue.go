package services

import (
	"context"
	"fmt"
	"log"
	"nexus/config"
	"nexus/models"
	"nexus/utils/chttp"
	"sync"
	"time"
)

type JudgeTask struct {
	SubmissionID int64
	ProblemID    string // 题库题目ID（比赛时为 ContestProblem.ID）
	UserID       uint64
	Code         string
	Language     string
	TestCases    []models.JudgeTestCase
	JudgeConfig  models.JudgeConfig
	ContestID    string // 空=题库提交，非空=比赛提交
	ProblemLabel string // 比赛时使用（A/B/C...）
}

type JudgeQueue struct {
	taskChan   chan *JudgeTask
	workerNum  int
	wg         sync.WaitGroup
	judgeQueue *JudgeServerQueue
	results    map[int64]*models.JudgeOutputResult
	resultMap  sync.Map // submissionID -> *JudgeOutputResult
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
		results:    make(map[int64]*models.JudgeOutputResult),
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

// SubmitTaskSync 提交任务并等待结果
func (q *JudgeQueue) SubmitTaskSync(task *JudgeTask) (*models.JudgeOutputResult, error) {
	// 创建一个接收结果的 channel
	resultChan := make(chan *models.JudgeOutputResult, 1)
	errChan := make(chan error, 1)

	// 存储 channel 以便 worker 可以发送结果
	q.resultMap.Store(task.SubmissionID, resultChan)

	// 提交任务到队列
	if err := q.SubmitTask(task); err != nil {
		q.resultMap.Delete(task.SubmissionID)
		return nil, err
	}

	// 等待结果或超时
	timeout := time.After(120 * time.Second) // 2分钟超时
	select {
	case result := <-resultChan:
		q.resultMap.Delete(task.SubmissionID)
		return result, nil
	case err := <-errChan:
		q.resultMap.Delete(task.SubmissionID)
		return nil, err
	case <-timeout:
		q.resultMap.Delete(task.SubmissionID)
		return nil, fmt.Errorf("judge timeout after 120 seconds")
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

		// 如果有等待的 channel，发送错误
		if value, ok := q.resultMap.Load(task.SubmissionID); ok {
			if resultChan, ok := value.(chan *models.JudgeOutputResult); ok {
				close(resultChan)
			}
		}
		return
	}

	// 保存判题结果
	if err := q.saveResult(task, result); err != nil {
		log.Printf("Worker %d: failed to save result for submission %d - %v", workerID, task.SubmissionID, err)

		// 如果有等待的 channel，发送错误
		if value, ok := q.resultMap.Load(task.SubmissionID); ok {
			if resultChan, ok := value.(chan *models.JudgeOutputResult); ok {
				close(resultChan)
			}
		}
		return
	}

	// 更新题目统计
	q.updateProblemStats(task, result.Verdict)

	// 如果有等待的 channel，发送结果
	if value, ok := q.resultMap.Load(task.SubmissionID); ok {
		if resultChan, ok := value.(chan *models.JudgeOutputResult); ok {
			resultChan <- result
		}
	}

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
	if task.ContestID != "" {
		// 比赛提交 → 写入 contest_record
		record := &models.ContestRecord{
			ID:           task.SubmissionID,
			ContestID:    task.ContestID,
			UserID:       task.UserID,
			ProblemLabel: task.ProblemLabel,
			Code:         task.Code,
			Language:     task.Language,
			Verdict:      models.SystemError,
		}
		if createErr := models.CreateContestRecord(record); createErr != nil {
			log.Printf("Failed to create error contest_record for submission %d: %v", task.SubmissionID, createErr)
		}
		return
	}

	// 题库提交 → 写入 record
	record := &models.Record{
		ID:        task.SubmissionID,
		UserId:    task.UserID,
		ProblemId: task.ProblemID,
		Code:      task.Code,
		Language:  task.Language,
		Verdict:   models.SystemError,
	}

	if createErr := models.CreateRecord(record); createErr != nil {
		log.Printf("Failed to create error record for submission %d: %v", task.SubmissionID, createErr)
	}
}

func (q *JudgeQueue) saveResult(task *JudgeTask, result *models.JudgeOutputResult) error {
	if task.ContestID != "" {
		// 比赛提交 → 写入 contest_record
		record := &models.ContestRecord{
			ID:           task.SubmissionID,
			ContestID:    task.ContestID,
			UserID:       task.UserID,
			ProblemLabel: task.ProblemLabel,
			Code:         task.Code,
			Language:     task.Language,
			MaxTime:      result.MaxTime,
			MaxMemory:    result.MaxMemory,
			Verdict:      result.Verdict,
			JudgeResult:  result.Result,
		}
		return models.CreateContestRecord(record)
	}

	// 题库提交 → 写入 record
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

func (q *JudgeQueue) updateProblemStats(task *JudgeTask, verdict models.JudgeVerdict) {
	if task.ContestID != "" {
		// 比赛提交 → 更新 contest_problem 统计
		cp, err := (models.ContestProblem{}).GetContestProblemByLabel(task.ContestID, task.ProblemLabel)
		if err != nil {
			log.Printf("Failed to get contest_problem for stats update: %v", err)
			return
		}
		cp.IncrSubmission(cp.ID, verdict == models.Accepted)
	} else {
		// 题库提交 → 更新 problem 统计
		problem, err := models.Problem{}.GetProblemInfoWithoutUsername(task.ProblemID)
		if err != nil {
			log.Printf("Failed to get problem %s for stats update: %v", task.ProblemID, err)
			return
		}

		problem.Submission++
		if verdict == models.Accepted {
			problem.Accept++
		}

		problemModel := models.Problem{}
		if updateErr := problemModel.UpdateProblem(&problem); updateErr != nil {
			log.Printf("Failed to update stats for problem %s: %v", task.ProblemID, updateErr)
		}
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
