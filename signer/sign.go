package signer

import "github.com/liteseed/goar/crypto"

func (s *Signer) Sign(data []byte) ([]byte, error) {
	return crypto.Sign(data, s.PrivateKey)
}
