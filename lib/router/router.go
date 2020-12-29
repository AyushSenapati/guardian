package router

import "net/http"

// MiddlewareFunc defines the middleware type
type MiddlewareFunc func(http.Handler) http.Handler

// Router interface defines basic functionality of a Guardian router
type Router interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	RegisterPath(path string, handler http.HandlerFunc, mwfs []MiddlewareFunc)
	Use(handlers ...MiddlewareFunc)
}
