package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func SuccessResponse(w http.ResponseWriter, httpStatus int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	err := json.NewEncoder(w).Encode(&Response{Data: data})
	if err != nil {
		log.Fatalf("Coundn't encode the data: " + err.Error())
		return
	}
}
func ErrorResponse(w http.ResponseWriter, httpStatus int, error string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	err := json.NewEncoder(w).Encode(&Response{Error: error})
	if err != nil {
		log.Fatalf("Coundn't encode the data: " + err.Error())
		return
	}
}
