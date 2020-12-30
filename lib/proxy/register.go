package proxy

import (
	"log"
	"net/http"
	"regexp"

	"github.com/AyushSenapati/guardian/lib/router"
)

var matcher = regexp.MustCompile(`(\/\*(.+)?)`)

// Register is the register of the proxy, which manages the choosen router
type Register struct {
	Router router.Router
}

// NewRegister returns an instance of proxy register initialised with provided router
func NewRegister(rtr router.Router) *Register {
	return &Register{Router: rtr}
}

// Add registers the provided proxy definition in the register
func (r *Register) Add(def *RouterDefinition) {
	lb, err := NewLB(def.Upstreams.Strategy)
	if err != nil {
		log.Fatalln("error:", err)
	}
	reverseProxy := newRevesedProxy(def.Definition, lb)
	if reverseProxy.Transport == nil {
		reverseProxy.Transport = http.DefaultTransport
	}

	if matcher.MatchString(def.ListenPath) {
		r.doRegister(
			matcher.ReplaceAllString(def.ListenPath, ""),
			reverseProxy.ServeHTTP,
			def.ListMiddlewareFuncs(),
		)
	} else {
		log.Println("warn: not a valid listen_path... skiping")
	}
}

func (r *Register) doRegister(
	listenPath string, handler http.HandlerFunc, mwfs []router.MiddlewareFunc) {
	// listenPath = listenPath + "/{[A-Za-z0-9_@./#&+-]+}"
	listenPath += "/"
	r.Router.RegisterPath(listenPath, handler, mwfs)
}
