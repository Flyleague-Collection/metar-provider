// Package config
package config

import "fmt"

type HSTSConfig struct {
	Enable            bool `yaml:"enable"`
	MaxAge            int  `yaml:"max_age"`
	IncludeSubdomains bool `yaml:"include_subdomains"`
}

func defaultHSTSConfig() *HSTSConfig {
	return &HSTSConfig{
		Enable:            false,
		MaxAge:            5184000,
		IncludeSubdomains: false,
	}
}

func (h *HSTSConfig) Verify() (bool, error) {
	if !h.Enable {
		return true, nil
	}
	if h.MaxAge <= 0 {
		return false, fmt.Errorf("max_age must larger than 0")
	}
	return true, nil
}

type SSLConfig struct {
	Enable     bool        `yaml:"enable"`
	Cert       string      `yaml:"cert"`
	Key        string      `yaml:"key"`
	ForceHttps bool        `yaml:"force_https"`
	HSTSConfig *HSTSConfig `yaml:"hsts"`
}

func defaultSSLConfig() *SSLConfig {
	return &SSLConfig{
		Enable:     false,
		Cert:       "",
		Key:        "",
		ForceHttps: false,
		HSTSConfig: defaultHSTSConfig(),
	}
}

func (s *SSLConfig) Verify() (bool, error) {
	if !s.Enable {
		return true, nil
	}
	if s.Cert == "" {
		return false, fmt.Errorf("cert is empty")
	}
	if s.Key == "" {
		return false, fmt.Errorf("key is empty")
	}
	if ok, err := s.HSTSConfig.Verify(); !ok {
		return false, err
	}
	return true, nil
}

type HttpServerConfig struct {
	Enable    bool       `yaml:"enable"`
	Host      string     `yaml:"host"`
	Port      int        `yaml:"port"`
	BodyLimit string     `yaml:"body_limit"`
	RateLimit int        `yaml:"rate_limit"`
	ProxyType int        `yaml:"proxy_type"`
	TrustIps  []string   `yaml:"trust_ips"`
	SSLConfig *SSLConfig `yaml:"ssl"`
}

func defaultHttpServerConfig() *HttpServerConfig {
	return &HttpServerConfig{
		Enable:    true,
		Host:      "0.0.0.0",
		Port:      8080,
		BodyLimit: "5M",
		RateLimit: 20,
		ProxyType: 0,
		TrustIps:  []string{"0.0.0.0/0"},
		SSLConfig: defaultSSLConfig(),
	}
}

func (h *HttpServerConfig) Verify() (bool, error) {
	if !h.Enable {
		return true, nil
	}
	if h.Host == "" {
		return false, fmt.Errorf("host is empty")
	}
	if h.Port <= 0 {
		return false, fmt.Errorf("port must larger than 0")
	}
	if h.Port > 65535 {
		return false, fmt.Errorf("port must less than 65535")
	}
	if h.ProxyType < 0 || h.ProxyType > 2 {
		return false, fmt.Errorf("proxy_type must be 0, 1 or 2")
	}
	if h.RateLimit < 0 {
		return false, fmt.Errorf("rate_limit must be positive")
	}
	if ok, err := h.SSLConfig.Verify(); !ok {
		return false, err
	}
	return true, nil
}
