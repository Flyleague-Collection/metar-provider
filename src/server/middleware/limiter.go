package middleware

import (
	"metar-provider/src/interfaces/server/dto"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	windowSize     time.Duration
	maxRequests    int
	requestRecords map[string][]time.Time
	mu             sync.RWMutex
}

// NewSlidingWindowLimiter 创建滑动窗口限流器
func NewSlidingWindowLimiter(windowSize time.Duration, maxRequests int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		windowSize:     windowSize,
		maxRequests:    maxRequests,
		requestRecords: make(map[string][]time.Time),
	}
}

// Allow 检查是否允许请求
func (l *SlidingWindowLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	if _, exists := l.requestRecords[key]; !exists {
		l.requestRecords[key] = make([]time.Time, 0, l.maxRequests*2)
	}

	windowStart := now.Add(-l.windowSize)
	records := l.requestRecords[key]
	for len(records) > 0 && records[0].Before(windowStart) {
		records = records[1:]
	}

	if len(records) >= l.maxRequests {
		l.requestRecords[key] = records
		return false
	}

	records = append(records, now)
	l.requestRecords[key] = records
	return true
}

func (l *SlidingWindowLimiter) StartCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			l.cleanup()
		}
	}()
}

func (l *SlidingWindowLimiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	threshold := now.Add(-2 * l.windowSize)

	for key, records := range l.requestRecords {
		if len(records) > 0 && records[len(records)-1].Before(threshold) {
			delete(l.requestRecords, key)
		}
	}
}

func CombinedKeyFunc(c echo.Context) string {
	return c.RealIP() + "|" + c.Path()
}

// RateLimitMiddleware Echo 限流中间件
func RateLimitMiddleware(limiter *SlidingWindowLimiter, keyFunc func(c echo.Context) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := keyFunc(c)

			if !limiter.Allow(key) {
				return dto.ErrorResponse(c, dto.ErrRateLimitExceeded)
			}

			return next(c)
		}
	}
}
