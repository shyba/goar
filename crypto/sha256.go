package crypto

import "crypto/sha256"

// SHA256 - Convert raw binary data to SHA256 hashed bytes. This is convenience function for the library.
func SHA256(data []byte) []byte {
	r := sha256.Sum256(data)
	return r[:]
}
