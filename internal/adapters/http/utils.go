package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func writeAnswer(writer http.ResponseWriter, status int, message string) {
	var errorFlag bool
	if status >= 400 {
		errorFlag = true
	}
	msg := Message{
		StatusCode: status,
		Message:    message,
		IsError:    errorFlag,
	}
	writer.WriteHeader(status)
	err := json.NewEncoder(writer).Encode(msg)
	if err != nil {
		log.Println("BAD json") //TODO
	}
}
