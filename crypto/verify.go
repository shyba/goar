package crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
)

// Verify to verify any data using the provided Arweave RSA Public Key
func Verify(data []byte, signature []byte, publicKey *rsa.PublicKey) error {
	hashed := sha256.Sum256(data)

	return rsa.VerifyPSS(publicKey, crypto.SHA256, hashed[:], signature, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	})
}
