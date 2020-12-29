package config

import (
	"errors"

	"github.com/spf13/viper"
)

// Specification defines basic application configs
type Specification struct {
	Port     int
	AddReqID bool

	// in second(s)
	GraceTimeout int
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

func init() {
	viper.SetDefault("port", "8080")
	viper.SetDefault("gracetimeout", 15)
	viper.SetDefault("readtimeout", 15)
	viper.SetDefault("writetimeout", 15)
	viper.SetDefault("idletimeout", 15)

	viper.SetDefault("addreqid", true)
}

// Load reads the config file and returns read configs
func Load(configFile string) (*Specification, error) {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.New("No config file found")
	}

	var config Specification
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
