package requests

import (
	"encoding/json"
	"github.com/worlvlhole/maladapt/internal/quarantine"
	"net/http"
	"time"
)

// Response is the base of all other responses
type Response struct {
	Date    time.Time   `json:"date"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

type responseError struct {
	Error Response `json:"error"`
}

type responseSuccess struct {
	Success Response `json:"success"`
}

const InvalidKeySupplied = "Invalid multipart/form-data key. Only accepts `file`"

func NewResponseError(code int, message string) responseError {
	return responseError{
		Error: Response{
			Code:    code,
			Date:    time.Now(),
			Message: message,
		},
	}
}

func NewResponseSuccess(message quarantine.ScanResponse) responseSuccess {
	return responseSuccess{
		Success: Response{
			Code:    http.StatusOK,
			Date:    time.Now(),
			Message: message,
		},
	}
}
func WriteError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(NewResponseError(code, message))
}

func WriteSuccess(w http.ResponseWriter, message quarantine.ScanResponse) {
	json.NewEncoder(w).Encode(NewResponseSuccess(message))
}
