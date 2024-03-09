package logx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

//TODO:增加使用的便捷性,自定义logger

var (
	Logger  *zap.Logger
	options *Options
)

type (
	FlyLogger struct {
		opt    *Options
		logger *zap.Logger
	}
)

func init() {
	options = &Options{Level: zapcore.DebugLevel}
	var w = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	var enc = zapcore.NewConsoleEncoder(newEncoderConfig())
	var core = zapcore.NewCore(
		enc,
		w,
		options.Level,
	)

	Logger = zap.New(core).Named("fly-im")
}

func Enable(level zapcore.Level) bool {
	return options.Level >= level
}

func EnableDebug() bool {
	return zapcore.DebugLevel >= options.Level
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
