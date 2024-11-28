package ingestor

import (
	"context"
	"errors"
	"fmt"
	jetstreamclient "github.com/bluesky-social/jetstream/pkg/client"
	"github.com/bluesky-social/jetstream/pkg/client/schedulers/parallel"
	"log"
	"log/slog"
	"time"
)

type Client struct {
	ctx    context.Context
	cancel context.CancelFunc

	notifier *KafkaNotifier
}

func NewIngestorClient(notifier *KafkaNotifier) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		ctx:    ctx,
		cancel: cancel,

		notifier: notifier,
	}
}

func (c *Client) Start(_ context.Context) error {
	scheduler := parallel.NewScheduler(250, "jetstream-prod", slog.Default(), c.handleEvent)

	conf := jetstreamclient.DefaultClientConfig()
	conf.WantedCollections = []string{"app.bsky.feed.post"}
	conf.WebsocketURL = "wss://jetstream.atproto.tools/subscribe"
	conf.Compress = true

	jetstreamClient, err := jetstreamclient.NewClient(conf, slog.Default(), scheduler)
	if err != nil {
		return fmt.Errorf("failed to create jetstream client: %w", err)
	}

	go func() {
		cursorTime := time.Now().Add(-time.Hour * 8)
		micros := cursorTime.UnixMicro()
		log.Printf("Starting jetstream client at time: %d", micros)

		err = jetstreamClient.ConnectAndRead(c.ctx, &micros)
		if !errors.Is(err, context.Canceled) {
			log.Fatalf("HandleRepoStream returned unexpectedly: %+v...", err)
		} else {
			log.Printf("HandleRepoStream closed on context cancel...")
		}
	}()

	return nil
}

func (c *Client) Shutdown(_ context.Context) error {
	c.cancel()
	return nil
}
