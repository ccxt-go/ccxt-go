package ccxt

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogLevel 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Logger 日志接口
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	Sync() error
}

// LogManager 日志管理器
type LogManager struct {
	logger     *zap.Logger
	level      LogLevel
	filePath   string
	maxSize    int
	maxAge     int
	maxBackups int
	compress   bool
	mu         sync.RWMutex
}

var (
	globalLogManager *LogManager
	logOnce          sync.Once
)

// GetLogManager 获取全局日志管理器
func GetLogManager() *LogManager {
	logOnce.Do(func() {
		globalLogManager = NewLogManager()
	})
	return globalLogManager
}

// NewLogManager 创建新的日志管理器
func NewLogManager() *LogManager {
	lm := &LogManager{
		level:      INFO,
		filePath:   "ccxt-go.log",
		maxSize:    100, // MB
		maxAge:     30,  // days
		maxBackups: 10,
		compress:   true,
	}

	lm.initLogger()
	return lm
}

// initLogger 初始化日志器
func (lm *LogManager) initLogger() {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	// 创建日志目录
	dir := filepath.Dir(lm.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(fmt.Sprintf("创建日志目录失败: %v", err))
	}

	// 配置日志轮转
	writer := &lumberjack.Logger{
		Filename:   lm.filePath,
		MaxSize:    lm.maxSize,
		MaxAge:     lm.maxAge,
		MaxBackups: lm.maxBackups,
		Compress:   lm.compress,
	}

	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建编码器
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 创建核心
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(writer), lm.getZapLevel()),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), lm.getZapLevel()),
	)

	// 创建日志器
	lm.logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

// getZapLevel 获取zap日志级别
func (lm *LogManager) getZapLevel() zapcore.Level {
	switch lm.level {
	case DEBUG:
		return zapcore.DebugLevel
	case INFO:
		return zapcore.InfoLevel
	case WARN:
		return zapcore.WarnLevel
	case ERROR:
		return zapcore.ErrorLevel
	case FATAL:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// SetLevel 设置日志级别
func (lm *LogManager) SetLevel(level LogLevel) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.level = level
	lm.initLogger() // 重新初始化
}

// SetFilePath 设置日志文件路径
func (lm *LogManager) SetFilePath(path string) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.filePath = path
	lm.initLogger() // 重新初始化
}

// SetMaxSize 设置日志文件最大大小
func (lm *LogManager) SetMaxSize(size int) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.maxSize = size
	lm.initLogger() // 重新初始化
}

// SetMaxAge 设置日志文件最大保存天数
func (lm *LogManager) SetMaxAge(age int) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.maxAge = age
	lm.initLogger() // 重新初始化
}

// SetMaxBackups 设置最大备份文件数
func (lm *LogManager) SetMaxBackups(backups int) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.maxBackups = backups
	lm.initLogger() // 重新初始化
}

// SetCompress 设置是否压缩备份文件
func (lm *LogManager) SetCompress(compress bool) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.compress = compress
	lm.initLogger() // 重新初始化
}

// Debug 记录调试日志
func (lm *LogManager) Debug(msg string, fields ...zap.Field) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if lm.logger != nil {
		lm.logger.Debug(msg, fields...)
	}
}

// Info 记录信息日志
func (lm *LogManager) Info(msg string, fields ...zap.Field) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if lm.logger != nil {
		lm.logger.Info(msg, fields...)
	}
}

// Warn 记录警告日志
func (lm *LogManager) Warn(msg string, fields ...zap.Field) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if lm.logger != nil {
		lm.logger.Warn(msg, fields...)
	}
}

// Error 记录错误日志
func (lm *LogManager) Error(msg string, fields ...zap.Field) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if lm.logger != nil {
		lm.logger.Error(msg, fields...)
	}
}

// Fatal 记录致命错误日志
func (lm *LogManager) Fatal(msg string, fields ...zap.Field) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if lm.logger != nil {
		lm.logger.Fatal(msg, fields...)
	}
}

// With 创建带字段的日志器
func (lm *LogManager) With(fields ...zap.Field) Logger {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if lm.logger != nil {
		return &LogManager{logger: lm.logger.With(fields...)}
	}
	return lm
}

// Sync 同步日志
func (lm *LogManager) Sync() error {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if lm.logger != nil {
		return lm.logger.Sync()
	}
	return nil
}

// Close 关闭日志管理器
func (lm *LogManager) Close() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if lm.logger != nil {
		err := lm.logger.Sync()
		lm.logger = nil
		return err
	}
	return nil
}

// 便捷函数
func Debug(msg string, fields ...zap.Field) {
	GetLogManager().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	GetLogManager().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GetLogManager().Warn(msg, fields...)
}

func LogError(msg string, fields ...zap.Field) {
	GetLogManager().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GetLogManager().Fatal(msg, fields...)
}

// 带调用者信息的日志函数
func DebugWithCaller(msg string, fields ...zap.Field) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fields = append(fields, zap.String("caller", fmt.Sprintf("%s:%d", filepath.Base(file), line)))
	}
	Debug(msg, fields...)
}

func InfoWithCaller(msg string, fields ...zap.Field) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fields = append(fields, zap.String("caller", fmt.Sprintf("%s:%d", filepath.Base(file), line)))
	}
	Info(msg, fields...)
}

func WarnWithCaller(msg string, fields ...zap.Field) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fields = append(fields, zap.String("caller", fmt.Sprintf("%s:%d", filepath.Base(file), line)))
	}
	Warn(msg, fields...)
}

func LogErrorWithCaller(msg string, fields ...zap.Field) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fields = append(fields, zap.String("caller", fmt.Sprintf("%s:%d", filepath.Base(file), line)))
	}
	LogError(msg, fields...)
}

// 交易所专用日志函数
func LogExchangeRequest(exchangeId, method, url string, duration time.Duration, statusCode int, fields ...zap.Field) {
	allFields := []zap.Field{
		zap.String("exchange", exchangeId),
		zap.String("method", method),
		zap.String("url", url),
		zap.Duration("duration", duration),
		zap.Int("status_code", statusCode),
	}
	allFields = append(allFields, fields...)
	Info("Exchange request completed", allFields...)
}

func LogExchangeError(exchangeId, method, url string, err error, fields ...zap.Field) {
	allFields := []zap.Field{
		zap.String("exchange", exchangeId),
		zap.String("method", method),
		zap.String("url", url),
		zap.Error(err),
	}
	allFields = append(allFields, fields...)
	LogError("Exchange request failed", allFields...)
}

func LogWebSocketEvent(exchangeId, eventType, message string, fields ...zap.Field) {
	allFields := []zap.Field{
		zap.String("exchange", exchangeId),
		zap.String("event_type", eventType),
		zap.String("message", message),
	}
	allFields = append(allFields, fields...)
	Debug("WebSocket event", allFields...)
}

// 性能监控日志
func LogPerformance(operation string, duration time.Duration, fields ...zap.Field) {
	allFields := []zap.Field{
		zap.String("operation", operation),
		zap.Duration("duration", duration),
	}
	allFields = append(allFields, fields...)

	if duration > time.Second {
		Warn("Slow operation", allFields...)
	} else {
		Debug("Operation completed", allFields...)
	}
}

// 内存使用监控
func LogMemoryUsage(operation string, fields ...zap.Field) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	allFields := []zap.Field{
		zap.String("operation", operation),
		zap.Uint64("alloc_mb", m.Alloc/1024/1024),
		zap.Uint64("total_alloc_mb", m.TotalAlloc/1024/1024),
		zap.Uint64("sys_mb", m.Sys/1024/1024),
		zap.Uint32("num_gc", m.NumGC),
	}
	allFields = append(allFields, fields...)

	Debug("Memory usage", allFields...)
}
