/*
 * @Author: x_wq3337 854541540@qq.com
 * @Date: 2025-07-29 16:28:52
 * @LastEditors: x_wq3337 854541540@qq.com
 * @LastEditTime: 2026-01-15 13:07:15
 * @FilePath: /backend/models/judge.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package models

import (
	"encoding/json"
	"pkg/utils/jsonx"
)

type JudgeResourcesLimits struct {
	CpuTime     int `json:"cpu_time"`
	MemoryBytes int `json:"memory_bytes"`
	StackBytes  int `json:"stack_bytes"`
	OutputBytes int `json:"output_bytes"`
}

type JudgeOutputResult struct {
	SubmissionID int64                 `json:"submission_id"`
	Language     string                `json:"language"`
	Code         string                `json:"code"`
	Verdict      JudgeVerdict          `json:"verdict"`
	MaxTime      float32               `json:"max_time"`
	MaxMemory    float32               `json:"max_memory"`
	Result       []JudgeTestCaseResult `json:"result"`
}

func (o *JudgeOutputResult) String() string {
	jsonBytes, _ := json.MarshalIndent(o, "", "  ")
	return string(jsonBytes)
}

type JudgeTestCaseResult struct {
	CaseID   string       `json:"case_id"`
	Stdin    string       `json:"stdin"`
	Stdout   string       `json:"stdout"`
	Stderr   string       `json:"stderr"`
	Status   JudgeVerdict `json:"status"`
	Time     float32      `json:"time"`
	Memory   float32      `json:"memory"`
	Expected string       `json:"expected"`
}

type JudgeVerdict string

const (
	Accepted             JudgeVerdict = "Accepted"
	WrongAnswer          JudgeVerdict = "WrongAnswer"
	TimeLimitExceeded    JudgeVerdict = "TimeLimitExceeded"
	MemoryLimitExceeded  JudgeVerdict = "MemoryLimitExceeded"
	RuntimeError         JudgeVerdict = "RuntimeError"
	RestrictedSystemCall JudgeVerdict = "RestrictedSystemCall"
	CompilationError     JudgeVerdict = "CompilationError"
	SystemError          JudgeVerdict = "SystemError"
)

type JudgeLanguageConfig struct {
	CompileCmd    string   `json:"compile_cmd"`
	RunCmd        string   `json:"run_cmd"`
	AllowSysCalls []string `json:"allow_sys_calls"`
}

type JudgeInputStruct struct {
	SubmissionID    int64                `json:"submission_id"`
	Language        string               `json:"language"`
	Code            string               `json:"code"`
	TestCases       []JudgeTestCase      `json:"test_cases"`
	ResourcesLimits JudgeResourcesLimits `json:"resources_limits"`
	Message         string               `json:"message"`
	SeccompProfile  string               `json:"seccomp_profile"`
}
type ProblemJudgeStruct struct {
	ProblemID string `json:"problem_id"`
	Code      string `json:"code"`
	Language  string `json:"language"`
}

func (i JudgeInputStruct) String() string {
	jsonBytes, _ := jsonx.MarshalIndent(
		i,
		"",
		"  ",
	)
	return string(jsonBytes)
}

type JudgeTestCase struct {
	CaseID   int    `json:"case_id"`
	Stdin    string `json:"stdin"`
	Expected string `json:"expected"`
}
