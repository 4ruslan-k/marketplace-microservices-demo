package dto

type CreateUserInput struct {
	Name     string `json:"name"  binding:"required"`
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required"`
}
