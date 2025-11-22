// Package service
package service

import "metar-provider/src/interfaces/server/dto"

type MetarInterface interface {
	QueryMetar(icao string) *dto.ApiResponse[[]string]
	BatchQueryMetar(icaos []string) *dto.ApiResponse[[]string]
	QueryTaf(icao string) *dto.ApiResponse[[]string]
	BatchQueryTaf(icaos []string) *dto.ApiResponse[[]string]
}
