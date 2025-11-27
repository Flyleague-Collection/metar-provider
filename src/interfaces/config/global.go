// Package config
package config

import (
	"fmt"
	gl "metar-service/src/interfaces/global"

	"half-nothing.cn/service-core/interfaces/config"
	"half-nothing.cn/service-core/interfaces/global"
)

type GlobalConfig struct {
	Name      string            `yaml:"name"`
	Version   string            `yaml:"version"`
	LogConfig *config.LogConfig `yaml:"log"`
}

func (g *GlobalConfig) InitDefaults() {
	g.Name = "metar-service"
	g.Version = gl.ConfigVersion
	g.LogConfig = &config.LogConfig{}
	g.LogConfig.InitDefaults()
}

func (g *GlobalConfig) Verify() (bool, error) {
	if g.Name == "" {
		return false, fmt.Errorf("global name is empty")
	}
	if g.Version == "" {
		return false, fmt.Errorf("global version is empty")
	}
	configVersion, err := global.NewVersion(g.Version)
	if err != nil {
		return false, fmt.Errorf("global version is invalid: %s", err)
	}
	targetConfigVersion, _ := global.NewVersion(gl.ConfigVersion)
	if targetConfigVersion.CheckVersion(configVersion) != global.AllMatch {
		return false, fmt.Errorf("config version mismatch, expected %s, got %s", gl.ConfigVersion, g.Version)
	}
	if g.LogConfig == nil {
		return false, fmt.Errorf("log config is empty")
	}
	if ok, err := g.LogConfig.Verify(); !ok {
		return false, err
	}
	return true, nil
}
