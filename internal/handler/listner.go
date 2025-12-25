package handler

import "net/http"

func Listner(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}
