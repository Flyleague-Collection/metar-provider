// Package global
package global

import (
	"flag"
	"time"
)

var (
	QueryThread        = flag.Int("thread", 16, "Number of query threads")
	CacheCleanInterval = flag.Duration("cache_clean_interval", 30*time.Minute, "cache cleanup interval")
)

const (
	AppVersion    = "0.3.1"
	ConfigVersion = "0.4.0"

	EnvQueryThread        = "QUERY_THREAD"
	EnvCacheCleanInterval = "CACHE_CLEAN_INTERVAL"
)
