package crypto

import (
	"encoding/base64"
)

// Base64URL encode bytes to string
func Base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// Base64URL decode string to bytes
func Base64URLDecode(data string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(data)
}
