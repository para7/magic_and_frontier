package custom_validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	model "maf_command_editor/app/domain/model"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

var slugIDPattern = regexp.MustCompile(`^[a-z0-9_-]+$`)

func init() {
	Validate = validator.New()

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	Validate.RegisterValidation("trimmed_required", func(fl validator.FieldLevel) bool {
		return NormalizeText(fl.Field().String()) != ""
	})
	Validate.RegisterValidation("trimmed_min", func(fl validator.FieldLevel) bool {
		value := NormalizeText(fl.Field().String())
		want, ok := parseIntParam(fl.Param())
		return ok && len([]rune(value)) >= want
	})
	Validate.RegisterValidation("trimmed_max", func(fl validator.FieldLevel) bool {
		value := NormalizeText(fl.Field().String())
		want, ok := parseIntParam(fl.Param())
		return ok && len([]rune(value)) <= want
	})
	Validate.RegisterValidation("trimmed_oneof", func(fl validator.FieldLevel) bool {
		value := NormalizeText(fl.Field().String())
		for _, candidate := range strings.Fields(fl.Param()) {
			if value == candidate {
				return true
			}
		}
		return false
	})
	Validate.RegisterValidation("maf_slug_id", func(fl validator.FieldLevel) bool {
		return slugIDPattern.MatchString(fl.Field().String())
	})
}

// NormalizeText は CRLF/CR を LF に統一し、前後の空白を除去する。
func NormalizeText(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.TrimSpace(s)
}

func parseIntParam(param string) (int, bool) {
	var n int
	_, err := fmt.Sscanf(param, "%d", &n)
	return n, err == nil
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
	case "trimmed_min", "min":
		return fmt.Sprintf("最小文字数は%sです (%s=%s)", e.Param, e.Tag, e.Param)
	case "trimmed_max", "max":
		return fmt.Sprintf("最大文字数は%sです (%s=%s)", e.Param, e.Tag, e.Param)
	case "trimmed_oneof", "oneof":
		return fmt.Sprintf("次のいずれかの値が必要です: %s", e.Param)
	case "maf_slug_id":
		return "半角小文字英数字、_、- のみ使用できます"
	case "gte":
		return fmt.Sprintf("%s以上の値が必要です (gte=%s)", e.Param, e.Param)
	case "lte":
		return fmt.Sprintf("%s以下の値が必要です (lte=%s)", e.Param, e.Param)
	case "gt":
		return fmt.Sprintf("%sより大きい値が必要です (gt=%s)", e.Param, e.Param)
	case "lt":
		return fmt.Sprintf("%sより小さい値が必要です (lt=%s)", e.Param, e.Param)
	case "dive":
		return "配列要素にエラーがあります"
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
