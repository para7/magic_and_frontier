package application

import (
	"strings"

	"tools2/app/internal/domain/common"
)

func appendSaveIssues[T any](report *ValidationReport, entity, id string, result common.SaveResult[T]) {
	if result.OK {
		return
	}
	if len(result.FieldErrors) == 0 {
		report.Issues = append(report.Issues, ValidationIssue{
			Entity:  entity,
			ID:      id,
			Message: result.FormError,
		})
		return
	}
	for field, message := range result.FieldErrors {
		report.Issues = append(report.Issues, ValidationIssue{
			Entity:  entity,
			ID:      id,
			Field:   field,
			Message: message,
		})
	}
}

func entryIDs[T any](entries []T, idOf func(T) string) map[string]struct{} {
	out := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		id := strings.TrimSpace(idOf(entry))
		if id != "" {
			out[id] = struct{}{}
		}
	}
	return out
}
