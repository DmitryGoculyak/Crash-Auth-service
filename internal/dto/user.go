package dto

type RegistrationInput struct {
	FullName     string `validate:"required,min=2,max=150,fullname"`
	Email        string `validate:"required,email,max=255"`
	Password     string `validate:"required,min=6,max=20"`
	CurrencyCode string `validate:"required,len=3"`
}

type UserRequest struct {
	FullName     string `json:"fullName"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	CurrencyCode string `json:"currencyCode"`
}
