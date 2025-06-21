package web

import "net/http"

func Param(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
