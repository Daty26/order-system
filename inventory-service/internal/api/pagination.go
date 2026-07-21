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
		value, err := strconv.Atoi(raw)
		if err != nil {
			return 0, 0, false
		}

		if value > maxLimit {
			limit = maxLimit
		} else {
			limit = value
		}
	}
	if raw := r.URL.Query().Get("offset"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			return 0, 0, false
		}
		if value < 0 {
			return 0, 0, false
		}
		offset = value
	}
	return limit, offset, true
}
