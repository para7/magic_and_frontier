package common

import (
	"regexp"
	"strings"
)

var (
	resourceIDPattern   = regexp.MustCompile(`^[a-z0-9_.-]+:[a-z0-9_./-]+$`)
	resourcePathPattern = regexp.MustCompile(`^[a-z0-9_./-]+$`)
)

func RequireNonEmptyID(errs FieldErrors, field, value string) string {
	id := NormalizeText(value)
	if id == "" {
		errs.Add(field, "Required.")
		return ""
	}
	return id
}

func NormalizeResourcePath(value string) string {
	value = strings.ReplaceAll(value, "\\", "/")
	value = NormalizeText(value)
	return strings.Trim(value, "/")
}

func IsNamespacedResourceID(value string) bool {
	return resourceIDPattern.MatchString(NormalizeText(value))
}

func HasSafeResourcePathSegments(value string) bool {
	path := NormalizeResourcePath(value)
	if path == "" {
		return false
	}
	for _, segment := range strings.Split(path, "/") {
		if segment == "" || segment == "." || segment == ".." {
			return false
		}
	}
	return true
}

func IsSafeNamespacedResourcePath(value string) bool {
	value = NormalizeText(value)
	if !resourceIDPattern.MatchString(value) {
		return false
	}
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return false
	}
	return HasSafeResourcePathSegments(parts[1])
}

func IsRelativeResourcePath(value string) bool {
	return resourcePathPattern.MatchString(NormalizeResourcePath(value))
}
