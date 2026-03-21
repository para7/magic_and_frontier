package common

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	resourceIDPattern   = regexp.MustCompile(`^[a-z0-9_.-]+:[a-z0-9_./-]+$`)
	resourcePathPattern = regexp.MustCompile(`^[a-z0-9_./-]+$`)
)

func IsPrefixedSequenceID(value, prefix string) bool {
	value = NormalizeText(value)
	if !strings.HasPrefix(value, prefix) {
		return false
	}
	suffix := strings.TrimPrefix(value, prefix)
	if suffix == "" {
		return false
	}
	n, err := strconv.Atoi(suffix)
	return err == nil && n >= 1
}

func RequirePrefixedSequenceID(errs FieldErrors, field, value, prefix string) string {
	id := NormalizeText(value)
	if id == "" {
		errs.Add(field, "Required.")
		return ""
	}
	if !IsPrefixedSequenceID(id, prefix) {
		errs.Add(field, "Invalid ID format.")
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
