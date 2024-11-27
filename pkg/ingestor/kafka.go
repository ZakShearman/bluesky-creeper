package ingestor

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/zakshearman/bluesky-creeper/pkg/bskytypes"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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
	msg, err := p.ToPostCreatedEvent()
	if err != nil {
		return err
	}

	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	return k.w.WriteMessages(ctx,
		kafka.Message{
			Key:     []byte(p.Did),
			Headers: []kafka.Header{{Key: "X-Proto-Type", Value: []byte(msg.ProtoReflect().Descriptor().FullName())}},
			Value:   bytes,
		},
	)
}

func (k *KafkaNotifier) Shutdown(_ context.Context) error {
	return k.w.Close()
}
