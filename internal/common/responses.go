package common

import (
	"encoding/json"
	"errors"
	"net/http"
)

// RespondJSON отправляет JSON-ответ с заданным статусом и данными.
func RespondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// RespondError отправляет JSON-ответ с ошибкой.
// Если err реализует *AppError — используется его статус и сообщение.
// Для обычных ошибок — 500 Internal Server Error.
func RespondError(w http.ResponseWriter, err error) {
	if err == nil {
		RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "unknown error"})
		return
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		RespondJSON(w, appErr.Status, map[string]string{"error": appErr.Message})
	} else {
		RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
}
