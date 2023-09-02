package dto

type LoginWithTotpInput struct {
	PasswordVerificationTokenID string `json:"tokenId"  binding:"required"`
	Code                        string `json:"code"  binding:"required"`
}
