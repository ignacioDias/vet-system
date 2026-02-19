package utils

import (
	"net/http"
	"strconv"
)

func Pagination(r *http.Request) (limit, offset int) {
	page := min(parseIntOrDefault(r.URL.Query().Get("page"), 1), 10000)
	limit = min(parseIntOrDefault(r.URL.Query().Get("limit"), 20), 100)
	offset = (page - 1) * limit
	return
}

func parseIntOrDefault(value string, def int) int {
	n, err := strconv.Atoi(value)
	if err != nil || n < 1 {
		return def
	}
	return n
}
