package common

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

type FieldViolation struct {
	Field string
	Tag   string
	Param string
}

type ValidationMessageFunc func(FieldViolation) string

var (
	validateOnce sync.Once
	validateInst *validator.Validate
)

func Validator() *validator.Validate {
	validateOnce.Do(func() {
		v := validator.New()
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			tag := fld.Tag.Get("json")
			if tag == "" {
				return fld.Name
			}
			name := strings.Split(tag, ",")[0]
			if name == "" || name == "-" {
				return fld.Name
			}
			return name
		})
		mustRegisterValidation(v, "trimmed_required", validateTrimmedRequired)
		mustRegisterValidation(v, "trimmed_min", validateTrimmedMin)
		mustRegisterValidation(v, "trimmed_max", validateTrimmedMax)
		mustRegisterValidation(v, "trimmed_oneof", validateTrimmedOneOf)
		mustRegisterValidation(v, "uuid_any", validateUUIDAny)
		validateInst = v
	})
	return validateInst
}

func ValidateStruct(input any) []FieldViolation {
	err := Validator().Struct(input)
	if err == nil {
		return nil
	}
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return []FieldViolation{{Field: "", Tag: "invalid", Param: err.Error()}}
	}
	out := make([]FieldViolation, 0, len(validationErrs))
	for _, fe := range validationErrs {
		out = append(out, FieldViolation{
			Field: normalizeNamespace(fe.Namespace()),
			Tag:   fe.Tag(),
			Param: fe.Param(),
		})
	}
	return out
}

func ViolationsToFieldErrors(violations []FieldViolation, messageFn ValidationMessageFunc) FieldErrors {
	errs := FieldErrors{}
	for _, violation := range violations {
		if violation.Field == "" {
			continue
		}
		if _, exists := errs[violation.Field]; exists {
			continue
		}
		errs[violation.Field] = messageFn(violation)
	}
	return errs
}

func DefaultValidationMessage(v FieldViolation) string {
	switch v.Tag {
	case "required", "trimmed_required":
		return "Required."
	case "uuid", "uuid_any":
		return "Must be a UUID."
	case "trimmed_min", "trimmed_max", "min", "max":
		return "Must be within allowed length."
	case "gte", "lte":
		return fmt.Sprintf("Must satisfy %s %s.", v.Tag, v.Param)
	case "oneof", "trimmed_oneof":
		return "Invalid value."
	default:
		return "Invalid value."
	}
}

func mustRegisterValidation(v *validator.Validate, tag string, fn validator.Func) {
	if err := v.RegisterValidation(tag, fn); err != nil {
		panic(err)
	}
}

func validateTrimmedRequired(fl validator.FieldLevel) bool {
	return NormalizeText(fl.Field().String()) != ""
}

func validateTrimmedMin(fl validator.FieldLevel) bool {
	value := NormalizeText(fl.Field().String())
	want, ok := parseIntParam(fl.Param())
	return ok && len([]rune(value)) >= want
}

func validateTrimmedMax(fl validator.FieldLevel) bool {
	value := NormalizeText(fl.Field().String())
	want, ok := parseIntParam(fl.Param())
	return ok && len([]rune(value)) <= want
}

func validateTrimmedOneOf(fl validator.FieldLevel) bool {
	value := NormalizeText(fl.Field().String())
	for _, candidate := range strings.Fields(fl.Param()) {
		if value == candidate {
			return true
		}
	}
	return false
}

func validateUUIDAny(fl validator.FieldLevel) bool {
	return IsUUID(fl.Field().String())
}

func parseIntParam(param string) (int, bool) {
	var n int
	_, err := fmt.Sscanf(param, "%d", &n)
	return n, err == nil
}

func normalizeNamespace(ns string) string {
	if ns == "" {
		return ""
	}
	if i := strings.IndexByte(ns, '.'); i >= 0 {
		ns = ns[i+1:]
	}
	ns = strings.ReplaceAll(ns, "[", ".")
	ns = strings.ReplaceAll(ns, "]", "")
	return ns
}
