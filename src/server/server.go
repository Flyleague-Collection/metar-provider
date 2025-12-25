// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package server
package server

import (
	"io"
	"metar-service/src/interfaces/content"
	controllerImpl "metar-service/src/server/controller"
	serviceImpl "metar-service/src/server/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	h "half-nothing.cn/service-core/http"
	"half-nothing.cn/service-core/interfaces/logger"
)

func StartHttpServer(content *content.ApplicationContent) {
	lg := logger.NewLoggerAdapter(content.Logger(), "http-server")
	c := content.ConfigManager().GetConfig()

	lg.Info("Http server starting...")

	e := echo.New()
	e.Logger.SetOutput(io.Discard)
	e.Logger.SetLevel(log.OFF)

	h.SetEchoConfig(lg, e, c.ServerConfig.HttpServerConfig, nil)

	if c.TelemetryConfig.HttpServerTrace {
		h.SetTelemetry(e, c.TelemetryConfig, h.SkipperHealthCheck)
	}

	metarController := controllerImpl.NewMetar(lg, serviceImpl.NewMetar(lg, content.MetarManager(), content.TafManager()))

	h.SetHealthPoint(e)

	apiGroup := e.Group("/api/v1")
	apiGroup.GET("/metar", metarController.QueryMetar)
	apiGroup.GET("/taf", metarController.QueryTaf)

	h.SetUnmatchedRoute(e)
	h.SetCleaner(content.Cleaner(), e)

	h.Serve(lg, e, c.ServerConfig.HttpServerConfig)
}
