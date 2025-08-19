package validation

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

var Validate *validator.Validate

func InitValidation() (*validator.Validate, error) {
	Validate = validator.New()

	_ = Validate.RegisterValidation("fullname", func(name validator.FieldLevel) bool {
		parts := strings.Fields(name.Field().String())
		return len(parts) >= 2
	})
	return Validate, nil
}
