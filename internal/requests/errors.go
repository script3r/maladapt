package requests

import (
	"encoding/json"
	"net/http"
	"time"
)

// Response is the base of all other responses
type Error struct {
	Date    time.Time `json:"date"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
}
type RequestError struct {
	Error Error `json:"error"`
}

const InvalidKeySupplied = "Invalid multipart/form-data key. Only accepts `file`"

func NewRequestError(code int, message string) RequestError {
	return RequestError{
		Error: Error{
			Date:    time.Now(),
			Code:    code,
			Message: message,
		},
	}
}

func WriteError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(NewRequestError(code, message))
}
