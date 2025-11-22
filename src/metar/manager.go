// Package metar
package metar

import (
	"metar-provider/src/interfaces/cache"
	"metar-provider/src/interfaces/config"
	"metar-provider/src/interfaces/global"
	"metar-provider/src/interfaces/logger"
	"metar-provider/src/interfaces/metar"
	"metar-provider/src/utils"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type Manager struct {
	logger       logger.Interface
	providers    []metar.ProviderInterface
	cache        cache.Interface[*string]
	requestGroup singleflight.Group
}

func NewManager(
	lg logger.Interface,
	providerConfigs []*config.ProviderConfig,
	cache cache.Interface[*string],
) *Manager {
	manager := &Manager{
		logger:    logger.NewLoggerAdapter(lg, "SourceManager"),
		providers: make([]metar.ProviderInterface, 0),
		cache:     cache,
	}

	utils.ForEach(providerConfigs, func(index int, providerConfig *config.ProviderConfig) {
		manager.providers = append(manager.providers, NewProvider(lg, providerConfig))
	})

	return manager
}
func (m *Manager) Query(icao string) (string, error) {
	if icao == "" || len(icao) != 4 {
		return "", metar.ErrICAOInvalid
	}

	if data, ok := m.cache.Get(icao); ok {
		if data != nil {
			return *data, nil
		}
		return "", metar.ErrTargetNotFound
	}

	result, err, _ := m.requestGroup.Do(icao, func() (interface{}, error) {
		for _, provider := range m.providers {
			data, err := provider.Get(icao)
			if err != nil {
				continue
			}
			m.setCache(icao, &data)
			return data, nil
		}
		m.setCache(icao, nil)
		return "", metar.ErrTargetNotFound
	})

	if err != nil {
		return "", err
	}

	// 将结果按行分割，去除空行，再重新组合成单行字符串
	lines := strings.Split(result.(string), "\n")
	nonEmptyLines := utils.Filter(lines, func(line string) bool {
		return strings.TrimSpace(line) != ""
	})
	// 对每一行进行trim处理后再连接
	utils.Map(nonEmptyLines, func(line string) string {
		return strings.TrimSpace(line)
	})
	data := strings.Join(nonEmptyLines, " ")

	return data, nil
}

func (m *Manager) BatchQuery(icaos []string) []string {
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}
	data := make([]string, 0, len(icaos))
	limiter := make(chan struct{}, *global.QueryThread)

	for _, icao := range icaos {
		wg.Add(1)
		limiter <- struct{}{}
		go func() {
			defer func() {
				<-limiter
				wg.Done()
			}()
			d, err := m.Query(icao)
			if err != nil {
				return
			}
			lock.Lock()
			data = append(data, d)
			lock.Unlock()
		}()
	}
	wg.Wait()

	return data
}

func (m *Manager) setCache(icao string, metar *string) {
	currentTime := time.Now()
	minute := currentTime.Minute()
	var addMinutes int
	if minute < 30 {
		addMinutes = 30 - minute
	} else {
		addMinutes = 60 - minute
	}
	m.cache.SetWithTTL(icao, metar, time.Duration(addMinutes)*time.Minute)
}
