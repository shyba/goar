package bundle

import (
	"errors"

	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/data_item"
)

func Decode(data []byte) (*Bundle, error) {
	// length must more than 32
	if len(data) < 32 {
		return nil, errors.New("binary length must more than 32")
	}
	headers, N := decodeBundleHeader(&data)
	bundle := &Bundle{
		Items:   make([]data_item.DataItem, N),
		RawData: crypto.Base64Encode(data),
	}
	bundleStart := 32 + 64*N
	for i := 0; i < N; i++ {
		header := (*headers)[i]
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

func New(dataItems *[]data_item.DataItem) (*Bundle, error) {
	bundle := &Bundle{}

	headers, err := generateBundleHeader(dataItems)
	if err != nil {
		return nil, err
	}

	bundle.Headers = *headers
	bundle.Items = *dataItems
	N := len(*dataItems)

	var sizeBytes []byte
	var headersBytes []byte
	var dataItemsBytes []byte

	for i := 0; i < N; i++ {
		headersBytes = append(headersBytes, (*headers)[i].Raw...)
		dataItemsBytes = append(dataItemsBytes, (*headers)[i].Raw...)
	}

	bundle.RawData = crypto.Base64Encode(append(sizeBytes, append(headersBytes, dataItemsBytes...)...))
	return bundle, nil
}

func Verify(data []byte) (bool, error) {
	// length must more than 32
	if len(data) < 32 {
		return false, errors.New("binary length must more than 32")
	}
	headers, N := decodeBundleHeader(&data)
	dataItemSize := 0
	for i := 0; i < N; i++ {
		dataItemSize += (*headers)[i].Size
	}
	return len(data) == dataItemSize+32+64*N, nil
}
