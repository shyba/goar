package bundle

import (
	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/transaction/data_item"
)

func generateBundleHeader(d *[]data_item.DataItem) (*[]BundleHeader, error) {
	headers := []BundleHeader{}

	for _, dataItem := range *d {
		idBytes, err := crypto.Base64Decode(dataItem.ID)
		if err != nil {
			return nil, err
		}

		id := byteArrayToLong(idBytes)
		size := len(dataItem.Raw)
		headers = append(headers, BundleHeader{ID: id, Size: size})
	}
	return &headers, nil
}

func decodeBundleHeader(data []byte) (*[]BundleHeader, int) {
	N := byteArrayToLong(data[:32])
	headers := []BundleHeader{}
	for i := 32; i < 32+64*N; i += 64 {
		size := byteArrayToLong(data[i : i+32])
		id := byteArrayToLong(data[i+32 : i+64])
		headers = append(headers, BundleHeader{ID: id, Size: size})
	}
	return &headers, N
}

func longTo32ByteArray(long int) []byte {
	byteArray := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < len(byteArray); i++ {
		byt := long & 0xff
		byteArray[i] = byte(byt)
		long = (long - byt) / 256
	}
	return byteArray
}
func byteArrayToLong(b []byte) int {
	value := 0
	for i := len(b) - 1; i >= 0; i-- {
		value = value*256 + int(b[i])
	}
	return value
}
