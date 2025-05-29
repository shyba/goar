// Package signer provides cryptographic signing functionality for Arweave transactions.
//
// This package handles RSA key management and transaction signing operations
// used in the Arweave protocol. It supports loading keys from JWK format,
// generating new keys, and creating signatures for transactions.
//
// Example usage:
//
//	// Create a new signer with generated key
//	signer, err := New()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Load signer from JWK file
//	signer, err := FromPath("wallet.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Get wallet address
//	address := signer.Address
//	fmt.Printf("Wallet address: %s\n", address)
package signer

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/everFinance/gojwk"
	"github.com/liteseed/goar/crypto"
)

// Signer represents an Arweave wallet signer with RSA key pair.
//
// A Signer contains the complete cryptographic identity for an Arweave wallet,
// including the wallet address, public key, and private key. It provides
// methods for signing transactions and managing key data.
type Signer struct {
	Address    string          // The Arweave wallet address derived from the public key
	PublicKey  *rsa.PublicKey  // RSA public key for verification operations
	PrivateKey *rsa.PrivateKey // RSA private key for signing operations
}

// New creates a new Signer with a randomly generated RSA key pair.
//
// This function generates a new 4096-bit RSA key pair suitable for use
// with the Arweave protocol. The generated key is automatically converted
// to JWK format and then loaded into a Signer instance.
//
// Returns a new Signer instance with a fresh key pair, or an error if
// key generation fails.
//
// Example:
//
//	signer, err := New()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Generated new wallet: %s\n", signer.Address)
func New() (*Signer, error) {
	bitSize := 4096
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}
	jwk, err := gojwk.PrivateKey(key)
	if err != nil {
		return nil, err
	}
	data, err := gojwk.Marshal(jwk)
	if err != nil {
		return nil, err
	}
	return FromJWK(data)
}

// FromPath creates a Signer from a JWK file on disk.
//
// This function reads a JSON Web Key (JWK) file from the specified path
// and creates a Signer instance from the contained RSA private key.
// The file should contain a JWK-formatted RSA private key as typically
// exported by Arweave wallet software.
//
// Parameters:
//   - path: The file system path to the JWK file
//
// Returns a Signer instance loaded with the key from the file, or an error
// if the file cannot be read or contains invalid key data.
//
// Example:
//
//	signer, err := FromPath("./wallet.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Loaded wallet: %s\n", signer.Address)
func FromPath(path string) (*Signer, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return FromJWK(b)
}

// FromJWK creates a Signer from JWK data in memory.
//
// This function parses JSON Web Key (JWK) data and extracts the RSA
// private key to create a Signer instance. The JWK data should contain
// a valid RSA private key in the standard JWK format.
//
// Parameters:
//   - b: The JWK data as bytes (should be valid JSON)
//
// Returns a Signer instance with the loaded key and computed address,
// or an error if the JWK data is invalid or cannot be parsed.
//
// Example:
//
//	jwkData := []byte(`{"kty":"RSA","n":"...","e":"AQAB",...}`)
//	signer, err := FromJWK(jwkData)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Loaded wallet: %s\n", signer.Address)
func FromJWK(b []byte) (*Signer, error) {
	key, err := gojwk.Unmarshal(b)
	if err != nil {
		return nil, err
	}
	rsaPublicKey, err := key.DecodePublicKey()
	if err != nil {
		return nil, err
	}
	publicKey, ok := rsaPublicKey.(*rsa.PublicKey)
	if !ok {
		err = fmt.Errorf("public key type error")
		return nil, err
	}

	rsaPrivateKey, err := key.DecodePrivateKey()
	if err != nil {
		return nil, err
	}
	privateKey, ok := rsaPrivateKey.(*rsa.PrivateKey)
	if !ok {
		err = fmt.Errorf("private key type error")
		return nil, err
	}

	return &Signer{
		Address:    crypto.GetAddressFromPublicKey(publicKey),
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

// FromPrivateKey creates a Signer from an existing RSA private key.
//
// This function takes an RSA private key and creates a Signer instance,
// automatically deriving the public key and wallet address. This is useful
// when you already have an RSA private key object from another source.
//
// Parameters:
//   - privateKey: An RSA private key instance
//
// Returns a Signer instance with the provided key and computed address.
//
// Example:
//
//	// Assuming you have an *rsa.PrivateKey from elsewhere
//	signer := FromPrivateKey(existingKey)
//	fmt.Printf("Wallet address: %s\n", signer.Address)
func FromPrivateKey(privateKey *rsa.PrivateKey) *Signer {
	p := &privateKey.PublicKey
	address := crypto.GetAddressFromPublicKey(p)
	return &Signer{
		Address:    address,
		PublicKey:  p,
		PrivateKey: privateKey,
	}
}

// Owner returns the base64url-encoded public key modulus.
//
// This method returns the owner field value as used in Arweave transactions.
// The owner field contains the public key modulus (N) encoded in base64url
// format, which uniquely identifies the transaction signer.
//
// Returns the base64url-encoded public key modulus.
//
// Example:
//
//	owner := signer.Owner()
//	fmt.Printf("Transaction owner: %s\n", owner)
func (s *Signer) Owner() string {
	return crypto.Base64URLEncode(s.PublicKey.N.Bytes())
}

// Generate creates a new Arweave-compatible RSA private key in JWK format.
//
// This function generates a new 4096-bit RSA key pair and returns it
// as JWK-formatted JSON bytes. This is useful for creating new wallet
// files that can be used with Arweave wallet software.
//
// Returns the JWK-formatted private key as bytes, or an error if
// key generation fails.
//
// Example:
//
//	jwkData, err := Generate()
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = os.WriteFile("new-wallet.json", jwkData, 0600)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("New wallet saved to new-wallet.json")
func Generate() ([]byte, error) {
	bitSize := 4096
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}
	jwk, err := gojwk.PrivateKey(key)
	if err != nil {
		return nil, err
	}
	data, err := gojwk.Marshal(jwk)
	if err != nil {
		return nil, err
	}
	return data, nil
}
