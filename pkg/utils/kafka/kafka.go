package kafkaUtils

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/zakshearman/bluesky-creeper/pkg/config"
	"go.uber.org/zap"
	"time"
)

func CreateWriter(cfg config.KafkaConfig, log *zap.SugaredLogger, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		Async:        true,
		BatchTimeout: time.Millisecond * 500,
		ErrorLogger:  kafka.LoggerFunc(log.Errorw),
	}
}
