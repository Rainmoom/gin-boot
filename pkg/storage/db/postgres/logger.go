package postgres

import (
	"context"
	"errors"
	"fmt"
	"gin-boot/pkg/logger"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type Logger struct {
	LogLevel                  gormLogger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

func NewPostgresLogger() *Logger {
	return &Logger{
		LogLevel:                  gormLogger.Error,
		IgnoreRecordNotFoundError: false,
		SlowThreshold:             time.Second * 10,
	}
}

func (l *Logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	l.LogLevel = level
	return &newLogger
}
func (l *Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		logger.Out.Info(fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}
func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		logger.Out.Warn(fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}
func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		logger.Out.Error(fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!errors.Is(err, gormLogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			logger.Out.Error(fmt.Sprintf(utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql))
		} else {
			logger.Out.Error(fmt.Sprintf(utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql))
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			logger.Out.Warn(fmt.Sprintf(utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql))
		} else {
			logger.Out.Warn(fmt.Sprintf(utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql))
		}
	case l.LogLevel == gormLogger.Info:
		sql, rows := fc()
		if rows == -1 {
			logger.Out.Warn(fmt.Sprintf(utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql))
		} else {
			logger.Out.Warn(fmt.Sprintf(utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql))
		}
	}
}
