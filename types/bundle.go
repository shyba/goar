package types

type BundleHeader struct {
	ID   int
	Size int
	Raw  []byte
}

type Bundle struct {
	Headers []BundleHeader `json:"bundle_header"`
	Items   []DataItem     `json:"items"`
	RawData string         `json:"raw_data"`
}
