// Package metar
package metar

import (
	"bytes"
	"fmt"
	"io"
	"metar-provider/src/interfaces/config"
	"metar-provider/src/interfaces/logger"
	"metar-provider/src/interfaces/metar"
	"net/http"
)

type Provider struct {
	config  *config.ProviderConfig
	logger  logger.Interface
	decoder metar.DecoderInterface
}

func NewProvider(
	lg logger.Interface,
	c *config.ProviderConfig,
) *Provider {
	return &Provider{
		config:  c,
		logger:  logger.NewLoggerAdapter(lg, fmt.Sprintf("provider-%s", c.Name)),
		decoder: config.DecoderTypes.GetEnum(c.Decoder).Data,
	}
}

func (p *Provider) Get(icao string) (string, error) {
	if icao == "" || len(icao) != 4 {
		return "", metar.ErrICAOInvalid
	}
	url := fmt.Sprintf(p.config.Target, icao)
	p.logger.Debugf("Getting Data from %s", url)
	response, err := http.Get(url)
	if err != nil {
		p.logger.Errorf("Error getting data from %s: %s", url, err.Error())
		return "", err
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(response.Body)

	if response.StatusCode != http.StatusOK {
		p.logger.Errorf("Error getting data from %s, status code %s", url, response.Status)
		return "", fmt.Errorf("error getting data from %s, status code %s", url, response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		p.logger.Errorf("Error getting data from %s, reading response body fail: %s", url, err.Error())
		return "", err
	}

	data = bytes.TrimRight(data, "\n")

	ok, decodeData, err := p.decoder.Decode(data, p.config.Selector, p.config.Reverse, p.config.Multiline)
	if err != nil {
		p.logger.Errorf("Error getting data from %s, decoding fail: %s", url, err.Error())
		return "", err
	}
	if !ok {
		p.logger.Errorf("Error getting data from %s, data not found", url)
		return "", metar.ErrTargetNotFound
	}
	return decodeData, nil
}
