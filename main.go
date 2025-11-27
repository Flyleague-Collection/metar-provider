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
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"half-nothing.cn/service-core/cache"
	"half-nothing.cn/service-core/cleaner"
	"half-nothing.cn/service-core/config"
	"half-nothing.cn/service-core/discovery"
	"half-nothing.cn/service-core/interfaces/global"
	"half-nothing.cn/service-core/logger"
	"half-nothing.cn/service-core/utils"
)

func main() {
	global.CheckFlags()
	global.CheckIntEnv(g.EnvQueryThread, g.QueryThread, 16)
	global.CheckDurationEnv(g.EnvCacheCleanInterval, g.CacheCleanInterval)
	global.CheckDurationEnv(g.EnvRequestTimeout, g.RequestTimeout)
	global.CheckIntEnv(g.EnvGzipLevel, g.GzipLevel, 5)

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
	lg.Infof("                     Copyright © %d-%d Half_nothing", global.BeginYear, time.Now().Year())
	lg.Infof("                                   MetarProvider v%s", g.AppVersion)

	cl := cleaner.NewCleaner(lg)
	cl.Init()

	metarManagerMemoryCache := cache.NewMemoryCache[*string](*g.CacheCleanInterval)
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

	tafManagerMemoryCache := cache.NewMemoryCache[*string](*g.CacheCleanInterval)
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

	go server.StartServer(applicationContent)

	if applicationConfig.ServerConfig.GrpcServerConfig.Enable {
		started := make(chan struct{})
		go func() {
			address := fmt.Sprintf("%s:%d", applicationConfig.ServerConfig.GrpcServerConfig.Host, applicationConfig.ServerConfig.GrpcServerConfig.Port)
			lis, err := net.Listen("tcp", address)
			if err != nil {
				lg.Fatalf("gRPC fail to listen: %v", err)
				return
			}
			s := grpc.NewServer()
			grpcServer := grpcImpl.NewMetarServer(lg, metarManager, tafManager)
			pb.RegisterMetarServer(s, grpcServer)
			reflection.Register(s)
			cl.Add("gRPC Server", func(ctx context.Context) error {
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
			lg.Infof("gRPC server listening at %v", lis.Addr())
			close(started)
			if err := s.Serve(lis); err != nil {
				lg.Fatalf("gRPC failed to serve: %v", err)
				return
			}
		}()

		<-started
		version, _ := global.NewVersion(g.AppVersion)
		service := discovery.NewServiceDiscovery(
			lg,
			"metar-service",
			applicationConfig.ServerConfig.GrpcServerConfig.Port,
			version,
		)
		if err := service.Start(); err != nil {
			lg.Fatalf("fail to start service discovery: %v", err)
			cl.Clean()
			return
		}
		cl.Add("Service Discovery", service.Stop)
	}

	cl.Wait()
}
