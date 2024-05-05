package option

import (
	"flag"
	"github.com/linkyaa/fly-im/pkg/appruntime"
	"go.uber.org/zap/zapcore"
)

type (
	LogOption struct {
		Level    string        `json:"level"`
		ZapLevel zapcore.Level `json:"-"`
	}
)

var (
	defaultLevel                       = zapcore.InfoLevel
	_            appruntime.AppOptions = (*LogOption)(nil)
)

func (l *LogOption) AddFlags() {
	flag.StringVar(&l.Level, "log.level", defaultLevel.String(), "日志级别:debug,info,warn,err,panic")
}

func (l *LogOption) Validate() []error {
	var errors []error
	level, err := zapcore.ParseLevel(l.Level)
	if err != nil {
		return append(errors, err)
	}

	l.ZapLevel = level
	return errors
}

func NewLogOption() *LogOption {
	res := &LogOption{
		Level:    defaultLevel.String(),
		ZapLevel: defaultLevel,
	}
	return res
}
