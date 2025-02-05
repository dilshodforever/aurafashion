package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App   `yaml:"app"`
		HTTP  `yaml:"http"`
		Log   `yaml:"logger"`
		PG    `yaml:"postgres"`
		JWT   `yaml:"jwt"`
		Redis `yaml:"redis"`
		Gmail `yaml:"gmail"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true"                 env:"PG_URL"`
	}

	// JWT -.
	JWT struct {
		Secret string `env-required:"true" yaml:"secret" env:"JWT_SECRET"`
	}

	// Redis -.
	Redis struct {
		RedisHost string `env-required:"true" yaml:"host" env:"REDIS_HOST"`
		RedisPort int    `env-required:"true" yaml:"port" env:"REDIS_PORT"`
	}

	// Gmail -.
	Gmail struct {
		Email     string `env-required:"true" yaml:"email" env:"EMAIL"`
		EmailPass string `env-required:"true" yaml:"email_pass" env:"EMAIL_PASS"`
		Host      string `env-required:"true" yaml:"host" env:"SMTP_HOST"`
		Port      string    `env-required:"true" yaml:"port" env:"SMTP_PORT"`
	}

)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}




// ---------------------------------- for local ---------------------------------

// package config

// import (
// 	"fmt"
// 	"os"

// 	"github.com/joho/godotenv"
// )

// type (
// 	// Config -.
// 	Config struct {
// 		App   `yaml:"app"`
// 		HTTP  `yaml:"http"`
// 		Log   `yaml:"logger"`
// 		PG    `yaml:"postgres"`
// 		JWT   `yaml:"jwt"`
// 		Redis `yaml:"redis"`
// 		Gmail `yaml:"gmail"`
// 	}

// 	// App -.
// 	App struct {
// 		Name    string `env-required:"true" env:"APP_NAME"`
// 		Version string `env-required:"true" env:"APP_VERSION"`
// 	}

// 	// HTTP -.
// 	HTTP struct {
// 		Port string `env-required:"true" env:"HTTP_PORT"`
// 	}

// 	// Log -.
// 	Log struct {
// 		Level string `env-required:"true" env:"LOG_LEVEL"`
// 	}

// 	// PG -.
// 	PG struct {
// 		PoolMax int    `env-required:"true" env:"PG_POOL_MAX"`
// 		URL     string `env-required:"true" env:"PG_URL"`
// 	}

// 	// JWT -.
// 	JWT struct {
// 		Secret string `env-required:"true" env:"JWT_SECRET"`
// 	}

// 	// Redis -.
// 	Redis struct {
// 		RedisHost string `env-required:"true" env:"REDIS_HOST"`
// 		RedisPort int    `env-required:"true" env:"REDIS_PORT"`
// 	}

// 	// Gmail -.
// 	Gmail struct {
// 		Email     string `env-required:"true" env:"EMAIL"`
// 		EmailPass string `env-required:"true" env:"EMAIL_PASS"`
// 		Host      string `env-required:"true" env:"SMTP_HOST"`
// 		Port      string `env-required:"true" env:"SMTP_PORT"`
// 	}
// )

// // NewConfig returns app config.
// func NewConfig() (*Config, error) {
// 	// Load environment variables from .env file
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		return nil, fmt.Errorf("error loading .env file: %w", err)
// 	}

// 	cfg := &Config{
// 		App: App{
// 			Name:    os.Getenv("APP_NAME"),
// 			Version: os.Getenv("APP_VERSION"),
// 		},
// 		HTTP: HTTP{
// 			Port: os.Getenv("HTTP_PORT"),
// 		},
// 		Log: Log{
// 			Level: os.Getenv("LOG_LEVEL"),
// 		},
// 		PG: PG{
// 			PoolMax: mustGetInt("PG_POOL_MAX"),
// 			URL:     os.Getenv("PG_URL"),
// 		},
// 		JWT: JWT{
// 			Secret: os.Getenv("JWT_SECRET"),
// 		},
// 		Redis: Redis{
// 			RedisHost: os.Getenv("REDIS_HOST"),
// 			RedisPort: mustGetInt("REDIS_PORT"),
// 		},
// 		Gmail: Gmail{
// 			Email:     os.Getenv("EMAIL"),
// 			EmailPass: os.Getenv("EMAIL_PASS"),
// 			Host:      os.Getenv("SMTP_HOST"),
// 			Port:      os.Getenv("SMTP_PORT"),
// 		},
// 	}

// 	return cfg, nil
// }

// // mustGetInt reads an environment variable and converts it to an integer.
// func mustGetInt(key string) int {
// 	value := os.Getenv(key)
// 	if value == "" {
// 		panic(fmt.Sprintf("missing required environment variable: %s", key))
// 	}
// 	var intValue int
// 	_, err := fmt.Sscanf(value, "%d", &intValue)
// 	if err != nil {
// 		panic(fmt.Sprintf("invalid value for environment variable %s: %s", key, value))
// 	}
// 	return intValue
// }
