package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger
var atom zap.AtomicLevel

func init() {
	atom = zap.NewAtomicLevel()
	atom.SetLevel(zapcore.InfoLevel)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	l := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer logSync()

	log = l
}

func SetDebug() {
	atom.SetLevel(zapcore.DebugLevel)
}

func Error(err error) {
	defer logSync()
	sugar := log.Sugar()
	sugar.Error(err)
}

func Errorf(template string, args ...interface{}) {
	defer logSync()
	sugar := log.Sugar()
	sugar.Errorf(template, args...)
}

func Info(msg string, fields ...zap.Field) {
	defer logSync()
	sugar := log.Sugar()
	sugar.Info(msg, fields)
}

func Infof(template string, args ...interface{}) {
	defer logSync()
	sugar := log.Sugar()
	sugar.Infof(template, args...)
}

func Debug(msg string, fields ...zap.Field) {
	defer logSync()
	sugar := log.Sugar()
	sugar.Debug(msg, fields)
}

func Debugf(template string, args ...interface{}) {
	defer logSync()
	sugar := log.Sugar()
	sugar.Debugf(template, args...)
}

func Warnf(template string, args ...interface{}) {
	defer logSync()
	sugar := log.Sugar()
	sugar.Warnf(template, args...)
}

func logSync() {
	_ = log.Sync()
}
