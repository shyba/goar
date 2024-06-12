package data_item

import "github.com/liteseed/goar/tag"

type DataItem struct {
	ID            string    `json:"id"`
	Signature     string    `json:"signature"`
	SignatureType int       `json:"signature_type"`
	Owner         string    `json:"owner"`
	Target        string    `json:"target"`
	Anchor        string    `json:"anchor"`
	Tags          []tag.Tag `json:"tags"`
	Data          string    `json:"data"`
	Raw           []byte
}
