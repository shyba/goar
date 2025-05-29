// Package crypto provides cryptographic utilities for the Arweave protocol.
//
// This package includes functions for base64url encoding/decoding, SHA256 hashing,
// deep hash computation, and RSA signature operations as used in Arweave transactions.
package crypto

import (
	"encoding/base64"
)

// Base64URLEncode encodes bytes to a Base64URL string.
//
// This function uses base64 URL encoding without padding as specified
// by RFC 4648. This encoding is used throughout the Arweave protocol
// for transaction IDs, signatures, and other binary data.
//
// Parameters:
//   - data: The byte data to encode
//
// Returns the base64url-encoded string representation.
//
// Example:
//
//	data := []byte("Hello, Arweave!")
//	encoded := Base64URLEncode(data)
//	fmt.Printf("Encoded: %s\n", encoded)
//	// Output: SGVsbG8sIEFyd2VhdmUh
func Base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// Base64URLDecode decodes a Base64URL string to bytes.
//
// This function decodes base64 URL-encoded strings without padding
// as used throughout the Arweave protocol. It's the inverse operation
// of Base64URLEncode.
//
// Parameters:
//   - data: The base64url-encoded string to decode
//
// Returns the decoded bytes or an error if the string is invalid.
//
// Example:
//
//	encoded := "SGVsbG8sIEFyd2VhdmUh"
//	decoded, err := Base64URLDecode(encoded)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Decoded: %s\n", string(decoded))
//	// Output: Hello, Arweave!
func Base64URLDecode(data string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(data)
}
