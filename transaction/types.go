package transaction

import (
	"errors"
	"fmt"

	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/tag"
)

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
	ByteRange    int
	MaxByteRange int
	Type         NodeType
	LeftChild    *Node
	RightChild   *Node
}

type Transaction struct {
	Format    int       `json:"format"`
	ID        string    `json:"id"`
	LastTx    string    `json:"last_tx"`
	Owner     string    `json:"owner"`
	Tags      []tag.Tag `json:"tags"`
	Target    string    `json:"target"`
	Quantity  string    `json:"quantity"`
	Data      []byte    `json:"data"`
	Reward    string    `json:"reward"`
	Signature string    `json:"signature"`
	DataSize  string    `json:"data_size"`
	DataRoot  string    `json:"data_root"`

	ChunkData *ChunkData
}

type TransactionOffset struct {
	Size   int64 `json:"size"`
	Offset int64 `json:"offset"`
}
type TransactionChunk struct {
	Chunk    string `json:"chunk"`
	DataPath string `json:"data_path"`
	TxPath   string `json:"tx_path"`
}


type GetChunkResult struct {
	DataRoot string `json:"data_root"`
	DataSize string `json:"data_size"`
	DataPath string `json:"data_path"`
	Offset   string `json:"offset"`
	Chunk    string `json:"chunks"`
}

func (t *Transaction) GetChunk(i int, data []byte) (*GetChunkResult, error) {
	if t.ChunkData != nil {
		return nil, errors.New("chunks have not been prepared")
	}
	proof := t.ChunkData.Proofs[i]
	chunk := t.ChunkData.Chunks[i]

	return &GetChunkResult{
		DataRoot: t.DataRoot,
		DataSize: t.DataSize,
		DataPath: crypto.Base64Encode(proof.Proof),
		Offset:   fmt.Sprint(proof.Offset),
		Chunk:    crypto.Base64Encode(data[chunk.MinByteRange:chunk.MaxByteRange]),
	}, nil
}

// Note: we *do not* use `t.Data`, the caller may be
// operating on a transaction with an zero length data field.
// This function computes the chunks for the data passed in and
// assigns the result to this transaction. It should not read the
// data *from* this transaction.
func (t *Transaction) PrepareChunks(data []byte) error {
	if len(data) > 0 {
		chunks, err := generateTransactionChunks(data)
		if err != nil {
			return err
		}
		t.DataSize = fmt.Sprint(len(data))
		t.ChunkData = chunks
		t.DataRoot = (*chunks).DataRoot
	}
	return nil
}
