// Package global
package global

import (
	"flag"
	"time"
)

var (
	QueryThread        = flag.Int("thread", 16, "Number of query threads")
	CacheCleanInterval = flag.Duration("cache_clean_interval", 30*time.Minute, "cache cleanup interval")
	RequestTimeout     = flag.Duration("request_timeout", 30*time.Second, "Request timeout")
	GzipLevel          = flag.Int("gzip_level", 5, "GZip level")
)

const (
	AppVersion    = "0.3.0"
	ConfigVersion = "0.3.0"

	EnvQueryThread        = "QUERY_THREAD"
	EnvCacheCleanInterval = "CACHE_CLEAN_INTERVAL"
	EnvRequestTimeout     = "REQUEST_TIMEOUT"
	EnvGzipLevel          = "GZIP_LEVEL"
)
