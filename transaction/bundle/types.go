package bundle

import "github.com/liteseed/goar/transaction/data_item"

type BundleHeader struct {
	ID   int
	Size int
}

type Bundle struct {
	Headers []BundleHeader       `json:"bundle_header"`
	Items   []data_item.DataItem `json:"items"`
	Raw     []byte
}
