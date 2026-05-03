package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level      string `mapstructure:"level"`
	Encoding   string `mapstructure:"encoding"`
	Output     string `mapstructure:"output"`
	Caller     bool   `mapstructure:"caller"`
	Stacktrace bool   `mapstructure:"stacktrace"`
}

type Logger struct {
	sugar *zap.SugaredLogger
	zap   *zap.Logger
}

type Field = zap.Field

var (
	global *Logger
	once   sync.Once
	mu     sync.RWMutex
)

func Init(cfg Config) error {
	var initErr error
	once.Do(func() {
		l, err := build(cfg)
		if err != nil {
			initErr = fmt.Errorf("logger.Init: %w", err)
			return
		}
		mu.Lock()
		global = l
		mu.Unlock()
	})
	return initErr
}

func build(cfg Config) (*Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	encoderCfg := buildEncoderConfig(cfg.Encoding)

	var sink zapcore.WriteSyncer
	switch strings.ToLower(cfg.Output) {
	case "stdout", "":
		sink = zapcore.AddSync(os.Stdout)
	case "stderr":
		sink = zapcore.AddSync(os.Stderr)
	default:
		f, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("cannot open log file %q: %w", cfg.Output, err)
		}
		sink = zapcore.AddSync(f)
	}

	var encoder zapcore.Encoder
	if strings.ToLower(cfg.Encoding) == "json" {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, sink, level)

	opts := []zap.Option{zap.WithCaller(cfg.Caller)}
	if cfg.Stacktrace {
		opts = append(opts, zap.AddStacktrace(zapcore.WarnLevel))
	}

	z := zap.New(core, opts...)
	return &Logger{zap: z, sugar: z.Sugar()}, nil
}

func buildEncoderConfig(encoding string) zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "time"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	cfg.EncodeCaller = zapcore.ShortCallerEncoder

	if strings.ToLower(encoding) != "json" {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	return cfg
}

func parseLevel(s string) (zapcore.Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info", "":
		return zapcore.InfoLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unknown log level %q", s)
	}
}

func g() *Logger {
	mu.RLock()
	defer mu.RUnlock()
	if global == nil {
		panic("logger: Init() must be called before using the global logger")
	}
	return global
}

func (l *Logger) With(args ...any) *Logger {
	child := l.sugar.With(args...)
	return &Logger{sugar: child, zap: child.Desugar()}
}

func (l *Logger) WithFields(fields ...Field) *Logger {
	child := l.zap.With(fields...)
	return &Logger{zap: child, sugar: child.Sugar()}
}

func (l *Logger) Debug(msg string, args ...any) { l.sugar.Debugw(msg, args...) }
func (l *Logger) Info(msg string, args ...any)  { l.sugar.Infow(msg, args...) }
func (l *Logger) Warn(msg string, args ...any)  { l.sugar.Warnw(msg, args...) }
func (l *Logger) Error(msg string, args ...any) { l.sugar.Errorw(msg, args...) }
func (l *Logger) Fatal(msg string, args ...any) { l.sugar.Fatalw(msg, args...) }

func (l *Logger) Sync() { _ = l.zap.Sync() }

func (l *Logger) Zap() *zap.Logger { return l.zap }

func With(args ...any) *Logger { return g().With(args...) }

func WithFields(fields ...Field) *Logger { return g().WithFields(fields...) }

func Debug(msg string, args ...any) { g().Debug(msg, args...) }
func Info(msg string, args ...any)  { g().Info(msg, args...) }
func Warn(msg string, args ...any)  { g().Warn(msg, args...) }
func Error(msg string, args ...any) { g().Error(msg, args...) }
func Fatal(msg string, args ...any) { g().Fatal(msg, args...) }

func Sync() { g().Sync() }

func Zap() *zap.Logger { return g().Zap() }

func Err(err error) Field { return zap.Error(err) }

func F(key string, val any) Field { return zap.Any(key, val) }

func String(key, val string) Field { return zap.String(key, val) }

func Int(key string, val int) Field { return zap.Int(key, val) }

func Bool(key string, val bool) Field { return zap.Bool(key, val) }
