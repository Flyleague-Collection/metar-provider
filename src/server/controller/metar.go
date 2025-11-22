// Package controller
package controller

import (
	"fmt"
	"metar-provider/src/interfaces/logger"
	"metar-provider/src/interfaces/server/dto"
	"metar-provider/src/interfaces/server/service"
	"strings"

	"github.com/labstack/echo/v4"
)

type Metar struct {
	logger  logger.Interface
	service service.MetarInterface
}

func NewMetar(
	lg logger.Interface,
	service service.MetarInterface,
) *Metar {
	return &Metar{
		logger:  logger.NewLoggerAdapter(lg, "MetarController"),
		service: service,
	}
}

func (m *Metar) QueryMetar(ctx echo.Context) error {
	data := &dto.QueryMetar{}

	if err := ctx.Bind(data); err != nil {
		m.logger.Errorf("QueryMetar bind param error: %+v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}

	m.logger.Debugf("QueryMetar called with args: %+v", data)

	icaos := strings.Split(data.ICAO, ",")

	if len(icaos) == 0 {
		m.logger.Errorf("QueryMetar param error icao(%s)", data.ICAO)
		return dto.ErrorResponse(ctx, dto.ErrInvalidParam)
	}

	var res *dto.ApiResponse[[]string]

	if len(icaos) == 1 {
		res = m.service.QueryMetar(icaos[0])
	} else {
		res = m.service.BatchQueryMetar(icaos)
	}

	if !data.Raw || res.Data == nil {
		return res.Response(ctx)
	}

	return dto.TextResponse(ctx, res.HttpCode, fmt.Sprintf("<pre>%s</pre>", strings.Join(res.Data, "</pre>\n<pre>")))
}

func (m *Metar) QueryTaf(ctx echo.Context) error {
	data := &dto.QueryTaf{}

	if err := ctx.Bind(data); err != nil {
		m.logger.Errorf("QueryTaf bind param error: %+v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}

	m.logger.Debugf("QueryTaf called with args: %+v", data)

	icaos := strings.Split(data.ICAO, ",")

	if len(icaos) == 0 {
		m.logger.Errorf("QueryTaf param error icao(%s)", data.ICAO)
		return dto.ErrorResponse(ctx, dto.ErrInvalidParam)
	}

	var res *dto.ApiResponse[[]string]

	if len(icaos) == 1 {
		res = m.service.QueryTaf(icaos[0])
	} else {
		res = m.service.BatchQueryTaf(icaos)
	}

	if !data.Raw || res.Data == nil {
		return res.Response(ctx)
	}

	return dto.TextResponse(ctx, res.HttpCode, fmt.Sprintf("<pre>%s</pre>", strings.Join(res.Data, "</pre>\n<pre>")))
}
