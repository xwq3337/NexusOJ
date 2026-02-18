package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"pkg/dao"
	"pkg/models"
	"pkg/utils"
	"pkg/utils/jsonx"
	"pkg/utils/logger"
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

type MessageStruct struct {
	Timestamp   int64   `json:"timestamp" gorm:"primaryKey;autoIncrement"`
	CurrentTime int64   `json:"currenttime" gorm:"autoUpdateTime:nano"`
	Sender      string  `json:"sender"`
	Receiver    string  `json:"receiver"`
	Text        *string `json:"text"`
	Type        uint    `json:"type"` // 1私信 2群发 3心跳检测
}

type Client struct {
	Conn          *websocket.Conn
	User          string
	Addr          string
	HeartBeatTime uint64
	DataQueue     chan []byte
	GroupSets     sync.Map
	mu            sync.Mutex // 防止并发写
}

var (
	clients   = make(map[string]*Client)
	clientsMu sync.RWMutex
)

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Conn.Close()
	close(c.DataQueue)
}

func (ChatController) Handler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(fmt.Errorf("WebSocket 升级失败: %w", err))
		logger.Error("WebSocket 升级失败", err)
		return
	}
	sender_id := c.Query("id")
	if sender_id == "" {
		conn.WriteMessage(websocket.CloseMessage, []byte("缺少用户ID"))
		return
	}
	currentTime := uint64(time.Now().Unix())
	client := &Client{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(),
		User:          sender_id,
		HeartBeatTime: currentTime,
		DataQueue:     make(chan []byte, 1024),
	}
	go addClient(sender_id, client)
	go client.writePump()
	client.readPump()
	defer func() {
		client.Close()
		removeClient(client)
	}()
}
func addClient(userID string, client *Client) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	clients[userID] = client
}
func findClient(id string) *Client {
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
			client.mu.Lock()
			err := client.Conn.WriteMessage(websocket.TextMessage, data)
			client.mu.Unlock()
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
			client.mu.Lock()
			client.Conn.WriteMessage(websocket.TextMessage, []byte("pong"))
			client.mu.Unlock()
			logger.Debug("收到心跳")
		} else {
			var msg MessageStruct
			err = json.Unmarshal([]byte(message), &msg)
			if err != nil {
				logger.Error("解码失败", err)
				continue
			}
			logger.Debugf("收到消息: %s , 来自 %s", message, msg.Sender)
			dispatch(msg)
		}
		client.HeartBeatTime = uint64(time.Now().Unix())

	}
}
func dispatch(msg MessageStruct) {
	switch msg.Type {
	case 1:
		sendPrivateMsg(msg)
	case 2:
		sendGroupMsg(msg)
	default:
		logger.Debug("未知消息类型: %d", msg.Type)
	}
}
func broadMsg(msg MessageStruct) error { // 局域网广播
	return nil
}

func sendPrivateMsg(msg MessageStruct) { //私聊
	// --- TODO
	channel := "unread_record:" + msg.Receiver
	ctx := context.Background()
	count, _ := models.QueryUnReadRecord(msg.Receiver)
	err := dao.RedisClient.Publish(ctx, channel, strconv.Itoa(count+1)).Err()
	if err != nil {
		logger.Error("redis publish error", err)
	}
	logger.Debug(count, "未读消息")
	// ----
	target_client := findClient(msg.Receiver)
	sender_client := findClient(msg.Sender)
	var online bool
	if target_client != nil {
		online = true
	} else {
		online = false
	}
	go func() {
		chat_record := &models.ChatRecord{
			SenderID:   msg.Sender,
			ReceiverID: msg.Receiver,
			Status:     online,
			Content:    *msg.Text,
			CreatedAt:  time.Now(),
		}
		err := models.CreateChatRecord(chat_record)
		if err != nil {
			return
		} else {
			message, _ := jsonx.Marshal(chat_record)
			if target_client != nil && target_client != sender_client {
				target_client.DataQueue <- []byte(message)
			}
			sender_client.DataQueue <- []byte(message) //告诉发送者，是否发送成功
		}
	}()
}
func sendGroupMsg(msg MessageStruct) { //群聊
	data := *msg.Text
	// sender_client := findClient(msg.Sender)
	target_id := msg.Receiver
	logger.Debug("群发消息[%s]到%s", data, target_id)
}

func (ChatController) GetChatRecord(c *gin.Context) {
	userID, _ := ParserToken(c)
	id1 := userID
	id2 := c.Query("id")
	chatRecords, err := models.QueryChatRecord(id1, id2)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", chatRecords)
}

func (ChatController) GetUnReadRecord(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		utils.ReturnError(c, http.StatusBadRequest, "缺少用户ID")
		return
	}
	channel := "unread_record:" + id
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
