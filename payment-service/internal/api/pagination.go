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
		if val > maxLimit {
			limit = maxLimit
		} else if val <= 0 {
			return 0, 0, false
		} else {
			limit = val
		}
	}

	if raw := r.URL.Query().Get("offset"); raw != "" {
		val, err := strconv.Atoi(raw)
		if err != nil {
			return 0, 0, false
		}
		if val < 0 {
			return 0, 0, false
		}
		offset = val
	}
	return limit, offset, true
}
