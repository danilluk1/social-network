package conf

import (
	"strings"

	"github.com/danilluk1/social-network/libs/go/pkg/conf"
	"github.com/danilluk1/social-network/libs/go/pkg/conf/environment"
	typeMapper "github.com/danilluk1/social-network/libs/go/pkg/reflection/type_mapper"
	"github.com/iancoleman/strcase"
)

type AppOptions struct {
	ServiceName string `mapstructure:"serviceName"  env:"ServiceName"`
}

func NewAppOptions(environment environment.Environment) (*AppOptions, error) {
	optionName := strcase.ToLowerCamel(typeMapper.GetTypeNameByT[AppOptions]())
	cfg, err := conf.BindConfigKey[*AppOptions](optionName, environment)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *AppOptions) GetMicroserviceNameUpper() string {
	return strings.ToUpper(cfg.ServiceName)
}

func (cfg *AppOptions) GetMicroserviceName() string {
	return cfg.ServiceName
}
