package main

import (
	"context"
	"flag"
	"fmt"
	"metar-provider/src/cache"
	cleanerImpl "metar-provider/src/cleaner"
	configImpl "metar-provider/src/config"
	grpcImpl "metar-provider/src/grpc"
	"metar-provider/src/interfaces/config"
	"metar-provider/src/interfaces/content"
	"metar-provider/src/interfaces/global"
	pb "metar-provider/src/interfaces/grpc"
	loggerImpl "metar-provider/src/logger"
	"metar-provider/src/metar"
	"metar-provider/src/server"
	"metar-provider/src/utils"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	flag.Parse()

	global.CheckBoolEnv(global.EnvNoLogs, global.NoLogs)
	global.CheckStringEnv(global.EnvConfigFilePath, global.ConfigFilePath)
	global.CheckIntEnv(global.EnvQueryThread, global.QueryThread, 16)
	global.CheckDurationEnv(global.EnvCacheCleanInterval, global.CacheCleanInterval)
	global.CheckDurationEnv(global.EnvRequestTimeout, global.RequestTimeout)
	global.CheckIntEnv(global.EnvGzipLevel, global.GzipLevel, 5)

	configManager := configImpl.NewManager()
	if err := configManager.Init(); err != nil {
		fmt.Printf("fail to initialize configuration file: %v", err)
		return
	}

	applicationConfig := configManager.GetConfig()
	logger := loggerImpl.NewLogger()
	logger.Init(
		applicationConfig.GlobalConfig.LogConfig.Path,
		global.LogName,
		applicationConfig.GlobalConfig.LogConfig.Level,
		applicationConfig.GlobalConfig.LogConfig,
	)

	logger.Info(" _____     _           _____             _   _")
	logger.Info("|     |___| |_ ___ ___|  _  |___ ___ _ _|_|_| |___ ___")
	logger.Info("| | | | -_|  _| .'|  _|   __|  _| . | | | | . | -_|  _|")
	logger.Info("|_|_|_|___|_| |__,|_| |__|  |_| |___|\\_/|_|___|___|_|")
	logger.Infof("                     Copyright © %d-%d Half_nothing", global.BeginYear, time.Now().Year())
	logger.Infof("                                   MetarProvider v%s", global.AppVersion)

	cleaner := cleanerImpl.NewCleaner(logger)
	cleaner.Init()
	defer cleaner.Clean()

	metarManagerMemoryCache := cache.NewMemoryCache[*string](*global.CacheCleanInterval)
	cleaner.Add("Metar Cache", func(ctx context.Context) error {
		metarManagerMemoryCache.Close()
		return nil
	})
	metarManager := metar.NewManager(
		logger,
		utils.Filter(applicationConfig.ProviderConfigs, func(providerConfig *config.ProviderConfig) bool {
			return providerConfig.Type == config.ProviderTypeMetar.Value
		}),
		metarManagerMemoryCache,
	)

	tafManagerMemoryCache := cache.NewMemoryCache[*string](*global.CacheCleanInterval)
	cleaner.Add("Taf Cache", func(ctx context.Context) error {
		tafManagerMemoryCache.Close()
		return nil
	})
	tafManager := metar.NewManager(
		logger,
		utils.Filter(applicationConfig.ProviderConfigs, func(providerConfig *config.ProviderConfig) bool {
			return providerConfig.Type == config.ProviderTypeTaf.Value
		}),
		tafManagerMemoryCache,
	)

	applicationContent := content.NewApplicationContentBuilder().
		SetConfigManager(configManager).
		SetCleaner(cleaner).
		SetLogger(logger).
		SetMetarManager(metarManager).
		SetTafManager(tafManager).
		Build()

	if applicationConfig.ServerConfig.GrpcServerConfig.Enable {
		go func() {
			address := fmt.Sprintf("%s:%d", applicationConfig.ServerConfig.GrpcServerConfig.Host, applicationConfig.ServerConfig.GrpcServerConfig.Port)
			lis, err := net.Listen("tcp", address)
			if err != nil {
				logger.Fatalf("gRPC fail to listen: %v", err)
				return
			}
			s := grpc.NewServer()
			grpcServer := grpcImpl.NewMetarServer(logger, metarManager, tafManager)
			pb.RegisterMetarServer(s, grpcServer)
			reflection.Register(s)
			cleaner.Add("gRPC Server", func(ctx context.Context) error {
				timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()
				cleanOver := make(chan struct{})
				go func() {
					s.GracefulStop()
					cleanOver <- struct{}{}
				}()
				select {
				case <-timeoutCtx.Done():
					s.Stop()
				case <-cleanOver:
				}
				return nil
			})
			logger.Infof("gRPC server listening at %v", lis.Addr())
			if err := s.Serve(lis); err != nil {
				logger.Fatalf("gRPC failed to serve: %v", err)
				return
			}
		}()
	}

	server.StartServer(applicationContent)
}
