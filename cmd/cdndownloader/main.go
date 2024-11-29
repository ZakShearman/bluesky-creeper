package cdndownloader

import (
	"github.com/zakshearman/bluesky-creeper/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		// Config
		fx.Supply(config.LoadCommonConfig(config.ComponentKafka, config.ComponentMinio)),

		// Logging
		fx.Provide(
			newZapLogger,
			newZapSugared,
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),

		// Kafka
		fx.Provide(newKafkaNotifier),

		// Ingestor listener
		fx.Invoke(newIngestorClient),
	).Run()
}

func newZapLogger(conf config.CommonConfig) (*zap.Logger, error) {
	if conf.Env == config.EnvValueProd {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

func newZapSugared(log *zap.Logger) *zap.SugaredLogger {
	zap.ReplaceGlobals(log)
	return log.Sugar()
}
