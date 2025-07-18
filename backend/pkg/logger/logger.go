package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is a global logger
var Log *zap.Logger

// Init initializes the logger
func Init(isDev bool) {
    var cfg zap.Config
    if isDev {
        cfg = zap.NewDevelopmentConfig()
    } else {
        cfg = zap.NewProductionConfig()
        cfg.EncoderConfig.TimeKey = "timestamp"
        cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // "2006-01-02T15:04:05Z07:00"
    }

    var err error
    Log, err = cfg.Build()
    if err != nil {
        panic("cannot initialize zap logger: " + err.Error())
    }
    zap.ReplaceGlobals(Log) // Optional: make zap.L() use this logger
}
