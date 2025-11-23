// Package cleaner 提供清理器接口，用于管理可清理资源
package cleaner

import (
	"context"
)

type ShutdownCallback func(ctx context.Context) error

// Interface 定义了清理器接口，用于初始化、添加和执行清理操作
type Interface interface {
	Init()
	Wait()
	Add(name string, callback ShutdownCallback)
	Clean()
}
