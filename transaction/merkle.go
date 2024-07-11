package transaction

import (
	"errors"
	"log"
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
	var chunks []Chunk

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
		dataSha := crypto.SHA256(chunk)

		cursor += len(chunk)
		chunks = append(chunks, Chunk{
			DataHash:     dataSha[:],
			MinByteRange: cursor - len(chunk),
			MaxByteRange: cursor,
		})

		rest = rest[chunkSize:]
	}

	hash := crypto.SHA256(rest)
	chunks = append(chunks, Chunk{
		DataHash:     hash[:],
		MinByteRange: cursor,
		MaxByteRange: cursor + len(rest),
	})
	return chunks, nil
}

func generateLeaves(chunks []Chunk) ([]Node, error) {
	var leaves []Node
	for _, chunk := range chunks {
		ID := crypto.SHA256(append(crypto.SHA256(chunk.DataHash), crypto.SHA256(intToByteArray(chunk.MaxByteRange))...))
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

	var nextLayer []Node
	for i := 0; i < len(nodes); i += 2 {
		var next *Node
		if i+1 < len(nodes) {
			next = &nodes[i+1]
		}
		node, err := hashBranch(&nodes[i], next)
		if err != nil {
			return nil, err
		}
		nextLayer = append(nextLayer, *node)
	}
	return buildLayer(nextLayer, level+1)
}

func hashBranch(left *Node, right *Node) (*Node, error) {
	if right == nil {
		return left, nil
	}
	ID := crypto.SHA256(
		append(crypto.SHA256(left.ID),
			append(
				crypto.SHA256(right.ID),
				crypto.SHA256(intToByteArray(left.MaxByteRange))...,
			)...,
		),
	)
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
	var proofs []Proof
	if node.Type == Leaf {
		var p []byte
		p = append(p, proof...)
		p = append(p, node.DataHash...)
		p = append(p, intToByteArray(node.MaxByteRange)...)
		proofs = append(proofs, Proof{Offset: node.MaxByteRange - 1, Proof: p})
	}
	if node.Type == Branch {
		var partialProof []byte
		partialProof = append(partialProof, proof...)
		partialProof = append(partialProof, node.LeftChild.ID...)
		partialProof = append(partialProof, node.RightChild.ID...)
		partialProof = append(partialProof, intToByteArray(node.ByteRange)...)
		proofs = append(proofs, generateProofs(node.LeftChild, partialProof, depth+1)...)
		proofs = append(proofs, generateProofs(node.RightChild, partialProof, depth+1)...)
	}

	return proofs
}

func validatePath(id []byte, dest int, leftBound int, rightBound int, path []byte) (*ValidatePathResult, error) {
	log.Println(crypto.Base64URLEncode(id), dest, leftBound, rightBound, len(path))
	if rightBound <= 0 {
		return nil, errors.New("right bound < 0")
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
		h := crypto.SHA256(append(crypto.SHA256(pathData), crypto.SHA256(endOffsetBuffer)...))
		if reflect.DeepEqual(id, h) {
			return &ValidatePathResult{
				Offset:     rightBound - 1,
				LeftBound:  leftBound,
				RightBound: rightBound,
				ChunkSize:  rightBound - leftBound,
			}, nil
		}
		return nil, errors.New("invalid path")
	}
	left := path[0:HASH_SIZE]
	right := path[len(left) : len(left)+HASH_SIZE]
	offsetBuffer := path[len(left)+len(right) : len(left)+len(right)+NOTE_SIZE]
	offset := byteArrayToInt(offsetBuffer)
	remainder := path[len(left)+len(right)+len(offsetBuffer):]

	var p []byte
	p = append(p, crypto.SHA256(left)...)
	p = append(p, crypto.SHA256(right)...)
	p = append(p, crypto.SHA256(offsetBuffer)...)

	if reflect.DeepEqual(id, crypto.SHA256(p)) {
		if dest < offset {
			return validatePath(
				left,
				dest,
				leftBound,
				min(rightBound, offset),
				remainder,
			)
		}
		return validatePath(
			right,
			dest,
			max(leftBound, offset),
			rightBound,
			remainder,
		)
	}
	return nil, errors.New("no valid path")
}

func flatten[T any](v []any) []T {
	var proofs []T
	for _, val := range v {
		if isSlice(val) {
			proofs = append(proofs, flatten[T](val.([]any))...)
		} else {
			proofs = append(proofs, val.(T))
		}
	}
	return proofs
}
