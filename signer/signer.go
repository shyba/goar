package signer

import (
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

	address, err := crypto.GetAddressFromPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return &Signer{
		Address:    address,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

func FromPrivateKey(privateKey *rsa.PrivateKey) (*Signer, error) {
	p := &privateKey.PublicKey
	address, err := crypto.GetAddressFromPublicKey(p)
	if err != nil {
		return nil, err
	}
	return &Signer{
		Address:    address,
		PublicKey:  p,
		PrivateKey: privateKey,
	}, nil
}

func (s *Signer) Owner() string {
	return crypto.Base64Encode(s.PublicKey.N.Bytes())
}
