package responses

import (
	"encoding/json"
	"net/http"
)

type ResponseError struct {
	Message []string `json:"messages"`
}

func Error(w http.ResponseWriter, code int, messages []string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ResponseError{messages})
}

func Success(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func CheckError(w http.ResponseWriter, isError bool, messages []string) {
	if isError {
		if messages == nil {
			messages = []string{"Invalid request"}
		}
		Error(w, http.StatusBadRequest, messages)
		return
	}
}
