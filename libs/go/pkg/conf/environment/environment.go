package environment

import (
	"log"
	"os"

	constants "github.com/danilluk1/social-network/libs/go/pkg/constants"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Environment string

var (
	Development = Environment(constants.Dev)
	Test        = Environment(constants.Test)
	Production  = Environment(constants.Production)
)

func ConfigAppEnv(environments ...Environment) Environment {
	environment := Environment("")
	if len(environments) > 0 {
		environment = environments[0]
	} else {
		environment = Development
	}

	viper.AutomaticEnv()

	err := godotenv.Load()
	if err != nil {
		log.Println(".env file cannot be found.")
	}

	manualEnv := os.Getenv(constants.AppEnv)

	if manualEnv != "" {
		environment = Environment(manualEnv)
	}

	return environment
}

func (env Environment) IsDev() bool {
	return env == Development
}

func (env Environment) IsProd() bool {
	return env == Production
}

func (env Environment) GetEnvironmentName() string {
	return string(env)
}
