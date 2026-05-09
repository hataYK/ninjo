package handler

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator は Echo 用のバリデーター。
// go-playground/validator をラップして Echo の Validator interface を満たす。
type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
