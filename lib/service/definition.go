package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/AyushSenapati/guardian/lib/proxy"
)

// Plugin defines the plugin data structure that each service definition can have
type Plugin struct {
	Name   string
	Enable bool
	Config map[string]interface{}
}

// Definition defines service config, which needs to be proxied
type Definition struct {
	Name    string
	Active  bool
	Proxy   *proxy.Definition
	Plugins []Plugin
}

// Validate returns error if proxy/plugin validation fails
func (d *Definition) Validate() (err error) {
	err = d.Proxy.Validate()
	// TODO: implement plugin validation
	return
}

// NewDefinition returns new instance of service definition initialised with defaults
func NewDefinition() *Definition {
	return &Definition{
		Active:  true,
		Proxy:   proxy.NewDefinition(),
		Plugins: make([]Plugin, 0),
	}
}

// Configuration holds the service configurations
type Configuration struct {
	Definitions []*Definition
}

// UnmarshalJSON tells unmarshaller how to unmarshal configuration
func (c *Configuration) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &c.Definitions)
}

// Validate validates all the service definitions
func (c *Configuration) Validate() error {
	for _, def := range c.Definitions {
		if err := def.Validate(); err != nil {
			return fmt.Errorf("svc: %s err: %s", def.Name, err)
		}
	}
	return nil
}

// ParseAndLoad parses and loads raw config
func ParseAndLoad(rawConfig []byte) []*Definition {
	config := Configuration{}
	if err := json.Unmarshal(rawConfig, &config); err != nil {
		log.Println(err)
	}
	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}
	return config.Definitions
}
