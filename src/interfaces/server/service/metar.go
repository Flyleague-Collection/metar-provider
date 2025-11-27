// Package service
package service

import (
	"half-nothing.cn/service-core/interfaces/http/dto"
)

type MetarInterface interface {
	QueryMetar(icao string) *dto.ApiResponse[[]string]
	BatchQueryMetar(icaos []string) *dto.ApiResponse[[]string]
	QueryTaf(icao string) *dto.ApiResponse[[]string]
	BatchQueryTaf(icaos []string) *dto.ApiResponse[[]string]
}
