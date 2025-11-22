// Package config
package config

import (
	"errors"
	"fmt"
	"metar-provider/src/interfaces/metar"
	decoderImpl "metar-provider/src/metar/decoder"
	"metar-provider/src/utils"
	"strings"
)

type ProviderConfig struct {
	Type      string `yaml:"type"`
	Name      string `yaml:"name"`
	Target    string `yaml:"target"`
	Decoder   string `yaml:"decoder"`
	Selector  string `yaml:"selector"`
	Reverse   bool   `yaml:"reverse"`
	Multiline string `yaml:"multiline"`
}

type ProviderType *utils.Enum[string, string]

var (
	ProviderTypeMetar ProviderType = utils.NewEnum("metar", "METAR")
	ProviderTypeTaf   ProviderType = utils.NewEnum("taf", "TAF")
)

var ProviderTypes = utils.NewEnums(ProviderTypeMetar, ProviderTypeTaf)

type DecoderType *utils.Enum[string, metar.DecoderInterface]

var (
	DecoderTypeRaw  DecoderType = utils.NewEnum[string, metar.DecoderInterface]("raw", &decoderImpl.RawDecoder{})
	DecoderTypeHtml DecoderType = utils.NewEnum[string, metar.DecoderInterface]("html", &decoderImpl.HtmlDecoder{})
	DecoderTypeJson DecoderType = utils.NewEnum[string, metar.DecoderInterface]("json", &decoderImpl.JsonDecoder{})
)

var DecoderTypes = utils.NewEnums(DecoderTypeRaw, DecoderTypeHtml, DecoderTypeJson)

func (p *ProviderConfig) Verify() (bool, error) {
	if p.Type == "" {
		return false, fmt.Errorf("type is required")
	}
	providerType := strings.ToLower(p.Type)
	if !ProviderTypes.IsValidEnum(providerType) {
		return false, errors.New("type is not supported")
	}
	if p.Name == "" {
		return false, fmt.Errorf("name is required")
	}
	if p.Target == "" {
		return false, fmt.Errorf("target is required")
	}
	if p.Decoder == "" {
		return false, fmt.Errorf("parser is required")
	}
	parser := strings.ToLower(p.Decoder)
	if !DecoderTypes.IsValidEnum(parser) {
		return false, errors.New("parser is not supported")
	}
	switch parser {
	case DecoderTypeHtml.Value:
		if p.Selector == "" {
			return false, fmt.Errorf("provider %s with type 'html' need a selector", p.Name)
		}
	case DecoderTypeJson.Value:
		if p.Selector == "" {
			return false, fmt.Errorf("provider %s with type 'json' need a selector", p.Name)
		}
	}
	return true, nil
}
