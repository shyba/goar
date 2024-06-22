package data_item

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/tag"
)

const (
	Arweave  = 1
	ED25519  = 2
	Ethereum = 3
	Solana   = 4
)

type SignatureMeta struct {
	SignatureLength int
	PublicKeyLength int
	Name            string
}

var SignatureConfig = map[int]SignatureMeta{
	Arweave: {
		SignatureLength: 512,
		PublicKeyLength: 512,
		Name:            "arweave",
	},
	ED25519: {
		SignatureLength: 64,
		PublicKeyLength: 32,
		Name:            "ed25519",
	},
	Ethereum: {
		SignatureLength: 65,
		PublicKeyLength: 65,
		Name:            "ethereum",
	},
	Solana: {
		SignatureLength: 64,
		PublicKeyLength: 32,
		Name:            "solana",
	},
}

func getTarget(data *[]byte, position int) (string, int) {
	target := ""
	if (*data)[position] == 1 {
		target = base64.RawURLEncoding.EncodeToString((*data)[position+1 : position+1+32])
		position += 32
	}
	return target, position + 1
}

func getAnchor(data *[]byte, position int) (string, int) {
	anchor := ""
	if (*data)[position] == 1 {
		anchor = string((*data)[position+1 : position+1+32])
		position += 32
	}
	return anchor, position + 1
}
func getSignatureMetadata(data []byte) (SignatureType int, SignatureLength int, PublicKeyLength int, err error) {
	SignatureType = int(binary.LittleEndian.Uint16(data))
	signatureMeta, ok := SignatureConfig[SignatureType]
	if !ok {
		return -1, -1, -1, fmt.Errorf("unsupported signature type:%d", SignatureType)
	}
	SignatureLength = signatureMeta.SignatureLength
	PublicKeyLength = signatureMeta.PublicKeyLength
	err = nil
	return
}

// This function assembles DataItem data in a format specified by ANS-104 and hashes is it using DeepHash
func (d *DataItem) getDataItemChunk() ([]byte, error) {
	rawOwner, err := crypto.Base64Decode(d.Owner)
	if err != nil {
		return nil, err
	}

	rawTarget, err := crypto.Base64Decode(d.Target)
	if err != nil {
		return nil, err
	}
	rawAnchor := []byte(d.Anchor)

	rawTags, err := tag.Serialize(d.Tags)
	if err != nil {
		return nil, err
	}
	rawData, err := crypto.Base64Decode(d.Data)
	if err != nil {
		return nil, err
	}

	chunks := [][]byte{
		[]byte("dataitem"),
		[]byte("1"),
		[]byte("1"),
		rawOwner,
		rawTarget,
		rawAnchor,
		rawTags,
		rawData,
	}
	deepHashChunk := crypto.DeepHash(chunks)
	return deepHashChunk[:], nil
}
