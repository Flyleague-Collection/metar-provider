package main

import (
	"context"
	"fmt"
	grpcImpl "metar-service/src/grpc"
	c "metar-service/src/interfaces/config"
	"metar-service/src/interfaces/content"
	g "metar-service/src/interfaces/global"
	pb "metar-service/src/interfaces/grpc"
	"metar-service/src/metar"
	"metar-service/src/server"
	"time"

	"google.golang.org/grpc"
	"half-nothing.cn/service-core/cache"
	"half-nothing.cn/service-core/cleaner"
	"half-nothing.cn/service-core/config"
	"half-nothing.cn/service-core/discovery"
	grpcUtils "half-nothing.cn/service-core/grpc"
	"half-nothing.cn/service-core/interfaces/global"
	"half-nothing.cn/service-core/logger"
	"half-nothing.cn/service-core/telemetry"
	"half-nothing.cn/service-core/utils"
)

func main() {
	global.CheckFlags()
	utils.CheckIntEnv(g.EnvQueryThread, g.QueryThread)
	utils.CheckDurationEnv(g.EnvCacheCleanInterval, g.CacheCleanInterval)

	configManager := config.NewManager[*c.Config]()
	if err := configManager.Init(); err != nil {
		fmt.Printf("fail to initialize configuration file: %v", err)
		return
	}

	applicationConfig := configManager.GetConfig()
	lg := logger.NewLogger()
	lg.Init(
		global.LogName,
		applicationConfig.GlobalConfig.LogConfig,
	)

	lg.Info(" _____     _           _____             _")
	lg.Info("|     |___| |_ ___ ___|   __|___ ___ _ _|_|___ ___")
	lg.Info("| | | | -_|  _| .'|  _|__   | -_|  _| | | |  _| -_|")
	lg.Info("|_|_|_|___|_| |__,|_| |_____|___|_|  \\_/|_|___|___|")
	lg.Info(fmt.Sprintf("%51s", fmt.Sprintf("Copyright © %d-%d Half_nothing", global.BeginYear, time.Now().Year())))
	lg.Info(fmt.Sprintf("%51s", fmt.Sprintf("MetarService v%s", g.AppVersion)))

	cl := cleaner.NewCleaner(lg)
	cl.Init()

	if applicationConfig.TelemetryConfig.Enable {
		sdk := telemetry.NewSDK(lg, applicationConfig.TelemetryConfig)
		shutdown, err := sdk.SetupOTelSDK(context.Background())
		if err != nil {
			lg.Fatalf("fail to initialize telemetry: %v", err)
			return
		}
		cl.Add("Telemetry", shutdown)
	}

	metarManagerMemoryCache := cache.NewMemoryCache[string, *string](*g.CacheCleanInterval)
	cl.Add("Metar Cache", func(ctx context.Context) error {
		metarManagerMemoryCache.Close()
		return nil
	})
	metarManager := metar.NewManager(
		lg,
		utils.Filter(applicationConfig.ProviderConfigs, func(providerConfig *c.ProviderConfig) bool {
			return providerConfig.Type == c.ProviderTypeMetar.Value
		}),
		metarManagerMemoryCache,
	)

	tafManagerMemoryCache := cache.NewMemoryCache[string, *string](*g.CacheCleanInterval)
	cl.Add("Taf Cache", func(ctx context.Context) error {
		tafManagerMemoryCache.Close()
		return nil
	})
	tafManager := metar.NewManager(
		lg,
		utils.Filter(applicationConfig.ProviderConfigs, func(providerConfig *c.ProviderConfig) bool {
			return providerConfig.Type == c.ProviderTypeTaf.Value
		}),
		tafManagerMemoryCache,
	)

	applicationContent := content.NewApplicationContentBuilder().
		SetConfigManager(configManager).
		SetCleaner(cl).
		SetLogger(lg).
		SetMetarManager(metarManager).
		SetTafManager(tafManager).
		Build()

	go server.StartHttpServer(applicationContent)

	if applicationConfig.ServerConfig.GrpcServerConfig.Enable {
		started := make(chan bool)
		initFunc := func(s *grpc.Server) {
			grpcServer := grpcImpl.NewMetarServer(lg, metarManager, tafManager)
			pb.RegisterMetarServer(s, grpcServer)
		}
		if applicationConfig.TelemetryConfig.Enable && applicationConfig.TelemetryConfig.GrpcServerTrace {
			go grpcUtils.StartGrpcServerWithTrace(lg, cl, applicationConfig.ServerConfig.GrpcServerConfig, started, initFunc)
		} else {
			go grpcUtils.StartGrpcServer(lg, cl, applicationConfig.ServerConfig.GrpcServerConfig, started, initFunc)
		}
		go discovery.StartServiceDiscovery(lg, cl, started, utils.NewVersion(g.AppVersion),
			"metar-service", applicationConfig.ServerConfig.GrpcServerConfig.Port)
	}

	cl.Wait()
}
