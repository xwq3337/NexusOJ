package router

import (
	"net/http"
	"pkg/config"
	"pkg/controllers"
	"pkg/middleware"
	"pkg/middleware/jwt"
	"pkg/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func noCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}
func Router() *gin.Engine {
	r := gin.Default()
	assets := r.Group("/assets")
	// assets.Use(noCacheMiddleware())
	{
		assets.Static("/avatar", config.AvatarDir)
		assets.Static("/images", config.ImagesDir)
		assets.Static("/videos", config.VideosDir)
		assets.Static("/markdown/images", config.MarkdownImagesDir)
	}
	r.SetTrustedProxies(nil)
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Real-IP"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           50 * time.Second,
	}))
	// r.Use(ipinterceptor.Interceptor())
	r.Use(middleware.RequestLogger(middleware.DefaultRequestLoggerConfig()))
	wsGroup := r.Group("/ws")
	{
		wsGroup.GET("/chat", controllers.ChatController{}.Handler)
	}
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/ip", func(c *gin.Context) {
			utils.ReturnSuccess(c, http.StatusOK, "success", c.ClientIP())
		})
		apiGroup.GET("/ping", func(c *gin.Context) {
			utils.ReturnSuccess(c, http.StatusOK, "success", "pong")
		})
		apiGroup.GET("/time", func(c *gin.Context) {
			utils.ReturnSuccess(c, http.StatusOK, "success", time.Now().Unix())
		})
	}
	// TODO: 写出中间件判断用户角色是否为管理员
	adminGroup := r.Group("/admin")
	{
		adminGroup.POST("/login", controllers.AdminController{}.AdminLogin) // 管理员登录
	}
	userGroup := r.Group("/user")
	{
		userGroup.GET("/validate-token", controllers.UserController{}.ValidateToken) // 验证token
		userGroup.POST("/create", controllers.UserController{}.CreateUser)           // 创建用户
		userGroup.POST("/login", controllers.UserController{}.UserLogin)             // 用户登录
		userGroup.GET("/count", controllers.UserController{}.GetNumber)              // 用户数量
		userGroup.GET("/search", controllers.UserController{}.FuzzyQuery)            // 模糊查询用户
	}
	chatGroup := r.Group("/chat")
	{
		chatGroup.GET("/record", controllers.ChatController{}.GetChatRecord)   // 获取聊天记录
		chatGroup.GET("/unread", controllers.ChatController{}.GetUnReadRecord) // 获取未读消息
	}
	//#####################################################################
	// -----------------------以下为JWT验证相关代码---------------------------

	r.Use(jwt.Auth())
	{
		userGroup.POST("/refresh", controllers.UserController{}.GetAccessToken)                    // 刷新token
		userGroup.POST("/update", controllers.UserController{}.UpdateUser)                         // 更新用户信息
		userGroup.POST("/update-avatar", controllers.UserController{}.UpdateAvatar)                // 更改头像
		userGroup.POST("/update-password", controllers.UserController{}.UpdatePassword)            // 更改密码
		userGroup.GET("/friend-list", controllers.UserController{}.GetAllFriends)                  // 获取所有好友
		userGroup.POST("/friend-request", controllers.UserController{}.FirendRequest)              // 添加好友，发送好友请求
		userGroup.POST("/handle-friend-request", controllers.UserController{}.HandleFriendRequest) // 处理新的好友请求，拒绝或者接受
		userGroup.GET("/friend-request-list", controllers.UserController{}.GetFriendRequestList)   // 获取新的好友请求
		userGroup.GET("/:id", controllers.UserController{}.GetUserInfo)                            // 根据id获取用户信息
		adminGroup.GET("/user-list", controllers.AdminController{}.GetUserList)                    // 获取用户列表
	}
	problem := r.Group("/problem")
	{
		problem.POST("/create", controllers.ProblemController{}.CreateProblem)                  // 创建题目
		problem.GET("/list", controllers.ProblemController{}.GetList)                           // 所有题目列表
		problem.POST("/update", controllers.ProblemController{}.UpdateProblem)                  // 修改题目
		problem.GET("/search", controllers.ProblemController{}.SearchProblem)                   // 根据id获取题目详情
		problem.POST("/submit", controllers.ProblemController{}.SubmitProblem)                  // 提交代码
		problem.GET("/count", controllers.ProblemController{}.GetNumber)                        // 所有题目数量
		problem.GET("/:id", controllers.ProblemController{}.GetProblemInfo)                     // 根据id获取题目详情
		problem.GET("/submission/:id", controllers.ProblemController{}.GetSubmissionStatus)     // 查询提交状态
		problem.GET("/judge-queue/status", controllers.ProblemController{}.GetJudgeQueueStatus) // 查询判题队列状态
	}
	record := r.Group("/record")
	{
		record.GET("/list", controllers.RecodeController{}.GetList)                 // 所有提交记录
		record.GET("/user/:id", controllers.RecodeController{}.GetRecodeListByUser) // 个人的提交记录
		record.GET("/:id", controllers.RecodeController{}.GetRecodeInfo)            // 根据id获取提交记录详情
	}
	ide := r.Group("/ide")
	{
		ide.POST("/submit", controllers.IDEController{}.JudgeCode) // 提交运行代码
	}
	blog := r.Group("/blog")
	{
		blog.POST("/create", controllers.BlogController{}.CreateBlog)                  // 创建博客
		blog.GET("/delete", controllers.BlogController{}.DeleteBlog)                   // 删除博客
		blog.POST("/update", controllers.BlogController{}.UpdateBlog)                  // 修改博客
		blog.GET("/count", controllers.BlogController{}.GetNumber)                     // 所有博客数量
		blog.GET("/available-list", controllers.BlogController{}.GetAvailableBlogList) // 所有可见的博客列表
		blog.GET("/full-list", controllers.BlogController{}.GetFullList)               // 所有博客列表
		blog.GET("/personal-list", controllers.BlogController{}.GetUserBlogList)       // 用户个人的所有博客列表
		blog.GET("/recycle-list", controllers.BlogController{}.RecycleBlog)            // 所有已经标记删除的博客
		blog.GET("/verify-list", controllers.BlogController{}.GetVerifyList)           // 待审核博客列表
		blog.GET("/:id", controllers.BlogController{}.GetBlogInfo)                     // 根据id获取博客详情
	}

	training := r.Group("/training")
	{
		training.GET("/list", controllers.TrainingController{}.GetList)        // 所有训练(题单)列表
		training.GET("/:id", controllers.TrainingController{}.GetTrainingInfo) // 根据id获取训练(题单)详情
	}
	contest := r.Group("/contest")
	{
		contest.GET("/list", controllers.ContestController{}.GetList)       // 所有比赛列表
		contest.GET("/:id", controllers.ContestController{}.GetContestInfo) // 根据id获取比赛题目
	}
	file := r.Group("/file")
	{
		file.POST("/upload", controllers.FileUploadController{}.UploadFile)            // 文件上传
		file.GET("/delete", controllers.FileUploadController{}.DeleteFile)             // 文件删除
		file.GET("/get-share", controllers.FileUploadController{}.GetShareFile)        // 获取分享的文件
		file.POST("/create-share", controllers.FileUploadController{}.CreateShareFile) // 创建分享的文件
	}
	fileUpload := r.Group("/upload")
	{
		fileUpload.POST("/merge", controllers.FileUploadController{}.MergeFileChunk)         // 合并文件块
		fileUpload.POST("/directory", controllers.FileUploadController{}.UploadFolder)       // 上传文件夹
		fileUpload.POST("/markdown-image", controllers.FileUploadController{}.MarkdownImage) // 上传markdown图片

	}
	fileDownload := r.Group("/download")
	{
		fileDownload.GET("/test", controllers.FileDownloadController{}.GetTest)   // 测试下载接口
		fileDownload.GET("/chunk", controllers.FileDownloadController{}.GetChunk) // 分块下载接口
	}
	r.GET("/cloud-storage", controllers.FileUploadController{}.GetDirStruct) // 获取云存储目录结构

	return r
}
