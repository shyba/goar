package bundle

import (
	"log"

	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/transaction/data_item"
)

func generateBundleHeader(d *[]data_item.DataItem) (*[]Header, error) {
	var headers []Header

	for _, dataItem := range *d {
		idBytes, err := crypto.Base64URLDecode(dataItem.ID)
		if err != nil {
			return nil, err
		}

		size := len(dataItem.Raw)
		raw := append(idBytes, longTo32ByteArray(size)...)
		headers = append(headers, Header{ID: dataItem.ID, Size: size, Raw: raw})
	}
	return &headers, nil
}

func decodeBundleHeader(data []byte) ([]Header, int) {
	N := byteArrayToLong(data[:32])
	var headers []Header
	for i := 32; i < 32+64*N; i += 64 {
		log.Println(i, i+32, i+32, i+64)
		log.Println(len(data[i:i+32]), len(data[i+32:i+64]))
		size := byteArrayToLong(data[i : i+32])
		id := crypto.Base64URLEncode(data[i+32 : i+64])
		headers = append(headers, Header{ID: id, Size: size, Raw: data[i : i+64]})
	}
	return headers, N
}

func longTo32ByteArray(long int) []byte {
	byteArray := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < len(byteArray); i++ {
		byt := long & 255
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
