package types

import (
	"crypto/sha256"

	"github.com/liteseed/goar/crypto"
)

const (
	MAX_CHUNK_SIZE = 256 * 1024
	MIN_CHUNK_SIZE = 32 * 1024
	NOTE_SIZE      = 32
	HASH_SIZE      = 3

	// Node Type

	Leaf = "Leaf"

	Branch = "Branch"
)

/**
 * Generates the data_root, chunks & proofs
 * needed for a transaction.
 *
 * This also checks if the last chunk is a zero-length
 * chunk and discards that chunk and proof if so.
 * (we do not need to upload this zero length chunk)
 */
func generateTransactionChunks(data []byte) (*ChunkData, error) {
	chunks, err := chunkData(data)
	if err != nil {
		return nil, err
	}

	leaves, err := generateLeaves(chunks)
	if err != nil {
		return nil, err
	}

	root, err := buildLayer(leaves, 0) // leaf node level == 0
	if err != nil {
		return nil, err
	}

	proofs := generateProofs(root, make([]byte, 0), 0)

	// Discard the last chunk & proof if it's zero length.
	lastChunk := chunks[len(chunks)-1]
	if lastChunk.MaxByteRange-lastChunk.MinByteRange == 0 {
		chunks = chunks[:len(chunks)-1]
		proofs = proofs[:len(proofs)-1]
	}

	return &ChunkData{
		DataRoot: string(root.ID),
		Chunks:   chunks,
		Proofs:   proofs,
	}, nil

}

func chunkData(data []byte) ([]Chunk, error) {
	chunks := []Chunk{}

	rest := data
	cursor := 0

	for len(rest) >= MAX_CHUNK_SIZE {
		chunkSize := MAX_CHUNK_SIZE
		byteLength := len(rest)

		nextChunkSize := byteLength - MAX_CHUNK_SIZE

		if nextChunkSize > 0 && nextChunkSize < MIN_CHUNK_SIZE {
			chunkSize = (byteLength + byteLength&1) / 2
		}

		chunk := rest[:chunkSize]
		dataHash, err := crypto.SHA256(chunk)
		if err != nil {
			return nil, err
		}

		cursor += len(chunk)
		chunks = append(chunks, Chunk{
			DataHash:     dataHash[:],
			MinByteRange: cursor - len(chunk),
			MaxByteRange: cursor,
		})

		rest = rest[chunkSize:]
	}

	hash := sha256.Sum256(rest)
	chunks = append(chunks, Chunk{
		DataHash:     hash[:],
		MinByteRange: cursor,
		MaxByteRange: cursor + len(rest),
	})
	return chunks, nil
}

func generateLeaves(chunks []Chunk) ([]Node, error) {
	leaves := []Node{}
	for _, chunk := range chunks {
		hashDataHash, err := crypto.SHA256(chunk.DataHash)
		if err != nil {
			return nil, err
		}
		hashMaxByteRange, err := crypto.SHA256(encodeUint(uint64(chunk.MaxByteRange)))
		if err != nil {
			return nil, err
		}

		id := append(hashDataHash, hashMaxByteRange...)
		leaves = append(leaves, Node{
			ID:           id,
			DataHash:     chunk.DataHash,
			MinByteRange: chunk.MinByteRange,
			MaxByteRange: chunk.MaxByteRange,
			LeftChild:    nil,
			RightChild:   nil,
			Type:         Leaf,
		})
	}
	return leaves, nil
}

// buildLayer
func buildLayer(nodes []Node, level int) (*Node, error) {
	if len(nodes) == 1 {
		return &nodes[0], nil
	}
	node, err := hashBranch(&nodes[level], &nodes[level+1])
	if err != nil {
		return nil, err
	}
	return node, nil
}

func hashBranch(left *Node, right *Node) (*Node, error) {
	if right == nil {
		return left, nil
	}
	leftIdHash, err := crypto.SHA256(left.ID)
	if err != nil {
		return nil, err
	}
	rightIdHash, err := crypto.SHA256(right.ID)
	if err != nil {
		return nil, err
	}
	leftMaxByteRangeHash, err := crypto.SHA256(encodeUint(uint64(left.MaxByteRange)))
	if err != nil {
		return nil, err
	}

	id := append(leftIdHash, append(rightIdHash, leftMaxByteRangeHash...)...)
	return &Node{
		ID:           id,
		MinByteRange: left.MinByteRange,
		MaxByteRange: left.MaxByteRange,
		LeftChild:    left,
		RightChild:   right,
		Type:         Branch,
	}, nil
}

func generateProofs(node *Node, proofData []byte, depth int) []Proof {
	if node.Type == Leaf {
		return []Proof{{
			Offset: node.MaxByteRange - 1,
			Proof:  append(proofData, append(node.DataHash, encodeUint(uint64(node.MaxByteRange))...)...),
		}}
	}
	if node.Type == Branch {
		left := node.LeftChild
		right := node.RightChild
		partialProofData := append(proofData, left.ID...)
		partialProofData = append(partialProofData, right.ID...)
		partialProofData = append(partialProofData, encodeUint(uint64(node.MaxByteRange))...)

		leftProofs := generateProofs(left, partialProofData, depth+1)
		rightProofs := generateProofs(right, partialProofData, depth+1)

		return append(leftProofs, rightProofs...)
	}
	return []Proof{}
}
