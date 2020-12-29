package router

import (
	"log"
	"net/http"
)

// HTTPServeMux is an adapater to http default servemux,
// which implements router interface
type HTTPServeMux struct {
	mux         *http.ServeMux
	middlewares []MiddlewareFunc
}

// NewHTTPServeMux instantiates new http serve mux and returns the adapater to it
func NewHTTPServeMux() *HTTPServeMux {
	return &HTTPServeMux{mux: http.NewServeMux()}
}

func (r *HTTPServeMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// RegisterPath registers the path and its handlerfunc to the router
func (r *HTTPServeMux) RegisterPath(path string, handler http.HandlerFunc, mwfs []MiddlewareFunc) {
	// append the mwfs to existing r.middlewares.
	// This is important because middleware funcs provided by using Use()
	// are intended to be used for all the services registed in the gateway.
	// so they must be registered first in order. Then the service specific mwfs
	// should be registered
	finalMiddlewares := append(r.middlewares, mwfs...)

	r.mux.Handle(path, middleware(handler, finalMiddlewares...))
	log.Printf(
		"debug: path `%s` [%d global & %d local middlewares are registered]",
		path, len(r.middlewares), len(mwfs),
	)
}

// Use can be used to chain of global middlewares
func (r *HTTPServeMux) Use(mwf ...MiddlewareFunc) {
	for _, m := range mwf {
		r.middlewares = append(r.middlewares, m)
	}
}

// This applies middlewares in order they were registered
func middleware(h http.Handler, mwf ...MiddlewareFunc) http.Handler {
	for i := len(mwf) - 1; i >= 0; i-- {
		h = mwf[i](h)
	}
	return h
}
