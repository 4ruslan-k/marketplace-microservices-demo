package dto

type ChangeCurrentPasswordInput struct {
	UserID                  string
	CurrentPassword         string
	NewPassword             string
	NewPasswordConfirmation string
}
