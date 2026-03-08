package common

import (
	"regexp"
	"strings"
)

var uuidPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// CRLF/CR を LF に統一し、前後の空白を除去する。
func NormalizeText(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.TrimSpace(s)
}

func IsUUID(s string) bool {
	return uuidPattern.MatchString(strings.ToLower(strings.TrimSpace(s)))
}

// 空文字または UUID 不正時は FieldErrors に追加し、空文字を返す。
func RequireUUID(errs FieldErrors, field, value string) string {
	v := NormalizeText(value)
	if v == "" {
		errs.Add(field, "Required.")
		return ""
	}
	if !IsUUID(v) {
		errs.Add(field, "Must be a UUID.")
		return ""
	}
	return v
}

// 正規化後の文字数が範囲外なら FieldErrors に追加し、空文字を返す。
func RequireText(errs FieldErrors, field, value string, min, max int) string {
	v := NormalizeText(value)
	if len(v) < min || len(v) > max {
		errs.Add(field, "Must be within allowed length.")
		return ""
	}
	return v
}

func OptionalText(value string) string {
	return NormalizeText(value)
}

// 範囲外なら FieldErrors に追加して nil、範囲内なら値のコピーへのポインタを返す。
func RequireNumberInRange(errs FieldErrors, field string, value, min, max float64) *float64 {
	if value < min || value > max {
		errs.Add(field, "Must be between range limits.")
		return nil
	}
	v := value
	return &v
}

// nil はそのまま返し、値ありで範囲外の場合のみ FieldErrors に追加して nil を返す。
func OptionalNumberInRange(errs FieldErrors, field string, value *float64, min, max float64) *float64 {
	if value == nil {
		return nil
	}
	if *value < min || *value > max {
		errs.Add(field, "Must be between range limits.")
		return nil
	}
	v := *value
	return &v
}
