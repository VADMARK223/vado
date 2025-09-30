package logger

import (
	"fmt"
	"vado/internal/util"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init инициализирует глобальный логгер
func Init() (*zap.Logger, error) {
	idDevMode := util.IsDevMode()
	var cfg zap.Config
	if idDevMode {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	log = l
	return log, err
}

// L возвращает текущий логгер (или заглушку)
func L() *zap.Logger {
	if log == nil {
		fmt.Println("Logger not initialized.")
		return zap.NewNop()
	}
	return log
}

// Sync закрывает logger
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
