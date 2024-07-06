package bundle

import (
	"errors"

	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/transaction/data_item"
)

// Create a data bundle from a group of data items
// Learn more: // Learn more: https://github.com/ArweaveTeam/arweave-standards/blob/master/ans/ANS-104.md
func New(ds *[]data_item.DataItem) (*Bundle, error) {
	b := &Bundle{}

	headers, err := generateBundleHeader(ds)
	if err != nil {
		return nil, err
	}

	b.Headers = *headers
	b.Items = *ds
	N := len(*ds)

	var headersBytes []byte
	var dataItemsBytes []byte

	for i := 0; i < N; i++ {
		h := (*headers)[i]
		sizeBytes := longTo32ByteArray(h.Size)
		idBytes, err := crypto.Base64URLDecode(h.ID)
		if err != nil {
			return nil, err
		}
		headersBytes = append(headersBytes, sizeBytes...)
		headersBytes = append(headersBytes, idBytes...)

		d := (*ds)[i]
		dataItemsBytes = append(dataItemsBytes, d.Raw...)
	}

	raw := make([]byte, 0)
	raw = append(raw, longTo32ByteArray(N)...)
	raw = append(raw, headersBytes...)
	raw = append(raw, dataItemsBytes...)
	b.Raw = raw
	return b, nil
}

// Decode raw bytes into a Bundle
func Decode(data []byte) (*Bundle, error) {
	// length must more than 32
	if len(data) < 32 {
		return nil, errors.New("binary length must more than 32")
	}
	headers, N := decodeBundleHeader(data)
	bundle := &Bundle{
		Items: make([]data_item.DataItem, N),
		Raw:   data,
	}
	bundleStart := 32 + 64*N
	for i := 0; i < N; i++ {
		header := headers[i]
		bundleEnd := bundleStart + header.Size
		dataItem, err := data_item.Decode(data[bundleStart:bundleEnd])
		if err != nil {
			return nil, err
		}
		bundle.Items[i] = *dataItem
		bundleStart = bundleEnd
	}
	return bundle, nil
}

func Verify(data []byte) (bool, error) {
	// length must more than 32
	if len(data) < 32 {
		return false, errors.New("binary length must more than 32")
	}
	headers, N := decodeBundleHeader(data)
	dataItemSize := 0
	for i := 0; i < N; i++ {
		dataItemSize += headers[i].Size
	}
	return len(data) == dataItemSize+32+64*N, nil
}
