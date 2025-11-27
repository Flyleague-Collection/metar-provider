// Package config
package config

import "fmt"

type Config struct {
	GlobalConfig    *GlobalConfig     `yaml:"global"`
	ServerConfig    *ServerConfig     `yaml:"server"`
	ProviderConfigs []*ProviderConfig `yaml:"provider"`
}

func (c *Config) InitDefaults() {
	c.GlobalConfig = &GlobalConfig{}
	c.GlobalConfig.InitDefaults()
	c.ServerConfig = &ServerConfig{}
	c.ServerConfig.InitDefaults()
	c.ProviderConfigs = []*ProviderConfig{{}}
	c.ProviderConfigs[0].InitDefaults()
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
	if c.ProviderConfigs == nil {
		return false, fmt.Errorf("provider config is nil")
	}
	for _, provider := range c.ProviderConfigs {
		if ok, err := provider.Verify(); !ok {
			return false, err
		}
	}
	return true, nil
}
