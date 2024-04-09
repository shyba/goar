package tx

const (
	Arweave  = 1
	ED25519  = 2
	Ethereum = 3
	Solana   = 4
)

type SignatureMeta struct {
	SignatureLength int
	PublicKeyLength int
	Name            string
}

var SignatureConfig = map[int]SignatureMeta{
	Arweave: {
		SignatureLength: 512,
		PublicKeyLength: 512,
		Name:            "arweave",
	},
	ED25519: {
		SignatureLength: 64,
		PublicKeyLength: 32,
		Name:            "ed25519",
	},
	Ethereum: {
		SignatureLength: 65,
		PublicKeyLength: 65,
		Name:            "ethereum",
	},
	Solana: {
		SignatureLength: 64,
		PublicKeyLength: 32,
		Name:            "solana",
	},
}

type DataItem struct {
	ID            string `json:"id"`
	Signature     string `json:"signature"`
	SignatureType int    `json:"signature_type"`
	Owner         string `json:"owner"`
	Target        string `json:"target"`
	Anchor        string `json:"anchor"`
	Tags          []Tag  `json:"tags"`
	Data          string `json:"data"`
	Raw           []byte
}

type BundleHeader struct {
	id   int
	size int
	raw  []byte
}

type Bundle struct {
	Headers []BundleHeader `json:"bundle_header"`
	Items   []DataItem     `json:"items"`
	RawData string         `json:"raw_data"`
}

type Transaction struct {
	Format    int    `json:"format"`
	ID        string `json:"id"`
	LastTx    string `json:"last_tx"`
	Owner     string `json:"owner"`
	Tags      []Tag  `json:"tags"`
	Target    string `json:"target"`
	Quantity  string `json:"quantity"`
	Data      []byte `json:"data"`
	Reward    string `json:"reward"`
	Signature string `json:"signature"`
	DataSize  string `json:"data_size"`
	DataRoot  string `json:"data_root"`

	Chunks Chunks
}

type Chunk struct {
  DataHash []byte `json:"data_hash"`
  MinByteRange int `json:"min_byte_range"`
  MaxByteRange int `json:"max_byte_range"`
}

type Chunks struct {
	DataRoot string `json:"data_root"`
	DataSize string `json:"data_size"`
	DataPath string `json:"data_path"`
	Offset   string `json:"offset"`
	Chunks   []Chunk `json:"chunks"`
}
