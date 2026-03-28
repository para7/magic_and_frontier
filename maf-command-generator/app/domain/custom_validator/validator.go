package custom_validator

import (
	"fmt"
	"reflect"
	"strings"

	model "maf_command_editor/app/domain/model"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// trimmed_required などは自前実装が必要なので、ここで定義。
	Validate.RegisterValidation("trimmed_required", func(fl validator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})
}

// 各パッケージで github.com/go-playground/validator/v10 を個別にimportする必要がないように再エクスポート

type ValidationErrors = validator.ValidationErrors
type FieldError = validator.FieldError

func NewValidationError(entity, id string, fe FieldError) model.ValidationError {
	return model.ValidationError{
		Entity: entity,
		ID:     id,
		Field:  fe.Field(),
		Tag:    fe.Tag(),
		Param:  fe.Param(),
	}
}

func formatMessage(e model.ValidationError) string {
	switch e.Tag {
	case "trimmed_required":
		return "値が空です"
	case "required":
		return "値が必要です"
	case "gte":
		return fmt.Sprintf("%s以上の値が必要です (gte=%s)", e.Param, e.Param)
	case "lte":
		return fmt.Sprintf("%s以下の値が必要です (lte=%s)", e.Param, e.Param)
	case "gt":
		return fmt.Sprintf("%sより大きい値が必要です (gt=%s)", e.Param, e.Param)
	case "lt":
		return fmt.Sprintf("%sより小さい値が必要です (lt=%s)", e.Param, e.Param)
	case "min":
		return fmt.Sprintf("最小値は%sです (min=%s)", e.Param, e.Param)
	case "max":
		return fmt.Sprintf("最大値は%sです (max=%s)", e.Param, e.Param)
	default:
		if e.Param != "" {
			return fmt.Sprintf("'%s=%s' ルールに違反しています", e.Tag, e.Param)
		}
		return fmt.Sprintf("'%s' ルールに違反しています", e.Tag)
	}
}

func FormatValidationError(e model.ValidationError) string {
	return fmt.Sprintf("%s【%s】%s: %s", e.Entity, e.ID, e.Field, formatMessage(e))
}
