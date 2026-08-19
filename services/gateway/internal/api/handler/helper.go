package handler

import (
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/api/middleware"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getUserID(c *gin.Context) (uuid.UUID, bool) {
	value, exists := c.Get(middleware.UserIDKey)
	if !exists {
		return uuid.Nil, false
	}

	userID, ok := value.(uuid.UUID)
	return userID, ok
}

func errorResponse(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, dto.ErrorResponse{
		Message: message,
	})
}
