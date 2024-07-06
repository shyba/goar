package transaction

import (
	"encoding/binary"
	"reflect"
)

func encodeUint(x uint64) []byte {
	buf := make([]byte, 32)

	// byteOffset by 24
	// JS implementation assumes a 32 byte length Uint8Array
	binary.BigEndian.PutUint64(buf[24:], x)
	return buf
}
func isSlice(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}


func byteArrayToLong(b []byte) int {
	value := 0
	for i := len(b) - 1; i >= 0; i-- {
		value = value*256 + int(b[i])
	}
	return value
}
