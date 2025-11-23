// Package logger
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"metar-provider/src/interfaces/config"
	"metar-provider/src/interfaces/global"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LevelFatal 自定义日志级别 FATAL
const (
	LevelFatal slog.Level = 12
)

// AsyncHandler 异步日志处理器，支持并发安全的日志写入
type AsyncHandler struct {
	ch       chan []byte    // 日志数据通道
	logName  string         // 日志名称
	writer   io.Writer      // 日志写入目标
	attrs    []slog.Attr    // 全局属性列表
	group    string         // 当前分组名称
	logLevel slog.Level     // 日志级别
	wg       sync.WaitGroup // 用于等待后台协程结束
}

// NewAsyncHandler 创建一个新的异步日志处理器
//
// 参数:
//   - logPath: 日志文件路径，用于指定日志文件的存储位置
//   - logName: 日志记录器名称，会在每条日志中显示
//   - logLevel: 日志级别，用于过滤不同级别的日志输出
//   - logConfig: 日志配置项，用于指定日志的轮转、压缩等选项
//
// 返回值:
//   - *AsyncHandler: 返回初始化后的异步日志处理器实例
//
// 注意: 此函数为内部函数，不建议直接调用。如需创建日志记录器，请使用 NewLogger 函数创建 Logger 结构体并调用 Init 方法进行初始化
func NewAsyncHandler(logPath, logName string, logLevel slog.Level, logConfig *config.LogConfig) *AsyncHandler {
	h := &AsyncHandler{
		ch:       make(chan []byte, 1024),
		logLevel: logLevel,
		logName:  strings.ToUpper(logName),
	}

	if *global.NoLogs {
		h.writer = os.Stdout
	} else if logConfig.Rotate {
		h.writer = io.MultiWriter(os.Stdout, &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    logConfig.MaxSize,
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,
			Compress:   logConfig.Compress,
			LocalTime:  logConfig.LocalTime,
		})
	} else {
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, global.DefaultFilePermissions)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to open log file %s: %v\n", logPath, err)
			h.writer = os.Stdout
		}
		h.writer = io.MultiWriter(os.Stdout, file)
	}

	go h.startWorker()
	return h
}

// startWorker 后台工作者，从通道读取日志数据并写入
func (h *AsyncHandler) startWorker() {
	h.wg.Add(1)
	defer h.wg.Done()

	for data := range h.ch {
		_, _ = h.writer.Write(data)
	}
}

// Enabled 判断指定级别的日志是否应该被处理
func (h *AsyncHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.logLevel
}

// Handle 处理日志记录
func (h *AsyncHandler) Handle(_ context.Context, r slog.Record) error {
	// 根据日志级别设置不同的颜色显示
	level := r.Level.String()
	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	case LevelFatal:
		level = color.HiRedString("FATAL")
	}

	// 构造日志行格式：时间 | 记录器 | 级别 | 消息
	line := fmt.Sprintf(
		"%s | %-5s | %-5s | %s",
		color.GreenString(r.Time.Format("2006-01-02T15:04:05")),
		h.logName,
		level,
		color.CyanString(r.Message),
	)

	for _, attr := range h.attrs {
		line += color.CyanString(fmt.Sprintf(" %s=%v", attr.Key, attr.Value))
	}

	r.Attrs(func(attr slog.Attr) bool {
		line += color.CyanString(fmt.Sprintf(" %s=%v", attr.Key, attr.Value))
		return true
	})

	line += "\n"

	h.Write([]byte(line))
	return nil
}

// WithAttrs 返回带有附加属性的新处理器
func (h *AsyncHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	newAttrs = append(newAttrs, h.attrs...)
	newAttrs = append(newAttrs, attrs...)

	return &AsyncHandler{
		writer:   h.writer,
		attrs:    newAttrs,
		group:    h.group,
		logLevel: h.logLevel,
	}
}

// WithGroup 返回带有指定分组名称的新处理器
func (h *AsyncHandler) WithGroup(name string) slog.Handler {
	return &AsyncHandler{
		writer:   h.writer,
		attrs:    h.attrs,
		group:    name,
		logLevel: h.logLevel,
	}
}

// Write 将日志数据写入通道
func (h *AsyncHandler) Write(p []byte) {
	pb := make([]byte, len(p))
	copy(pb, p)
	h.ch <- pb
}

// Close 关闭日志处理器，确保所有日志都被写入
func (h *AsyncHandler) Close() error {
	close(h.ch)
	h.wg.Wait()

	if f, ok := h.writer.(*os.File); ok {
		_ = f.Sync()
	}
	return nil
}

// Logger 日志记录器封装结构体
type Logger struct {
	handler *AsyncHandler
	logger  *slog.Logger
}

// NewLogger 创建新的日志记录器实例
func NewLogger() *Logger {
	return &Logger{
		logger: nil,
	}
}

// Init 初始化日志记录器
// logPath: 日志文件路径
// logName: 日志记录器名称
// debug: 是否启用调试模式
// noLogs: 是否只输出到标准输出
func (lg *Logger) Init(logPath, logName, logLevel string, logConfig *config.LogConfig) {
	lg.handler = NewAsyncHandler(logPath, logName, slog.LevelInfo, logConfig)

	switch strings.ToLower(logLevel) {
	case "debug":
		lg.handler.logLevel = slog.LevelDebug
	case "info":
		lg.handler.logLevel = slog.LevelInfo
	case "warn":
		lg.handler.logLevel = slog.LevelWarn
	case "error":
		lg.handler.logLevel = slog.LevelError
	case "fatal":
		lg.handler.logLevel = LevelFatal
	}

	lg.logger = slog.New(lg.handler)

	lg.Debugf("%s logger initialized", strings.ToUpper(logName))
}

// ShutdownCallback 获取关闭回调
func (lg *Logger) ShutdownCallback(context.Context) error {
	return lg.handler.Close()
}

// LogHandler 获取底层 slog.Logger 实例
func (lg *Logger) LogHandler() *slog.Logger {
	return lg.logger
}

func (lg *Logger) Debug(msg string) {
	lg.logger.Debug(msg)
}

func (lg *Logger) Debugf(msg string, v ...interface{}) {
	lg.logger.Debug(fmt.Sprintf(msg, v...))
}

func (lg *Logger) Info(msg string) {
	lg.logger.Info(msg)
}

func (lg *Logger) Infof(msg string, v ...interface{}) {
	lg.logger.Info(fmt.Sprintf(msg, v...))
}

func (lg *Logger) Warn(msg string) {
	lg.logger.Warn(msg)
}

func (lg *Logger) Warnf(msg string, v ...interface{}) {
	lg.logger.Warn(fmt.Sprintf(msg, v...))
}

func (lg *Logger) Error(msg string) {
	lg.logger.Error(msg)
}

func (lg *Logger) Errorf(msg string, v ...interface{}) {
	lg.logger.Error(fmt.Sprintf(msg, v...))
}

func (lg *Logger) Fatal(msg string) {
	lg.logger.Log(context.Background(), LevelFatal, msg)
}

func (lg *Logger) Fatalf(msg string, v ...interface{}) {
	lg.logger.Log(context.Background(), LevelFatal, fmt.Sprintf(msg, v...))
}
