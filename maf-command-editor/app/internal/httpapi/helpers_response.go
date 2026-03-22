package httpapi

import (
	"encoding/json"
	"net/http"

	"tools2/app/internal/domain/common"
)

type apiErrorResponse struct {
	OK        bool   `json:"ok"`
	Code      string `json:"code,omitempty"`
	FormError string `json:"formError,omitempty"`
	Details   string `json:"details,omitempty"`
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dest any) bool {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		writeFormError(w, http.StatusBadRequest, "Invalid JSON request body.")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeFormError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, apiErrorResponse{
		OK:        false,
		FormError: message,
	})
}

func writeCodedError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, apiErrorResponse{
		OK:        false,
		Code:      code,
		FormError: message,
	})
}

func writeInternalError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusInternalServerError, apiErrorResponse{
		OK:        false,
		FormError: http.StatusText(http.StatusInternalServerError),
		Details:   err.Error(),
	})
}

func writeDuplicateIDValidationError[T any](w http.ResponseWriter) {
	writeJSON(w, http.StatusBadRequest, common.SaveValidationError[T](
		common.FieldErrors{"id": "この ID は既に使用されています。"},
		"Validation failed. Fix the highlighted fields.",
	))
}
