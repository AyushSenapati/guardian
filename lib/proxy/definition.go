package proxy

import (
	"github.com/AyushSenapati/guardian/lib/router"
)

// Definition defines proxy definition
type Definition struct {
	ListenPath   string `json:"listen_path"`
	Upstream     string `json:"upstream"`
	PreserveHost bool   `json:"preserve_host"`
	StripPath    bool   `json:"strip_path"`
}

// RouterDefinition defines the proxy router
// It is helpful to hold service's proxy specific plugins
type RouterDefinition struct {
	*Definition
	middlewareFuncs []router.MiddlewareFunc
}

// NewDefinition returns new instance of proxy definition initialised with defaults
func NewDefinition() *Definition {
	return &Definition{}
}

// NewRouterDefinition takes proxy definition and
// returns a new instance of RouterDefintion
func NewRouterDefinition(definition *Definition) *RouterDefinition {
	return &RouterDefinition{
		Definition: definition,
	}
}

// AddMiddleware adds middlewares
func (rd *RouterDefinition) AddMiddleware(mw router.MiddlewareFunc) {
	rd.middlewareFuncs = append(rd.middlewareFuncs, mw)
}

// ListMiddlewareFuncs retuns list of registered middleware functions
func (rd *RouterDefinition) ListMiddlewareFuncs() []router.MiddlewareFunc {
	return rd.middlewareFuncs
}
