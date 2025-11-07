package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func SuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(Response{Data: data})
	if err != nil {
		log.Fatalf("Couldn't encode the data: %s", err.Error())
		return
	}
}
func ErrorResponse(w http.ResponseWriter, status int, error string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(Response{Error: error})
	if err != nil {
		log.Fatalf("Couldn't encode the error: %s", err.Error())
		return
	}
}
