package transaction

import (
	"reflect"
)

func intToByteArray(n int) []byte {
	byteArray := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := len(byteArray) - 1; i >= 0; i-- {
		byt := n % 256
		byteArray[i] = byte(byt)
		n = (n - byt) / 256
	}
	return byteArray
}

func isSlice(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

func byteArrayToInt(b []byte) int {
	value := 0
	for i := 0; i < len(b); i++ {
		value = value*256 + int(b[i])
	}
	return value
}
