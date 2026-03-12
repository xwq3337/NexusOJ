package models

import "time"

type FriendshipDTO struct {
	ID             uint      `json:"id"` // 好友关系ID
	UserID         string    `json:"user_id"`
	FriendID       string    `json:"friend_id"`
	FriendAvatar   string    `json:"friend_avatar"`   // 好友头像
	FriendUsername string    `json:"friend_username"` // 好友用户名
	FriendNickname string    `json:"friend_nickname"` // 好友昵称
	Remark         string    `json:"remark"`          // 好友备注
	UnreadCount    int       `json:"unread_count"`    // 未读消息数量
	LatestMessage  *string   `json:"latest_message"`  // 最新消息内容
	CreatedAt      time.Time `json:"created_at"`      // 好友添加时间
}

type FriendshipRequestDTO struct {
	ID             uint      `json:"id"` // 好友请求ID
	UserID         string    `json:"user_id"`
	FriendID       string    `json:"friend_id"`
	FriendAvatar   string    `json:"friend_avatar"`   // 好友头像
	FriendUsername string    `json:"friend_username"` // 好友用户名
	FriendNickname string    `json:"friend_nickname"` // 好友昵称
	Status         string    `json:"status"`          // 请求状态
	Message        string    `json:"message"`         // 请求消息
	CreatedAt      time.Time `json:"created_at"`      // 请求发送时间
}
