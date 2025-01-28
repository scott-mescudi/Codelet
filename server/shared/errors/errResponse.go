package errors

import (
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ErrorWithJson(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := map[string]interface{}{
		"error": message,
		"code":  code,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode JSON response: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
