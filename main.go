package main

import (
	"context"
	"encoding/json"
	"errors"
	jetstreamclient "github.com/bluesky-social/jetstream/pkg/client"
	"github.com/bluesky-social/jetstream/pkg/client/schedulers/parallel"
	"github.com/bluesky-social/jetstream/pkg/models"
	"log"
	"log/slog"
	"time"
)

func main() {
	scheduler := parallel.NewScheduler(6000, "jetstream-prod", slog.Default(), func(ctx context.Context, event *models.Event) error {
		if event.Commit == nil || event.Commit.Record == nil {
			return nil
		}

		record := event.Commit.Record
		var post Post
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

		parsedEvent := PostEvent{
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

		log.Printf("Received post: %+v", parsedEvent)
		return nil
	})

	conf := jetstreamclient.DefaultClientConfig()
	conf.WantedCollections = []string{"app.bsky.feed.post"}
	conf.WebsocketURL = "wss://jetstream.atproto.tools/subscribe"
	conf.Compress = true

	jetstreamClient, err := jetstreamclient.NewClient(conf, slog.Default(), scheduler)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	go func() {
		//cursor := time.Now().Add(-time.Hour * 8).UnixMicro()

		ctx := context.Background()
		go func() {
			err = jetstreamClient.ConnectAndRead(ctx, nil)
			if !errors.Is(err, context.Canceled) {
				log.Fatalf("HandleRepoStream returned unexpectedly, killing the consumer: %+v...", err)
			} else {
				log.Printf("HandleRepoStream closed on context cancel...")
			}
		}()
	}()

	select {}
}

type PostEvent struct {
	// Did is a Decentralized Identifier. This is an ID for the entity that created the post. PLC DIDs can be resolved here: https://web.plc.directory/resolve
	Did    string
	TimeUS time.Time
	Post   Post
}

type Post struct {
	Type      string    `json:"$type"`           // The type of the post, e.g., "app.bsky.feed.post".
	CreatedAt time.Time `json:"createdAt"`       // Timestamp of the post creation.
	Text      string    `json:"text"`            // The content of the post.
	Languages []string  `json:"langs,omitempty"` // Languages of the post.
	//Embed     *Embed    `json:"embed,omitempty"`  // Embedding data (if present).
	//Facets    []Facet   `json:"facets,omitempty"` // Rich text features like mentions and links.
	//Reply     *Reply    `json:"reply,omitempty"`  // Reply context (if the post is a reply).
}

// Embed represents an embedded record or media within a post.
type Embed struct {
	Type     string         `json:"$type"`              // Type of the embed (e.g., "app.bsky.embed.record", "app.bsky.embed.images").
	Record   *RecordEmbed   `json:"record,omitempty"`   // Details of the record embed.
	Images   []ImageEmbed   `json:"images,omitempty"`   // Details of image embeds (if applicable).
	External *ExternalEmbed `json:"external,omitempty"` // External link details (if present).
}

// RecordEmbed represents an embedded record.
type RecordEmbed struct {
	CID string `json:"cid"` // Content ID of the embedded record.
	URI string `json:"uri"` // URI of the embedded record.
}

// ImageEmbed represents a single image in an embedded media.
type ImageEmbed struct {
	Alt         string     `json:"alt"`         // Alt text for the image.
	AspectRatio Dimensions `json:"aspectRatio"` // Dimensions of the image.
	Image       Blob       `json:"image"`       // Blob data for the image.
}

// ExternalEmbed represents an external link embed.
type ExternalEmbed struct {
	URI         string `json:"uri"`         // The URI of the external content.
	Title       string `json:"title"`       // Title of the external content.
	Description string `json:"description"` // Description of the external content.
	Thumb       Blob   `json:"thumb"`       // Thumbnail image for the external content.
}

// Dimensions represents the aspect ratio or size of an image.
type Dimensions struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

// Blob represents binary data for an image or file.
type Blob struct {
	Type     string `json:"$type"`    // Blob type (e.g., "blob").
	Ref      Ref    `json:"ref"`      // Reference to the blob.
	MimeType string `json:"mimeType"` // MIME type of the blob.
	Size     int    `json:"size"`     // Size of the blob in bytes.
}

// Ref represents a reference to a blob.
type Ref struct {
	Link string `json:"$link"` // Link to the blob.
}

// Facet represents a rich text feature such as mentions, hashtags, or links.
type Facet struct {
	Type     string    `json:"$type"`    // Type of the facet.
	Features []Feature `json:"features"` // Features within the facet.
	Index    ByteIndex `json:"index"`    // Byte range of the feature within the text.
}

// Feature represents a specific rich text feature like a mention or hashtag.
type Feature struct {
	Type string `json:"$type"` // Type of the feature (e.g., "app.bsky.richtext.facet#mention").
	DID  string `json:"did"`   // DID of the mentioned entity (if applicable).
	Tag  string `json:"tag"`   // Tag name (if applicable, e.g., for hashtags).
	URI  string `json:"uri"`   // URI for linked content (if applicable).
}

// ByteIndex specifies the byte range of a feature in the post's text.
type ByteIndex struct {
	ByteStart int `json:"byteStart"` // Start byte position.
	ByteEnd   int `json:"byteEnd"`   // End byte position.
}

// Reply represents a reply context for a post.
type Reply struct {
	Parent Reference `json:"parent"` // Parent post details.
	Root   Reference `json:"root"`   // Root post details.
}

// Reference represents a reference to another post.
type Reference struct {
	CID string `json:"cid"` // Content ID of the referenced post.
	URI string `json:"uri"` // URI of the referenced post.
}
