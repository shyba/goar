package crypto

import (
	"encoding/base64"
)

// Base64URLEncode Encode bytes to Base64URL string
func Base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// Base64URLDecode Decode Base64URL string to bytes
func Base64URLDecode(data string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(data)
}
