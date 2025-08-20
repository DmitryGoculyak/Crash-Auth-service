package dto

type PasswordInput struct {
	Email       string `validate:"required,email"`
	OldPassword string `validate:"required,min=6,max=20"`
	NewPassword string `validate:"required,min=6,max=20"`
}

type ChangePasswordRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
