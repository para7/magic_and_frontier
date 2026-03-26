package custom_validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

// trimmed_required などは自前実装が必要なので、ここで定義。
func init() {
	Validate = validator.New()

	Validate.RegisterValidation("trimmed_required", func(fl validator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})

	// Validate.RegisterValidation("trimmed_max", func(fl validator.FieldLevel) bool {
	// 	max, err := strconv.Atoi(fl.Param())
	// 	if err != nil {
	// 		return false
	// 	}
	// 	return len([]rune(strings.TrimSpace(fl.Field().String()))) <= max
	// })
}

type ValidationErrors = validator.ValidationErrors
