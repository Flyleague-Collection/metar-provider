// Package content
package content

import (
	"metar-provider/src/interfaces/cleaner"
	"metar-provider/src/interfaces/config"
	"metar-provider/src/interfaces/logger"
	"metar-provider/src/interfaces/metar"
)

// ApplicationContent 应用程序上下文结构体，包含所有核心组件的接口
type ApplicationContent struct {
	configManager config.ManagerInterface // 配置管理器
	cleaner       cleaner.Interface       // 清理器
	logger        logger.Interface        // 日志
	metarManager  metar.ManagerInterface  // METAR气象数据管理器
}

func (app *ApplicationContent) ConfigManager() config.ManagerInterface {
	return app.configManager
}

func (app *ApplicationContent) Cleaner() cleaner.Interface { return app.cleaner }

func (app *ApplicationContent) Logger() logger.Interface { return app.logger }

func (app *ApplicationContent) MetarManager() metar.ManagerInterface { return app.metarManager }
