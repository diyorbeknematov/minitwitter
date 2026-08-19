package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/security"
	"github.com/diyorbeknematov/minitwitter/services/gateway/pkg/traceid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	UserIDKey   = "user-id"
	UsernameKey = "username"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := security.ValidateToken(tokenString, cfg.AccessTokenSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			return
		}

		ctx.Set(UserIDKey, userID)
		ctx.Set(UsernameKey, claims.Username)

		ctx.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow localhost development URLs
		allowedOrigins := map[string]bool{}

		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func TraceMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := traceid.NewContext(c.Request.Context())
		
		c.Request = c.Request.WithContext(ctx)

		logger.InfoContext(ctx, "so'rov qabul qilindi",
			slog.String("trace_id", traceid.FromContext(ctx)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		)

		c.Next()
	}
}
