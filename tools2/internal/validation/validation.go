package validation

import (
	"regexp"
	"strconv"
	"strings"

	"tools2/internal/form"
)

var phoneDigitsOnly = regexp.MustCompile(`^[0-9]+$`)

func Validate(state form.State) form.Errors {
	err := form.Errors{}
	trimmedName := strings.TrimSpace(state.Name)
	trimmedPhone := strings.TrimSpace(state.Phone)

	if len(trimmedName) < 1 {
		err.Name = "名前は1文字以上で入力してください"
	}

	if len(trimmedPhone) < 1 {
		err.Phone = "電話番号は1文字以上で入力してください"
	} else if !phoneDigitsOnly.MatchString(trimmedPhone) {
		err.Phone = "電話番号は数字のみで入力してください"
	}

	if !form.IsValidMode(state.Mode) {
		err.Mode = "入力種別を選択してください"
		return err
	}

	switch state.Mode {
	case form.ModeLatLng:
		if !isNumber(state.Latitude) {
			err.Latitude = "緯度は数値で入力してください"
		}
		if !isNumber(state.Longitude) {
			err.Longitude = "経度は数値で入力してください"
		}
	case form.ModeBirthdate:
		if strings.TrimSpace(state.Birthdate) == "" {
			err.Birthdate = "生年月日を入力してください"
		}
	case form.ModeHeightWeight:
		if !isNumber(state.Height) {
			err.Height = "身長は数値で入力してください"
		}
		if !isNumber(state.Weight) {
			err.Weight = "体重は数値で入力してください"
		}
	}

	return err
}

func isNumber(v string) bool {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return false
	}
	_, err := strconv.ParseFloat(trimmed, 64)
	return err == nil
}
