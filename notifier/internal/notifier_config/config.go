package notifier_config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type KafkaConfig struct {
	Host            string `yaml:"host" validate:"required"`
	Port            string `yaml:"port" validate:"required,number,gt=0,lte=65535"`
	OrderTopic      string `yaml:"order_topic" validate:"required"`
	ConsumerGroupID string `yaml:"consumer_group_id" validate:"required"`
	BrokersPath     string `yaml:"brokers" validate:"required"`
}

type Config struct {
	Kafka KafkaConfig `yaml:"kafka"`
}

func LoadNotifierConfig(filename string) (*Config, error) {
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
