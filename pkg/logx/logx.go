package logx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// 全局日志
var (
	Logger      *zap.Logger
	atomicLevel zap.AtomicLevel
)

func init() {
	atomicLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	var w = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	var enc = zapcore.NewConsoleEncoder(newEncoderConfig())
	var core = zapcore.NewCore(
		enc,
		w,
		atomicLevel,
	)

	//默认加上caller和stackTrac
	//zap.AddCaller()
	//zap.AddStacktrace 考虑作为参数配置
	Logger = zap.New(core).Named("fly-im").WithOptions(zap.AddStacktrace(zapcore.ErrorLevel))
}

func Enable(level zapcore.Level) bool {
	return atomicLevel.Enabled(level)
}

func EnableDebug() bool {
	return atomicLevel.Enabled(zapcore.DebugLevel)
}

func SetLogLevel(level zapcore.Level) {
	atomicLevel.SetLevel(level)
}

// WithName 使用前缀创建一个logger
func WithName(name string) *zap.Logger {
	return Logger.Named(name)
}

func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	Logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

func newEncoderConfig() zapcore.EncoderConfig {
	var el = zapcore.CapitalColorLevelEncoder ///	对level进行大写,并且添加颜色

	return zapcore.EncoderConfig{
		MessageKey:     "M",
		LevelKey:       "L",
		TimeKey:        "T",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    "F",
		StacktraceKey:  "ST",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    el,
		EncodeTime:     timeEncoder, ///	time 的格式化
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func timeEncoder(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(time.Format("2006-01-02 15:04:05.000"))
}
