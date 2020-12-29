package service

import (
	"io/ioutil"
	"log"

	"github.com/AyushSenapati/guardian/lib/plugin"
	"github.com/AyushSenapati/guardian/lib/proxy"
)

// Loader is responsible for loading service config and configuring proxy register
type Loader struct {
	Register *proxy.Register
}

// NewLoader returns a service loader by initialising it with the given proxy register
func NewLoader(r *proxy.Register) *Loader {
	return &Loader{Register: r}
}

// LoadServiceDefinitions reads provided config file and returns service definitions
func (l *Loader) LoadServiceDefinitions(filePath string) []*Definition {
	config, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("srv: error reading service definition file", filePath)
		return []*Definition{}
	}

	// Parse and load the service definitions
	definitions := ParseAndLoad(config)
	for _, def := range definitions {
		log.Printf("debug: PROXY CONFIGURATION: %+v", def.Proxy)
		log.Printf("debug: PLUGIN CONFIGURATION: %+v", def.Plugins)
	}

	return definitions
}

// RegisterServices registers the provided service defitions
func (l *Loader) RegisterServices(definitions []*Definition) {
	for _, def := range definitions {
		l.registerService(def)
	}
}

func (l *Loader) registerService(svcDef *Definition) {
	log.Printf("debug: register service: %s... started\n", svcDef.Name)

	if !svcDef.Active {
		log.Printf("warn: service %s is not active. skipping registraion\n", svcDef.Name)
		return
	}

	routerDefinition := proxy.NewRouterDefinition(svcDef.Proxy)

	// Configure the service specific plugins
	for _, plg := range svcDef.Plugins {
		log.Printf("debug: registering plugin: %s for service: %s...", plg.Name, svcDef.Name)
		if !plg.Enable {
			log.Printf(" warn: plugin `%s` is not enabled!!! skipping...", plg.Name)
			continue
		}

		// Get the setup function from the plugin registry by its name
		setupFunc, err := plugin.GetSetupFunc(plg.Name)
		if err != nil {
			log.Printf("error: could not load plugin `%s`", plg.Name)
			continue
		}

		err = setupFunc(routerDefinition, plg.Config)
		if err != nil {
			log.Printf("error: failed configuring plugin:`%s` [%s]", plg.Name, err)
		}
	}

	// Register the proxy configs along with plugins in the proxy register
	l.Register.Add(routerDefinition)

	log.Printf("debug: register service: %s... Completed\n", svcDef.Name)
}
