package dto

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}
