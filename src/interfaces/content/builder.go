// Package content
package content

import (
	c "metar-service/src/interfaces/config"
	"metar-service/src/interfaces/metar"

	"half-nothing.cn/service-core/interfaces/cleaner"
	"half-nothing.cn/service-core/interfaces/config"
	"half-nothing.cn/service-core/interfaces/logger"
)

type ApplicationContentBuilder struct {
	content *ApplicationContent
}

func NewApplicationContentBuilder() *ApplicationContentBuilder {
	return &ApplicationContentBuilder{
		content: &ApplicationContent{},
	}
}

func (builder *ApplicationContentBuilder) SetConfigManager(configManager config.ManagerInterface[*c.Config]) *ApplicationContentBuilder {
	builder.content.configManager = configManager
	return builder
}

func (builder *ApplicationContentBuilder) SetCleaner(cleaner cleaner.Interface) *ApplicationContentBuilder {
	builder.content.cleaner = cleaner
	return builder
}

func (builder *ApplicationContentBuilder) SetLogger(logger logger.Interface) *ApplicationContentBuilder {
	builder.content.logger = logger
	return builder
}

func (builder *ApplicationContentBuilder) SetMetarManager(metarManager metar.ManagerInterface) *ApplicationContentBuilder {
	builder.content.metarManager = metarManager
	return builder
}

func (builder *ApplicationContentBuilder) SetTafManager(tafManager metar.ManagerInterface) *ApplicationContentBuilder {
	builder.content.tafManager = tafManager
	return builder
}

func (builder *ApplicationContentBuilder) Build() *ApplicationContent {
	return builder.content
}
