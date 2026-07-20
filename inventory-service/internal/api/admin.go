package api

import "net/http"

func isAdmin(r *http.Request) bool {
	role, ok := r.Context().Value("role").(string)
	return ok && role == "ADMIN"
}
