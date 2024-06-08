package crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
)

func Verify(data []byte, signature []byte, owner string) error {
	hashed := sha256.Sum256(data)

	publicKey, err := GetPublicKeyFromOwner(owner)
	if err != nil {
		return err
	}
	return rsa.VerifyPSS(publicKey, crypto.SHA256, hashed[:], signature, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	})
}
