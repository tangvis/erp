package logutil

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"

	ctxUtil "github.com/tangvis/erp/pkg/context"
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

func (l *Logger) Debugf(template string, args ...any) {
	l.sugar.Debugf(template, args...)
}

func (l *Logger) Infof(template string, args ...any) {
	l.sugar.Infof(template, args...)
}

func (l *Logger) Warnf(template string, args ...any) {
	l.sugar.Warnf(template, args...)
}

func (l *Logger) Errorf(template string, args ...any) {
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

func DebugF(template string, args ...any) { logger.Debugf(template, args...) }
func Info(msg string, fields ...Field)    { logger.Info(msg, fields...) }
func CtxInfo(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, zap.String(ctxUtil.TraceIDKey, ctxUtil.GetTraceID(ctx)))
	logger.Info(msg, fields...)
}
func InfoF(template string, args ...any) { logger.Infof(template, args...) }
func CtxInfoF(ctx context.Context, template string, args ...any) {
	logger.Infof(fmt.Sprintf("[%s]%s", ctxUtil.GetTraceID(ctx), template), args...)
}
func Warn(msg string, fields ...Field)   { logger.Warn(msg, fields...) }
func WarnF(template string, args ...any) { logger.Warnf(template, args...) }
func Error(msg string, fields ...Field)  { logger.Error(msg, fields...) }
func CtxError(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, zap.String(ctxUtil.TraceIDKey, ctxUtil.GetTraceID(ctx)))
	logger.Error(msg, fields...)
}
func ErrorF(template string, args ...any) { logger.Errorf(template, args...) }
func CtxErrorF(ctx context.Context, template string, args ...any) {
	logger.Errorf(fmt.Sprintf("[%s]%s", ctxUtil.GetTraceID(ctx), template), args...)
}
func Panic(msg string, fields ...Field) { logger.Panic(msg, fields...) }
func Fatal(msg string, fields ...Field) { logger.Fatal(msg, fields...) }

func Sync() error { return logger.Sync() }
