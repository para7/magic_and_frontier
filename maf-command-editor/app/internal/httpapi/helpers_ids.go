package httpapi

import (
	"strings"
)

func findEntry[T any](entries []T, id string, idOf func(T) string) (T, bool) {
	var zero T
	id = strings.TrimSpace(id)
	if id == "" {
		return zero, false
	}
	for _, entry := range entries {
		if strings.TrimSpace(idOf(entry)) == id {
			return entry, true
		}
	}
	return zero, false
}
