package ingestor

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zakshearman/bluesky-creeper/pkg/utils/runtime"
	"strings"
)

const (
	envFlag = "env"

	kafkaAddressFlag = "kafka-address"
	kafkaPortFlag    = "kafka-port"
)

type IngestorConfig struct {
	Env   Env
	Kafka KafkaConfig
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

func LoadIngestorConfig() IngestorConfig {
	viper.SetDefault(envFlag, EnvValueDev)
	viper.SetDefault(kafkaAddressFlag, "localhost")
	viper.SetDefault(kafkaPortFlag, 9092)

	pflag.String(envFlag, viper.GetString(envFlag), "Environment - dev or prod")
	pflag.String(kafkaAddressFlag, viper.GetString(kafkaAddressFlag), "Kafka address")
	pflag.Int(kafkaPortFlag, viper.GetInt(kafkaPortFlag), "Kafka port")
	pflag.Parse()

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	runtime.Must(viper.BindEnv(envFlag))
	runtime.Must(viper.BindEnv(kafkaAddressFlag))
	runtime.Must(viper.BindEnv(kafkaPortFlag))

	return IngestorConfig{
		Env: envs[viper.GetString(envFlag)],
		Kafka: KafkaConfig{
			Address: viper.GetString(kafkaAddressFlag),
			Port:    viper.GetInt(kafkaPortFlag),
		},
	}
}
