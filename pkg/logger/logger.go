package logger

import (
	"context"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
	"strings"

	"github.com/Rainmoom/gin-boot/pkg/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Out         *zap.Logger
	atomicLevel = zap.NewAtomicLevel()
)

func Init(ctx context.Context, cfg *conf.LogConfig) (err error) {
	err = initLogger(cfg)
	if err != nil {
		return
	}

	Out = zap.L()

	go func() {
		<-ctx.Done()
		err = Out.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}()

	return nil
}

func SetLevelHTTP(w http.ResponseWriter, r *http.Request) {
	atomicLevel.ServeHTTP(w, r)
}

func SetLevel(level zapcore.Level) {
	atomicLevel.SetLevel(level)
}

func getEncoder(format string) zapcore.Encoder {

	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConfig.TimeKey = "time"
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder

	if strings.ToUpper(format) == "JSON" {
		return zapcore.NewJSONEncoder(encodeConfig)
	} else {
		encodeConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encodeConfig)
	}
}

func getLogWriter(cfg *conf.LogConfig) zapcore.Core {
	var cores []zapcore.Core

	if cfg.Path != "" {
		logRotate := &lumberjack.Logger{
			Filename:   cfg.Path,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		fileEncoder := getEncoder(cfg.Format)
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(logRotate), atomicLevel))
	}

	if cfg.ConsoleEnable {
		consoleEncoder := getEncoder("console")
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), atomicLevel))
	}

	return zapcore.NewTee(cores...)
}

func initLogger(cfg *conf.LogConfig) (err error) {

	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return err
	}
	atomicLevel.SetLevel(level.Level())

	core := getLogWriter(cfg)

	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)

	return
}
