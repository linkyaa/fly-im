package appruntime

import "flag"

//解析flag

type (
	AppOptions interface {
		AddFlags()
		Validate() []error
	}
)

var (
	cfg = "" //配置文件地址

)

func initBasicFlag() {
	return
	flag.StringVar(&cfg, "cfg", cfg, "配置文件解析地址")
}
