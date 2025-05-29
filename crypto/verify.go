package crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
)

// Verify validates an RSA-PSS signature using an Arweave public key.
//
// This function implements the signature verification algorithm used by Arweave.
// It verifies RSA-PSS signatures created with SHA256 hashing and automatic salt
// length, matching the signature format used in Arweave transactions.
//
// The verification process:
// 1. Computes SHA256 hash of the input data
// 2. Verifies the signature against the hash using RSA-PSS
// 3. Uses automatic salt length matching the signing process
//
// Parameters:
//   - data: The original data that was signed
//   - signature: The signature bytes to verify
//   - publicKey: The RSA public key to verify against
//
// Returns nil if the signature is valid, or an error if verification fails.
//
// Example:
//
//	// Verify a transaction signature
//	tx, err := client.GetTransactionByID("txid...")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Get the signature data and verify
//	err = tx.Verify()
//	if err != nil {
//		log.Printf("Invalid signature: %v", err)
//	} else {
//		fmt.Println("Signature is valid")
//	}
func Verify(data []byte, signature []byte, publicKey *rsa.PublicKey) error {
	hashed := sha256.Sum256(data)

	return rsa.VerifyPSS(publicKey, crypto.SHA256, hashed[:], signature, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	})
}
