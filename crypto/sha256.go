package crypto

import "crypto/sha256"

func SHA256(data []byte) []byte {
	r := sha256.Sum256(data)
	return r[:]
}
