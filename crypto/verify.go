package crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"math/big"
)

func getPublicKeyFromOwner(owner string) (*rsa.PublicKey, error) {
	data, err := Base64Decode(owner)
	if err != nil {
		return nil, err
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(data),
		E: 65537, //"AQAB"
	}, nil
}

func Verify(data []byte, signature []byte, owner string) error {
	hashed := sha256.Sum256(data)

	publicKey, err := getPublicKeyFromOwner(owner)
	if err != nil {
		return err
	}
	return rsa.VerifyPSS(publicKey, crypto.SHA256, hashed[:], signature, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	})
}
