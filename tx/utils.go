package tx

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/liteseed/goar/types"
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

func generateBundleHeader(d *[]types.DataItem) (*[]types.BundleHeader, error) {
	headers := []types.BundleHeader{}

	for _, dataItem := range *d {
		idBytes, err := base64.RawURLEncoding.DecodeString(dataItem.ID)
		if err != nil {
			return nil, err
		}

		id := int(binary.LittleEndian.Uint16(idBytes))
		size := len(dataItem.Raw)
		raw := make([]byte, 64)
		binary.LittleEndian.PutUint16(raw, uint16(size))
		binary.LittleEndian.AppendUint16(raw, uint16(id))
		headers = append(headers, types.BundleHeader{ID: id, Size: size, Raw: raw})
	}
	return &headers, nil
}

func decodeBundleHeader(data *[]byte) (*[]types.BundleHeader, int) {
	N := int(binary.LittleEndian.Uint32((*data)[:32]))
	headers := []types.BundleHeader{}
	for i := 32; i < 32+64*N; i += 64 {
		size := int(binary.LittleEndian.Uint16((*data)[i : i+32]))
		id := int(binary.LittleEndian.Uint16((*data)[i+32 : i+64]))
		headers = append(headers, types.BundleHeader{ID: id, Size: size})
	}
	return &headers, N
}
