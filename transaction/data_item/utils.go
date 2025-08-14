package data_item

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"

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

// This function assembles DataItem data in a format specified by ANS-104 and hashes it using DeepHash
func (d *DataItem) getDataItemChunk() ([]byte, error) {
	rawOwner, err := crypto.Base64URLDecode(d.Owner)
	if err != nil {
		return nil, err
	}

	rawTarget, err := crypto.Base64URLDecode(d.Target)
	if err != nil {
		return nil, err
	}
	rawAnchor := []byte(d.Anchor)

	rawTags, err := tag.Serialize(d.Tags)
	if err != nil {
		return nil, err
	}

	// Use streaming approach for large data
	if d.DataReader != nil && d.DataSize > 0 {
		return d.getDataItemChunkStreaming(rawOwner, rawTarget, rawAnchor, rawTags)
	}

	// Handle in-memory data
	var rawData []byte
	if d.Data != "" {
		rawData, err = crypto.Base64URLDecode(d.Data)
		if err != nil {
			return nil, err
		}
	} else {
		rawData = []byte{}
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

// getDataItemChunkStreaming computes the DataItem hash using streaming for large data
func (d *DataItem) getDataItemChunkStreaming(rawOwner, rawTarget, rawAnchor, rawTags []byte) ([]byte, error) {
	// Prepare the chunks that come before the data
	chunks := [][]byte{
		[]byte("dataitem"),
		[]byte("1"),
		[]byte("1"),
		rawOwner,
		rawTarget,
		rawAnchor,
		rawTags,
	}

	// Get a reader for the data
	reader, err := d.getDataReader()
	if err != nil {
		return nil, err
	}
	// Note: We don't close the reader - it's the caller's responsibility

	// Seek to beginning to ensure we read from start
	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("failed to seek to beginning: %v", err)
	}

	// Use streaming DeepHash for the mixed case
	deepHashChunk, err := crypto.DeepHashMixed(chunks, reader, d.DataSize)
	if err != nil {
		return nil, err
	}

	return deepHashChunk[:], nil
}
