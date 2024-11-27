package ingestor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bluesky-social/jetstream/pkg/models"
	"github.com/zakshearman/bluesky-creeper/pkg/bskytypes"
	"log"
	"strings"
	"time"
)

var keywords = []string{"uwu", "owo"}

func (c *Client) handleEvent(_ context.Context, event *models.Event) error {
	if event.Commit == nil || event.Commit.Record == nil {
		return nil
	}

	record := event.Commit.Record
	var post bskytypes.Post
	if err := json.Unmarshal(record, &post); err != nil {
		log.Fatalf("Failed to unmarshal record: %v", err)
	}

	containsEn := false
	for _, lang := range post.Languages {
		if lang == "en" {
			containsEn = true
			break
		}
	}

	if !containsEn {
		return nil
	}

	parsedEvent := bskytypes.PostEvent{
		Did:    event.Did,
		TimeUS: time.UnixMicro(event.TimeUS),
		Post:   post,
	}

	//jsonOutput, err := json.Marshal(event.Commit.Record)
	//if err != nil {
	//	log.Fatalf("Failed to marshal event: %v", err)
	//}
	//
	//log.Printf("Received event: %s", jsonOutput)
	//
	//log.Printf("Received post: %+v", parsedEvent)

	contains := false
	for _, keyword := range keywords {
		if contains = strings.Contains(strings.ToLower(post.Text), keyword); contains {
			break
		}
	}

	if contains {
		log.Printf("Received post: %+v", parsedEvent.Post.Text)
	}

	// Log to Kafka
	if err := c.notifier.NotifyPostCreated(c.ctx, parsedEvent); err != nil {
		return fmt.Errorf("failed to notify post created: %w", err)
	}

	return nil
}
