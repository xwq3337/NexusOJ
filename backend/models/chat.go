package models

import (
	"context"
	"nexus/dao"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type ChatRecord struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SenderID    uint64             `json:"sender_id" bson:"sender_id"`
	ReceiverID  uint64             `json:"receiver_id" bson:"receiver_id"`
	Status      bool               `json:"status" bson:"status"`             // 0 未读 1 已读
	MessageType string             `json:"message_type" bson:"message_type"` // 文本(text)、图片(image)、语音(voice)、视频(video)、文件(file)
	Content     string             `json:"content" bson:"content"`           // 消息内容, 可以是文本内容或者文件URL等
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	DeletedAt   *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func CreateChatRecord(chatRecord *ChatRecord) error {
	return dao.InsertDocument("chat", "chat_record", chatRecord)
}

func QueryChatRecord(user_id uint64, friend_id uint64, page int) ([]ChatRecord, error) {
	const pageSize = 20

	// 构建查询条件：查询两个方向的聊天记录
	var filter bson.M
	if user_id != friend_id {
		filter = bson.M{
			"$or": []bson.M{
				{"sender_id": user_id, "receiver_id": friend_id},
				{"sender_id": friend_id, "receiver_id": user_id},
			},
		}
	} else {
		filter = bson.M{"sender_id": user_id, "receiver_id": friend_id}
	}

	// 按时间逆序排序并分页
	skip := int64((page - 1) * pageSize)
	results, err := dao.QueryDocumentWithOptions("chat", "chat_record", filter, bson.M{"created_at": -1}, skip, int64(pageSize))
	if err != nil {
		return nil, err
	}

	// 转换为 ChatRecord 结构
	var chatRecords []ChatRecord
	for _, m := range results {
		var chatRecord ChatRecord
		Bytes, _ := bson.Marshal(m)
		bson.Unmarshal(Bytes, &chatRecord)
		chatRecords = append(chatRecords, chatRecord)
	}

	// 异步更新消息状态（标记为已读）
	go func() {
		dao.UpdateDocument("chat", "chat_record", bson.M{"sender_id": friend_id, "receiver_id": user_id}, bson.M{"status": true})
		// 重置好友关系的未读消息数并触发推送
		totalCount, _ := ResetFriendshipUnreadCount(user_id, friend_id)
		// 触发 Redis 推送更新后的未读数
		channel := "unread_record:" + string(user_id)
		ctx := context.Background()
		dao.RedisClient.Publish(ctx, channel, strconv.Itoa(totalCount))
	}()

	return chatRecords, nil
}

func QueryUnReadRecord(id uint64) (int, error) {
	results, err := dao.QueryDocument("chat", "chat_record", bson.M{"receiver_id": id, "status": false})
	return len(results), err
}

// UpdateFriendshipForNewMessage 当有新消息时更新好友关系
func UpdateFriendshipForNewMessage(senderID, receiverID uint64, content string) error {
	// 更新接收者的好友关系：增加未读数，更新最新消息
	err := dao.MysqlClient.Model(&FriendShips{}).
		Where("user_id = ? AND friend_id = ?", receiverID, senderID).
		Updates(map[string]interface{}{
			"unread_count":   gorm.Expr("unread_count + ?", 1),
			"latest_message": content,
		}).Error

	// 更新发送者的好友关系：仅更新最新消息（不增加未读数）
	err = dao.MysqlClient.Model(&FriendShips{}).
		Where("user_id = ? AND friend_id = ?", senderID, receiverID).
		Update("latest_message", content).Error

	return err
}

// ResetFriendshipUnreadCount 重置好友关系的未读消息数，并返回更新后的总未读数
func ResetFriendshipUnreadCount(userID, friendID uint64) (int, error) {
	err := dao.MysqlClient.Model(&FriendShips{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Update("unread_count", 0).Error

	if err != nil {
		return 0, err
	}

	// 获取更新后的总未读数
	return GetTotalUnreadCount(userID)
}
