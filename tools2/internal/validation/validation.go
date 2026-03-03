package validation

import (
	"regexp"
	"strings"

	"tools2/internal/form"
)

var phoneDigitsOnly = regexp.MustCompile(`^[0-9]+$`)

func Validate(name, phone string) form.Errors {
	err := form.Errors{}
	trimmedName := strings.TrimSpace(name)
	trimmedPhone := strings.TrimSpace(phone)

	if len(trimmedName) < 1 {
		err.Name = "名前は1文字以上で入力してください"
	}

	if len(trimmedPhone) < 1 {
		err.Phone = "電話番号は1文字以上で入力してください"
	} else if !phoneDigitsOnly.MatchString(trimmedPhone) {
		err.Phone = "電話番号は数字のみで入力してください"
	}

	return err
}
