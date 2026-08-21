package handler

import (
	"net/http"
	"strconv"

	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// @Description Get Tweet
// @Summary Get Tweet by ID
// @Tags Tweet
// @Accept json
// @Produce json
// @Param tweet_id path string true "Tweet ID"
// @Success 200 {object} dto.Tweet
// @Failure 400,401,404,500 dto.ErrorResponse
// @Router /{tweet_id} [get]
// @Security BearerAuth
func (h *Handler) GetTweetByID(c *gin.Context) {
	ctx := c.Request.Context()

	tweetIDStr := c.Param("tweet_id")

	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		h.logger.WarnContext(ctx,
			"invalid tweet id",
			"tweet_id", tweetIDStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid tweet id")
		return
	}

	resp, err := h.service.Tweet.GetTweet(ctx, tweetID)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"get tweet request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description Update Tweet
// @Summary Update Tweet
// @Tags Tweet
// @Accept json
// @Param tweet_id path string true "Tweet ID"
// @Param tweet body dto.UpdateTweetReq true "Update Tweet"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router /{tweet_id}update [put]
// @Security BearerAuth
func (h *Handler) UpdateTweet(c *gin.Context) {
	ctx := c.Request.Context()

	tweetIDStr := c.Param("tweet_id")

	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		h.logger.WarnContext(ctx,
			"invalid tweet id",
			"tweet_id", tweetIDStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid tweet id")
		return
	}

	var in dto.UpdateTweetReq
	if err := c.ShouldBindJSON(&in); err != nil {
		h.logger.WarnContext(ctx,
			"invalid update tweet request",
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	err = h.service.Tweet.UpdateTweet(ctx, tweetID, in)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"update tweet request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Tweet updated successfully",
	})
}

// @Description Delete Tweet
// @Summary Delete Tweet
// @Tags Tweet
// @Accept json
// @Produce json
// @Param tweet_id path string true "Tweet ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400,401,404,500 dto.ErrorResponse
// @Router /{tweet_id} [delete]
// @Security BearerAuth
func (h *Handler) DeleteTweet(c *gin.Context) {
	ctx := c.Request.Context()

	tweetIDStr := c.Param("tweet_id")

	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		h.logger.WarnContext(ctx,
			"invalid tweet id",
			"tweet_id", tweetIDStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid tweet id")
		return
	}

	err = h.service.Tweet.DeleteTweet(ctx, tweetID)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"delete tweet request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "tweet deleted successfully",
	})
}

// @Description Get Tweets
// @Description Get Tweets by user
// @Tags Tweet
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number per page" default(20)
// @Success 200 {object} dto.GetTweetByUserResp
// @Failure 400,401,404,500 {object} dto.ErrorRespnse
// @Router / [get]
// @Security BearerAuth
func (h *Handler) GetTweetsByUser(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		h.logger.WarnContext(ctx, "unauthorized")
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		h.logger.WarnContext(ctx,
			"invalid page parameter",
			"page", pageStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid page")
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		h.logger.WarnContext(ctx,
			"invalid limit parameter",
			"limit", limitStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.service.Tweet.GetTweetsByUser(ctx,
		userID,
		dto.GetTweetByUserReq{
			Page:  int32(page),
			Limit: int32(limit),
		},
	)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"get tweets by user request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description Get Timeline
// @Summary Get Timeline
// @Tags Tweet
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number per page" default(20)
// @Success 200 {object} dto.GetTimelineResp
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router /timeline [get]
// @Security BearerAuth
func (h *Handler) GetTimeline(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		h.logger.WarnContext(ctx, "unauthorized")

		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		h.logger.WarnContext(ctx,
			"invalid page parameter",
			"page", pageStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid page")
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		h.logger.WarnContext(ctx,
			"invalid limit parameter",
			"limit", limitStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.service.Tweet.GetTimeline(ctx, userID, int32(page), int32(limit))
	if err != nil {
		h.logger.ErrorContext(ctx,
			"get timeiline request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description Like Tweet
// @Summary Like Tweet
// @Tags Like
// @Accept json
// @Produce json
// @Param tweet_id path string true "Tweet ID"
// @Success 200 {objcet} dto.SuccessResponse
// @Failure 400,401,404,405 {object} dto.ErrorResponse
// @Router /like [post]
// @Security BearerAuth
func (h *Handler) LikeTweet(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		h.logger.WarnContext(ctx, "unauthorized")

		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	tweetIDStr := c.Param("tweet_id")

	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		h.logger.WarnContext(ctx,
			"invalid tweet id",
			"tweet_id", tweetIDStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid tweet id")
		return
	}

	err = h.service.Tweet.LikeTweet(ctx, userID, tweetID)
	if err != nil {
		h.logger.ErrorContext(ctx, 
			"like tweet request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "tweet like successfully",
	})
}

// @Description Unlike Tweet
// @Summary unlike Tweet
// @Tags Like
// @Accept json
// @Produce json
// @Param tweet_id path string true "Tweet ID"
// @Success 200 {objcet} dto.SuccessResponse
// @Failure 400,401,404,405 {object} dto.ErrorResponse
// @Router /like [delete]
// @Security BearerAuth
func (h *Handler) UnlikeTweet(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		h.logger.WarnContext(ctx, "unauthorized")

		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	tweetIDStr := c.Param("tweet_id")

	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		h.logger.WarnContext(ctx,
			"invalid tweet id",
			"tweet_id", tweetIDStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid tweet id")
		return
	}

	err = h.service.Tweet.UnlikeTweet(ctx, userID, tweetID)
	if err != nil {
		h.logger.ErrorContext(ctx, 
			"unlike tweet request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "tweet unlike successfully",
	})
}

// @Description Retweet Tweet 
// @Summary Retweet Tweet 
// @Tags Retweet 
// @Accept json 
// @Produce json 
// @Param tweet_id path string true "Tweet ID"
// @Success 200 {objcet} dto.SuccessResponse
// @Failure 400,401,404,405 {object} dto.ErrorResponse
// @Router /retweet [post]
// @Security BearerAuth
func (h *Handler) RetweetTweet(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		h.logger.WarnContext(ctx, "unauthorized")

		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	tweetIDStr := c.Param("tweet_id")

	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		h.logger.WarnContext(ctx,
			"invalid tweet id",
			"tweet_id", tweetIDStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid tweet id")
		return
	}

	err = h.service.Tweet.Retweet(ctx, userID, tweetID)
	if err != nil {
		h.logger.ErrorContext(ctx, 
			"retweet request failed", 
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "tweet retweet successfully",
	})
}

// @Description Retweet Tweet 
// @Summary Retweet Tweet 
// @Tags Retweet 
// @Accept json 
// @Produce json 
// @Param tweet_id path string true "Tweet ID"
// @Success 200 {objcet} dto.SuccessResponse
// @Failure 400,401,404,405 {object} dto.ErrorResponse
// @Router /retweet [delete]
// @Security BearerAuth
func (h *Handler) UndoRetweetTweet(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		h.logger.WarnContext(ctx, "unauthorized")

		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	tweetIDStr := c.Param("tweet_id")

	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		h.logger.WarnContext(ctx,
			"invalid tweet id",
			"tweet_id", tweetIDStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid tweet id")
		return
	}

	err = h.service.Tweet.Retweet(ctx, userID, tweetID)
	if err != nil {
		h.logger.ErrorContext(ctx, 
			"undo retweet request failed", 
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "undo retweet successfully",
	})
}
