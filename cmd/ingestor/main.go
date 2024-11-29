package main

import (
	"github.com/zakshearman/bluesky-creeper/pkg/config"
	"github.com/zakshearman/bluesky-creeper/pkg/ingestor"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		// Config
		fx.Supply(config.LoadCommonConfig(config.ComponentKafka)),

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

func newKafkaNotifier(cfg config.CommonConfig, log *zap.SugaredLogger, lc fx.Lifecycle) *ingestor.KafkaNotifier {
	notifier := ingestor.NewKafkaNotifier(cfg.Kafka, log)
	lc.Append(fx.Hook{
		OnStop: notifier.Shutdown,
	})
	return notifier
}

func newIngestorClient(notifier *ingestor.KafkaNotifier, lc fx.Lifecycle) *ingestor.Client {
	client := ingestor.NewIngestorClient(notifier)
	lc.Append(fx.Hook{
		OnStart: client.Start,
		OnStop:  client.Shutdown,
	})
	return client
}

//func newElasticClient() {
//	es, err := elasticsearch.NewTypedClient(elasticsearch.Config{})
//	es.Index().Do()
//	es.
//}
