package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		App   App    `yaml:"app" validate:"required"`
		HTTP  HTTP   `yaml:"http" validate:"required"`
		PgSDN string `yaml:"pg_dsn"  validate:"required"`
	}
	App struct {
		Name    string `yaml:"name" validate:"required"`
		Version string `yaml:"version" validate:"required"`
	}

	HTTP struct {
		Port string `yaml:"port"  validate:"required"`
	}
)

func (c Config) Validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return validationErrors
	}
	return nil
}

func NewConfig() (*Config, error) {
	envFilePath := os.Getenv("ENV_FILE_PATH")
	godotenv.Load(envFilePath)
	projectRoot := os.Getenv("PROJECT_ROOT")
	config, err := ioutil.ReadFile(projectRoot + "/config/config.yml")
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	config = []byte(os.ExpandEnv(string(config)))

	cfg := &Config{}
	err = yaml.Unmarshal(config, cfg)

	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("cfg.Validate %w", err)
	}
	return cfg, nil
}

func NewTestConfig() (*Config, error) {
	cfg := &Config{}
	return cfg, nil
}
