package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type Resposne struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func SuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(&Resposne{Data: data})
	if err != nil {
		log.Fatalln("Couldn't encode request body: " + err.Error())
		return
	}
}
func ErrorResponse(w http.ResponseWriter, status int, error string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(&Resposne{Error: error})
	if err != nil {
		log.Fatalln("Couldn't encode request body: " + err.Error())
		return
	}

}
