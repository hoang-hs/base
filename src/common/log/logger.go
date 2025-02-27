package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

const callerSkip = 2

type logger struct {
	zap *zap.SugaredLogger
	*zap.Logger
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func NewLogger() {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeCaller: zapcore.FullCallerEncoder,
		EncodeTime:   SyslogTimeEncoder,
		EncodeLevel:  CustomLevelEncoder,
	}

	var encoder zapcore.Encoder
	var level zapcore.Level

	encoder = zapcore.NewConsoleEncoder(encoderConfig)
	level = zap.DebugLevel

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), level)
	//set log
	globalLogger = &logger{
		zap: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(callerSkip)).Sugar(),
	}
	return
}

func (l *logger) Infof(msg string, args ...interface{}) {
	l.zap.Infof(msg, args...)
}

func (l *logger) Debugf(msg string, args ...interface{}) {
	l.Error("debug", zap.Error(nil))
	l.zap.Debugf(msg, args...)
}

func (l *logger) Warnf(msg string, args ...interface{}) {
	l.zap.Warnf(msg, args...)
}

func (l *logger) Errorf(msg string, args ...interface{}) {
	l.zap.Errorf(msg, args...)
}

func (l *logger) Fatalf(msg string, args ...interface{}) {
	l.zap.Fatalf(msg, args...)
}

func (l *logger) GetZap() *zap.SugaredLogger {
	return l.zap
}

var (
	Any        = zap.Any
	Bool       = zap.Bool
	Duration   = zap.Duration
	Float64    = zap.Float64
	Int        = zap.Int
	Int32      = zap.Int32
	Int64      = zap.Int64
	Skip       = zap.Skip
	String     = zap.String
	Stringer   = zap.Stringer
	Time       = zap.Time
	Uint       = zap.Uint
	Uint32     = zap.Uint32
	Uint64     = zap.Uint64
	Uintptr    = zap.Uintptr
	ByteString = zap.ByteString
)

func Err(err error) zapcore.Field {
	if err == nil {
		return Skip()
	}
	return String("error", err.Error())
}
