package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"nexus/dao"
	"nexus/models"
	"nexus/utils"
	"nexus/utils/jsonx"
	"nexus/utils/logger"
	"strconv"
	"sync"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatController struct{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	HandshakeTimeout: 10 * time.Second,
}

type Message struct {
	SenderID    uint64 `json:"sender_id"`
	ReceiverID  uint64 `json:"receiver_id"`  // 群聊时为群ID
	MessageType string `json:"message_type"` // 消息类型 text/image/file/video/voice
	Content     string `json:"content"`      // 消息内容
	Timestamp   int64  `json:"timestamp"`    // 消息发送时间戳
}

type Client struct {
	Conn          *websocket.Conn
	User          uint64
	Addr          string
	HeartBeatTime uint64
	DataQueue     chan []byte
	GroupSets     sync.Map
	mu            sync.Mutex // 防止并发写
}

var (
	clients   = make(map[uint64]*Client)
	clientsMu sync.RWMutex
)

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Conn.Close()
	close(c.DataQueue)
}
func (c *Client) Send(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.Conn.WriteMessage(websocket.TextMessage, data)
	return err
}

func (ChatController) Handler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(fmt.Errorf("WebSocket 升级失败: %w", err))
		logger.Error("WebSocket 升级失败", err)
		return
	}
	UserID, err := ParserTokenByString(c.Query("token"))
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, []byte("用户ID解析失败"+err.Error()))
		return
	}
	// 创建客户端实例
	client := &Client{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(),
		User:          UserID,
		HeartBeatTime: uint64(time.Now().Unix()),
		DataQueue:     make(chan []byte, 1024),
	}
	go addClient(UserID, client)
	go client.writePump() // 写消息
	client.readPump()     // 读消息
	defer func() {
		client.Close()
		removeClient(client)
	}()
}
func addClient(userID uint64, client *Client) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	clients[userID] = client
}
func findClient(id uint64) *Client {
	clientsMu.RLock()
	defer clientsMu.RUnlock()
	return clients[id]
}
func removeClient(client *Client) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	delete(clients, client.User)
	logger.Debug("客户端断开", client.User)
}
func (client *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second) // 心跳间隔
	defer ticker.Stop()
	for {
		select {
		case data := <-client.DataQueue:
			err := client.Send(data)
			if err != nil {
				logger.Error("发送消息错误", err)
				return
			}
		case <-ticker.C:
			// 发送心跳
			client.mu.Lock()
			err := client.Conn.WriteMessage(websocket.PingMessage, nil)
			client.mu.Unlock()

			if err != nil {
				logger.Error("发送心跳失败", err)
				return
			}
		}
	}
}
func (client *Client) readPump() {
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			logger.Error("传输错误", err)
			return
		}
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			logger.Error("连接断开", err)
			return
		}
		if bytes.Equal(message, []byte("ping")) {
			client.Send([]byte("pong"))
			logger.Debug("收到心跳")
		} else {
			var msg Message
			err = json.Unmarshal([]byte(message), &msg)
			if err != nil {
				logger.Error("解码失败", err)
				continue
			}
			logger.Debugf("收到消息: %s , 来自 %s", message, msg.SenderID)
			dispatch(msg)
		}
		client.HeartBeatTime = uint64(time.Now().Unix())
	}
}
func dispatch(msg Message) {
	sendPrivateMsg(msg)
	// sendGroupMsg(msg)
}
func broadMsg(_ Message) error { // 局域网广播
	return nil
}

func sendPrivateMsg(msg Message) { //私聊
	channel := "unread_record:" + string(msg.ReceiverID)
	ctx := context.Background()
	count, _ := models.QueryUnReadRecord(msg.ReceiverID)
	err := dao.RedisClient.Publish(ctx, channel, strconv.Itoa(count+1)).Err()
	if err != nil {
		logger.Error("redis publish error", err)
	}
	logger.Debug(count, "未读消息")
	// ----
	target_client := findClient(msg.ReceiverID)
	sender_client := findClient(msg.SenderID)
	var online bool
	if target_client != nil {
		online = true
	} else {
		online = false
	}
	go func() {
		chat_record := &models.ChatRecord{
			SenderID:    msg.SenderID,
			ReceiverID:  msg.ReceiverID,
			Status:      online,
			Content:     msg.Content,
			MessageType: "text",
			CreatedAt:   time.Now(),
		}
		err := models.CreateChatRecord(chat_record)
		if err != nil {
			return
		} else {
			// 更新好友关系的未读数和最新消息
			models.UpdateFriendshipForNewMessage(msg.SenderID, msg.ReceiverID, msg.Content)

			message, _ := jsonx.Marshal(chat_record)
			if target_client != nil && target_client != sender_client {
				target_client.DataQueue <- []byte(message)
			}
			sender_client.DataQueue <- []byte(message) //告诉发送者，是否发送成功
		}
	}()
}
func sendGroupMsg(msg Message) { //群聊
	data := msg.Content
	// sender_client := findClient(msg.SenderID)
	target_id := msg.ReceiverID
	logger.Debug("群发消息[%s]到%s", data, target_id)
}

func (ChatController) GetChatRecord(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}
	friend_id, err := strconv.ParseUint(c.Query("friend_id"), 10, 64)
	if err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "缺少好友ID")
		return
	}
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // 默认第一页
	}
	chatRecords, err := models.QueryChatRecord(userID, friend_id, page)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", chatRecords)
}

func (ChatController) MarkMessagesAsRead(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}
	friendID, err := strconv.ParseUint(c.Query("friend_id"), 10, 64)
	if err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "缺少好友ID")
		return
	}

	totalCount, err := models.ResetFriendshipUnreadCount(userID, friendID)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}

	// 触发 Redis 推送更新后的未读数
	channel := "unread_record:" + string(userID)
	ctx := context.Background()
	err = dao.RedisClient.Publish(ctx, channel, strconv.Itoa(totalCount)).Err()
	if err != nil {
		logger.Error("redis publish error", err)
	}

	utils.ReturnSuccess(c, http.StatusOK, "标记成功", nil)
}

func (ChatController) GetUnReadRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "缺少用户ID")
		return
	}
	channel := "unread_record:" + string(id)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	pubsub := dao.RedisClient.Subscribe(c.Request.Context(), channel)
	ticker := time.NewTicker(20 * time.Second)

	msgChan := make(chan int, 10)
	defer func() {
		close(msgChan)
		ticker.Stop()
		pubsub.Close()
	}()
	count, _ := models.QueryUnReadRecord(id)
	msgChan <- count
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			if count, err := strconv.Atoi(msg.Payload); err == nil {
				select {
				case msgChan <- count:
				default:
				}
			}
		}
	}()

	// Flush the initial headers
	c.Writer.Flush()
	for {
		select {
		case count := <-msgChan:
			sse.Encode(c.Writer, sse.Event{
				Event: "message",
				Data:  strconv.Itoa(count),
			})
			c.Writer.Flush()

		case <-ticker.C:
			sse.Encode(c.Writer, sse.Event{
				Event: "heartbeat",
				Data:  "pong",
			})
			c.Writer.Flush()

		case <-c.Request.Context().Done():
			return
		}
	}
}
