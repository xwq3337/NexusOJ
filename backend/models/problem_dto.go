package models

type ProblemDTO struct {
	Problem
	Status string `json:"status"` // 用户对该题目的状态, 例如 "solved", "attempted", "unattempted"
}
