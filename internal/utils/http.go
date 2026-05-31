package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondJSON отправляет HTTP-ответ с указанным HTTP-статусом и телом в формате JSON.
// Устанавливает заголовок Content-Type: application/json.
func RespondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}

// RespondError отправляет HTTP-ответ с ошибкой в формате JSON.
// Тело ответа содержит поле "error" с текстом ошибки.
func RespondError(w http.ResponseWriter, status int, err error) {
	RespondJSON(w, status, map[string]string{"error": err.Error()})
}
