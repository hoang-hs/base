package common

import (
	"github.com/hoang-hs/base/src/configs"
)

const (
	AppEnvDev  = "dev"
	AppEnvProd = "prod"

	TraceIdName = "trace_id"
)

var IsProdEnv bool

func SetMode(cf *configs.Config) {
	switch cf.Mode {
	case AppEnvProd:
		IsProdEnv = true
	case AppEnvDev:
		IsProdEnv = false
	default:
		panic("mode invalid")
	}
}
