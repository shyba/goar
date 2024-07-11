package bundle

import "github.com/liteseed/goar/transaction/data_item"

type Header struct {
	ID   string
	Size int
	Raw  []byte
}

type Bundle struct {
	Headers []Header             `json:"bundle_header"`
	Items   []data_item.DataItem `json:"items"`
	Raw     []byte
}
