// Package global
package global

import (
	"flag"
)

var (
	NoLogs         = flag.Bool("no_logs", false, "Disable logging to file")
	ConfigFilePath = flag.String("config", "./config.yaml", "Path to configuration file")
	QueryThread    = flag.Int("thread", 16, "Number of query threads")
)

const (
	AppVersion    = "1.0.0"
	ConfigVersion = "1.0.0"

	SigningMethod = "HS512"

	BeginYear = 2025

	DefaultFilePermissions     = 0644
	DefaultDirectoryPermission = 0755

	LogName = "MAIN"

	EnvNoLogs         = "NO_LOGS"
	EnvConfigFilePath = "CONFIG_FILE_PATH"
	EnvQueryThread    = "QUERY_THREAD"
)
