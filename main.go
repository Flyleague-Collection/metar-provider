package main

import (
	"flag"
	"fmt"
	cleanerImpl "metar-provider/src/cleaner"
	configImpl "metar-provider/src/config"
	"metar-provider/src/interfaces/global"
	loggerImpl "metar-provider/src/logger"
	"time"
)

func main() {
	flag.Parse()

	global.CheckBoolEnv(global.EnvNoLogs, global.NoLogs)
	global.CheckStringEnv(global.EnvConfigFilePath, global.ConfigFilePath)
	global.CheckIntEnv(global.EnvQueryThread, global.QueryThread, 16)

	configManager := configImpl.NewManager()
	if err := configManager.Init(); err != nil {
		fmt.Printf("fail to initialize configuration file: %v", err)
		return
	}

	config := configManager.GetConfig()
	logger := loggerImpl.NewLogger()
	logger.Init(
		config.GlobalConfig.LogConfig.Path,
		global.LogName,
		config.GlobalConfig.LogConfig.Level,
		config.GlobalConfig.LogConfig,
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

}
