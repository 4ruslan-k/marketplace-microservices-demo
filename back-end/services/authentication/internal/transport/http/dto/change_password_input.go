package dto

// Change Password
type ChangePasswordInput struct {
	CurrentPassword         string `json:"currentPassword"  binding:"required"`
	NewPassword             string `json:"newPassword"  binding:"required"`
	NewPasswordConfirmation string `json:"newPasswordConfirmation"  binding:"required"`
}
