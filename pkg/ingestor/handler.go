package ingestor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bluesky-social/jetstream/pkg/models"
	"github.com/zakshearman/bluesky-creeper/pkg/bskytypes"
	"log"
	"time"
)

// imageURIFormat placeholders: DID, link ref
const imageURIFormat = "https://cdn.bsky.app/img/feed_fullsize/plain/%s/%s"

func (c *Client) handleEvent(_ context.Context, event *models.Event) error {
	if event.Commit == nil || event.Commit.Record == nil {
		return nil
	}

	record := event.Commit.Record
	var post bskytypes.Post
	if err := json.Unmarshal(record, &post); err != nil {
		log.Fatalf("Failed to unmarshal record: %v", err)
	}

	parsedEvent := bskytypes.PostEvent{
		Did:    event.Did,
		TimeUS: time.UnixMicro(event.TimeUS),
		Post:   post,
	}

	//log.Printf("Received post: %+v", parsedEvent)

	if post.Embed != nil && len(post.Embed.Images) > 0 {
		for i, image := range post.Embed.Images {
			log.Printf("Image(%d): %s", i, fmt.Sprintf(imageURIFormat, parsedEvent.Did, image.Image.Ref.Link))
		}
	}

	//if err := c.notifier.NotifyPostCreated(c.ctx, parsedEvent); err != nil {
	//	return fmt.Errorf("failed to notify post created: %w", err)
	//}

	return nil
}
