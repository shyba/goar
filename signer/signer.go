package signer

import (
	"crypto/rsa"
	"crypto/sha256"
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

func FromPath(path string) (*Signer, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return New(b)
}

func New(b []byte) (*Signer, error) {
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
		err = fmt.Errorf("pubKey type error")
		return nil, err
	}

	rsaPrivateKey, err := key.DecodePrivateKey()
	if err != nil {
		return nil, err
	}
	privateKey, ok := rsaPrivateKey.(*rsa.PrivateKey)
	if !ok {
		err = fmt.Errorf("prvKey type error")
		return nil, err
	}
	addr := sha256.Sum256(publicKey.N.Bytes())
	return &Signer{
		Address:    crypto.Base64Encode(addr[:]),
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

func FromPrivateKey(privateKey *rsa.PrivateKey) *Signer {
	pub := &privateKey.PublicKey
	addr := sha256.Sum256(pub.N.Bytes())
	return &Signer{
		Address:    crypto.Base64Encode(addr[:]),
		PublicKey:  pub,
		PrivateKey: privateKey,
	}
}

func (s *Signer) Owner() string {
	return crypto.Base64Encode(s.PublicKey.N.Bytes())
}
