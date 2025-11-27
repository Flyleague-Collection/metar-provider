// Package service
package service

import (
	"errors"
	"metar-service/src/interfaces/metar"

	"half-nothing.cn/service-core/interfaces/http/dto"
	"half-nothing.cn/service-core/interfaces/logger"
)

type Metar struct {
	logger       logger.Interface
	metarManager metar.ManagerInterface
	tafManager   metar.ManagerInterface
}

func NewMetar(
	lg logger.Interface,
	metarManager metar.ManagerInterface,
	tafManager metar.ManagerInterface,
) *Metar {
	return &Metar{
		logger:       logger.NewLoggerAdapter(lg, "MetarService"),
		metarManager: metarManager,
		tafManager:   tafManager,
	}
}

var ErrMetarNotFound = dto.NewApiStatus("NOT_FOUND", "Metar not found", dto.HttpCodeNotFound)

func (m *Metar) QueryMetar(icao string) *dto.ApiResponse[[]string] {
	data, err := m.metarManager.Query(icao)
	if errors.Is(err, metar.ErrTargetNotFound) {
		return dto.NewApiResponse[[]string](ErrMetarNotFound, nil)
	}
	if errors.Is(err, metar.ErrICAOInvalid) {
		return dto.NewApiResponse[[]string](dto.ErrErrorParam, nil)
	}
	if err != nil {
		return dto.NewApiResponse[[]string](dto.ErrServerError, nil)
	}
	return dto.NewApiResponse[[]string](dto.SuccessHandleRequest, []string{data})
}

func (m *Metar) BatchQueryMetar(icaos []string) *dto.ApiResponse[[]string] {
	if len(icaos) == 0 {
		return dto.NewApiResponse[[]string](dto.ErrInvalidParam, nil)
	}
	data := m.metarManager.BatchQuery(icaos)
	return dto.NewApiResponse[[]string](dto.SuccessHandleRequest, data)
}

var ErrTafNotFound = dto.NewApiStatus("NOT_FOUND", "Taf not found", dto.HttpCodeNotFound)

func (m *Metar) QueryTaf(icao string) *dto.ApiResponse[[]string] {
	data, err := m.tafManager.Query(icao)
	if errors.Is(err, metar.ErrTargetNotFound) {
		return dto.NewApiResponse[[]string](ErrTafNotFound, nil)
	}
	if errors.Is(err, metar.ErrICAOInvalid) {
		return dto.NewApiResponse[[]string](dto.ErrErrorParam, nil)
	}
	if err != nil {
		return dto.NewApiResponse[[]string](dto.ErrServerError, nil)
	}
	return dto.NewApiResponse[[]string](dto.SuccessHandleRequest, []string{data})
}

func (m *Metar) BatchQueryTaf(icaos []string) *dto.ApiResponse[[]string] {
	if len(icaos) == 0 {
		return dto.NewApiResponse[[]string](dto.ErrInvalidParam, nil)
	}
	data := m.tafManager.BatchQuery(icaos)
	return dto.NewApiResponse[[]string](dto.SuccessHandleRequest, data)
}
