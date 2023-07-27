package conf

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	AppEnv string `json:"app_env" split_words:"true" required:"true"`
	SMTP   SMTPConfiguration
	Kafka  KafkaConfiguration
}

type SMTPConfiguration struct {
	Host         string        `json:"host" split_words:"true" required:"true"`
	Port         int           `json:"port,omitempty" split_words:"true" required:"true"`
	Pass         string        `json:"pass,omitempty" split_words:"true" required:"true"`
	User         string        `json:"user" split_words:"true" required:"true"`
	AdminEmail   string        `json:"admin_email" split_words:"true"`
	SenderName   string        `json:"sender_name" split_words:"true"`
	MaxFrequency time.Duration `json:"max_frequency" split_words:"true"`
}

func (sc *SMTPConfiguration) Validate() error {
	return nil
}

type KafkaConfiguration struct {
	SchemaRegistryUrl string `json:"schema_registry_url" split_words:"true" required:"true"`
	KafkaUrl          string `json:"kafka_url" split_words:"true" required:"true"`
	GroupID           string `json:"group_id" split_words:"true" required:"true"`
}

func (kc *KafkaConfiguration) Validate() error {
	return nil
}

func (c *Configuration) ApplyDefaults() error {
	if c.AppEnv == "" {
		c.AppEnv = "development"
	}

	if c.Kafka.KafkaUrl == "" {
		c.Kafka.KafkaUrl = "localhost:9092"
	}

	if c.Kafka.SchemaRegistryUrl == "" {
		c.Kafka.SchemaRegistryUrl = "localhost:18085"
	}

	if c.Kafka.GroupID == "" {
		c.Kafka.GroupID = "mailer"
	}

	if c.SMTP.Host == "" {
		c.SMTP.Host = "localhost"
	}

	if c.SMTP.Port == 0 {
		c.SMTP.Port = 1025
	}

	if c.SMTP.AdminEmail == "" {
		c.SMTP.AdminEmail = "admin@socialy.ru"
	}

	if c.SMTP.MaxFrequency == 0 {
		c.SMTP.MaxFrequency = 10
	}

	if c.SMTP.Pass == "" {
		c.SMTP.Pass = "admin"
	}

	if c.SMTP.User == "" {
		c.SMTP.User = "admin"
	}

	return nil
}

func (c *Configuration) Validate() error {
	validatables := []interface {
		Validate() error
	}{
		&c.SMTP,
		&c.Kafka,
	}

	for _, validatable := range validatables {
		if err := validatable.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func loadEnv(filename string) error {
	var err error
	if filename != "" {
		err = godotenv.Overload(filename)
	} else {
		godotenv.Load()
		if os.IsNotExist(err) {
			return nil
		}
	}

	return err
}

func Load(filename string) (*Configuration, error) {
	if err := loadEnv(filename); err != nil {
		return nil, err
	}

	config := new(Configuration)
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}

	if err := config.ApplyDefaults(); err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}
