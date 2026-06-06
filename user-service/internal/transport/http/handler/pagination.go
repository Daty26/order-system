package transport_http_handler

import (
	"net/http"
	"strconv"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

func parsePagination (r *http.Request) (int,int, bool){
	limit := defaultLimit
	offset := 0
	if limitRaw := r.URL.Query().Get("limit"); limitRaw != ""{
		val, err := strconv.Atoi(limitRaw)
		if err != nil{
			return 0, 0, false
		}
		limit = val
	}
	if limit > maxLimit{
		limit = maxLimit
	}
	if offsetRaw := r.URL.Query().Get("offset"); offsetRaw != ""{
		val, err := strconv.Atoi(offsetRaw)
		if err != nil{
			return 0,0,false
		}
		offset = val
	}
	return limit, offset, true
}

