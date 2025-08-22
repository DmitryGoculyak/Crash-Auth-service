package dto

type PasswordInput struct {
	CurrentPassword string `validate:"required,min=6,max=20"`
	NewPassword     string `validate:"required,min=6,max=20"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
