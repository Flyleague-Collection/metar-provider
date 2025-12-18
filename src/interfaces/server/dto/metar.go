// Package dto
package dto

type QueryMetar struct {
	ICAO string `query:"icao" valid:"required"`
	Raw  bool   `query:"raw"`
}

type QueryTaf struct {
	ICAO string `query:"icao" valid:"required"`
	Raw  bool   `query:"raw"`
}
