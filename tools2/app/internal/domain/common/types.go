package common

type FieldErrors map[string]string

func (e FieldErrors) Add(field, message string) {
	if e == nil {
		return
	}
	e[field] = message
}

func (e FieldErrors) Any() bool {
	return len(e) > 0
}

type SaveMode string

const (
	SaveModeCreated SaveMode = "created"
	SaveModeUpdated SaveMode = "updated"
)

type SaveResult[T any] struct {
	OK          bool        `json:"ok"`
	Entry       *T          `json:"entry,omitempty"`
	Mode        SaveMode    `json:"mode,omitempty"`
	FieldErrors FieldErrors `json:"fieldErrors,omitempty"`
	FormError   string      `json:"formError,omitempty"`
}

func SaveValidationError[T any](errs FieldErrors, fallback string) SaveResult[T] {
	return SaveResult[T]{
		OK:          false,
		FieldErrors: errs,
		FormError:   fallback,
	}
}

func SaveSuccess[T any](entry T, mode SaveMode) SaveResult[T] {
	return SaveResult[T]{
		OK:    true,
		Entry: &entry,
		Mode:  mode,
	}
}

type DeleteResult struct {
	OK        bool   `json:"ok"`
	DeletedID string `json:"deletedId,omitempty"`
	FormError string `json:"formError,omitempty"`
	Code      string `json:"code,omitempty"`
}

func DeleteSuccess(id string) DeleteResult {
	return DeleteResult{OK: true, DeletedID: id}
}

func DeleteNotFound(entity string) DeleteResult {
	return DeleteResult{OK: false, Code: "NOT_FOUND", FormError: entity + " not found."}
}

type EntryState[T any] struct {
	Entries []T `json:"entries"`
}

func EmptyEntryState[T any]() EntryState[T] {
	return EntryState[T]{Entries: []T{}}
}
