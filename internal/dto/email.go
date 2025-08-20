package dto

type EmailInput struct {
	Password string `validate:"required,min=6,max=20"`
	NewEmail string `validate:"required,email"`
}

type ChangeEmailRequest struct {
	Password string `json:"password"`
	NewEmail string `json:"newEmail"`
}
