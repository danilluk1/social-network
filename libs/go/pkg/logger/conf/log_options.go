package conf

import (
	"github.com/danilluk1/social-network/libs/go/pkg/conf"
	"github.com/danilluk1/social-network/libs/go/pkg/conf/environment"
	"github.com/danilluk1/social-network/libs/go/pkg/logger/models"
	typeMapper "github.com/danilluk1/social-network/libs/go/pkg/reflection/type_mapper"
	"github.com/iancoleman/strcase"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetTypeNameByT[LogOptions]())

type LogOptions struct {
	LogLevel      string         `mapstructure:"level"`
	LogType       models.LogType `mapstructure:"logType"`
	CallerEnabled bool           `mapstructure:"callerEnabled"`
}

func ProvideLogConfig(env environment.Environment) (*LogOptions, error) {
	return conf.BindConfigKey[*LogOptions](optionName, env)
}
