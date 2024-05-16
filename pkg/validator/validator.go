package validator

import (
	"github.com/go-playground/validator/v10"
	"unicode"
)

func New() (*validator.Validate, error) {
	val := validator.New()
	err := val.RegisterValidation("contains_special", func(fl validator.FieldLevel) bool {
		for _, char := range fl.Field().String() {
			if unicode.IsPunct(char) || unicode.IsSymbol(char) {
				return true
			}
		}
		return false
	})
	if err != nil {
		return nil, err
	}
	err = val.RegisterValidation("contains_uppercase", func(fl validator.FieldLevel) bool {
		for _, char := range fl.Field().String() {
			if unicode.IsUpper(char) {
				return true
			}
		}
		return false
	})
	if err != nil {
		return val, err
	}
	return val, nil
}
