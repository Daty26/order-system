package api

import (
	"net/http"
	"strconv"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

func parsePagination(r *http.Request) (int, int, bool) {
	limit := defaultLimit
	offset := 0
	if raw := r.URL.Query().Get("limit"); raw != "" {
		val, err := strconv.Atoi(raw)
		if err != nil {
			return 0, 0, false
		}
		if limit > 100 {
			limit = maxLimit
		} else {
			limit = val
		}
	}
	if raw := r.URL.Query().Get("offset"); raw != "" {
		val, err := strconv.Atoi(raw)
		if err != nil {
			return 0, 0, false
		}
		offset = val
	}
	return limit, offset, true
}
