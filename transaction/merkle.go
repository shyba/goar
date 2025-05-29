package transaction

import (
	"errors"
	"math"
	"reflect"

	"github.com/liteseed/goar/crypto"
)

// ValidatePathResult contains the result of validating a Merkle path.
//
// This structure is returned when validating that a specific chunk
// belongs to a Merkle tree and provides information about the chunk's
// position and boundaries within the complete dataset.
type ValidatePathResult struct {
	Offset     int // The offset of the chunk within the complete dataset
	LeftBound  int // The left boundary of the chunk's byte range
	RightBound int // The right boundary of the chunk's byte range
	ChunkSize  int // The size of the chunk in bytes
}

// Merkle tree and chunking constants used by Arweave protocol
const (
	MAX_CHUNK_SIZE = 256 * 1024 // Maximum size of a single chunk (256KB)
	MIN_CHUNK_SIZE = 32 * 1024  // Minimum size of a single chunk (32KB)
	NOTE_SIZE      = 32         // Size of note/offset information in bytes
	HASH_SIZE      = 32         // Size of SHA256 hash in bytes

	// Node types for Merkle tree structure
	Leaf   = "Leaf"   // Leaf node containing actual data chunk
	Branch = "Branch" // Branch node containing child node references
)

// generateTree creates a complete Merkle tree from the provided data.
//
// This function implements Arweave's chunking and Merkle tree generation
// algorithm. It splits the data into appropriate-sized chunks and builds
// a binary Merkle tree structure for efficient verification.
//
// Parameters:
//   - data: The raw data to create a Merkle tree from
//
// Returns the root node of the generated Merkle tree, or an error if
// tree generation fails.
//
// Example:
//
//	data := []byte("Large dataset to be stored on Arweave")
//	rootNode, err := generateTree(data)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Root hash: %x\n", rootNode.ID)
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

// generateTransactionChunks generates the complete chunk data needed for an Arweave transaction.
//
// This function creates all the components required for a transaction with data:
// - Data root hash (Merkle tree root)
// - Individual chunks with their hashes and byte ranges
// - Merkle proofs for each chunk
//
// The function also handles the Arweave protocol requirement to discard any
// zero-length chunks that may be generated at the end of the chunking process.
//
// Parameters:
//   - data: The raw data to be chunked and processed
//
// Returns ChunkData containing the data root, chunks, and proofs, or an error
// if processing fails.
//
// Example:
//
//	data := []byte("Data to be uploaded to Arweave")
//	chunkData, err := generateTransactionChunks(data)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Generated %d chunks with root: %s\n",
//		len(chunkData.Chunks), chunkData.DataRoot)
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

// chunkData splits transaction data into chunks according to Arweave's chunking algorithm.
//
// This function implements the specific chunking strategy used by Arweave:
//   - Preferred chunk size is MAX_CHUNK_SIZE (256KB)
//   - If the remaining data would create a chunk smaller than MIN_CHUNK_SIZE (32KB),
//     the current chunk is split to create two roughly equal chunks
//   - Each chunk includes a SHA256 hash and byte range information
//
// Parameters:
//   - data: The raw data to be chunked
//
// Returns a slice of Chunk structs containing hash and range information
// for each chunk, or an error if chunking fails.
//
// Example:
//
//	data := []byte("Data to be chunked")
//	chunks, err := chunkData(data)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for i, chunk := range chunks {
//		fmt.Printf("Chunk %d: bytes %d-%d, hash: %x\n",
//			i, chunk.MinByteRange, chunk.MaxByteRange, chunk.DataHash)
//	}
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

// generateLeaves creates leaf nodes for the Merkle tree from data chunks.
//
// Each leaf node represents a single chunk of data and contains:
// - A unique ID (hash of the chunk hash and chunk size)
// - The original data hash
// - The maximum byte range (end position) of the chunk
// - Type set to "Leaf"
//
// Parameters:
//   - chunks: Slice of chunks to convert to leaf nodes
//
// Returns a slice of Node structs representing the leaf level of the Merkle tree.
//
// Example:
//
//	chunks := []Chunk{{DataHash: hash1, MaxByteRange: 1024}}
//	leaves, err := generateLeaves(chunks)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Created %d leaf nodes\n", len(leaves))
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

// buildLayer recursively builds the Merkle tree from a layer of nodes.
//
// This function creates parent nodes by pairing adjacent nodes and continues
// recursively until only one root node remains. It handles odd numbers of
// nodes by promoting the last node to the next layer.
//
// Parameters:
//   - nodes: The current layer of nodes to build upon
//   - level: The current level in the tree (used for tracking depth)
//
// Returns the root node of the tree when construction is complete.
//
// Example:
//
//	leaves := []Node{leaf1, leaf2, leaf3, leaf4}
//	root, err := buildLayer(leaves, 0)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Tree built with root ID: %x\n", root.ID)
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

// hashBranch creates a branch node from two child nodes.
//
// This function implements Arweave's branch hashing algorithm:
// - Hash = SHA256(SHA256(left.ID) + SHA256(right.ID) + SHA256(left.MaxByteRange))
// - If right node is nil, returns left node unchanged
//
// Parameters:
//   - left: The left child node
//   - right: The right child node (can be nil for odd numbers of nodes)
//
// Returns a new branch node containing the two children, or an error if
// hashing fails.
//
// Example:
//
//	leftNode := &Node{ID: leftHash, MaxByteRange: 1024}
//	rightNode := &Node{ID: rightHash, MaxByteRange: 2048}
//	branch, err := hashBranch(leftNode, rightNode)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Branch node created with ID: %x\n", branch.ID)
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

// generateProofs recursively generates Merkle proofs for all chunks in the tree.
//
// A Merkle proof allows verification that a specific chunk belongs to the
// complete dataset without requiring the entire dataset. The proof contains
// the path from the chunk to the root of the Merkle tree.
//
// Parameters:
//   - node: The current node being processed
//   - proof: The accumulated proof data from parent nodes
//   - depth: The current depth in the tree (used for tracking)
//
// Returns a slice of Proof structs, one for each leaf node reachable
// from the given node.
//
// Example:
//
//	proofs := generateProofs(rootNode, nil, 0)
//	fmt.Printf("Generated %d proofs\n", len(proofs))
//	for i, proof := range proofs {
//		fmt.Printf("Proof %d: offset=%d, size=%d bytes\n",
//			i, proof.Offset, len(proof.Proof))
//	}
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

// validatePath verifies that a Merkle path is valid for a given chunk.
//
// This function verifies that a provided Merkle proof correctly proves
// that a chunk at a specific destination belongs to a dataset with the
// given root hash. It recursively validates the path through the tree.
//
// Parameters:
//   - id: The root hash of the Merkle tree
//   - dest: The byte offset of the chunk being verified
//   - leftBound: The left boundary of the current search range
//   - rightBound: The right boundary of the current search range
//   - path: The Merkle proof data to validate
//
// Returns ValidatePathResult with chunk information if the path is valid,
// or an error if validation fails.
//
// Example:
//
//	result, err := validatePath(rootHash, 1024, 0, 4096, proofBytes)
//	if err != nil {
//		log.Printf("Invalid proof: %v", err)
//	} else {
//		fmt.Printf("Valid chunk at offset %d, size %d\n",
//			result.Offset, result.ChunkSize)
//	}
func validatePath(id []byte, dest int, leftBound int, rightBound int, path []byte) (*ValidatePathResult, error) {
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

// flatten is a generic utility function that flattens nested slices into a single slice.
//
// This function recursively processes nested slice structures and flattens them
// into a single-level slice of the specified type T. It's used internally for
// processing proof data structures.
//
// Type Parameters:
//   - T: The type of elements in the flattened slice
//
// Parameters:
//   - v: A slice of any type that may contain nested slices
//
// Returns a flattened slice containing all elements of type T found in the
// nested structure.
//
// Example:
//
//	nested := []any{[]int{1, 2}, 3, []int{4, 5}}
//	flat := flatten[int](nested)
//	// Result: [1, 2, 3, 4, 5]
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
