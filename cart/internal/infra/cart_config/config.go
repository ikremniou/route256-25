package cart_config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" validate:"required"`
		Port string `yaml:"port" validate:"required,number,gt=0,lte=65535"`
	} `yaml:"service"`

	Products struct {
		Host  string `yaml:"host" validate:"required"`
		Port  string `yaml:"port" validate:"required,number,gt=0,lte=65535"`
		Token string `yaml:"token" validate:"required"`
	} `yaml:"product_service"`

	Loms struct {
		Host string `yaml:"host" validate:"required"`
		Port string `yaml:"port" validate:"required,number,gt=0,lte=65535"`
	} `yaml:"loms_service"`

	Jaeger struct {
		Host string `yaml:"host" validate:"required"`
		Port string `yaml:"port" validate:"required,number,gt=0,lte=65535"`
	} `yaml:"jaeger"`
}

func LoadCartConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(file, config); err != nil {
		return nil, err
	}

	var validator = validator.New(validator.WithRequiredStructEnabled())
	if err := validator.Struct(config); err != nil {
		return nil, err
	}

	return config, nil
}
