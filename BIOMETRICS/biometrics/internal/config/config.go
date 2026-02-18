package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`

	Server struct {
		Port         string `mapstructure:"SERVER_PORT"`
		ReadTimeout  int    `mapstructure:"SERVER_READ_TIMEOUT"`
		WriteTimeout int    `mapstructure:"SERVER_WRITE_TIMEOUT"`
		IdleTimeout  int    `mapstructure:"SERVER_IDLE_TIMEOUT"`
	} `mapstructure:"server"`

	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Worker    WorkerConfig    `mapstructure:"worker"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"DB_HOST"`
	Port            string `mapstructure:"DB_PORT"`
	User            string `mapstructure:"DB_USER"`
	Password        string `mapstructure:"DB_PASSWORD"`
	Database        string `mapstructure:"DB_NAME"`
	MaxOpenConns    int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime int    `mapstructure:"DB_CONN_MAX_LIFETIME"`
	SSLMode         string `mapstructure:"DB_SSL_MODE"`
}

func (d DatabaseConfig) URL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Database, d.SSLMode)
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	Database int    `mapstructure:"REDIS_DATABASE"`
	PoolSize int    `mapstructure:"REDIS_POOL_SIZE"`
}

func (r RedisConfig) URL() string {
	if r.Password != "" {
		return fmt.Sprintf("redis://:%s@%s:%s/%d", r.Password, r.Host, r.Port, r.Database)
	}
	return fmt.Sprintf("redis://%s:%s/%d", r.Host, r.Port, r.Database)
}

type AuthConfig struct {
	JWTSecret          string        `mapstructure:"JWT_SECRET"`
	JWTExpire          time.Duration `mapstructure:"JWT_EXPIRE"`
	RefreshTokenExpire time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRE"`
	OAuth              OAuthConfig   `mapstructure:"oauth"`
}

type OAuthConfig struct {
	Google OAuthProviderConfig `mapstructure:"google"`
	GitHub OAuthProviderConfig `mapstructure:"github"`
	Apple  OAuthProviderConfig `mapstructure:"apple"`
}

type OAuthProviderConfig struct {
	ClientID     string   `mapstructure:"CLIENT_ID"`
	ClientSecret string   `mapstructure:"CLIENT_SECRET"`
	RedirectURL  string   `mapstructure:"REDIRECT_URL"`
	Scopes       []string `mapstructure:"SCOPES"`
}

type WorkerConfig struct {
	Count   int           `mapstructure:"WORKER_COUNT"`
	Captcha CaptchaConfig `mapstructure:"captcha"`
	Survey  SurveyConfig  `mapstructure:"survey"`
}

type CaptchaConfig struct {
	Enabled    bool `mapstructure:"ENABLED"`
	MaxRetries int  `mapstructure:"MAX_RETRIES"`
	RetryDelay int  `mapstructure:"RETRY_DELAY_SECONDS"`
	Timeout    int  `mapstructure:"TIMEOUT_SECONDS"`
}

type SurveyConfig struct {
	Enabled         bool `mapstructure:"ENABLED"`
	MaxRetries      int  `mapstructure:"MAX_RETRIES"`
	RetryDelay      int  `mapstructure:"RETRY_DELAY_SECONDS"`
	ParallelWorkers int  `mapstructure:"PARALLEL_WORKERS"`
}

type RateLimitConfig struct {
	RequestsPerMinute int `mapstructure:"REQUESTS_PER_MINUTE"`
	BurstSize         int `mapstructure:"BURST_SIZE"`
}

type LoadOptions struct {
	SkipValidation bool
	RequireDB      bool
	RequireRedis   bool
}

func Load(opts ...LoadOptions) (*Config, error) {
	var options LoadOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.biometrics/")
	viper.AddConfigPath("/etc/biometrics/")
	viper.AutomaticEnv()

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_READ_TIMEOUT", 30)
	viper.SetDefault("SERVER_WRITE_TIMEOUT", 30)
	viper.SetDefault("SERVER_IDLE_TIMEOUT", 120)
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("ENVIRONMENT", "development")

	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 5)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", 300)
	viper.SetDefault("DB_SSL_MODE", "disable")

	viper.SetDefault("REDIS_POOL_SIZE", 10)
	viper.SetDefault("REDIS_DATABASE", 0)

	viper.SetDefault("JWT_EXPIRE", 24*time.Hour)
	viper.SetDefault("REFRESH_TOKEN_EXPIRE", 7*24*time.Hour)

	viper.SetDefault("WORKER_COUNT", 5)
	viper.SetDefault("WORKER_CAPTCHA_MAX_RETRIES", 3)
	viper.SetDefault("WORKER_CAPTCHA_RETRY_DELAY_SECONDS", 5)
	viper.SetDefault("WORKER_CAPTCHA_TIMEOUT_SECONDS", 30)
	viper.SetDefault("WORKER_SURVEY_MAX_RETRIES", 3)
	viper.SetDefault("WORKER_SURVEY_RETRY_DELAY_SECONDS", 10)
	viper.SetDefault("WORKER_SURVEY_PARALLEL_WORKERS", 3)

	viper.SetDefault("RATE_LIMIT_REQUESTS_PER_MINUTE", 100)
	viper.SetDefault("RATE_LIMIT_BURST_SIZE", 20)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if !options.SkipValidation {
		if err := validateConfig(&cfg, options); err != nil {
			return nil, fmt.Errorf("invalid configuration: %w", err)
		}
	}

	return &cfg, nil
}

func validateConfig(cfg *Config, opts LoadOptions) error {
	if cfg.Environment == "" {
		return fmt.Errorf("ENVIRONMENT is required")
	}

	if opts.RequireDB || opts.SkipValidation == false {
		if cfg.Database.Host == "" {
			return fmt.Errorf("database host is required")
		}
	}

	if opts.RequireRedis || opts.SkipValidation == false {
		if cfg.Redis.Host == "" {
			return fmt.Errorf("redis host is required")
		}
	}

	if cfg.Auth.JWTSecret == "" && cfg.Environment != "development" {
		return fmt.Errorf("JWT_SECRET is required in non-development environments")
	}

	return nil
}
