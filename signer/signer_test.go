// Package signer tests - verifies key management and signing functionality
package signer

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew verifies that new signers can be created with generated keys
func TestNew(t *testing.T) {
	signer, err := New()
	require.NoError(t, err)
	assert.NotNil(t, signer)
	assert.NotEmpty(t, signer.Address)
	assert.NotNil(t, signer.PrivateKey)
	assert.NotNil(t, signer.PublicKey)
	assert.Equal(t, 4096, signer.PrivateKey.Size()*8) // Should be 4096-bit key
}

// TestFromPath verifies loading signers from JWK files
func TestFromPath(t *testing.T) {
	signer, err := FromPath("../test/signer.json")
	require.NoError(t, err)
	assert.NotNil(t, signer)
	assert.NotEmpty(t, signer.Address)
	assert.NotNil(t, signer.PrivateKey)
	assert.NotNil(t, signer.PublicKey)
}

// TestFromPathInvalidFile verifies error handling for invalid file paths
func TestFromPathInvalidFile(t *testing.T) {
	_, err := FromPath("nonexistent.json")
	assert.Error(t, err)
}

// TestFromJWK verifies creating signers from JWK data
func TestFromJWK(t *testing.T) {
	// Load test JWK data
	data, err := os.ReadFile("../test/signer.json")
	require.NoError(t, err)

	signer, err := FromJWK(data)
	require.NoError(t, err)
	assert.NotNil(t, signer)
	assert.NotEmpty(t, signer.Address)
	assert.NotNil(t, signer.PrivateKey)
	assert.NotNil(t, signer.PublicKey)
}

// TestFromJWKInvalidData verifies error handling for invalid JWK data
func TestFromJWKInvalidData(t *testing.T) {
	invalidData := []byte("{invalid json}")
	_, err := FromJWK(invalidData)
	assert.Error(t, err)
}

// TestFromPrivateKey verifies creating signers from existing private keys
func TestFromPrivateKey(t *testing.T) {
	// First create a signer to get a private key
	originalSigner, err := New()
	require.NoError(t, err)

	// Create new signer from the private key
	newSigner := FromPrivateKey(originalSigner.PrivateKey)
	assert.NotNil(t, newSigner)
	assert.Equal(t, originalSigner.Address, newSigner.Address)
	assert.Equal(t, originalSigner.PrivateKey, newSigner.PrivateKey)
	assert.Equal(t, originalSigner.PublicKey, newSigner.PublicKey)
}

// TestOwner verifies that Owner() returns correct base64url-encoded modulus
func TestOwner(t *testing.T) {
	signer, err := FromPath("../test/signer.json")
	require.NoError(t, err)

	owner := signer.Owner()
	assert.NotEmpty(t, owner)
	// Owner should be base64url encoded, so no padding and URL-safe characters
	assert.NotContains(t, owner, "+")
	assert.NotContains(t, owner, "/")
	assert.NotContains(t, owner, "=")
}

// TestGenerate verifies that Generate() creates valid JWK data
func TestGenerate(t *testing.T) {
	jwkData, err := Generate()
	require.NoError(t, err)
	assert.NotEmpty(t, jwkData)

	// Verify the generated JWK can be used to create a signer
	signer, err := FromJWK(jwkData)
	require.NoError(t, err)
	assert.NotNil(t, signer)
	assert.NotEmpty(t, signer.Address)

	// Verify it's valid JSON
	var jwkMap map[string]interface{}
	err = json.Unmarshal(jwkData, &jwkMap)
	require.NoError(t, err)
	assert.Equal(t, "RSA", jwkMap["kty"])
}

// TestSignerConsistency verifies that the same private key produces the same address
func TestSignerConsistency(t *testing.T) {
	// Load same signer twice
	signer1, err := FromPath("../test/signer.json")
	require.NoError(t, err)

	signer2, err := FromPath("../test/signer.json")
	require.NoError(t, err)

	// Should have identical properties
	assert.Equal(t, signer1.Address, signer2.Address)
	assert.Equal(t, signer1.Owner(), signer2.Owner())
	assert.Equal(t, signer1.PrivateKey.N, signer2.PrivateKey.N)
	assert.Equal(t, signer1.PublicKey.N, signer2.PublicKey.N)
}
