package conf

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const defaultMinPasswordLength int = 5
const defaultChallengeExpiryDuration float64 = 300
const defaultFlowStateExpiryDuration time.Duration = 300 * time.Second

type OAuthProviderConfiguration struct {
	ClientID    []string `json:"client_id" split_words:"true"`
	Secret      string   `json:"secret"`
	RedirectURI string   `json:"redirect_uri" split_words:"true"`
	URL         string   `json:"url"`
	ApiURL      string   `json:"api_url" split_words:"true"`
	Enabled     bool     `json:"enabled"`
}

type EmailProviderConfiguration struct {
	Enabled bool `json:"enabled" default:"true"`
}

type DBConfiguration struct {
	URL               string        `json:"url" envconfig:"DATABASE_URL" required:"true"`
	Namespace         string        `json:"namespace" envconfig:"DB_NAMESPACE" default:"public"`
	MaxPoolSize       int           `json:"max_pool_size" split_words:"true"`
	MaxIdlePoolSize   int           `json:"max_idle_pool_size" split_words:"true"`
	ConnMaxLifetime   time.Duration `json:"conn_max_lifetime,omitempty" split_words:"true"`
	ConnMaxIdleTime   time.Duration `json:"conn_max_idle_time,omitempty" split_words:"true"`
	HealthCheckPeriod time.Duration `json:"health_check_period" split_words:"true"`
	MigrationsPath    string        `json:"migrations_path" split_words:"true" default:"/home/danluki/Projects/social-network/apps/auth/internal/db/migration"`
	CleanupEnabled    bool          `json:"cleanup_enabled" split_words:"true" default:"false"`
}

func (c *DBConfiguration) Validate() error {
	return nil
}

type KafkaConfiguration struct {
	KafkaUrl          string `default:"localhost:29092" envconfig:"KAFKA_URL"`
	SchemaRegistryUrl string `default:"http://localhost:18085" envconfig:"SCHEMA_REGISTRY_URL"`
	SchemasPath       string `default:"/home/danluki/Projects/social-network/libs/kafka/schemas/" envconfig:"SCHEMAS_PATH"`
}

type PASETOConfiguration struct {
	Secret string   `json:"secret" envconfig:"PASETO_SECRET" required:"true"`
	Exp    int      `json:"exp"`
	Aud    string   `json:"aud"`
	Roles  []string `json:"roles" split_words:"true"`
}

type MFAConfiguration struct {
	Enabled                     bool    `default:"false"`
	ChallengeExpiryDuration     float64 `json:"challenge_expiry_duration" default:"300" split_words:"true"`
	RateLimitChallengeAndVerify float64 `split_words:"true" default:"15"`
	MaxEnrolledFactors          float64 `split_words:"true" default:"10"`
	MaxVerifiedFactors          int     `split_words:"true" default:"10"`
}

type GRPCConfiguration struct {
	Host string
	Port string `envconfig:"PORT" default:"50051"`
}

func (c *GRPCConfiguration) Validate() error {
	return nil
}

type GlobalConfiguration struct {
	GRPC                  GRPCConfiguration
	DB                    DBConfiguration
	Logging               LoggingConfig `envconfig:"LOG"`
	External              ProviderConfiguration
	Profiler              ProfilerConfig `envconfig:"PROFILER"`
	Tracing               TracingConfig
	Metrics               MetricsConfig
	SMTP                  SMTPConfiguration
	Kafka                 KafkaConfiguration
	RateLimitHeader       string  `split_words:"true"`
	RateLimitEmailSent    float64 `split_words:"true" default:"30"`
	RateLimitSmsSent      float64 `split_words:"true" default:"30"`
	RateLimitVerify       float64 `split_words:"true" default:"30"`
	RateLimitTokenRefresh float64 `split_words:"true" default:"30"`

	AppEnv            string              `json:"app_env" default:"development" envconfig:"APP_ENV"  required:"true"`
	SiteURL           string              `json:"site_url" default:"localhost:5173" envconfig:"SITE_URL" required:"true"`
	PasswordMinLength int                 `json:"password_min_length" split_words:"true"`
	PASETO            PASETOConfiguration `json:"paseto"`
	Mailer            MailerConfiguration `json:"mailer"`
	DisableSignup     bool                `json:"disable_signup" split_words:"true"`
	MFA               MFAConfiguration    `json:"MFA"`
	Cookie            struct {
		Key      string `json:"key"`
		Domain   string `json:"domain"`
		Duration int    `json:"duration"`
	} `json:"cookies"`
}

type EmailContentConfiguration struct {
	Invite           string `json:"invite"`
	Confirmation     string `json:"confirmation"`
	Recovery         string `json:"recovery"`
	EmailChange      string `json:"email_change" split_words:"true"`
	MagicLink        string `json:"magic_link" split_words:"true"`
	Reauthentication string `json:"reauthentication"`
}

type SMTPConfiguration struct {
	MaxFrequency time.Duration `json:"max_frequency" split_words:"true"`
	Host         string        `json:"host"`
	Port         int           `json:"port,omitempty" default:"587"`
	User         string        `json:"user"`
	Pass         string        `json:"pass,omitempty"`
	AdminEmail   string        `json:"admin_email" split_words:"true"`
	SenderName   string        `json:"sender_name" split_words:"true"`
}

type MailerConfiguration struct {
	Autoconfirm              bool                      `json:"autoconfirm"`
	Subjects                 EmailContentConfiguration `json:"subjects"`
	Templates                EmailContentConfiguration `json:"templates"`
	URLPaths                 EmailContentConfiguration `json:"url_paths"`
	SecureEmailChangeEnabled bool                      `json:"secure_email_change_enabled" split_words:"true" default:"true"`
	OtpExp                   uint                      `json:"otp_exp" split_words:"true"`
	OtpLength                int                       `json:"otp_length" split_words:"true"`
}

func (c *SMTPConfiguration) Validate() error {
	return nil
}

type ProviderConfiguration struct {
	Email EmailProviderConfiguration `json:"email"`
}

func loadEnvironment(filename string) error {
	var err error
	if filename != "" {
		err = godotenv.Overload(filename)
	} else {
		err = godotenv.Load()
		if os.IsNotExist(err) {
			return nil
		}
	}

	return err
}

func LoadGlobal(filename string) (*GlobalConfiguration, error) {
	if err := loadEnvironment(filename); err != nil {
		return nil, err
	}

	config := new(GlobalConfiguration)
	if err := envconfig.Process("auth", config); err != nil {
		return nil, err
	}

	if err := config.ApplyDefaults(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *GlobalConfiguration) ApplyDefaults() error {
	return nil
}

func (c *GlobalConfiguration) Validate() error {
	validatables := []interface {
		Validate() error
	}{
		&c.GRPC,
		&c.DB,
		&c.Tracing,
		&c.Metrics,
		&c.SMTP,
	}

	for _, validatable := range validatables {
		if err := validatable.Validate(); err != nil {
			return err
		}
	}

	return nil
}
