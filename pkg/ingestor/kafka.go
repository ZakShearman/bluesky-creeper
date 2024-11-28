package ingestor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/zakshearman/bluesky-creeper/pkg/bskytypes"
	"go.uber.org/zap"
	"log"
	"time"
)

const topic = "raw-posts"

type KafkaNotifier struct {
	w *kafka.Writer
}

func NewKafkaNotifier(cfg KafkaConfig, log *zap.SugaredLogger) *KafkaNotifier {
	w := &kafka.Writer{
		Addr:         kafka.TCP(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		Async:        true,
		BatchTimeout: time.Millisecond * 500,
		ErrorLogger:  kafka.LoggerFunc(log.Errorw),
	}

	return &KafkaNotifier{
		w: w,
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
