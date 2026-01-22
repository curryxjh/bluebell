package router

import (
	"bluebell/controller"
	_ "bluebell/docs"
	"bluebell/logger"
	"bluebell/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"                  // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"net/http"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 设置成开发模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler), middlewares.RateLimitMiddleware(2*time.Second, 1))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册路由
	v1 := r.Group("/api/v1")

	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	// 注意：具体的路由要放在动态路由前面，避免被匹配
	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/list", controller.CommunityListWithMembersHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)
	v1.GET("/posts", controller.GetPostListHandler)
	v1.GET("/posts2", controller.GetPostListHandler2)
	v1.GET("/post/:id", controller.GetPostDetailHandler)
	// 评论相关（公开接口）
	v1.GET("/comments/:id", controller.GetCommentsByPostIDHandler)

	v1.Use(middlewares.JWTAuthMiddleware())

	{
		v1.POST("/post", controller.CreatePostHandler)

		v1.POST("/vote", controller.PostVoteController)

		// 社区相关接口
		v1.POST("/community/join", controller.JoinCommunityHandler)
		v1.POST("/community/leave", controller.LeaveCommunityHandler)
		v1.GET("/community/user", controller.GetUserCommunitiesHandler)

		// 评论相关接口（需要登录）
		v1.POST("/comment", controller.CreateCommentHandler)
		v1.DELETE("/comment/:id", controller.DeleteCommentHandler)
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
