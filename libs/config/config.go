package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MigrationUrl      string `required:"true" default:"file://apps/auth/internal/db/migration"`
	AuthPostgresUrl   string `required:"true" default:"postgres://test:test@localhost:5432/sn_auth" envconfig:"AUTH_POSTGRES_URL"`
	AppEnv            string `required:"true" default:"development" envconfig:"APP_ENV"`
	JwtSecret         string `required:"true" default:"CoolSecretForJWT" envconfig:"JWT_SECRET"`
	RedisUrl          string `required:"true" default:"redis://localhost:6379" envconfig:"REDIS_URL"`
	MailHost          string `required:"true" default:"localhost" envconfig:"MAIL_HOST"`
	MailPort          int    `required:"true" default:"1025" envconfig:"MAIL_PORT"`
	MailUser          string `required:"false" default:"" envconfig:"MAIL_USER"`
	MailPass          string `required:"false" default:"" envconfig:"MAIL_PASS"`
	KafkaUrl          string `required:"true" default:"localhost:29092" envconfig:"KAFKA_URL"`
	SchemaRegistryUrl string `required:"true" default:"localhost:18085" envconfig:"SCHEMA_REGISTRY_URL"`
}

func New() (*Config, error) {
	var newCfg Config
	var err error

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	wd = filepath.Join(wd, "..", "..")
	envPath := filepath.Join(wd, ".env")
	_ = godotenv.Load(envPath)

	if err = envconfig.Process("", &newCfg); err != nil {
		return nil, err
	}

	return &newCfg, nil
}
