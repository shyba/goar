package tx

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"reflect"
)

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

func generateBundleHeader(d *[]DataItem) (*[]BundleHeader, error) {
	headers := []BundleHeader{}

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
		headers = append(headers, BundleHeader{id: id, size: size, raw: raw})
	}
	return &headers, nil
}

func decodeBundleHeader(data *[]byte) (*[]BundleHeader, int) {
	N := int(binary.LittleEndian.Uint32((*data)[:32]))
	headers := []BundleHeader{}
	for i := 32; i < 32+64*N; i += 64 {
		size := int(binary.LittleEndian.Uint16((*data)[i : i+32]))
		id := int(binary.LittleEndian.Uint16((*data)[i+32 : i+64]))
		headers = append(headers, BundleHeader{id: id, size: size})
	}
	return &headers, N
}

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
