package conf

import (
	"fmt"
	"time"
)

type GlobalConfiguration struct {
	App  *AppConfiguration
	DB   *DBConfiguration
	Mail *MailConfiguration
}

type AppConfiguration struct {
	Host              string `envconfig:"APP_HOST"`
	Port              string `envconfig:"APP_PORT" default:"8080"`
	PasswordMinLength int    `envconfig:"APP_PW_MIN_LEN" default:"6"`
	PasswordMaxLength int    `envconfig:"APP_PW_MAX_LEN" default:"30"`

	ConfirmationTokenExpiration time.Duration `envconfig:"APP_CONFIRMATION_TOKEN_EXPIRATION" default:"1h"`
}

// DBConfiguration holds all the database related configuration.
type DBConfiguration struct {
	Host     string `envconfig:"DATABASE_HOST"`
	Port     string `envconfig:"DATABASE_PORT" default:"5432"`
	UserName string `envconfig:"DATABASE_USERNAME"`
	Password string `envconfig:"DATABASE_PASSWORD"`
	DBName   string `envconfig:"DATABASE_NAME"`
	SSLmode  string `envconfig:"DATABASE_SSLMODE" default:"disable"`
	Timezone string `envconfig:"DATABASE_TIMEZONE" default:"UTC"`
	// MaxPoolSize defaults to 0 (unlimited).
	// MaxPoolSize       int
	// MaxIdlePoolSize   int
	// ConnMaxLifetime   time.Duration
	// ConnMaxIdleTime   time.Duration
	// HealthCheckPeriod time.Duration
	// MigrationsPath    string
	// CleanupEnabled    bool
}

func (c *DBConfiguration) Dsn() string {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.UserName, c.Password, c.DBName, c.SSLmode,
	)

	return dsn
}

type MailConfiguration struct {
	ResendApiKey              string        `envconfig:"RESEND_API_KEY"`
	SendConfirmationFrequency time.Duration `envconfig:"CONFIRMATION_FREQUENCY_MIN" default:"3m"` // 単位なしだとnsになる
	ResendFromEmail           string        `envconfig:"RESEND_EMAIL_FROM"`
}
