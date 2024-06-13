package bundle

import "github.com/liteseed/goar/transaction/data_item"

type BundleHeader struct {
	ID   int
	Size int
	Raw  []byte
}

type Bundle struct {
	Headers []BundleHeader       `json:"bundle_header"`
	Items   []data_item.DataItem `json:"items"`
	RawData string               `json:"raw_data"`
}
