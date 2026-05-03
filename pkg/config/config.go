package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"ai-gamification-elearning-server/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	App    App           `mapstructure:"app"`
	Logger logger.Config `mapstructure:"logger"`
	Server Server        `mapstructure:"server"`
	CORS   CORS          `mapstructure:"cors"`
}

type App struct {
	Name    string `mapstructure:"name"`
	Env     string `mapstructure:"env"`
	Port    int    `mapstructure:"port"`
	Version string `mapstructure:"version"`
	Debug   bool   `mapstructure:"debug"`
	BaseURL string `mapstructure:"base_url"`
}

type Server struct {
	ReadTimeout  int `mapstructure:"read_timeout"`
	WriteTimeout int `mapstructure:"write_timeout"`
}

type CORS struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAgeSeconds    int      `mapstructure:"max_age_seconds"`
}

func Load() (*Config, error) {
	if err := loadDotEnvLocal(".env.local"); err != nil {
		return nil, fmt.Errorf("config: reading .env.local: %w", err)
	}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AddConfigPath("./configs")
	v.AddConfigPath("configs")
	v.AddConfigPath("../configs")
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config: reading config.yaml: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshalling: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	if cfg.App.Port == 0 {
		return errors.New("app.port must be greater than 0")
	}

	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[cfg.App.Env] {
		return fmt.Errorf("app.env must be one of: development, staging, production (got %q)", cfg.App.Env)
	}

	return nil
}

func loadDotEnvLocal(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			return fmt.Errorf("line %d: invalid format (expected KEY=VALUE, got %q)", i+1, line)
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		if _, exists := os.LookupEnv(key); !exists {
			if err := os.Setenv(key, val); err != nil {
				return fmt.Errorf("line %d: setting %s: %w", i+1, key, err)
			}
		}
	}

	return nil
}
