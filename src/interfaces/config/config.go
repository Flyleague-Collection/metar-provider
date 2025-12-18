// Package config
package config

import (
	"fmt"

	"half-nothing.cn/service-core/interfaces/config"
)

type Config struct {
	GlobalConfig    *GlobalConfig           `yaml:"global"`
	ServerConfig    *config.ServerConfig    `yaml:"server"`
	ProviderConfigs []*ProviderConfig       `yaml:"provider"`
	TelemetryConfig *config.TelemetryConfig `yaml:"telemetry"`
}

func (c *Config) InitDefaults() {
	c.GlobalConfig = &GlobalConfig{}
	c.GlobalConfig.InitDefaults()
	c.ServerConfig = &config.ServerConfig{}
	c.ServerConfig.InitDefaults()
	c.ProviderConfigs = []*ProviderConfig{{}}
	c.ProviderConfigs[0].InitDefaults()
	c.TelemetryConfig = &config.TelemetryConfig{}
	c.TelemetryConfig.InitDefaults()
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
	if ok, err := c.TelemetryConfig.Verify(); !ok {
		return false, err
	}
	return true, nil
}
