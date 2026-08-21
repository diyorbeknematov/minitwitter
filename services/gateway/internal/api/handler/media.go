package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Description Upload Media
// @Summary Upload Media
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Media file"
// @Param folder formData string false "Folder name"
// @Param category formData string true "Media category" Enums(avatar,tweet)
// @Succes 200 {object} dto.Media
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router /upload [post]
// @Security BearerAuth
func (h *Handler) UploadMedia(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		h.logger.WarnContext(ctx, "unauthorized")

		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	category := c.PostForm("category")
	switch category {
	case "avatar", "tweet":
		// OK
	default:
		errorResponse(c, http.StatusBadRequest, "invalid media category")
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		h.logger.WarnContext(ctx,
			"failed to get uploaded file",
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "file is required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.logger.ErrorContext(ctx,
			"failed to open uploaded file",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "failed to process file")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"failed to read uploaded file",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "failed to read file")
		return
	}

	resp, err := h.service.Media.UploadMedia(
		ctx,
		userID,
		data,
		fileHeader.Filename,
		dto.UploadMediaReq{Category: dto.MediaCategory(category)},
	)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"upload media request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description Get Media
// @Summary Get Media
// @Tags Media
// @Accept json
// @Produce json
// @Param media_id path string true "Media ID"
// @Success 200 {object} dto.Media
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router /{media_id} [get]
// @Security BearerAuth
func (h *Handler) GetMedia(c *gin.Context) {
	ctx := c.Request.Context()

	mediaIDStr := c.Param("media_id")

	mediaID, err := uuid.Parse(mediaIDStr)
	if err != nil {
		h.logger.WarnContext(ctx,
			"invalid media id",
			"media_id", mediaIDStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid media id")
		return
	}

	resp, err := h.service.Media.GetMedia(ctx, mediaID)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"get media request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description Get Medias
// @Summary Get Medias
// @Tags Media
// @Accept json
// @Produce json
// @Param ids query string true "Comma-separated media IDs"
// @Success 200 {object} dto.GetMediasResp
// @Failure 400,401,404,500 {object} dto.ErrorResponse
// @Router / [get]
// @Security BearerAuth
func (h *Handler) GetMedias(c *gin.Context) {
	ctx := c.Request.Context()

	idsStr := c.Query("ids")
	parts := strings.Split(idsStr, ",")

	ids := make([]uuid.UUID, 0, len(parts))

	for _, idStr := range parts {
		id, err := uuid.Parse(strings.TrimSpace(idStr))
		if err != nil {
			h.logger.WarnContext(ctx,
				"invalid media ids",
				"media_ids", idStr,
				"error", err,
			)

			errorResponse(c, http.StatusBadRequest, "invalid media ids")
			return
		}

		ids = append(ids, id)
	}

	resp, err := h.service.Media.GetMedias(ctx, ids)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"get medias request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description Delete Media 
// @Summary Delete Media 
// @Tags Media 
// @Accept json 
// @Produce json 
// @Param media_id path string true "Media ID"
// @Success 200 {object} dto.SuccessResponse 
// @Failure 400,401,404,500 {object} dto.ErrorRespnose 
// @Router /{media_id} [delete]
// @Security BearerAuth 
func (h *Handler) DeleteMedia(c *gin.Context) {
	ctx := c.Request.Context()

	mediaIDStr := c.Param("media_id")

	mediaID, err := uuid.Parse(mediaIDStr)
	if err != nil {
		h.logger.WarnContext(ctx,
			"invalid media id",
			"media_id", mediaIDStr,
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid media id")
		return
	}

	err = h.service.Media.DeleteMedia(ctx, mediaID)
	if err != nil {
		h.logger.ErrorContext(ctx, 
			"delete media request failed",
			"error", err,
		)

		errorResponse(c, http.StatusInternalServerError, "unexpected error")
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "media deleted successfully",
	})
}
