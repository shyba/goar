package types

type Chunk struct {
	DataHash     []byte `json:"data_hash"`
	MinByteRange int    `json:"min_byte_range"`
	MaxByteRange int    `json:"max_byte_range"`
}

type Proof struct {
	Offset int    `json:"offset"`
	Proof  []byte `json:"proof"`
}

type ChunkData struct {
	DataRoot string  `json:"data_root"`
	Chunks   []Chunk `json:"chunks"`
	Proofs   []Proof `json:"proofs"`
}

type NodeType = string
type Node struct {
	ID           []byte
	DataHash     []byte
	MinByteRange int
	MaxByteRange int
	Type         NodeType
	LeftChild    *Node
	RightChild   *Node
}

type Transaction struct {
	Format    int    `json:"format"`
	ID        string `json:"id"`
	LastTx    string `json:"last_tx"`
	Owner     string `json:"owner"`
	Tags      []Tag  `json:"tags"`
	Target    string `json:"target"`
	Quantity  string `json:"quantity"`
	Data      string `json:"data"`
	Reward    string `json:"reward"`
	Signature string `json:"signature"`
	DataSize  string `json:"data_size"`
	DataRoot  string `json:"data_root"`

	ChunkData *ChunkData
}

