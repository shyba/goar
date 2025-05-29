package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// Sign creates an RSA-PSS signature for the given data using an Arweave private key.
//
// This function implements the signature algorithm used by Arweave for transaction
// signing. It uses RSA-PSS (Probabilistic Signature Scheme) with SHA256 hashing
// and automatic salt length as specified by the Arweave protocol.
//
// The signing process:
// 1. Computes SHA256 hash of the input data
// 2. Signs the hash using RSA-PSS with the provided private key
// 3. Uses automatic salt length for optimal security
//
// Parameters:
//   - data: The raw data to sign (typically transaction signature data)
//   - privateKey: The RSA private key to sign with (should be 4096-bit for Arweave)
//
// Returns the signature bytes or an error if signing fails.
//
// Example:
//
//	// Load your private key
//	signer, err := signer.FromPath("wallet.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Sign some data
//	data := []byte("Hello, Arweave!")
//	signature, err := Sign(data, signer.PrivateKey)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Signature: %x\n", signature)
func Sign(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hashed := sha256.Sum256(data)

	return rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hashed[:], &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	})
}
