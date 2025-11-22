// Package metar
package metar

import (
	"errors"
)

var (
	ErrICAOInvalid    = errors.New("invalid ICAO value")
	ErrTargetNotFound = errors.New("target not found")
)

type ManagerInterface interface {
	Query(icao string) (string, error)
	BatchQuery(icaos []string) []string
}

type ProviderInterface interface {
	Get(icao string) (string, error)
}

type DecoderInterface interface {
	Decode(raw []byte, selector string, reverse bool, multiline string) (bool, string, error)
}

type ParserInterface[T any] interface {
	Parse(data string) (T, error)
}
