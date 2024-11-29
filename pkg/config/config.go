package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zakshearman/bluesky-creeper/pkg/utils/runtime"
	"strings"
)

type Component string

const (
	ComponentKafka Component = "kafka"
	ComponentMinio Component = "minio"
)

const (
	envFlag = "env"

	kafkaAddressFlag = "kafka-address"
	kafkaPortFlag    = "kafka-port"

	minioAddressFlag     = "minio-address"
	minioPortFlag        = "minio-port"
	minioAccessKeyIDFlag = "minio-access-key-id"
	minioAccessSecret    = "minio-access-secret"
)

var componentFlags = map[Component][]configVar{
	ComponentKafka: {
		{kafkaAddressFlag, "localhost", "Kafka address"},
		{kafkaPortFlag, 9092, "Kafka port"},
	},
	ComponentMinio: {
		{minioAddressFlag, "localhost", "Minio address"},
		{minioPortFlag, 9000, "Minio port"},
		{minioAccessKeyIDFlag, nil, "Minio access key ID"},
		{minioAccessSecret, nil, "Minio access secret"},
	},
}

type configVar struct {
	key          string
	defaultValue any

	usage string
}

type Env string

const (
	EnvValueProd Env = "prod"
	EnvValueDev  Env = "dev"
)

var envs = map[string]Env{"prod": EnvValueProd, "dev": EnvValueDev}

type KafkaConfig struct {
	Address string
	Port    int
}

type MinioConfig struct {
	Address string
	Port    int

	AccessKeyID     string
	SecretAccessKey string
}

type CommonConfig struct {
	Components []Component

	Env   Env
	Kafka KafkaConfig
	Minio MinioConfig
}

func LoadCommonConfig(components ...Component) CommonConfig {
	viper.SetDefault(envFlag, EnvValueDev)
	pflag.String(envFlag, viper.GetString(envFlag), "Environment - dev or prod")

	for _, component := range components {
		for _, flag := range componentFlags[component] {
			if flag.defaultValue != nil {
				viper.SetDefault(flag.key, flag.defaultValue)
			}

			pflag.String(flag.key, viper.GetString(flag.key), flag.usage)
		}
	}

	pflag.Parse()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	for _, component := range components {
		for _, flag := range componentFlags[component] {
			runtime.Must(viper.BindEnv(flag.key))
		}
	}

	cfg := CommonConfig{
		Components: components,
		Env:        envs[viper.GetString(envFlag)],
	}

	if hasComponent(components, ComponentKafka) {
		cfg.Kafka = KafkaConfig{
			Address: viper.GetString(kafkaAddressFlag),
			Port:    viper.GetInt(kafkaPortFlag),
		}
	}

	if hasComponent(components, ComponentMinio) {
		cfg.Minio = MinioConfig{
			Address:         viper.GetString(minioAddressFlag),
			Port:            viper.GetInt(minioPortFlag),
			AccessKeyID:     viper.GetString(minioAccessKeyIDFlag),
			SecretAccessKey: viper.GetString(minioAccessSecret),
		}
	}

	return cfg
}

func hasComponent(components []Component, component Component) bool {
	for _, c := range components {
		if c == component {
			return true
		}
	}
	return false
}
