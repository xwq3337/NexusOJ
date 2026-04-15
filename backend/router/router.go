package router

import (
	"net/http"
	"nexus/config"
	"nexus/controllers"
	"nexus/microapps"
	"nexus/middleware"
	"nexus/middleware/jwt"
	"nexus/utils"
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
	// 请求日志
	r.Use(middleware.RequestLogger())
	microapps.MicroRouter(r) // 微应用路由
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
	adminUserGroup := r.Group("/admin/user")
	{
		adminUserGroup.POST("/login", controllers.AdminController{}.AdminLogin) // 管理员登录
	}
	userGroup := r.Group("/user")
	{
		userGroup.GET("/validate-token", controllers.UserController{}.ValidateToken) // 验证token
		userGroup.POST("/create", controllers.UserController{}.CreateUser)           // 创建用户
		userGroup.POST("/login", controllers.UserController{}.UserLogin)             // 用户登录
		userGroup.GET("/count", controllers.UserController{}.GetNumber)              // 用户数量
		userGroup.GET("/search", controllers.UserController{}.FuzzyQuery)            // 模糊查询用户
		userGroup.GET("/top-rating", controllers.UserController{}.GetTopUsers)       // Rating排行榜Top10
	}
	chatGroup := r.Group("/chat")
	{
		chatGroup.GET("/record", controllers.ChatController{}.GetChatRecord)          // 获取聊天记录
		chatGroup.GET("/unread", controllers.ChatController{}.GetUnReadRecord)        // 获取未读消息
		chatGroup.POST("/mark-read", controllers.ChatController{}.MarkMessagesAsRead) // 标记消息已读
	}
	userGroup.GET("/refresh", controllers.UserController{}.RefreshToken) // 刷新token
	// 把mysql 数据同步到 redis 里
	r.GET("/refresh-redis", controllers.RedisCache{}.RefreshRedisCache) // 刷新redis缓存
	//#####################################################################
	// -----------------------以下为JWT验证相关代码---------------------------
	r.Use(jwt.Auth())

	{
		userGroup.POST("/update", controllers.UserController{}.UpdateUser)                                   // 更新用户信息
		userGroup.POST("/update-avatar", controllers.UserController{}.UpdateAvatar)                          // 更改头像
		userGroup.POST("/update-password", controllers.UserController{}.UpdatePassword)                      // 更改密码
		userGroup.GET("/friend-list", controllers.FriendshipController{}.GetFriendList)                      // 获取所有好友
		userGroup.GET("/friend-request-list", controllers.FriendshipController{}.GetFriendRequestList)       // 获取新的好友请求
		userGroup.POST("/friend-request", controllers.FriendshipController{}.CreateFriendship)               // 添加好友，发送好友请求
		userGroup.POST("/handle-friend-request", controllers.FriendshipController{}.HandleFriendshipRequest) // 处理新的好友请求，拒绝或者接受
		userGroup.GET("/:id", controllers.UserController{}.GetUserInfo)                                      // 根据id获取用户信息
		userGroup.GET("/homepage/:id", controllers.UserController{}.GetUserHomePage)                         // 根据id获取用户主页信息
		adminUserGroup.GET("/list", controllers.AdminController{}.GetUserList)                               // 获取用户列表
		adminUserGroup.POST("/update-role", controllers.UserController{}.UpdateRole)
	}
	problem := r.Group("/problem")
	{
		problem.POST("/create", controllers.ProblemController{}.CreateProblem) // 创建题目
		problem.GET("/list", controllers.ProblemController{}.GetList)          // 所有题目列表
		problem.POST("/update", controllers.ProblemController{}.UpdateProblem) // 修改题目
		problem.GET("/search", controllers.ProblemController{}.SearchProblem)  // 根据id获取题目详情
		problem.POST("/submit", controllers.ProblemController{}.SubmitProblem) // 提交代码
		problem.GET("/count", controllers.ProblemController{}.GetNumber)       // 所有题目数量
		problem.GET("/:id", controllers.ProblemController{}.GetProblemInfo)    // 根据id获取题目详情
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
		contest.GET("/list", controllers.ContestController{}.GetList)                                    // 所有比赛列表
		contest.GET("/:id", controllers.ContestController{}.GetContestInfo)                              // 比赛详情
		contest.POST("/register", controllers.ContestController{}.RegisterContest)                       // 报名比赛
		contest.GET("/:id/problems", controllers.ContestController{}.GetContestProblems)                 // 比赛题目
		contest.GET("/:id/problems/:label", controllers.ContestController{}.GetContestProblemDetail)     // 比赛单题详情
		contest.POST("/:id/submit", controllers.ContestController{}.SubmitContestProblem)                // 提交代码
		contest.GET("/:id/submissions", controllers.ContestController{}.GetContestSubmissions)           // 比赛提交列表
		contest.GET("/:id/submissions/:rid", controllers.ContestController{}.GetContestSubmissionDetail) // 提交详情
		contest.GET("/:id/ranking", controllers.ContestController{}.GetContestRanking)                   // 实时排名
		contest.GET("/:id/ranking/stream", controllers.ContestController{}.StreamContestRanking)         // SSE排名推送
		contest.GET("/:id/my-status", controllers.ContestController{}.GetMyContestStatus)                // 用户参赛状态
	}

	solution := r.Group("/solution")
	{
		solution.GET("/list", controllers.SolutionController{}.GetSolutions)          // 题解列表
		solution.GET("/:id", controllers.SolutionController{}.GetSolutionDetail)      // 题解详情
		solution.POST("/create", controllers.SolutionController{}.CreateSolution)     // 创建题解
		solution.POST("/update/:id", controllers.SolutionController{}.UpdateSolution) // 更新题解
		solution.POST("/delete/:id", controllers.SolutionController{}.DeleteSolution) // 删除题解
	}

	adminContest := r.Group("/admin/contest")
	{
		adminContest.POST("/create", controllers.ContestController{}.CreateContest)                      // 创建比赛
		adminContest.POST("/update", controllers.ContestController{}.UpdateContest)                      // 更新比赛
		adminContest.POST("/delete", controllers.ContestController{}.DeleteContest)                      // 删除比赛
		adminContest.GET("/list", controllers.ContestController{}.GetAdminContestList)                   // 管理员比赛列表
		adminContest.GET("/:id", controllers.ContestController{}.GetAdminContestDetail)                  // 管理员比赛详情
		adminContest.POST("/:id/problems", controllers.ContestController{}.SetContestProblems)           // 设置比赛题目
		adminContest.POST("/:id/participants", controllers.ContestController{}.ManageParticipant)        // 管理参赛者
		adminContest.GET("/:id/report", controllers.ContestController{}.GenerateReport)                  // 生成比赛报告
		adminContest.GET("/:id/import-preview", controllers.ContestController{}.GetImportPreview)        // 预览可导入题目
		adminContest.POST("/:id/import-problems", controllers.ContestController{}.ImportContestProblems) // 导入题目到题库
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
	logGroup := r.Group("/log")
	{
		logGroup.GET("/date", controllers.LogController{}.GetDate) // 获取所有日志日期
		logGroup.GET("/list", controllers.LogController{}.GetList) // 获取某一天的日志列表
	}

	recommend := r.Group("/recommend")
	{
		recommend.GET("/profile", controllers.RecommendationController{}.GetMyProfile)            // 当前用户画像
		recommend.GET("/profile/:id", controllers.RecommendationController{}.GetUserProfile)      // 其他用户画像
		recommend.GET("/ability", controllers.RecommendationController{}.GetAbilityAnalysis)       // 能力分析
		recommend.GET("/activity", controllers.RecommendationController{}.GetActivityStats)        // 活跃度统计
		recommend.GET("/problems", controllers.RecommendationController{}.GetRecommendations)      // 推荐题目
		recommend.POST("/refresh", controllers.RecommendationController{}.RefreshRecommendations)  // 刷新推荐
	}

	return r
}
