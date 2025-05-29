package crypto

import "crypto/sha256"

// SHA256 computes the SHA256 hash of the provided data.
//
// This is a convenience function that wraps Go's standard crypto/sha256
// package to return the hash as a byte slice instead of a fixed-size array.
// SHA256 is used extensively throughout the Arweave protocol for creating
// transaction IDs, chunk hashes, and Merkle tree nodes.
//
// Parameters:
//   - data: The raw binary data to hash
//
// Returns the SHA256 hash as a 32-byte slice.
//
// Example:
//
//	data := []byte("Hello, Arweave!")
//	hash := SHA256(data)
//	fmt.Printf("SHA256: %x\n", hash)
//	// Output: SHA256: a1b2c3d4e5f6... (32 bytes)
func SHA256(data []byte) []byte {
	r := sha256.Sum256(data)
	return r[:]
}
