package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"social-app/internal/config"
	"social-app/internal/domains/auth"
	"social-app/internal/domains/chat"
	"social-app/internal/domains/comment"
	"social-app/internal/domains/like"
	"social-app/internal/domains/llm"
	"social-app/internal/domains/media"
	"social-app/internal/domains/notification"
	"social-app/internal/domains/post"
	"social-app/internal/domains/profile"
	"social-app/internal/models"
	"social-app/pkg/middleware"
	"social-app/pkg/ws"
)

type Router struct {
	profileHandler      profile.Handler
	chatHandler         chat.Handler
	postHandler         post.Handler
	notificationHandler notification.Handler
	commentHandler      comment.Handler
	likeHandler         like.Handler
	mediaHandler        media.Handler
	wsHandler           ws.Handler
	llmHandler          llm.Handler
	authHandler         auth.Handler
	cfg                 config.Auth
}

func NewRouter(
	p post.Handler,
	c comment.Handler,
	pf profile.Handler,
	ch chat.Handler,
	lh like.Handler,
	nh notification.Handler,
	ah auth.Handler,
	mh media.Handler,
	cfg config.Auth,
	wsh ws.Handler,
	lmh llm.Handler,
) Router {
	return Router{
		postHandler:         p,
		commentHandler:      c,
		profileHandler:      pf,
		chatHandler:         ch,
		likeHandler:         lh,
		notificationHandler: nh,
		authHandler:         ah,
		mediaHandler:        mh,
		wsHandler:           wsh,
		llmHandler:          lmh,
		cfg:                 cfg,
	}
}

func (rt Router) SetupRoutes(r *gin.Engine) {
	// cors
	r.Use(middleware.Api(), middleware.LogMiddleware(), gin.Recovery(), middleware.CORSMiddleware())
	pub := r.Group("/")
	rt.publicRoutes(pub)

	authorized := pub.Group("/")
	authorized.Use(middleware.Auth(rt.cfg))
	rt.profileRoutes(authorized)
	rt.mediaRoutes(authorized)
	rt.authRoutes(authorized)
	rt.postRoutes(authorized)
	rt.chatRoutes(authorized)
	rt.notifRoutes(authorized)
	rt.wsRoutes(authorized)
}

func (rt Router) wsRoutes(authorized *gin.RouterGroup) gin.IRoutes {
	return authorized.GET("/ws", middleware.Verified(rt.wsHandler.HandleWebSocket))
}

func (rt Router) mediaRoutes(authorized *gin.RouterGroup) gin.IRoutes {
	return authorized.POST("/upload", middleware.Verified(rt.mediaHandler.UploadImage))
}

func (rt Router) profileRoutes(authorized *gin.RouterGroup) {
	authorized.GET("/users/:id", middleware.Verified(rt.profileHandler.GetProfile))
	authorized.PATCH("/users/:id", middleware.Verified(rt.profileHandler.UpdateProfile))
}

func (rt Router) publicRoutes(pub *gin.RouterGroup) {
	pub.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	pub.Static("/uploads", "./uploads")

	pub.POST("/register", rt.authHandler.Register)
	pub.POST("/login", rt.authHandler.Login)
	pub.GET("/oauth/google/login", rt.authHandler.OauthLogin)
	pub.GET("/oauth/google/callback", rt.authHandler.OauthCallback)
	pub.POST("/profile/verify", middleware.Profile(rt.authHandler.Verify))
	pub.POST("/profile/request-code", middleware.Profile(rt.authHandler.RequestVerification))

	pub.POST("/llm/chat", middleware.Profile(rt.llmHandler.Prompt))
	pub.GET("/admin/users", func(c *gin.Context) {
		var users []models.User
		c.JSON(200, users)
	})
}

func (rt Router) notifRoutes(authorized *gin.RouterGroup) {
	notif := authorized.Group("/notifications")
	notif.GET("", middleware.Verified(rt.notificationHandler.GetNotifications))
	notif.POST("/read", middleware.Verified(rt.notificationHandler.MarkRead))
	notif.POST("/all-read", middleware.Verified(rt.notificationHandler.MarkAllRead))
}

func (rt Router) postRoutes(authorized *gin.RouterGroup) {
	profileByID := authorized.Group("/profile")
	profileByID.GET("/:id/posts", middleware.Verified(rt.postHandler.GetPosts))
	posts := authorized.Group("/posts")
	posts.POST("", middleware.Verified(rt.postHandler.CreatePost))
	posts.GET("", middleware.Verified(rt.postHandler.GetPosts))
	posts.GET("/like/counts", middleware.Verified(rt.likeHandler.GetStats))
	posts.GET("/comments/counts", middleware.Verified(rt.commentHandler.CountByPosts))
	postByID := posts.Group("/:id")
	postByID.POST("/comments", middleware.Verified(rt.commentHandler.CreateComment))
	postByID.GET("/comments", middleware.Verified(rt.commentHandler.GetComments))
	postByID.POST("/like", middleware.Verified(rt.likeHandler.LikePost))
	postByID.DELETE("/like", middleware.Verified(rt.likeHandler.UnlikePost))
	postByID.GET("/likes", middleware.Verified(rt.likeHandler.GetLikes))
}

func (rt Router) chatRoutes(authorized *gin.RouterGroup) {
	authorized.GET("/chats", middleware.Verified(rt.chatHandler.GetChatList))
	authorized.POST("/messages", middleware.Verified(rt.chatHandler.CreatMessage))
	authorized.GET("/messages/:id", middleware.Verified(rt.chatHandler.GetMessages))
	authorized.POST("/messages/:id/read", middleware.Verified(rt.chatHandler.MarkMessagesAsRead))
}

func (rt Router) authRoutes(authorized *gin.RouterGroup) {
	authorized.POST("/refresh", middleware.Verified(rt.authHandler.RefreshToken))
	authorized.POST("/logout", middleware.Verified(rt.authHandler.Logout))
	authorized.GET("/users/:id/online", middleware.Verified(rt.authHandler.IsUserOnline))
}
