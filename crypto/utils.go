package crypto

import (
	"crypto/rsa"
	"math/big"
	"reflect"
)

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func unpackArray(s any) []any {
	v := reflect.ValueOf(s)
	r := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		r[i] = v.Index(i).Interface()
	}
	return r
}

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
	data, err := Base64Decode(owner)
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
	return Base64Encode(address), nil
}
