package transaction

import (
	"reflect"
)

// intToByteArray converts an integer to a 32-byte big-endian byte array.
//
// This function is used internally for converting integer values to bytes
// in the format expected by Arweave's hashing and encoding algorithms.
// The result is always exactly 32 bytes long, with leading zeros if necessary.
//
// Parameters:
//   - n: The integer to convert (must be non-negative)
//
// Returns a 32-byte array representing the integer in big-endian format.
//
// Example:
//   - intToByteArray(256) returns [0,0,...,0,1,0] (32 bytes total)
//   - intToByteArray(0) returns [0,0,...,0,0] (32 bytes total)
func intToByteArray(n int) []byte {
	byteArray := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := len(byteArray) - 1; i >= 0; i-- {
		byt := n % 256
		byteArray[i] = byte(byt)
		n = (n - byt) / 256
	}
	return byteArray
}

// isSlice checks if the provided value is a slice type.
//
// This utility function uses reflection to determine if a given interface{}
// value represents a slice. It's used internally for type checking during
// data processing operations.
//
// Parameters:
//   - v: The value to check (can be any type)
//
// Returns true if the value is a slice, false otherwise.
//
// Example:
//   - isSlice([]int{1,2,3}) returns true
//   - isSlice("hello") returns false
//   - isSlice(123) returns false
func isSlice(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

// byteArrayToInt converts a byte array to an integer using big-endian interpretation.
//
// This function is the inverse of intToByteArray and is used to convert
// byte arrays back to integer values. It interprets the bytes in big-endian
// format (most significant byte first).
//
// Parameters:
//   - b: The byte array to convert (can be any length)
//
// Returns the integer value represented by the byte array.
//
// Example:
//   - byteArrayToInt([]byte{0,0,1,0}) returns 256
//   - byteArrayToInt([]byte{0,0,0,0}) returns 0
//   - byteArrayToInt([]byte{1,2,3}) returns 66051 (1*256² + 2*256 + 3)
func byteArrayToInt(b []byte) int {
	value := 0
	for i := 0; i < len(b); i++ {
		value = value*256 + int(b[i])
	}
	return value
}
