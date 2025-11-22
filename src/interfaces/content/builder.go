// Package content
package content

import (
	"metar-provider/src/interfaces/cleaner"
	"metar-provider/src/interfaces/config"
	"metar-provider/src/interfaces/logger"
	"metar-provider/src/interfaces/metar"
)

type ApplicationContentBuilder struct {
	content *ApplicationContent
}

func NewApplicationContentBuilder() *ApplicationContentBuilder {
	return &ApplicationContentBuilder{
		content: &ApplicationContent{},
	}
}

func (builder *ApplicationContentBuilder) SetConfigManager(configManager config.ManagerInterface) *ApplicationContentBuilder {
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
