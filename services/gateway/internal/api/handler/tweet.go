package handler

import (
	"net/http"

	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/gin-gonic/gin"
)

// @Description Create Tweet
// @Summary Createa Tweet
// @Tags Tweet
// @Accept json 
// @Produce json 
// @Param tweet body dto.CreateTweetReq true "Tweet"
// @Success 201 {object} dto.Tweet
// @Success 400,401,404,500 {object} dto.ErrorResponse
// @Router /tweet [post]
// @Security BearerAuth
func (h *Handler) CreateTweet(c *gin.Context) {
	ctx := c.Request.Context()

	var in dto.CreateTweetReq
	if err := c.ShouldBindBodyWithJSON(&in); err != nil {
		h.logger.WarnContext(ctx, 
			"invalid create tweet request",
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid request input")
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		h.logger.WarnContext(ctx,
			"unauthorized",
		)

		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	} 

	resp, err := h.service.Tweet.CreateTweet(ctx, userID, in)
	if err != nil {
		h.logger.ErrorContext(ctx, 
			"failed to create tweet",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusCreated, resp)
}