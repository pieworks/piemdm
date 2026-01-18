package log

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

// PrettyHandler 是一个自定义的 slog.Handler
type PrettyHandler struct {
	opts slog.HandlerOptions
	w    io.Writer
	mu   *sync.Mutex
	// 用于处理 WithGroup 和 WithAttrs
	groups []string    // stores group names
	attrs  []slog.Attr // stores attributes from WithAttrs
}

// NewPrettyHandler 创建一个新的 PrettyHandler
func NewPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	if opts.Level == nil {
		opts.Level = slog.LevelInfo
	}
	return &PrettyHandler{
		opts:   *opts,
		w:      w,
		mu:     new(sync.Mutex),
		groups: make([]string, 0),
		attrs:  make([]slog.Attr, 0),
	}
}

// Enabled 检查日志级别是否应该被记录
func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

// Handle 处理并格式化日志记录
func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	// 在打印日志记录之前先写入一个空行
	_, err := h.w.Write([]byte("\n"))
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	// 1. 格式化第一行：时间、级别、消息
	buf.WriteString(r.Time.Format(time.DateTime))
	buf.WriteByte(' ')
	buf.WriteString(r.Level.String())
	buf.WriteByte(' ')
	buf.WriteString(r.Message)

	// 2. 添加来源 (source) 信息到第一行
	sourceAttrPrinted := false
	if h.opts.AddSource && r.PC != 0 {
		// 获取完整的调用栈，跳过log包内部的帧
		pcs := make([]uintptr, 20)
		n := runtime.Callers(0, pcs)
		if n > 0 {
			frames := runtime.CallersFrames(pcs[:n])
			for {
				f, more := frames.Next()
				// 跳过log包内部的调用帧，找到实际业务代码的调用位置
				if !strings.Contains(f.File, "/pkg/log/") && f.Function != "" &&
					!strings.Contains(f.Function, "runtime.") &&
					!strings.Contains(f.Function, "slog.") {
					source := fmt.Sprintf(" %s=%s:%d", slog.SourceKey, f.File, f.Line)
					buf.WriteString(source)
					sourceAttrPrinted = true
					break
				}
				if !more {
					break
				}
			}
		}
	}

	// Collect all attributes (including from WithAttrs and groups), excluding the source if already printed
	var allAttrs []slog.Attr

	// Process attributes from WithAttrs, applying groups
	currentAttrsFromWithAttrs := h.attrs

	// Process attributes from the record, applying groups and excluding source if already printed
	var currentAttrsFromRecord []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		// Only add source if it was not already added to the first line by AddSource
		if sourceAttrPrinted && a.Key == slog.SourceKey {
			return true // Skip this source attribute, it's already handled in the first line
		}
		currentAttrsFromRecord = append(currentAttrsFromRecord, a)
		return true
	})

	// Combine and flatten all attributes
	allCombinedAttrs := append(currentAttrsFromWithAttrs, currentAttrsFromRecord...)

	groupPrefix := ""
	for _, group := range h.groups {
		groupPrefix += group + "."
	}
	for _, a := range allCombinedAttrs {
		// Flatten groups within attributes
		if a.Value.Kind() == slog.KindGroup {
			for _, ga := range a.Value.Group() {
				allAttrs = append(allAttrs, slog.Attr{Key: groupPrefix + a.Key + "." + ga.Key, Value: ga.Value})
			}
		} else {
			allAttrs = append(allAttrs, slog.Attr{Key: groupPrefix + a.Key, Value: a.Value})
		}
	}

	// 3. If there are other attributes, add a newline and format them with indentation
	if len(allAttrs) > 0 {
		buf.WriteByte('\n') // Newline after the first line (time, level, msg, source)

		for _, a := range allAttrs {
			buf.WriteString("\t") // Indent
			buf.WriteString(a.Key)
			buf.WriteByte('=')

			// Format value based on type for better readability
			// For complex types, use %+v to print fields
			switch a.Value.Kind() {
			case slog.KindString:
				buf.WriteString(fmt.Sprintf("%q", a.Value.String()))
			case slog.KindAny:
				// Use %+v for structs/complex types for detailed output
				buf.WriteString(fmt.Sprintf("%+v", a.Value.Any()))
			default:
				buf.WriteString(a.Value.String())
			}
			buf.WriteByte('\n')
		}
	} else {
		// If no other attributes, just ensure a newline after the first line
		buf.WriteByte('\n')
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err = h.w.Write(buf.Bytes())
	return err
}

// WithAttrs 返回一个新的 Handler，其中包含已有的和新增的属性
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := *h
	newHandler.attrs = make([]slog.Attr, len(h.attrs)) // Create new slice for safety
	copy(newHandler.attrs, h.attrs)
	newHandler.attrs = append(newHandler.attrs, attrs...)
	return &newHandler
}

// WithGroup 返回一个新的 Handler，其中包含新的组前缀
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	newHandler := *h
	newHandler.groups = make([]string, len(h.groups)) // Create new slice for safety
	copy(newHandler.groups, h.groups)
	newHandler.groups = append(newHandler.groups, name)
	return &newHandler
}

const LOGGER_KEY = "slogLogger"

// Logger 封装 slog.Logger,提供与项目集成的功能
type Logger struct {
	*slog.Logger
}

// NewLog 创建新的 Logger 实例
func NewLog(conf *viper.Viper) *Logger {
	return initSlog(conf)
}

// initSlog 初始化 slog logger
func initSlog(conf *viper.Viper) *Logger {
	// 获取配置
	logPath := conf.GetString("log.log_path")
	logFileName := conf.GetString("log.log_file_name")
	logLevel := conf.GetString("log.log_level")
	encoding := conf.GetString("log.encoding")
	env := conf.GetString("env")

	// 解析日志级别
	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// 配置日志文件轮转
	var writers []io.Writer
	// 直接将 os.Stdout 添加到 writers，由 PrettyHandler 统一处理空行
	writers = append(writers, os.Stdout)

	if logPath != "" && logFileName != "" {
		// 处理日志路径
		// 1. 如果以 ~/ 开头,扩展为用户主目录
		if strings.HasPrefix(logPath, "~/") {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				logPath = filepath.Join(homeDir, logPath[2:])
			}
		} else if !filepath.IsAbs(logPath) {
			// 2. 如果是相对路径,则相对于配置文件所在目录
			configFile := conf.ConfigFileUsed()
			if configFile != "" {
				configDir := filepath.Dir(configFile)
				// 获取配置文件目录的绝对路径
				if absConfigDir, err := filepath.Abs(configDir); err == nil {
					logPath = filepath.Join(absConfigDir, logPath)
				}
			} else {
				// 如果无法获取配置文件路径,则相对于当前工作目录
				if absLogPath, err := filepath.Abs(logPath); err == nil {
					logPath = absLogPath
				}
			}
		}

		// 组合完整的日志文件路径
		logFile := filepath.Join(logPath, logFileName)

		// 确保日志目录存在
		if err := os.MkdirAll(logPath, 0755); err == nil {
			fileWriter := &lumberjack.Logger{
				Filename:   logFile,
				MaxSize:    conf.GetInt("log.max_size"), // MB
				MaxBackups: conf.GetInt("log.max_backups"),
				MaxAge:     conf.GetInt("log.max_age"), // days
				Compress:   conf.GetBool("log.compress"),
			}
			writers = append(writers, fileWriter)
		}
	}

	// 创建多写入器
	multiWriter := io.MultiWriter(writers...)

	// 创建 handler 选项
	handlerOpts := &slog.HandlerOptions{
		// PrettyHandler 会根据 AddSource 决定是否显示源代码位置
		AddSource: true, // Enable source for PrettyHandler and default JSON handler
		Level:     level,
	}

	// 根据环境和配置选择 handler
	var handler slog.Handler
	if encoding == "console" || env != "prod" {
		// 开发环境使用 PrettyHandler,便于阅读,且支持多行属性
		handler = NewPrettyHandler(multiWriter, handlerOpts)
	} else {
		// 生产环境使用 JSON 格式,便于日志收集
		handler = slog.NewJSONHandler(multiWriter, handlerOpts)
	}

	// 添加全局属性
	serviceName := conf.GetString("app.name")
	if serviceName == "" {
		serviceName = "piemdm"
	}
	version := conf.GetString("app.version")
	if version == "" {
		version = "1.0.0"
	}

	// 使用 With 添加全局属性
	logger := slog.New(handler).With(
	// slog.Group("service",
	// 	slog.String("name", serviceName),
	// 	slog.String("version", version),
	// ),
	)

	return &Logger{Logger: logger}
}

// NewContext 向指定的 Gin Context 添加日志字段
func (l *Logger) NewContext(ctx *gin.Context, args ...any) {
	ctx.Set(LOGGER_KEY, l.WithContext(ctx).With(args...))
}

// WithContext 从指定的 Gin Context 返回 logger 实例
func (l *Logger) WithContext(ctx *gin.Context) *Logger {
	if ctx == nil {
		return l
	}

	// 尝试从 context 获取 logger
	if ctxLogger, exists := ctx.Get(LOGGER_KEY); exists {
		if logger, ok := ctxLogger.(*Logger); ok {
			return logger
		}
	}

	return l
}

// With 返回带有附加属性的新 logger
func (l *Logger) With(args ...any) *Logger {
	return &Logger{Logger: l.Logger.With(args...)}
}

// WithGroup 返回带有指定分组的新 logger
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{Logger: l.Logger.WithGroup(name)}
}

// Debug 记录 debug 级别日志
func (l *Logger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

// Info 记录 info 级别日志
func (l *Logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

// Warn 记录 warn 级别日志
func (l *Logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

// Error 记录 error 级别日志
func (l *Logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

// DebugContext 使用 context 记录 debug 级别日志
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.Logger.DebugContext(ctx, msg, args...)
}

// InfoContext 使用 context 记录 info 级别日志
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.Logger.InfoContext(ctx, msg, args...)
}

// WarnContext 使用 context 记录 warn 级别日志
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.Logger.WarnContext(ctx, msg, args...)
}

// ErrorContext 使用 context 记录 error 级别日志
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.Logger.ErrorContext(ctx, msg, args...)
}
