package crypto

import (
	"crypto/rsa"
	"math/big"
)

func GetAddressFromOwner(owner string) (string, error) {
	publicKey, err := GetPublicKeyFromOwner(owner)
	if err != nil {
		return "", err
	}
	address, err := GetAddressFromPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	return address, nil
}

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

func GetAddressFromPublicKey(p *rsa.PublicKey) (string, error) {
	address, err := SHA256(p.N.Bytes())
	if err != nil {
		return "", err
	}
	return Base64URLEncode(address), nil
}
