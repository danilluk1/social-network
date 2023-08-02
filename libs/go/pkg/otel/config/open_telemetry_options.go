package config

import (
	"github.com/iancoleman/strcase"

	"github.com/danilluk1/social-network/libs/go/pkg/conf"
	"github.com/danilluk1/social-network/libs/go/pkg/conf/environment"
	typeMapper "github.com/danilluk1/social-network/libs/go/pkg/reflection/type_mapper"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetTypeNameByT[OpenTelemetryOptions]())

type OpenTelemetryOptions struct {
	Enabled               bool                   `mapstructure:"enabled"`
	ServiceName           string                 `mapstructure:"serviceName"`
	InstrumentationName   string                 `mapstructure:"instrumentationName"`
	Id                    int64                  `mapstructure:"id"`
	AlwaysOnSampler       bool                   `mapstructure:"alwaysOnSampler"`
	JaegerExporterOptions *JaegerExporterOptions `mapstructure:"jaegerExporterOptions"`
	OTelMetricsOptions    *OTelMetricsOptions    `mapstructure:"otelMetricsOptions"`
	UseStdout             bool                   `mapstructure:"useStdout"`
}

type JaegerExporterOptions struct {
	AgentHost string `mapstructure:"agentHost"`
	AgentPort string `mapstructure:"agentPort"`
}

type OTelMetricsOptions struct {
	Host             string `mapstructure:"host"`
	Port             string `mapstructure:"port"`
	Name             string `mapstructure:"name"`
	MetricsRoutePath string `mapstructure:"metricsRoutePath"`
}

func ProvideOtelConfig(environment environment.Environment) (*OpenTelemetryOptions, error) {
	return conf.BindConfigKey[*OpenTelemetryOptions](optionName, environment)
}
