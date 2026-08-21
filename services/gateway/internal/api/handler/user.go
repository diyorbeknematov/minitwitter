package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Description Get User
// @Summary Get User by ID
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} dto.User
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router /{user_id} [get]
// @Security BearerAuth
func (h *Handler) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr := c.Param("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.WarnContext(ctx,
			"invalid user id",
			"user_id", userIDStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	resp, err := h.service.User.GetUserByID(ctx, userID)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"get user by id request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description Get User profile
// @Summary Get User profile
// @Tags User
// @Accept json
// @Produce json
// @Succes 200 {object} dto.User
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router /profile [get]
// @Security BearerAuth
func (h *Handler) GetUserProfile(c *gin.Context) {
	ctx := c.Request.Context()

	username, ok := getUsername(c)
	if !ok {
		h.logger.WarnContext(ctx,
			"unauthorized",
		)
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp, err := h.service.User.GetProfile(ctx, username)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"get user profile request failed",
			"error", err,
		)
		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description Update User profile
// @Summary Update User profile
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.UpdateProfileReq true "Update profile"
// @Success 200 {object} dto.User
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router /profile [put]
// @Security BearerAuth
func (h *Handler) UpdateUserProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		h.logger.WarnContext(ctx,
			"unauthorized",
		)
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var in dto.UpdateProfileReq
	if err := c.ShouldBindJSON(&in); err != nil {
		h.logger.WarnContext(ctx,
			"invalid update user profiel request",
			"error", err,
		)
		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	resp, err := h.service.User.UpdateProfile(ctx, userID, in)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"update user profile request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description User follow
// @Summary User follow
// @Tags User
// @Accept json
// @Produce json
// @Param follow_id path string true "User ID to follow"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400,401,404,500 dto.ErrorResponse
// @Router /{follow_id}/follow [post]
// @Security BearerAuth
func (h *Handler) UserFollow(c *gin.Context) {
	ctx := c.Request.Context()

	// Kim follow qilyapti?
	userID, ok := getUserID(c)
	if !ok {
		h.logger.WarnContext(ctx, "unauthorized")

		errorResponse(
			c,
			http.StatusUnauthorized,
			"unauthorized",
		)
		return
	}

	// Kimni follow qilyapti?
	followIDStr := c.Param("follow_id")

	followID, err := uuid.Parse(followIDStr)
	if err != nil {
		h.logger.WarnContext(ctx,
			"invalid follow user id",
			"follow_id", followIDStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	err = h.service.User.Follow(ctx, userID, followID)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"failed to follow user",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "user follow successfully",
	})
}

// @Description User unfollow
// @Summary User unfollow
// @Tags User
// @Accept json
// @Produce json
// @Param following_id path string true "User ID to unfollow"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400,401,404,500 dto.ErrorResponse
// @Router /{following_id}/unfollow [delete]
// @Security BearerAuth
func (h *Handler) UserUnfollow(c *gin.Context) {
	ctx := c.Request.Context()

	// Kim unfollow qilyapti?
	userID, ok := getUserID(c)
	if !ok {
		h.logger.WarnContext(ctx, "unauthorized")

		errorResponse(
			c,
			http.StatusUnauthorized,
			"unauthorized",
		)
		return
	}

	// Kimni unfollow qilyapti?
	followingID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	err = h.service.User.Unfollow(ctx, userID, followingID)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"unfollow user request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "user unfollow successfully",
	})
}

// @Description Get User Followers
// @Summary Get User Followers
// @Tags User
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number per page" default(20)
// @Success 200 {object} dto.UserResp
// @Failure 400,401,404,500 dto.ErrorResponse
// @Router /followers [get]
// @Security BearerAuth
func (h *Handler) GetUserFollowers(c *gin.Context) {
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

	resp, err := h.service.User.GetFollowers(ctx, userID, int32(page), int32(limit))
	if err != nil {
		h.logger.ErrorContext(ctx,
			"get followers request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description Get User Following
// @Summary Get User Following
// @Tags User
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number per page" default(20)
// @Success 200 {object} dto.UserResp
// @Failure 400,401,404,500 dto.ErrorResponse
// @Router /following [get]
// @Security BearerAuth
func (h *Handler) UserUserFollowing(c *gin.Context) {
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

	resp, err := h.service.User.GetFollowing(ctx, userID, int32(page), int32(limit))
	if err != nil {
		h.logger.ErrorContext(ctx,
			"get user following request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description Search users
// @Summary Search users
// @Tags User
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number per page" default(20)
// @Success 200 {object} dto.UserListResponse
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router /search [get]
// @Security BearerAuth
func (h *Handler) GetSearchUsers(c *gin.Context) {
	ctx := c.Request.Context()

	query := strings.TrimSpace(c.Query("q"))

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

	resp, err := h.service.User.SearchUsers(ctx, dto.SearchUsersQuery{
		Query: query,
		Page:  int32(page),
		Limit: int32(limit),
	})

	if err != nil {
		h.logger.ErrorContext(ctx,
			"failed to search users",
			"query", query,
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}
