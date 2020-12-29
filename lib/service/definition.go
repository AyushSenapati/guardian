package service

import (
	"encoding/json"
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

// ParseAndLoad parses and loads raw config
func ParseAndLoad(rawConfig []byte) []*Definition {
	config := Configuration{}
	if err := json.Unmarshal(rawConfig, &config); err != nil {
		log.Println(err)
	}
	return config.Definitions
}
