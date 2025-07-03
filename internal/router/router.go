package router

import "net/http"

func NewRouter(
	userCreateHandler http.Handler,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/user", userCreateHandler)

	return mux
}
