package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"emperror.dev/errors"
	"github.com/caarlos0/env/v6"
	"github.com/danilluk1/social-network/libs/go/pkg/conf/environment"
	constants "github.com/danilluk1/social-network/libs/go/pkg/constants"
	typeMapper "github.com/danilluk1/social-network/libs/go/pkg/reflection/type_mapper"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

func BindConfig[T any](environments ...environment.Environment) (T, error) {
	return BindConfigKey[T]("", environments...)
}

func BindConfigKey[T any](configKey string, environments ...environment.Environment) (T, error) {
	var configPath string

	environment := environment.Environment("")
	if len(environments) > 0 {
		environment = environments[0]
	} else {
		environment = constants.Dev
	}

	configPathFromEnv := viper.Get(constants.ConfigPath)
	if configPathFromEnv != nil {
		configPath = configPathFromEnv.(string)
	} else {
		d, err := getConfigRootPath()
		if err != nil {
			return *new(T), err
		}

		configPath = d
	}

	cfg := typeMapper.GenericInstanceByT[T]()

	viper.SetConfigName(fmt.Sprintf("config.%s.json", environment))
	viper.AddConfigPath(configPath)
	viper.SetConfigType(constants.Json)

	if err := viper.ReadInConfig(); err != nil {
		return *new(T), errors.WrapIf(err, "viper.ReadInConfig")
	}

	if len(configKey) == 0 {
		if err := viper.Unmarshal(cfg); err != nil {
			return *new(T), errors.WrapIf(err, "viper.Unmarshal")
		}
	} else {
		if err := viper.UnmarshalKey(configKey, cfg); err != nil {
			return *new(T), errors.WrapIf(err, "viper.Unmarshal")
		}
	}

	viper.AutomaticEnv()

	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	defaults.SetDefaults(cfg)

	return cfg, nil
}

func getConfigRootPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	absCurrentDir, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(absCurrentDir, "conf")

	return configPath, nil
}
