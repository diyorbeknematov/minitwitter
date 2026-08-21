package api

import (
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/api/handler"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/api/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	// middleware
	r.Use(middleware.CORSMiddleware())

	r.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	RegisterRoutes(r, h)

	return r
}

func RegisterRoutes(r *gin.Engine, h *handler.Handler) {

	api := r.Group("/api")

	// Auth
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)
	api.POST("/refresh", h.Refresh)
	api.POST("/logout", h.Logout)

	// Users
	users := api.Group("/users")
	users.GET("/profile", h.GetUserProfile)
	users.PUT("/profile", h.UpdateUserProfile)
	users.GET("/search", h.GetSearchUsers)
	users.GET("/:user_id", h.GetUserByID)
	users.POST("/:user_id/follow", h.UserFollow)
	users.DELETE("/:user_id/follow", h.UserUnfollow)
	users.GET("/followers", h.GetUserFollowers)
	users.GET("/following", h.GetUserFollowing)

	// Tweets
	tweets := api.Group("/tweets")
	tweets.POST("", h.CreateTweet)
	tweets.GET("/:tweet_id", h.GetTweetByID)
	tweets.PUT("/:tweet_id", h.UpdateTweet)
	tweets.DELETE("/:tweet_id", h.DeleteTweet)
	tweets.GET("/", h.GetTweetsByUser)
	tweets.GET("/timeline", h.GetTimeline)

	// Likes
	likes := api.Group("/likes")
	likes.POST("/", h.LikeTweet)
	likes.DELETE("/:tweet_id", h.UnlikeTweet)

	// Retweets
	retweets := api.Group("/retweets")
	retweets.POST("/", h.RetweetTweet)
	retweets.DELETE("/:tweet_id", h.UndoRetweetTweet)

	// Media
	media := api.Group("/media")
	media.POST("/upload", h.UploadMedia)
	media.GET("/:media_id", h.GetMedia)
	media.GET("/", h.GetMedias)
	media.DELETE("/:media_id", h.DeleteMedia)

	// Notifications
	notifications := api.Group("/notifications")
	notifications.GET("", h.GetNotifications)
	notifications.GET("/:notif_id", h.GetNotification)
	notifications.GET("/unreadcount", h.GetUnreadCount)
	notifications.POST("/:notif_id", h.MarkAsRead)
	notifications.POST("/", h.MarkAllAsRead)

	// WebSocket
	api.GET("/ws/notifications", h.NotificationWS)
}
