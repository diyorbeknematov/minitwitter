package handler

import (
	"errors"
	"net/http"

	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/diyorbeknematov/minitwitter/services/gateway/pkg/apperror"
	"github.com/gin-gonic/gin"
)

// @Description Register User
// @Summary Register User
// @Tags Auth
// @Accept json
// @Produce json
// @Param signup body dto.RegisterReq true "Register"
// @Success 201 {object} dto.RegisterResp
// @Failure 400,404,409,500 {object} dto.ErrorResponse
// @Router /register [post]
func (h *Handler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var body dto.RegisterReq

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		h.logger.WarnContext(ctx,
			"invalid register request",
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Username == "" || body.Email == "" {
		h.logger.WarnContext(ctx,
			"register validation failed",
			"username_empty", body.Username == "",
			"email_empty", body.Email == "",
		)

		errorResponse(c, http.StatusBadRequest, "username or email cannot be empty")
		return
	}

	resp, err := h.service.Auth.Register(ctx, body)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"failed to register user",
			"error", err,
		)

		switch {
		case errors.Is(err, apperror.ErrEmailExists):
			errorResponse(c, http.StatusConflict, "Email already exists")
		case errors.Is(err, apperror.ErrUsernameExists):
			errorResponse(c, http.StatusConflict, "Username already exists")
		default:
			errorResponse(c, http.StatusInternalServerError, "unexpected error")
		}
		return
	}

	h.logger.InfoContext(ctx,
		"user registered successfully",
	)

	c.JSON(http.StatusCreated, resp)
}

// @Description Login User
// @Summary Login User
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body dto.LoginReq true "Login"
// @Succes 200 {object} dto.LoginResp
// @Failure 400,404,409,500 {object} dto.ErrorResponse
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var body dto.LoginReq
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		h.logger.WarnContext(ctx,
			"invalid login request",
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Email == "" || body.Password == "" {
		h.logger.WarnContext(ctx,
			"login validation failed",
			"email_empty", body.Email == "",
			"password_empty", body.Password == "",
		)

		errorResponse(c, http.StatusBadRequest, "Email and password are required")
	}

	resp, err := h.service.Auth.Login(ctx, body)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"failed to login user",
			"error", err,
		)

		switch {
		case errors.Is(err, apperror.ErrUserNotFound):
			errorResponse(c, http.StatusNotFound, "User not found")
		case errors.Is(err, apperror.ErrInvalidInput):
			errorResponse(c, http.StatusUnauthorized, "Invalid username or password")
		default:
			errorResponse(c, http.StatusInternalServerError, "Unexpected error")
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Description Logout User
// @Summary Logout User
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400,404,500 {object} dto.ErrorResponse
// @Router /logout [post]
// @Security BearerAuth
func (h *Handler) Logout(c *gin.Context) {
	// ctx := c.Request.Context()

	// userID, ok := getUserID(c)
	// if !ok {
	// 	h.logger.WarnContext(ctx, 
	// 		"unauthorized",
	// 	)
	// }

	// resp, err := h.service.Auth.Logout(userID)

	// c.JSON(http.StatusOK, dto.SuccessResponse{
	// 	Success: true,
	// 	Message: "",
	// })
}

// @Description Refresh Token
// @Summary Refresh Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body token string true "Refresh Token"
// @Success 200 {object} dto.LoginResp
// @Failure 400,404,500 {object} dto.ErrorResponse
// @Router /refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()

	var body dto.RefreshTokenReq
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		h.logger.WarnContext(ctx,
			"invalid refresh request",
			"error", err,
		)

		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.service.Auth.RefreshToken(ctx, body)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"failed refresh request",
			"error", err,
		)
		errorResponse(c, http.StatusBadRequest, "unexpected error")
	}

	c.JSON(http.StatusOK, resp)
}
