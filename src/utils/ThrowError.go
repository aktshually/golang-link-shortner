package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ThrowError(w http.ResponseWriter, code int, message string) {
	completeError := Error{code, message}
	response, err := json.Marshal(completeError)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	fmt.Fprintln(w, string(response))
}
