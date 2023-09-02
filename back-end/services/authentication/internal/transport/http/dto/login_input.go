package dto

type LoginInput struct {
	Email    string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}
