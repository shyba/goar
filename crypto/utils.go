package crypto

import (
	"crypto/rsa"
	"math/big"
)

// GetAddressFromOwner - Convert the 512 byte owner to the Arweave public address
func GetAddressFromOwner(owner string) (string, error) {
	publicKey, err := GetPublicKeyFromOwner(owner)
	if err != nil {
		return "", err
	}
	address := GetAddressFromPublicKey(publicKey)
	return address, nil
}

// GetPublicKeyFromOwner - Convert the 512 byte owner from the Public
func GetPublicKeyFromOwner(owner string) (*rsa.PublicKey, error) {
	data, err := Base64URLDecode(owner)
	if err != nil {
		return nil, err
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(data),
		E: 65537, //"AQAB"
	}, nil
}

// GetAddressFromPublicKey - Convert the RSA Public Key to the Arweave public address
func GetAddressFromPublicKey(p *rsa.PublicKey) string {
	return Base64URLEncode(SHA256(p.N.Bytes()))
}
