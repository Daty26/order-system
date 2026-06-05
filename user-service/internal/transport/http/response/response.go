package transport_http_response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct{
	Data any `json:"data"`
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, status int, body Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("failed to encode HTTP response", "error", err)
	}
}

func SuccessJSON(w http.ResponseWriter, status int, data any){
	JSON(w, status, Response{Data: data})
}
func ErrorJSON(w http.ResponseWriter, status int, errMsg string){
	JSON(w, status, Response{Error: errMsg})	
}