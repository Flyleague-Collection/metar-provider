// Package config
package config

import "fmt"

type GrpcServerConfig struct {
	Enable bool   `yaml:"enable"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

func (g *GrpcServerConfig) InitDefaults() {
	g.Enable = true
	g.Host = "0.0.0.0"
	g.Port = 8081
}

func (g *GrpcServerConfig) Verify() (bool, error) {
	if !g.Enable {
		return true, nil
	}
	if g.Host == "" {
		return false, fmt.Errorf("host is empty")
	}
	if g.Port <= 0 {
		return false, fmt.Errorf("port must larger than 0")
	}
	if g.Port > 65535 {
		return false, fmt.Errorf("port must less than 65535")
	}
	return true, nil
}
