package loms_config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required,number,gt=0,lte=65535"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	DbName   string `yaml:"db_name" validate:"required"`
}

type KafkaConfig struct {
	Host       string `yaml:"host" validate:"required"`
	Port       string `yaml:"port" validate:"required,number,gt=0,lte=65535"`
	OrderTopic string `yaml:"order_topic" validate:"required"`
	Brokers    string `yaml:"brokers" validate:"required"`
	PollMs     int64  `yaml:"poll" validate:"required,number,gt=0"`
}

type ServerConfig struct {
	Host         string `yaml:"host" validate:"required"`
	HttpPort     string `yaml:"http_port" validate:"required,number,gt=0,lte=65535"`
	GrpcPort     string `yaml:"grpc_port" validate:"required,number,gt=0,lte=65535"`
	AllowSwagger bool   `yaml:"allow_swagger"`
	IsInMemory   bool   `yaml:"in_memory"`
}

type JaegerConfig struct {
	Host string `yaml:"host" validate:"required"`
	Port string `yaml:"port" validate:"required,number,gt=0,lte=65535"`
}

type Config struct {
	Server    ServerConfig   `yaml:"service"`
	MasterDb  DatabaseConfig `yaml:"db_master"`
	ReplicaDb DatabaseConfig `yaml:"db_replica"`
	Kafka     KafkaConfig    `yaml:"kafka"`
	Jaeger    JaegerConfig   `yaml:"jaeger"`
}

func LoadLomsConfig(filename string) (*Config, error) {
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
