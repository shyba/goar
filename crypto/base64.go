package crypto

import (
	"encoding/base64"
)

// Base64 encode bytes to string
func Base64Encode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// Base64 decode string to bytes
func Base64Decode(data string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(data)
}
