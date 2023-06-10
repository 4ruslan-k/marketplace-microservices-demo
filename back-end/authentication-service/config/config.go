package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		App               App           `yaml:"app" validate:"required"`
		HTTP              HTTP          `yaml:"http" validate:"required"`
		MongoURI          string        `yaml:"mongo_url" validate:"required"`
		MongoDatabaseName string        `yaml:"mongo_database_name" validate:"required"`
		FrontendURL       string        `yaml:"frontend_url" validate:"required"`
		GatewayURL        string        `yaml:"gateway_url" validate:"required"`
		SessionSecret     string        `yaml:"session_secret" validate:"required,min=10"`
		SocialSignIn      SocialSignIn  `yaml:"social_sign_in"`
		OpenTelemetry     OpenTelemetry `yaml:"open_telemetry" validate:"required"`
	}
	App struct {
		Name    string `yaml:"name" validate:"required"`
		Version string `yaml:"version" validate:"required"`
	}

	HTTP struct {
		Port string `yaml:"port"  validate:"required"`
	}

	SocialSignIn struct {
		GithubKey    string `yaml:"github_key" validate:"required"`
		GithubSecret string `yaml:"github_secret" validate:"required"`
		GoogleKey    string `yaml:"google_key" validate:"required"`
		GoogleSecret string `yaml:"google_secret" validate:"required"`
	}

	OpenTelemetry struct {
		IsEnabled    *bool  `yaml:"is_enabled" validate:"required"`
		IsInsecure   *bool  `yaml:"is_insecure" validate:"required"`
		CollectorURL string `yaml:"collector_url" validate:"required"`
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
	godotenv.Load(".env.test")
	cfg := &Config{}
	absPath, _ := filepath.Abs("../../../config/config_test.yml")

	err := cleanenv.ReadConfig(absPath, cfg)

	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
