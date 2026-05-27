package common

import "net/http"

// AppError — ошибка с HTTP-статусом и произвольным сообщением.
type AppError struct {
	Status  int    `json:"-"`
	Message string `json:"error"`
}

func (e *AppError) Error() string { return e.Message }

// BadRequest — ошибка 400 с произвольным сообщением.
func BadRequest(msg ...string) *AppError {
	if len(msg) == 0 {
		msg = []string{"bad request"}
	}
	return &AppError{Status: http.StatusBadRequest, Message: msg[0]}
}

// NotFound — ошибка 404 с произвольным сообщением.
func NotFound(msg ...string) *AppError {
	if len(msg) == 0 {
		msg = []string{"not found"}
	}
	return &AppError{Status: http.StatusNotFound, Message: msg[0]}
}

// AlreadyExists — ошибка 409 с произвольным сообщением.
func AlreadyExists(msg ...string) *AppError {
	if len(msg) == 0 {
		msg = []string{"already exists"}
	}
	return &AppError{Status: http.StatusConflict, Message: msg[0]}
}

// InternalError — ошибка 500 с произвольным сообщением.
func InternalError(msg ...string) *AppError {
	if len(msg) == 0 {
		msg = []string{"internal error"}
	}
	return &AppError{Status: http.StatusInternalServerError, Message: msg[0]}
}
