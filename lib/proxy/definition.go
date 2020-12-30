package proxy

import (
	"fmt"
	"log"

	"github.com/AyushSenapati/guardian/lib/router"
	"github.com/asaskevich/govalidator"
)

// Definition defines proxy definition
type Definition struct {
	ListenPath   string     `json:"listen_path"`
	Upstreams    *Upstreams `json:"upstreams"`
	PreserveHost bool       `json:"preserve_host"`
	StripPath    bool       `json:"strip_path"`
}

// Validate returns error in case proxy definition validation fails
// currently it validates upstreams only
func (d *Definition) Validate() error {
	return d.Upstreams.Validate()
}

// RouterDefinition represents the proxy router
// It is helpful to hold service's proxy specific plugins
type RouterDefinition struct {
	*Definition
	middlewareFuncs []router.MiddlewareFunc
}

// Upstreams represents the targets where the requests would be forwarded to
type Upstreams struct {
	Strategy string   `json:"strategy"`
	Targets  []string `json:"targets" valid:"requrl,required"`
}

func (u *Upstreams) String() string {
	for _, t := range u.Targets {
		if !govalidator.IsRequestURL(t) {
			log.Fatalln("error: invalid URL:", t)
		}
	}
	return fmt.Sprintf("Strategy: %s, Targets: %v", u.Strategy, u.Targets)
}

// Validate returns error if invalid targets are configured
func (u *Upstreams) Validate() error {
	if len(u.Targets) == 0 {
		return ErrEmptyTargets
	}
	for _, t := range u.Targets {
		if !govalidator.IsRequestURL(t) {
			return fmt.Errorf("invalid URL: %s", t)
		}
	}
	return nil
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
