package config

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
)

var (
	AppMode           string
	Port              string
	Address           string
	Db                string
	DbHost            string
	DbPort            string
	DbUser            string
	DbPassword        string
	DbName            string
	RedisDb           string
	RedisAddr         string
	RedisPort         string
	RedisPwd          string
	RedisDbName       string
	MongoDBName       string
	MongoDBAddr       string
	MongoDBPwd        string
	MongoDBPort       string
	MongoDBUserName   string
	UploadDir         string
	ChunkDir          string
	ImagesDir         string
	AvatarDir         string
	VideosDir         string
	MarkdownImagesDir string
	LogDir            string
	DataDir           string
	JudgeServer       string
	JudgeServerPort   string
)

func init() {
	file, err := ini.Load("config.ini")
	if err != nil {
		panic(fmt.Sprintf("config.ini load failed, err: %s", err))
	}
	loadServer(file)
	loadMysql(file)
	loadMongoDB(file)
	loadRedis(file)
	loadPath(file)
	loadJudge(file)
}

func loadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	Port = file.Section("service").Key("Port").String()
	Address = file.Section("service").Key("Address").String()
}
func loadMysql(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassword = file.Section("mysql").Key("DbPassword").String()
	DbName = file.Section("mysql").Key("DbName").String()
}
func loadMongoDB(file *ini.File) {
	MongoDBName = file.Section("MongoDB").Key("MongoDBName").String()
	MongoDBAddr = file.Section("MongoDB").Key("MongoDBAddr").String()
	MongoDBPwd = file.Section("MongoDB").Key("MongoDBPwd").String()
	MongoDBPort = file.Section("MongoDB").Key("MongoDBPort").String()
	MongoDBUserName = file.Section("MongoDB").Key("MongoDBUserName").String()
}
func loadRedis(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPort = file.Section("redis").Key("RedisPort").String()
	RedisPwd = file.Section("redis").Key("RedisPwd").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}

func loadPath(file *ini.File) {

	UploadDir = file.Section("path").Key("UploadDir").String()
	ChunkDir = file.Section("path").Key("ChunkDir").String()
	ImagesDir = file.Section("path").Key("ImagesDir").String()
	VideosDir = file.Section("path").Key("VideosDir").String()
	AvatarDir = file.Section("path").Key("AvatarDir").String()
	MarkdownImagesDir = file.Section("path").Key("MarkdownImagesDir").String()
	LogDir = file.Section("path").Key("LogDir").String()
	DataDir = file.Section("path").Key("DataDir").String()
	if gin.Mode() == gin.ReleaseMode {
		UploadDir = DataDir + UploadDir
		ChunkDir = DataDir + ChunkDir
		ImagesDir = DataDir + ImagesDir
		VideosDir = DataDir + VideosDir
		AvatarDir = DataDir + AvatarDir
		MarkdownImagesDir = DataDir + MarkdownImagesDir
		LogDir = DataDir + LogDir
	}
}

func loadJudge(file *ini.File) {
	JudgeServer = file.Section("judge").Key("JudgeServer").String()
	JudgeServerPort = file.Section("judge").Key("JudgeServerPort").String()
}
