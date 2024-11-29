package ingestor

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/zakshearman/bluesky-creeper/pkg/bskytypes"
	"github.com/zakshearman/bluesky-creeper/pkg/config"
	kafkaUtils "github.com/zakshearman/bluesky-creeper/pkg/utils/kafka"
	"go.uber.org/zap"
	"log"
)

const topic = "raw-posts"

type KafkaNotifier struct {
	w *kafka.Writer
}

func NewKafkaNotifier(cfg config.KafkaConfig, log *zap.SugaredLogger) *KafkaNotifier {
	return &KafkaNotifier{
		w: kafkaUtils.CreateWriter(cfg, log, topic),
	}
}

func (k *KafkaNotifier) NotifyPostCreated(ctx context.Context, p bskytypes.PostEvent) error {
	bytes, err := json.Marshal(p)
	if err != nil {
		return err
	}

	log.Printf("Sending post: %+v", string(bytes))

	return k.w.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(p.Did),
			Value: bytes,
		},
	)
}

func (k *KafkaNotifier) Shutdown(_ context.Context) error {
	return k.w.Close()
}
