package dto

type EnableTotpMfaInput struct {
	Otp string `json:"code" binding:"required"`
}
