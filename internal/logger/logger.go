package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level  string
	Format string
}

var (
	logBuffer   []string
	logBufferMu sync.RWMutex
	maxLogLines = 500
)

func AddLog(msg string) {
	logBufferMu.Lock()
	defer logBufferMu.Unlock()
	timestamp := time.Now().Format("15:04:05.000")
	logBuffer = append(logBuffer, timestamp+" "+msg)
	if len(logBuffer) > maxLogLines {
		logBuffer = logBuffer[len(logBuffer)-maxLogLines:]
	}
}

func New(cfg Config) (*zap.Logger, error) {
	var level zapcore.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		return nil, fmt.Errorf("неверный уровень логирования: %s", cfg.Level)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	switch strings.ToLower(cfg.Format) {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	logger := zap.New(core, zap.AddCaller())

	AddLog("=== Сервер запущен ===")

	logger.Info("logger инициализирован",
		zap.String("level", cfg.Level),
		zap.String("format", cfg.Format),
	)

	return logger, nil
}
