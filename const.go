package base

const (
	AppEnvDev  = "dev"
	AppEnvProd = "prod"

	TraceIdName = "trace_id"
)

func IsProdEnv() bool {
	return Get().Mode == AppEnvProd
}
