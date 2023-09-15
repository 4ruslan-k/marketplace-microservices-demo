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
		App                      `yaml:"app"`
		HTTP                     `yaml:"http"`
		AccountsAppURL           string `yaml:"accounts_app_url"  validate:"required"`
		MarketplaceAppUrl        string `yaml:"marketplace_app_url"  validate:"required"`
		CatalogServiceURL        string `yaml:"catalog_service_url" validate:"required"`
		CartServiceURL           string `yaml:"cart_service_url" validate:"required"`
		ChatsServiceWebsocketURL string `yaml:"chat_service_websocket_url" validate:"required"`
		ChatsServiceURL          string `yaml:"chat_service_url" validate:"required"`
		AuthenticationServiceURL string `yaml:"authentication_service_url" validate:"required"`
		NotificationServiceURL   string `yaml:"notification_service_url" validate:"required"`
		SwaggerUIDomain          string `yaml:"swagger_ui_domain"`
		SwaggerEditorDomain      string `yaml:"swagger_editor_domain"`
		RedisAddress             string `yaml:"redis_address" validate:"required"`
	}
	App struct {
		Name    string `yaml:"name" validate:"required"`
		Version string `yaml:"version" validate:"required"`
	}

	HTTP struct {
		Port string `yaml:"port" env:"PORT" validate:"required"`
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
