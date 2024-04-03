package tx

import (
	"encoding/json"
	"os"
)

const (
	ArweaveSignType  = 1
	ED25519SignType  = 2
	EthereumSignType = 3
	SolanaSignType   = 4
)

type SigMeta struct {
	SigLength int
	PubLength int
	SigName   string
}

var SigConfigMap = map[int]SigMeta{
	ArweaveSignType: {
		SigLength: 512,
		PubLength: 512,
		SigName:   "arweave",
	},
	ED25519SignType: {
		SigLength: 64,
		PubLength: 32,
		SigName:   "ed25519",
	},
	EthereumSignType: {
		SigLength: 65,
		PubLength: 65,
		SigName:   "ethereum",
	},
	SolanaSignType: {
		SigLength: 64,
		PubLength: 32,
		SigName:   "solana",
	},
}

type Bundle struct {
	Items            []DataItem `json:"items"`
	BundleBinary     []byte
	BundleDataReader *os.File
}

type DataItem struct {
	SignatureType int    `json:"signatureType"`
	Signature     string `json:"signature"`
	Owner         string `json:"owner"`
	Target        string `json:"target"` // optional, if exist must length 32, and is base64 str
	Anchor        string `json:"anchor"` // optional, if exist must length 32, and is base64 str
	Data          string `json:"data"`
	Id            string `json:"id"`
	Tags          string `json:"tags"`

	RawData    []byte   `json:"-"`
	DataReader *os.File `json:"-"`
}

type Tag struct {
	Name  string `json:"name" avro:"name"`
	Value string `json:"value" avro:"value"`
}

type TransactionChunk struct {
	Chunk    string `json:"chunk"`
	DataPath string `json:"data_path"`
	TxPath   string `json:"tx_path"`
}

type TransactionOffset struct {
	Size   string `json:"size"`
	Offset string `json:"offset"`
}

type TxStatus struct {
	BlockHeight           int    `json:"block_height"`
	BlockIndepHash        string `json:"block_indep_hash"`
	NumberOfConfirmations int    `json:"number_of_confirmations"`
}

type BundlrResp struct {
	Id                  string   `json:"id"`
	Signature           string   `json:"signature"`
	N                   string   `json:"n"`
	Public              string   `json:"public"`
	Block               int64    `json:"block"`
	ValidatorSignatures []string `json:"validatorSignatures"`
}

type Chunks struct {
	DataRoot []byte   `json:"data_root"`
	Chunks   []Chunk  `json:"chunks"`
	Proofs   []*Proof `json:"proofs"`
}

type Chunk struct {
	DataHash     []byte
	MinByteRange int
	MaxByteRange int
}

// Node include leaf node and branch node
type Node struct {
	ID           []byte
	Type         string // "branch" or "leaf"
	DataHash     []byte // only leaf node
	MinByteRange int    // only leaf node
	MaxByteRange int
	ByteRange    int   // only branch node
	LeftChild    *Node // only branch node
	RightChild   *Node // only branch node
}

type Proof struct {
	Offest int
	Proof  []byte
}

type Transaction struct {
	Format     int      `json:"format"`
	ID         string   `json:"id"`
	LastTx     string   `json:"last_tx"`
	Owner      string   `json:"owner"` // utils.Base64Encode(wallet.PubKey.N.Bytes())
	Tags       []Tag    `json:"tags"`
	Target     string   `json:"target"`
	Quantity   string   `json:"quantity"`
	Data       string   `json:"data"` // base64.encode
	DataReader *os.File `json:"-"`    // when dataSize too big use dataReader, set Data = ""
	DataSize   string   `json:"data_size"`
	DataRoot   string   `json:"data_root"`
	Reward     string   `json:"reward"`
	Signature  string   `json:"signature"`

	// Computed when needed.
	Chunks *Chunks `json:"-"`
}

type GetChunk struct {
	DataRoot string `json:"data_root"`
	DataSize string `json:"data_size"`
	DataPath string `json:"data_path"`
	Offset   string `json:"offset"`
	Chunk    string `json:"chunk"`
}

func (gc *GetChunk) Marshal() ([]byte, error) {
	return json.Marshal(gc)
}
