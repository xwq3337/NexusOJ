package models

type UserDTO struct {
	User
	// TODO 热力图
	Solved int64 `json:"solved"` // 用户已经解决的题目数量
}
