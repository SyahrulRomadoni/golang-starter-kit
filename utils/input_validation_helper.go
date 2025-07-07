package utils

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// formatValidationError mengubah error validasi menjadi pesan yang mudah dibaca
func formatValidationError(fe validator.FieldError) string {
	field := fe.Field()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("Field '%s' wajib diisi", field)
	case "email":
		return fmt.Sprintf("Field '%s' harus berupa email yang valid", field)
	case "min":
		return fmt.Sprintf("Field '%s' minimal harus %s karakter", field, fe.Param())
	case "max":
		return fmt.Sprintf("Field '%s' maksimal %s karakter", field, fe.Param())
	default:
		return fmt.Sprintf("Field '%s' tidak valid", field)
	}
}

// InputValidation menggabungkan proses bind JSON dan validasi input
func InputValidation(c *gin.Context, input interface{}) (bool, APIResponse) {
	if err := c.ShouldBindJSON(input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			msg := formatValidationError(ve[0])
			return false, APIResponseError(msg, nil)
		}
		return false, APIResponseError("Input tidak valid", nil)
	}

	if err := validate.Struct(input); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			msg := formatValidationError(ve[0])
			return false, APIResponseError(msg, nil)
		}
		return false, APIResponseError("Input tidak valid", nil)
	}

	return true, APIResponse{}
}

// InputValidationPasswordCriteria memastikan password hanya mengandung huruf, angka, dan karakter @#$.
func InputValidationPasswordCriteria(password string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9@#$]+$`, password)
	return matched
}
