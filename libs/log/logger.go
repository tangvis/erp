package logutil

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"

	ctxUtil "github.com/tangvis/erp/libs/context"
)

var logger *Logger

func InitLogger(config *Config) {
	logger = New(config)
}

/*New New Logger*/
func New(config *Config) *Logger {
	lg := &Logger{
		Config: config,
	}
	lg.ApplyConfig()
	return lg
}

/*ApplyConfig 应用当前Config配置*/
func (l *Logger) ApplyConfig() {
	conf := l.Config
	var cores []zapcore.Core

	var encoder zapcore.Encoder

	if conf.jsonFormat {
		encoder = zapcore.NewJSONEncoder(getEncoder())
	} else {
		encoder = zapcore.NewConsoleEncoder(getEncoder())
	}

	conf.atomicLevel.SetLevel(getLevel(conf.defaultLogLevel))

	if conf.consoleOut {
		writer := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(encoder, writer, conf.atomicLevel)
		cores = append(cores, core)
	}

	if conf.fileOut.enable {
		fileWriter := getFileWriter(
			conf.fileOut.path,
			conf.fileOut.name,
			conf.fileOut.rotationTime,
			conf.fileOut.rotationCount,
		)
		writer := zapcore.AddSync(fileWriter)
		core := zapcore.NewCore(encoder, writer, conf.atomicLevel)
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)

	lg := zap.New(combinedCore,
		zap.AddCallerSkip(conf.callerSkip),
		zap.AddStacktrace(getLevel(conf.stacktraceLevel)),
		zap.AddCaller(),
	)

	if conf.projectName != "" {
		lg = lg.Named(conf.projectName)
	}

	// nolint:errcheck
	defer lg.Sync()

	l.l = lg
	l.sugar = lg.Sugar()
}

type Logger struct {
	l      *zap.Logger
	sugar  *zap.SugaredLogger
	Config *Config
}

type Field = zap.Field

func (l *Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}

func (l *Logger) DebugF(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

func (l *Logger) InfoF(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

func (l *Logger) WarnF(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

func (l *Logger) ErrorF(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Sync() error {
	if err := l.sugar.Sync(); err != nil {
		return err
	}
	return l.l.Sync()
}

func Debug(msg string, fields ...Field) { logger.Debug(msg, fields...) }

func DebugF(template string, args ...interface{}) { logger.DebugF(template, args...) }
func Info(msg string, fields ...Field)            { logger.Info(msg, fields...) }
func CtxInfo(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, zap.String(ctxUtil.TraceIDKey, ctxUtil.GetTranceID(ctx)))
	logger.Info(msg, fields...)
}
func InfoF(template string, args ...interface{}) { logger.InfoF(template, args...) }
func CtxInfoF(ctx context.Context, template string, args ...interface{}) {
	logger.InfoF(fmt.Sprintf("[%s]%s", ctxUtil.GetTranceID(ctx), template), args...)
}
func Warn(msg string, fields ...Field)           { logger.Warn(msg, fields...) }
func WarnF(template string, args ...interface{}) { logger.WarnF(template, args...) }
func Error(msg string, fields ...Field)          { logger.Error(msg, fields...) }
func CtxError(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, zap.String(ctxUtil.TraceIDKey, ctxUtil.GetTranceID(ctx)))
	logger.Error(msg, fields...)
}
func ErrorF(template string, args ...interface{}) { logger.ErrorF(template, args...) }
func CtxErrorF(ctx context.Context, template string, args ...interface{}) {
	logger.ErrorF(fmt.Sprintf("[%s]%s", ctxUtil.GetTranceID(ctx), template), args...)
}
func Panic(msg string, fields ...Field) { logger.Panic(msg, fields...) }
func Fatal(msg string, fields ...Field) { logger.Fatal(msg, fields...) }

func Sync() error { return logger.Sync() }
