package models

import (
	"pkg/dao"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type ChatRecord struct {
	SenderID    string
	ReceiverID  string
	Type        uint // 1 私发 2 群发
	Status      bool // 0 未读 1 已读
	MessageType string
	Content     string
	CreatedAt   time.Time
	DeletedAt   *time.Time
}

func CreateChatRecord(chatRecord *ChatRecord) error {
	return dao.InsertDocument("chat", "chat_record", chatRecord)
}

func QueryChatRecord(id1 string, id2 string) ([]ChatRecord, error) {
	results1, err := dao.QueryDocument("chat", "chat_record", bson.M{"senderid": id1, "receiverid": id2})
	results2, err := dao.QueryDocument("chat", "chat_record", bson.M{"senderid": id2, "receiverid": id1})
	var chatRecords []ChatRecord
	var results []bson.M
	if id1 != id2 {
		results = append(results1, results2...)
	} else {
		results = results1
	}
	for _, m := range results {
		var chatRecord ChatRecord
		Bytes, _ := bson.Marshal(m)
		bson.Unmarshal(Bytes, &chatRecord)
		chatRecords = append(chatRecords, chatRecord)
	}
	go func() {
		dao.UpdateDocument("chat", "chat_record", bson.M{"senderid": id2, "receiverid": id1}, bson.M{"status": true})
	}()
	return chatRecords, err
}

func QueryUnReadRecord(id string) (int, error) {
	results, err := dao.QueryDocument("chat", "chat_record", bson.M{"receiverid": id, "status": false})
	return len(results), err
}
