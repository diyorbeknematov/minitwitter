package handler

import (
	"net/http"
	"strconv"

	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Description Get Notification
// @Summary Get Notification
// @Tags Notification
// @Accept json
// @Produce json
// @Param notif_id path string true "Notification ID"
// @Success 200 {object} dto.Notification
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router /{notif_id} [get]
// @Security BearerAuth
func (h *Handler) GetNotification(c *gin.Context) {
	ctx := c.Request.Context()

	notifIDStr := c.Param("notif_id")

	notifID, err := uuid.Parse(notifIDStr)
	if err != nil {
		h.logger.WarnContext(ctx,
			"invalid notification id",
			"notification_id", notifIDStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid notification id")
		return
	}

	resp, err := h.service.Notification.GetNotification(ctx, notifID)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"get notification request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description Get Notifications
// @Summary Get Notifications
// @Tags Notification
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number per page" default(20)
// @Success 200 {object} dto.GetNotificationsResp
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router / [get]
// @Security BearerAuth
func (h *Handler) GetNotifications(c *gin.Context) {
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

	resp, err := h.service.Notification.GetNotifications(
		ctx,
		userID,
		int32(page),
		int32(limit),
	)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"get notifications request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description Get unread count
// @Summary Get unread nofitication count
// @Tags Notification
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetUnreadCountResp
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router /unreadcount [get]
// @Security BearerAuth
func (h *Handler) GetUnreadCount(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		h.logger.WarnContext(ctx, "unauthorized")

		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp, err := h.service.Notification.GetUnreadCount(ctx, userID)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"get unread notification count request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description Delete Notification
// @Summary Delete Notification
// @Tags Notification
// @Accept json
// @Produce json
// @Param notif_id path string true "Notification ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router /{notif_id} [delete]
// @Security BearerAuth
func (h *Handler) DeleteNotification(c *gin.Context) {
	ctx := c.Request.Context()

	notifIDStr := c.Param("notif_id")

	notifID, err := uuid.Parse(notifIDStr)
	if err != nil {
		h.logger.WarnContext(ctx,
			"invalid notification id",
			"notification_id", notifIDStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid notification id")
		return
	}

	err = h.service.Notification.DeleteNotification(ctx, notifID)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"delete notification request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "notification deleted successfully",
	})
}

// @Description Mark as read
// @Summary Mark notification as read
// @Tags Notification
// @Accept json
// @Produce json
// @Param notif_id path string true "Notification ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router /markasread/{notif_id} [post]
// @Security BearerAuth
func (h *Handler) MarkAsRead(c *gin.Context) {
	ctx := c.Request.Context()

	notifIDStr := c.Param("notif_id")

	notifID, err := uuid.Parse(notifIDStr)
	if err != nil {
		h.logger.WarnContext(ctx,
			"invalid notification id",
			"notification_id", notifIDStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid notification id")
		return
	}

	err = h.service.Notification.MarkAsRead(ctx, notifID)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"mark notification as read request failed",
			"error", err,
		)

		errorResponse(c, http.StatusOK, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "notificaton marked as read successfully",
	})
}

// @Description Mark all notification as read
// @Summary Mark all notification as read
// @Tags Notification
// @Accept json
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router /markasread [post]
// @Security BearerAuth
func (h *Handler) MarkAllAsRead(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		h.logger.WarnContext(ctx, "unauthorized")

		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := h.service.Notification.MarkAllAsRead(ctx, userID)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"mark all notification as read request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "all notification marked as read successfully",
	})
}
