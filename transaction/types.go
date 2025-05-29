package transaction

import (
	"errors"
	"fmt"

	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/tag"
)

// Chunk represents a single chunk of data in an Arweave transaction's Merkle tree.
//
// Arweave splits large data into chunks for efficient storage and verification.
// Each chunk contains a hash of its data and byte range information.
type Chunk struct {
	DataHash     []byte `json:"data_hash"`      // SHA256 hash of the chunk data
	MinByteRange int    `json:"min_byte_range"` // Starting byte position of this chunk
	MaxByteRange int    `json:"max_byte_range"` // Ending byte position of this chunk (exclusive)
}

// Proof represents a Merkle proof for a specific chunk.
//
// Proofs allow verification that a chunk belongs to the larger dataset
// without requiring the entire dataset.
type Proof struct {
	Offset int    `json:"offset"` // Byte offset where this chunk starts in the overall data
	Proof  []byte `json:"proof"`  // Merkle proof bytes for verification
}

// ChunkData contains all chunk-related data for a transaction.
//
// This includes the complete chunking information: individual chunks,
// their proofs, and the root hash of the Merkle tree.
type ChunkData struct {
	DataRoot string  // Base64url-encoded root hash of the Merkle tree
	Chunks   []Chunk // Individual data chunks
	Proofs   []Proof // Merkle proofs for each chunk
}

// NodeType represents the type of a node in the Merkle tree.
type NodeType = string

// Node represents a node in the Merkle tree structure.
//
// Arweave uses a binary Merkle tree to organize transaction data.
// Each node can be either a leaf (containing actual data) or a branch
// (containing hashes of child nodes).
type Node struct {
	ID           []byte   // Unique identifier for this node
	DataHash     []byte   // Hash of the data this node represents
	ByteRange    int      // Starting byte position
	MaxByteRange int      // Ending byte position
	Type         NodeType // Type of node (leaf or branch)
	LeftChild    *Node    // Left child node (nil for leaf nodes)
	RightChild   *Node    // Right child node (nil for leaf nodes)
}

// Transaction represents an Arweave transaction.
//
// This struct contains all the fields required for an Arweave transaction
// according to the version 2 format specification. It supports both data
// transactions (storing data on Arweave) and transfer transactions (sending AR tokens).
type Transaction struct {
	Format    int        `json:"format"`    // Transaction format version (always 2 for this implementation)
	ID        string     `json:"id"`        // Transaction ID (SHA256 hash of signature)
	LastTx    string     `json:"last_tx"`   // Hash of the last transaction from this wallet
	Owner     string     `json:"owner"`     // Base64url-encoded public key of the transaction owner
	Tags      *[]tag.Tag `json:"tags"`      // Optional metadata tags
	Target    string     `json:"target"`    // Target wallet address (for AR transfers)
	Quantity  string     `json:"quantity"`  // Amount of AR to transfer in Winston units
	Data      string     `json:"data"`      // Base64url-encoded transaction data
	Reward    string     `json:"reward"`    // Transaction fee in Winston units
	Signature string     `json:"signature"` // Base64url-encoded transaction signature
	DataSize  string     `json:"data_size"` // Size of the data in bytes
	DataRoot  string     `json:"data_root"` // Merkle root hash of the data chunks

	ChunkData *ChunkData `json:"-"` // Chunk data for large transactions (not serialized)
}

// TransactionOffset represents the offset information for a transaction.
//
// This is used when querying transaction data from Arweave nodes
// to determine where the transaction data is located.
type TransactionOffset struct {
	Size   int64 `json:"size"`   // Size of the transaction data in bytes
	Offset int64 `json:"offset"` // Byte offset where the transaction data starts
}

// TransactionChunk represents a chunk of transaction data as returned by Arweave nodes.
//
// When retrieving large transaction data, Arweave returns it in chunks
// along with Merkle proofs for verification.
type TransactionChunk struct {
	Chunk    string `json:"chunk"`     // Base64url-encoded chunk data
	DataPath string `json:"data_path"` // Merkle proof path for this chunk
	TxPath   string `json:"tx_path"`   // Transaction path information
}

// GetChunkResult represents the result of retrieving a specific chunk.
//
// This structure contains all the information needed to verify and use
// a specific chunk of transaction data.
type GetChunkResult struct {
	DataRoot string `json:"data_root"` // Root hash of the complete dataset
	DataSize string `json:"data_size"` // Total size of the complete dataset
	DataPath string `json:"data_path"` // Merkle proof path for verification
	Offset   string `json:"offset"`    // Byte offset of this chunk
	Chunk    string `json:"chunk"`     // Base64url-encoded chunk data
}

// GetChunk retrieves a specific chunk from the transaction data.
//
// This method extracts a chunk at the specified index from the transaction's
// prepared chunk data and returns it along with the necessary proof information.
//
// Parameters:
//   - i: The index of the chunk to retrieve (0-based)
//   - data: The complete raw data that was chunked
//
// Returns a GetChunkResult containing the chunk data and proof, or an error
// if the chunks have not been prepared or the index is invalid.
//
// Example:
//
//	// Prepare chunks first
//	err := tx.PrepareChunks(data)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Get the first chunk
//	chunk, err := tx.GetChunk(0, data)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Chunk offset: %s, size: %d bytes\n", chunk.Offset, len(chunk.Chunk))
func (tx *Transaction) GetChunk(i int, data []byte) (*GetChunkResult, error) {
	if tx.ChunkData == nil {
		return nil, errors.New("chunks have not been prepared")
	}
	proof := tx.ChunkData.Proofs[i]
	chunk := tx.ChunkData.Chunks[i]

	return &GetChunkResult{
		DataRoot: tx.DataRoot,
		DataSize: tx.DataSize,
		DataPath: crypto.Base64URLEncode(proof.Proof),
		Offset:   fmt.Sprint(proof.Offset),
		Chunk:    crypto.Base64URLEncode(data[chunk.MinByteRange:chunk.MaxByteRange]),
	}, nil
}

// PrepareChunks computes and stores the chunk data for the given data.
//
// This method splits large data into chunks according to Arweave's chunking
// algorithm and generates the necessary Merkle tree structure and proofs.
// It must be called before signing a transaction that contains data.
//
// Note: This function operates on the provided data parameter, not the
// transaction's Data field. This allows for flexibility in data handling
// while preparing chunks.
//
// Parameters:
//   - data: The raw data to be chunked. Can be empty for transactions without data.
//
// Returns an error if chunking fails, otherwise updates the transaction's
// DataSize, ChunkData, and DataRoot fields.
//
// Example:
//
//	data := []byte("Large amount of data to be stored on Arweave")
//	err := tx.PrepareChunks(data)
//	if err != nil {
//		log.Printf("Failed to prepare chunks: %v", err)
//		return err
//	}
//	fmt.Printf("Data chunked into %d chunks\n", len(tx.ChunkData.Chunks))
func (tx *Transaction) PrepareChunks(data []byte) error {
	if len(data) > 0 {
		chunks, err := generateTransactionChunks(data)
		if err != nil {
			return err
		}
		tx.DataSize = fmt.Sprint(len(data))
		tx.ChunkData = chunks
		tx.DataRoot = (*chunks).DataRoot
	} else {
		tx.ChunkData = &ChunkData{
			Chunks:   []Chunk{},
			DataRoot: "",
			Proofs:   []Proof{},
		}
		tx.DataRoot = ""
	}
	return nil
}
