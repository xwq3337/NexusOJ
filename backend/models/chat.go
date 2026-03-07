package models

import (
	"nexus/dao"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatRecord struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SenderID    string             `json:"sender_id" bson:"sender_id"`
	ReceiverID  string             `json:"receiver_id" bson:"receiver_id"`
	Status      bool               `json:"status" bson:"status"`             // 0 未读 1 已读
	MessageType string             `json:"message_type" bson:"message_type"` // 文本(text)、图片(image)、语音(voice)、视频(video)、文件(file)
	Content     string             `json:"content" bson:"content"`           // 消息内容, 可以是文本内容或者文件URL等
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	DeletedAt   *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func CreateChatRecord(chatRecord *ChatRecord) error {
	return dao.InsertDocument("chat", "chat_record", chatRecord)
}

func QueryChatRecord(user_id string, friend_id string, page int) ([]ChatRecord, error) {
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
	}()

	return chatRecords, nil
}

func QueryUnReadRecord(id string) (int, error) {
	results, err := dao.QueryDocument("chat", "chat_record", bson.M{"receiver_id": id, "status": false})
	return len(results), err
}
