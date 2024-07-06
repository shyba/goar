package transaction

import (
	"errors"
	"math"
	"reflect"

	"github.com/liteseed/goar/crypto"
)

type ValidatePathResult struct {
	Offset     int
	LeftBound  int
	RightBound int
	ChunkSize  int
}

const (
	MAX_CHUNK_SIZE = 256 * 1024
	MIN_CHUNK_SIZE = 32 * 1024
	NOTE_SIZE      = 32
	HASH_SIZE      = 32

	// Node Type

	Leaf = "Leaf"

	Branch = "Branch"
)

func generateTree(data []byte) (*Node, error) {
	chunks, err := chunkData(data)
	if err != nil {
		return nil, err
	}
	leaves, err := generateLeaves(chunks)
	if err != nil {
		return nil, err
	}
	rootNode, err := buildLayer(leaves, 0)
	if err != nil {
		return nil, err
	}
	return rootNode, err
}

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

	root, err := buildLayer(leaves, 0)
	if err != nil {
		return nil, err
	}

	proofs := generateProofs(root, nil, 0)

	// Discard the last chunk & proof if it's zero length.
	lastChunk := chunks[len(chunks)-1]
	if lastChunk.MaxByteRange-lastChunk.MinByteRange == 0 {
		chunks = chunks[:len(chunks)-1]
		proofs = proofs[:len(proofs)-1]
	}

	return &ChunkData{
		DataRoot: crypto.Base64URLEncode(root.ID),
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
			chunkSize = int(math.Ceil(float64(byteLength) / 2))
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

	hash, err := crypto.SHA256(rest)
	if err != nil {
		return nil, err
	}
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

		hashRange, err := crypto.SHA256(encodeUint(uint64(chunk.MaxByteRange)))
		if err != nil {
			return nil, err
		}

		ID, err := crypto.SHA256(append(hashDataHash, hashRange...))
		if err != nil {
			return nil, err
		}
		leaves = append(leaves, Node{
			ID:           ID,
			DataHash:     chunk.DataHash,
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
	if len(nodes) < 2 {
		return &nodes[0], nil
	}

	nextLayer := []Node{}
	for i := 0; i < len(nodes)-1; i += 2 {
		node, err := hashBranch(&nodes[i], &nodes[i+1])
		if err != nil {
			return nil, err
		}
		nextLayer = append(nextLayer, *node)
	}
	return buildLayer(nextLayer, level+1)
}

func hashBranch(left *Node, right *Node) (*Node, error) {
	if right == nil {
		return &Node{
			ID:           left.ID,
			DataHash:     left.DataHash,
			ByteRange:    left.ByteRange,
			MaxByteRange: left.MaxByteRange,
			Type:         Branch,
			LeftChild:    left.LeftChild,
			RightChild:   left.RightChild,
		}, nil
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
	ID, err := crypto.SHA256(append(leftIdHash, append(rightIdHash, leftMaxByteRangeHash...)...))
	if err != nil {
		return nil, err
	}
	return &Node{
		ID:           ID,
		ByteRange:    left.MaxByteRange,
		MaxByteRange: right.MaxByteRange,
		LeftChild:    left,
		RightChild:   right,
		Type:         Branch,
	}, nil
}

func generateProofs(node *Node, proof []byte, depth int) []Proof {
	proofs := []Proof{}
	if node.Type == Branch {
		partialProof := append(proof, append(node.LeftChild.ID, append(node.RightChild.ID, encodeUint(uint64(node.ByteRange))...)...)...)
		proofs = append(proofs, generateProofs(node.LeftChild, partialProof, depth+1)...)
		proofs = append(proofs, generateProofs(node.RightChild, partialProof, depth+1)...)
	} else if node.Type == Leaf {
		proofs = append(proofs, Proof{Offset: node.MaxByteRange - 1, Proof: append(append(proof, node.DataHash...), encodeUint(uint64(node.MaxByteRange))...)})
	}
	return proofs
}

func validatePath(id []byte, dest int, leftBound int, rightBound int, path []byte) (*ValidatePathResult, error) {
	if rightBound <= 0 {
		return nil, errors.New("out of bound right")
	}
	if dest >= rightBound {
		return validatePath(id, 0, rightBound-1, rightBound, path)
	}
	if dest < 0 {
		return validatePath(id, 0, 0, rightBound, path)
	}
	if len(path) == HASH_SIZE+NOTE_SIZE {
		pathData := path[0:HASH_SIZE]
		endOffsetBuffer := path[len(pathData) : len(pathData)+NOTE_SIZE]

		hash0, err := crypto.SHA256(pathData)
		if err != nil {
			return nil, err
		}

		hash1, err := crypto.SHA256(endOffsetBuffer)
		if err != nil {
			return nil, err
		}

		pathDataHash, err := crypto.SHA256(append(hash0, hash1...))
		if err != nil {
			return nil, err
		}

		if reflect.DeepEqual(id, pathDataHash) {
			return &ValidatePathResult{
				Offset:     rightBound - 1,
				LeftBound:  leftBound,
				RightBound: rightBound,
				ChunkSize:  rightBound - leftBound,
			}, nil
		}
		return nil, errors.New("invalid path")
	} else {
		left := path[0:HASH_SIZE]
		right := path[len(left) : len(left)+HASH_SIZE]
		offsetBuffer := path[len(left)+len(right) : len(left)+len(right)+NOTE_SIZE]
		offset := byteArrayToLong(offsetBuffer)
		remainder := path[len(left)+len(right)+len(offsetBuffer):]

		l, err := crypto.SHA256(left)
		if err != nil {
			return nil, err
		}
		r, err := crypto.SHA256(right)
		if err != nil {
			return nil, err
		}

		o, err := crypto.SHA256(offsetBuffer)
		if err != nil {
			return nil, err
		}

		p := []byte{}
		p = append(p, l...)
		p = append(p, r...)
		p = append(p, o...)
		pathDataHash, err := crypto.SHA256(p)
		if err != nil {
			return nil, err
		}

		if reflect.DeepEqual(id, pathDataHash) {
			if dest < int(offset) {
				return validatePath(
					left,
					dest,
					leftBound,
					min(rightBound, offset),
					remainder,
				)
			} else {
				return validatePath(
					right,
					dest,
					max(leftBound, offset),
					rightBound,
					remainder,
				)
			}
		}
	}
	return nil, errors.New("no valid path")
}

func flatten[T any](v []any) []T {
	proofs := []T{}
	for _, val := range v {
		if isSlice(val) {
			proofs = append(proofs, flatten[T](val.([]any))...)
		} else {
			proofs = append(proofs, val.(T))
		}
	}
	return proofs
}
