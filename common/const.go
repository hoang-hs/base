package common

import "github.com/hoang-hs/base/config"

const (
	AppEnvDev  = "dev"
	AppEnvProd = "prod"

	TraceIdName = "trace_id"
)

func IsProdEnv() bool {
	return config.Get().Mode == AppEnvProd
}
