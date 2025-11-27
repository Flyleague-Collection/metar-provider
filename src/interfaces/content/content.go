// Package content
package content

import (
	c "metar-service/src/interfaces/config"
	"metar-service/src/interfaces/metar"

	"half-nothing.cn/service-core/interfaces/cleaner"
	"half-nothing.cn/service-core/interfaces/config"
	"half-nothing.cn/service-core/interfaces/logger"
)

// ApplicationContent 应用程序上下文结构体，包含所有核心组件的接口
type ApplicationContent struct {
	configManager config.ManagerInterface[*c.Config] // 配置管理器
	cleaner       cleaner.Interface                  // 清理器
	logger        logger.Interface                   // 日志
	metarManager  metar.ManagerInterface             // METAR气象数据管理器
	tafManager    metar.ManagerInterface             // TAF天气预报数据管理器
}

func (app *ApplicationContent) ConfigManager() config.ManagerInterface[*c.Config] {
	return app.configManager
}

func (app *ApplicationContent) Cleaner() cleaner.Interface { return app.cleaner }

func (app *ApplicationContent) Logger() logger.Interface { return app.logger }

func (app *ApplicationContent) MetarManager() metar.ManagerInterface { return app.metarManager }

func (app *ApplicationContent) TafManager() metar.ManagerInterface { return app.tafManager }
