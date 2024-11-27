package bskytypes

import (
	"github.com/zakshearman/bluesky-creeper/pkg/proto/generated"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"
)

type PostEvent struct {
	// Did is a Decentralized Identifier. This is an ID for the entity that created the post. PLC DIDs can be resolved here: https://web.plc.directory/resolve
	Did    string
	TimeUS time.Time
	Post   Post
}

func (p *PostEvent) ToPostCreatedEvent() (*generated.PostCreatedEvent, error) {
	facets := make([]*generated.Facet, len(p.Post.Facets))
	for i, f := range p.Post.Facets {
		facets[i] = f.ToProto()
	}
	return &generated.PostCreatedEvent{
		Did:          p.Did,
		EventTime_US: timestamppb.New(p.TimeUS),
		PostTime_US:  timestamppb.New(p.Post.CreatedAt),
		Text:         p.Post.Text,
		Languages:    p.Post.Languages,
		Facets:       facets,
	}, nil
}

type Post struct {
	CreatedAt time.Time `json:"createdAt"`        // Timestamp of the post creation.
	Text      string    `json:"text"`             // The content of the post.
	Languages []string  `json:"langs,omitempty"`  // Languages of the post.
	Embed     *Embed    `json:"embed,omitempty"`  // Embedding data (if present).
	Facets    []Facet   `json:"facets,omitempty"` // Rich text features like mentions and links.
	Reply     *Reply    `json:"reply,omitempty"`  // Reply context (if the post is a reply).
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

func (f *Facet) ToProto() *generated.Facet {
	features := make([]*generated.Feature, len(f.Features))
	for i, f := range f.Features {
		features[i] = f.ToProto()
	}

	return &generated.Facet{
		Features:  features,
		ByteStart: f.Index.ByteStart,
		ByteEnd:   f.Index.ByteEnd,
	}
}

// Feature represents a specific rich text feature like a mention or hashtag.
type Feature struct {
	Type string `json:"$type"` // Type of the feature (e.g., "app.bsky.richtext.facet#mention").
	DID  string `json:"did"`   // DID of the mentioned entity (if applicable).
	Tag  string `json:"tag"`   // Tag name (if applicable, e.g., for hashtags).
	URI  string `json:"uri"`   // URI for linked content (if applicable).
}

func (f *Feature) ToProto() *generated.Feature {
	feature := &generated.Feature{}
	var fType generated.FeatureType
	switch strings.Split(f.Type, "#")[1] {
	case "link":
		fType = generated.FeatureType_LINK
		feature.LinkUri = &f.URI
	case "mention":
		fType = generated.FeatureType_MENTION
		feature.MentionTargetDid = &f.DID
	case "tag":
		fType = generated.FeatureType_TAG
		feature.TagValue = &f.Tag
	default:
		fType = generated.FeatureType_UNKNOWN
	}

	return &generated.Feature{
		Type: fType,
	}
}

// ByteIndex specifies the byte range of a feature in the post's text.
type ByteIndex struct {
	ByteStart int32 `json:"byteStart"` // Start byte position.
	ByteEnd   int32 `json:"byteEnd"`   // End byte position.
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
