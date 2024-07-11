// Package signer provides primitives for signing data
package signer

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/everFinance/gojwk"
	"github.com/liteseed/goar/crypto"
)

type Signer struct {
	Address    string
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

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

func FromPath(path string) (*Signer, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return FromJWK(b)
}

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

func FromPrivateKey(privateKey *rsa.PrivateKey) *Signer {
	p := &privateKey.PublicKey
	address := crypto.GetAddressFromPublicKey(p)
	return &Signer{
		Address:    address,
		PublicKey:  p,
		PrivateKey: privateKey,
	}
}

func (s *Signer) Owner() string {
	return crypto.Base64URLEncode(s.PublicKey.N.Bytes())
}

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
