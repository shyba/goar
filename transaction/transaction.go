// Package transaction provides functionality for creating, signing, and verifying Arweave transactions.
//
// This package implements the Arweave transaction format version 2 and provides
// utilities for working with transaction data, signatures, and verification.
//
// Example usage:
//
//	data := []byte("Hello, Arweave!")
//	tags := []tag.Tag{{Name: "Content-Type", Value: "text/plain"}}
//	tx := transaction.New(data, "", "0", &tags)
//
//	signer := wallet.Signer()
//	err := tx.Sign(signer)
//	if err != nil {
//		log.Fatal(err)
//	}
package transaction

import (
	"errors"

	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/tag"
)

// New creates a new Arweave transaction with the provided data and metadata.
//
// Parameters:
//   - data: The data to include in the transaction. Can be nil for transactions without data.
//   - target: The target wallet address for AR transfers. Use empty string for data-only transactions.
//   - quantity: The amount of AR to transfer in Winston units. Use "0" for data-only transactions.
//   - tags: Optional metadata tags for the transaction. Can be nil.
//
// Returns a new Transaction struct with format version 2, which is the current
// standard for Arweave transactions.
//
// Example:
//
//	// Data transaction with tags
//	tags := []tag.Tag{{Name: "Content-Type", Value: "application/json"}}
//	tx := New(jsonData, "", "0", &tags)
//
//	// AR transfer transaction
//	tx := New(nil, targetAddress, "1000000000000", nil) // 1 AR in Winston
func New(data []byte, target string, quantity string, tags *[]tag.Tag) *Transaction {
	if tags == nil {
		tags = &[]tag.Tag{}
	}
	if quantity == "" {
		quantity = "0"
	}
	if data == nil {
		data = []byte("")
	}
	return &Transaction{
		Format:   2,
		Data:     crypto.Base64URLEncode(data),
		Target:   target,
		Quantity: quantity,
		Tags:     tag.ConvertToBase64(tags),
		DataSize: "0",
	}
}

// Sign signs the transaction using the provided signer and generates the transaction ID.
//
// This method:
// 1. Generates the signature data from the transaction fields
// 2. Signs the data using the signer's private key
// 3. Sets the transaction ID as the SHA256 hash of the signature
// 4. Sets the signature field with the base64url-encoded signature
//
// Parameters:
//   - s: A signer containing the private key to sign with
//
// Returns an error if signing fails or if the transaction format is unsupported.
//
// Example:
//
//	signer := wallet.Signer()
//	err := tx.Sign(signer)
//	if err != nil {
//		log.Printf("Failed to sign transaction: %v", err)
//		return err
//	}
//	fmt.Printf("Transaction signed with ID: %s", tx.ID)
func (tx *Transaction) Sign(s *signer.Signer) error {
	payload, err := tx.getSignatureData()
	if err != nil {
		return err
	}
	rawSignature, err := crypto.Sign(payload, s.PrivateKey)
	if err != nil {
		return err
	}
	tx.ID = crypto.Base64URLEncode(crypto.SHA256(rawSignature))
	tx.Signature = crypto.Base64URLEncode(rawSignature)
	return nil
}

// Verify verifies the transaction signature against the transaction data.
//
// This method:
// 1. Regenerates the signature data from the transaction fields
// 2. Extracts the public key from the Owner field
// 3. Verifies the signature against the data using the public key
//
// Returns nil if the signature is valid, or an error if verification fails.
// This is useful for validating transactions received from other sources.
//
// Example:
//
//	err := tx.Verify()
//	if err != nil {
//		log.Printf("Transaction signature invalid: %v", err)
//		return err
//	}
//	fmt.Println("Transaction signature verified successfully")
func (tx *Transaction) Verify() error {
	signatureData, err := tx.getSignatureData()
	if err != nil {
		return err
	}
	rawSignature, err := crypto.Base64URLDecode(tx.Signature)
	if err != nil {
		return err
	}
	publicKey, err := crypto.GetPublicKeyFromOwner(tx.Owner)
	if err != nil {
		return err
	}
	return crypto.Verify(signatureData, rawSignature, publicKey)
}

// getSignatureData generates the data that should be signed for this transaction.
//
// This internal method implements the Arweave signature data format for version 2
// transactions. It creates a deep hash of the transaction components in the
// correct order as specified by the Arweave protocol.
//
// The signature data includes:
// - Format version ("2")
// - Owner (public key)
// - Target address
// - Quantity in Winston
// - Reward amount
// - Last transaction hash
// - Tags
// - Data size
// - Data root (Merkle root of data chunks)
//
// Returns the signature data as bytes, or an error if the transaction format
// is unsupported or if any field cannot be decoded.
func (tx *Transaction) getSignatureData() ([]byte, error) {
	if tx.Format != 2 {
		return nil, errors.New("only type 2 transaction supported")
	}
	rawOwner, err := crypto.Base64URLDecode(tx.Owner)
	if err != nil {
		return nil, err
	}
	rawTarget, err := crypto.Base64URLDecode(tx.Target)
	if err != nil {
		return nil, err
	}

	rawTags, err := tag.Decode(tx.Tags)
	if err != nil {
		return nil, err
	}

	rawLastTx, err := crypto.Base64URLDecode(tx.LastTx)
	if err != nil {
		return nil, err
	}

	data, err := crypto.Base64URLDecode(tx.Data)
	if err != nil {
		return nil, err
	}

	err = tx.PrepareChunks(data)
	if err != nil {
		return nil, err
	}

	rawDataRoot, err := crypto.Base64URLDecode(tx.DataRoot)
	if err != nil {
		return nil, err
	}

	chunks := []any{
		[]byte("2"),
		rawOwner,
		rawTarget,
		[]byte(tx.Quantity),
		[]byte(tx.Reward),
		rawLastTx,
		rawTags,
		[]byte(tx.DataSize),
		rawDataRoot,
	}

	deepHash := crypto.DeepHash(chunks)
	signatureData := deepHash[:]
	return signatureData, nil
}
