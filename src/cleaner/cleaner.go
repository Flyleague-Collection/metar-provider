package base

import (
	"context"
	"fmt"
	"metar-provider/src/interfaces/cleaner"
	"metar-provider/src/interfaces/logger"
	"metar-provider/src/utils"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type shutdownFunc struct {
	callback cleaner.ShutdownCallback
	name     string
}

// Cleaner 负责管理应用关闭时需要执行的清理任务
type Cleaner struct {
	cleaners       []*shutdownFunc          // 需要执行的清理函数列表
	mu             sync.Mutex               // 保护cleaners和cleaning状态的互斥锁
	cleaning       bool                     // 标识是否正在进行清理
	loggerShutdown cleaner.ShutdownCallback // 日志系统关闭回调
	logger         logger.Interface         // 日志接口
}

// NewCleaner 创建一个新的Cleaner实例
func NewCleaner(logger logger.Interface) *Cleaner {
	return &Cleaner{
		cleaners:       make([]*shutdownFunc, 0),
		loggerShutdown: logger.ShutdownCallback,
		logger:         logger,
	}
}

// Add 添加一个清理函数到清理队列中
// 如果清理已经开始，则忽略新的清理函数
func (c *Cleaner) Add(name string, callback cleaner.ShutdownCallback) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查是否已经在清理过程中
	if c.cleaning {
		c.logger.Warn("Cleaner is already shutting down, ignoring new cleaner")
		return
	}

	// 添加清理函数到队列
	c.cleaners = append(c.cleaners, &shutdownFunc{callback: callback, name: name})
	c.logger.Debugf("Adding cleaner #%d %s(%p)", len(c.cleaners), name, callback)
}

// Clean 执行所有已注册的清理函数
// 清理函数会以相反的顺序执行，确保后添加的先执行
func (c *Cleaner) Clean() {
	// 锁定并复制清理函数列表，防止在清理过程中被修改
	c.mu.Lock()
	if c.cleaning {
		c.mu.Unlock()
		return
	}
	c.cleaning = true // 标记为清理中，阻止后续Add操作
	cleanersCopy := make([]*shutdownFunc, len(c.cleaners))
	copy(cleanersCopy, c.cleaners)
	c.mu.Unlock()

	c.logger.Debugf("Starting cleanup of %d registered functions", len(cleanersCopy))

	// 执行所有清理函数并收集错误
	var errs []error
	utils.ReverseForEach(cleanersCopy, func(idx int, sf *shutdownFunc) {
		c.logger.Debugf("Invoking cleaner #%d (%s)", idx+1, sf.name)

		// 为每个清理函数设置10秒超时
		timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFunc()

		// 在单独的goroutine中执行清理函数，防止某个清理函数阻塞整个清理过程
		done := make(chan error, 1)
		go func() {
			done <- sf.callback(timeoutCtx)
		}()

		// 等待清理函数完成或超时
		select {
		case err := <-done:
			// 执行清理函数并处理错误
			if err != nil {
				c.logger.Errorf("Cleaner #%d (%s) failed: %v", idx+1, sf.name, err)
				errs = append(errs, err)
			}
		case <-timeoutCtx.Done():
			c.logger.Errorf("Cleaner #%d (%s) timed out", idx+1, sf.name)
			errs = append(errs, timeoutCtx.Err())
		}
	})

	// 输出清理结果
	if len(errs) > 0 {
		c.logger.Errorf("%d errors occurred during cleanup:", len(errs))
		for i, err := range errs {
			c.logger.Errorf("Error %d: %v", i+1, err)
		}
	} else {
		c.logger.Debug("All cleaners executed successfully")
	}
	c.logger.Info("Cleanup finished, server offline")

	// 关闭日志系统
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := c.loggerShutdown(shutdownCtx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LOGGER SHUTDOWN ERROR: %v\n", err)
	}

	syscall.Exit(0)
}

// Init 初始化信号监听器，监听中断信号(SIGINT, SIGTERM)
// 收到信号后自动触发清理流程
func (c *Cleaner) Init() {
	// 创建监听中断信号和终止信号的上下文
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// 启动goroutine监听信号
	go func() {
		<-ctx.Done() // 等待信号到达
		stop()       // 停止信号监听
		c.logger.Info("Received interrupt signal, shutting down")
		c.Clean() // 执行清理
	}()
}
