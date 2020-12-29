package plugin

import (
	"fmt"

	"github.com/AyushSenapati/guardian/lib/plugin/limiter"
	"github.com/AyushSenapati/guardian/lib/proxy"
)

// // Registry maintains system wide plugins
// type Registry struct {
// 	sync.RWMutex
// 	plugins map[string]SetupFunc
// }
type registry map[string]SetupFunc

// register keeps all the active plugins
var register = registry{
	"limiter": limiter.SetupLimiter,
}

// GetSetupFunc returns SetupFunc for the requested plugin name
func GetSetupFunc(pluginName string) (SetupFunc, error) {
	setupFunc, found := register[pluginName]
	if !found {
		return nil, fmt.Errorf("plugin `%s` not found", pluginName)
	}
	return setupFunc, nil
}

// SetupFunc defines how a plugin should configure itself using given raw config to the proxy
type SetupFunc func(def *proxy.RouterDefinition, rawConfig map[string]interface{}) error

// // Config is the raw config DS
// type Config map[string]interface{}

// // RegisterPlugin registers the given plugin in the plugin registry
// // NOTE: Every plugin that is intended to be used,
// // must register itself to the plugin registry
// func (r *Registry) RegisterPlugin(name string, plugin SetupFunc) error {
// 	r.Lock()
// 	defer r.Unlock()

// 	if name == "" {
// 		return errors.New("plugin name can not be blank")
// 	}

// 	if _, found := r.plugins[name]; found {
// 		return fmt.Errorf("plugin: %s already registered", name)
// 	}

// 	r.plugins[name] = plugin
// 	return nil
// }
