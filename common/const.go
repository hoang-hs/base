package common

import "github.com/hoang-hs/base/config"

const (
	AppEnvDev  = "dev"
	AppEnvProd = "prod"

	TraceIdName = "trace_id"
)

var IsProdEnv bool

func SetMode(cf *config.Config) {
	switch cf.Mode {
	case AppEnvProd:
		IsProdEnv = true
	case AppEnvDev:
		IsProdEnv = false
	default:
		panic("mode invalid")
	}
}
