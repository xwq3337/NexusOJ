package dao

import (
	"context"
	"errors"
	"fmt"
	"pkg/config"
	"pkg/utils/logger"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	MysqlClient   *gorm.DB
	MongoDBClient *mongo.Client
	RedisClient   *redis.Client
	ctx           = context.Background()
)

func init() {
	InitMongoDB()
	mySql()
}

func InitMongoDB() {
	// 1. 创建客户端选项
	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s", config.MongoDBAddr, config.MongoDBPort)).
		SetAuth(options.Credential{
			AuthSource: "admin", // 明确指定认证数据库
			Username:   config.MongoDBUserName,
			Password:   config.MongoDBPwd,
		}).
		SetConnectTimeout(10 * time.Second).         // 设置连接超时
		SetServerSelectionTimeout(10 * time.Second). // 设置服务器选择超时
		SetMaxPoolSize(20)
	// 2. 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var err error
	// 3. 连接 MongoDB
	MongoDBClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to MongoDB: %s", err))
	}
}
func mySql() {
	var err error
	gsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_0900_ai_ci&parseTime=True&loc=Local", config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	MysqlClient, err = gorm.Open(mysql.Open(gsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		logger.Error("mysql connect error", err)
		panic(err)
	}
	db, _ := MysqlClient.DB()
	db.SetMaxIdleConns(10)  // 空闲连接数
	db.SetMaxOpenConns(100) // 打开最大连接
	db.SetConnMaxLifetime(time.Hour)
}
func init() {
	db, _ := strconv.ParseUint(config.RedisDbName, 10, 64)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisAddr, config.RedisPort),
		Password: config.RedisPwd,
		DB:       int(db),
		PoolSize: 20,
	})
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		logger.Error("redis connect error", err)
		panic(err)
	}
}

func InsertDocument(db, collection string, document interface{}) error {
	Collection := MongoDBClient.Database(db).Collection(collection)
	_, err := Collection.InsertOne(ctx, document)
	return err
}
func QueryDocument(db, collection string, filter bson.M) ([]bson.M, error) {
	if MongoDBClient == nil {
		return nil, errors.New("MongoDB 客户端未初始化")
	}

	Collection := MongoDBClient.Database(db).Collection(collection)
	cur, err := Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	err = cur.All(ctx, &results)
	return results, err
}
func UpdateDocument(db, collection string, filter, update bson.M) error {
	Collection := MongoDBClient.Database(db).Collection(collection)
	_, err := Collection.UpdateMany(ctx, filter, bson.M{
		"$set": update,
	})
	return err
}
func DeleteDocuments(db, collection string, filter bson.M) error {
	Collection := MongoDBClient.Database(db).Collection(collection)
	_, err := Collection.DeleteMany(ctx, filter)
	return err
}
