package api

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func SuccessResp(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ApiResponse{Data: data})
}
func ErrorResponse(w http.ResponseWriter, status int, error string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ApiResponse{Error: error})

}
