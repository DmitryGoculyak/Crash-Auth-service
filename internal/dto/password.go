package dto

type PasswordInput struct {
	Email           string `validate:"required,email"`
	CurrentPassword string `validate:"required,min=6,max=20"`
	NewPassword     string `validate:"required,min=6,max=20"`
}

type ChangePasswordRequest struct {
	Email           string `json:"email"`
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
