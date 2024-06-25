package bundle

import (
	"encoding/base64"
	"encoding/binary"

	"github.com/liteseed/goar/transaction/data_item"
)

func generateBundleHeader(d *[]data_item.DataItem) (*[]BundleHeader, error) {
	headers := []BundleHeader{}

	for _, dataItem := range *d {
		idBytes, err := base64.RawURLEncoding.DecodeString(dataItem.ID)
		if err != nil {
			return nil, err
		}

		id := int(binary.LittleEndian.Uint16(idBytes))
		size := len(dataItem.Raw)
		raw := make([]byte, 0)
		binary.LittleEndian.AppendUint16(raw, uint16(size))
		binary.LittleEndian.AppendUint16(raw, uint16(id))
		headers = append(headers, BundleHeader{ID: id, Size: size, Raw: raw})
	}
	return &headers, nil
}

func decodeBundleHeader(data *[]byte) (*[]BundleHeader, int) {
	N := int(binary.LittleEndian.Uint32((*data)[:32]))
	headers := []BundleHeader{}
	for i := 32; i < 32+64*N; i += 64 {
		size := int(binary.LittleEndian.Uint16((*data)[i : i+32]))
		id := int(binary.LittleEndian.Uint16((*data)[i+32 : i+64]))
		headers = append(headers, BundleHeader{ID: id, Size: size})
	}
	return &headers, N
}
