package controllers

import (
	"nexus/models"
	"nexus/utils"

	"github.com/gin-gonic/gin"
)

type FriendshipController struct{}

func (FriendshipController) CreateFriendship(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, 401, "未授权")
		return
	}
	var param struct {
		FriendID uint64 `json:"friend_id"` // 接收好友请求的用户ID
		Message  string `json:"message"`   // 验证消息
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.ReturnError(c, 400, err.Error())
		return
	}
	// 2. 调用模型层方法创建好友请求记录
	err = models.FriendShipRequest{}.CreateRequest(userID, param.FriendID, param.Message)

	// 3. 返 回结果给客户端
	if err != nil {
		utils.ReturnError(c, 500, err.Error())
		return
	}
	utils.ReturnSuccess(c, 200, "好友请求已发送", nil)
}

func (FriendshipController) HandleFriendshipRequest(c *gin.Context) {
	// 1. 获取请求参数, 包括好友请求ID和处理结果等
	var param struct {
		RequestID uint   `json:"request_id"` // 好友请求ID
		Status    string `json:"status"`     // 处理结果: accepted 或 rejected
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.ReturnError(c, 400, err.Error())
		return
	}
	// 2. 调用模型层方法处理好友请求
	result, err := models.FriendShipRequest{}.HandleRequest(param.RequestID, param.Status)
	// 3. 返回结果给客户端
	if err != nil {
		utils.ReturnError(c, 500, err.Error())
		return
	}
	utils.ReturnSuccess(c, 200, result, nil)
}

func (FriendshipController) GetFriendList(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, 401, "未授权")
		return
	}
	friends, err := models.FriendShips{}.GetFriendList(userID)
	if err != nil {
		utils.ReturnError(c, 500, err.Error())
		return
	}
	utils.ReturnSuccess(c, 200, "success", friends)
}

func (FriendshipController) GetFriendRequestList(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, 401, "未授权")
		return
	}
	requests, err := models.FriendShipRequest{}.GetFriendRequestList(userID)
	if err != nil {
		utils.ReturnError(c, 500, err.Error())
		return
	}
	utils.ReturnSuccess(c, 200, "success", requests)
}
