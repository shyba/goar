package data_item

import (
	"io"

	"github.com/liteseed/goar/tag"
)

type DataItem struct {
	ID            string     `json:"id"`
	Signature     string     `json:"signature"`
	SignatureType int        `json:"signature_type"`
	Owner         string     `json:"owner"`
	Target        string     `json:"target"`
	Anchor        string     `json:"anchor"`
	Tags          *[]tag.Tag `json:"tags"`
	Data          string     `json:"data"` // Used only for serialization/deserialization
	Raw           []byte

	// Fields for streaming large data
	DataReader io.ReadSeeker `json:"-"` // Seekable reader for large data (required for multiple passes)
	DataSize   int64         `json:"-"` // Size of data for streaming
}
