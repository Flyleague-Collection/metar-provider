// Package config
package config

import "fmt"

type Verifiable interface {
	Verify() (bool, error)
}

type Config struct {
	GlobalConfig   *GlobalConfig     `yaml:"global"`
	ServerConfig   *ServerConfig     `yaml:"server"`
	ProviderConfig []*ProviderConfig `yaml:"provider"`
}

func NewConfig() *Config {
	return &Config{
		GlobalConfig: defaultGlobalConfig(),
		ServerConfig: defaultServerConfig(),
		ProviderConfig: []*ProviderConfig{
			{
				Type:      "metar",
				Name:      "aviationweather",
				Target:    "https://aviationweather.gov/api/data/metar?ids=%s",
				Decoder:   "raw",
				Selector:  "",
				Reverse:   false,
				Multiline: "",
			},
		},
	}
}

func (c *Config) Verify() (bool, error) {
	if c.GlobalConfig == nil {
		return false, fmt.Errorf("global config is nil")
	}
	if ok, err := c.GlobalConfig.Verify(); !ok {
		return false, err
	}
	if c.ServerConfig == nil {
		return false, fmt.Errorf("server config is nil")
	}
	if ok, err := c.ServerConfig.Verify(); !ok {
		return false, err
	}
	if c.ProviderConfig == nil {
		return false, fmt.Errorf("provider config is nil")
	}
	for _, provider := range c.ProviderConfig {
		if ok, err := provider.Verify(); !ok {
			return false, err
		}
	}
	return true, nil
}
