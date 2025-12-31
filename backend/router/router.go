package router

import (
	"net/http"
	"pkg/config"
	"pkg/controllers"
	"pkg/middleware"
	"pkg/middleware/jwt"
	jwtgo "pkg/middleware/jwt"
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
			c.JSON(http.StatusOK, gin.H{
				"ip": c.ClientIP(),
			})
		})
	}

	userGroup := r.Group("/user")
	{
		userGroup.POST("/create", controllers.UserController{}.CreateUser)
		userGroup.POST("/login", controllers.UserController{}.UserLogin)
		userGroup.POST("/admin-login", controllers.UserController{}.AdminLogin)
		userGroup.GET("/logout", controllers.UserController{}.UserLogout)
		userGroup.GET("/friend", controllers.UserController{}.GetAllFriends)
		userGroup.GET("/handle", controllers.UserController{}.HandleNewFriend)
		userGroup.POST("/changeAvatar", controllers.UserController{}.ChangeAvatar)
		userGroup.POST("/changePassword", controllers.UserController{}.ChangePassword)
		userGroup.GET("/newfriend", controllers.UserController{}.GetNewFriends)
		userGroup.GET("/search", controllers.UserController{}.FuzzyQuery)
		userGroup.POST("/add", controllers.UserController{}.AddFirend)
	}
	// -----------------------以下为JWT验证相关代码---------------------------
	chatGroup := r.Group("/chat")
	{
		chatGroup.GET("/record", controllers.ChatController{}.GetChatRecord)
		chatGroup.GET("/unread", controllers.ChatController{}.GetUnReadRecord)
	}
	r.Use(jwt.Auth())
	{
		userGroup.POST("/refresh", controllers.UserController{}.GetAccessToken)
		userGroup.GET("/list", controllers.UserController{}.GetList)
		userGroup.GET("/count", controllers.UserController{}.GetNumber)
		userGroup.GET("/:id", controllers.UserController{}.GetUserInfo)
	}
	problem := r.Group("/problem")
	{

		problem.GET("/list", controllers.ProblemController{}.GetList)          // 所有题目列表
		problem.POST("/change", controllers.ProblemController{}.ChangeProblem) // 修改题目
		problem.GET("/search", controllers.ProblemController{}.SearchProblem)  // 根据id获取题目详情
		problem.POST("/create", controllers.ProblemController{}.CreateProblem) // 创建题目
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
		blog.POST("/create", controllers.BlogController{}.CreateBlog)        // 创建博客
		blog.GET("/count", controllers.BlogController{}.GetNumber)           // 所有博客数量
		blog.GET("/all", controllers.BlogController{}.GetList)               // 所有博客列表
		blog.GET("/list", controllers.BlogController{}.GetAvailableBlogList) // 所有可见的博客列表
		blog.GET("/my-list", controllers.BlogController{}.GetUserBlogList)   // 用户个人的所有博客列表
		blog.POST("/change", controllers.BlogController{}.ChangeBlog)        // 修改博客
		blog.GET("/delete", controllers.BlogController{}.DeleteBlog)         // 删除博客
		blog.GET("/recycling", controllers.BlogController{}.RecycleBlog)     // 所有已经标记删除的博客
		blog.GET("/verify-list", controllers.BlogController{}.GetVerifyList)
		blog.GET("/:id", controllers.BlogController{}.GetBlogInfo) // 根据id获取博客详情
	}
	comment := r.Group("/comment") // TODO: 评论相关接口
	{
		comment.GET("", controllers.CommentController{}.GetComment)
		comment.POST("/main", controllers.CommentController{}.CreateComment)
		comment.POST("/sub", controllers.CommentController{}.CreateSubComment)
	}
	training := r.Group("/training")
	{
		training.POST("/create", controllers.TrainingController{}.CreateTraining) // 创建训练(题单)
		training.GET("/list", controllers.TrainingController{}.GetList)           // 所有训练(题单)列表
		training.GET("/:id", controllers.TrainingController{}.GetTrainingInfo)
	}
	contest := r.Group("/contest")
	{
		contest.GET("/list", controllers.ContestController{}.GetList)
		contest.GET("/:id", controllers.ContestController{}.GetContestProblems)
	}
	file := r.Group("/file")
	{
		file.POST("/upload", controllers.FileUploadController{}.UploadFile)
		file.GET("/delete", controllers.FileUploadController{}.DeleteFile)
		file.GET("/get_share", controllers.FileUploadController{}.GetShareFile)
		file.POST("/create_share", controllers.FileUploadController{}.CreateShareFile)
	}
	fileUpload := r.Group("/upload")
	{
		fileUpload.POST("/merge", controllers.FileUploadController{}.MergeFileChunk)
		fileUpload.DELETE("file")
		fileUpload.POST("/directory", controllers.FileUploadController{}.UploadFolder)
		fileUpload.POST("/markdown-image", controllers.FileUploadController{}.MarkdownImage)
		fileUpload.POST("/avatar") // TODO:

	}
	fileDownload := r.Group("/download")
	{
		fileDownload.GET("/test", controllers.FileDownloadController{}.GetTest)
		fileDownload.GET("/chunk", controllers.FileDownloadController{}.GetChunk)
	}
	r.GET("/cloud-storage", controllers.FileUploadController{}.GetDirStruct)

	sv1 := r.Group("/auth/")
	{
		sv1.GET("/time", func(c *gin.Context) {
			claims := c.MustGet("claims").(*jwtgo.CustomClaims)
			if claims != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 0,
					"msg":  "token有效",
					"info": claims,
				})
			}
		})
	}

	return r
}
