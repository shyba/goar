package crypto

import (
	"crypto/sha512"
	"fmt"
	"io"
	"reflect"
)

// DeepHash is a hash algorithm which takes a nested list of values as input
// and produces a 384 bit hash, where a change of any value or the structure
// will affect the hash.
// https://www.arweave.org/yellow-paper.pdf
func DeepHash(data any) [48]byte {
	if typeof(data) == "[]uint8" {
		tag := append([]byte("blob"), []byte(fmt.Sprint(len(data.([]byte))))...)
		tagHashed := sha512.Sum384(tag)
		dataHashed := sha512.Sum384(data.([]byte))
		r := append(tagHashed[:], dataHashed[:]...)
		rHashed := sha512.Sum384(r)
		return rHashed
	} else {
		d := unpackArray(data)
		tag := append([]byte("list"), []byte(fmt.Sprint(len(d)))...)
		return deepHashChunk(d, sha512.Sum384(tag))
	}
}

// DeepHashStream is a streaming version of DeepHash for large data that won't fit in memory.
// It takes a reader and the data size, and computes the same hash as DeepHash would
// for the equivalent []byte, but without loading all data into memory.
func DeepHashStream(reader io.Reader, dataSize int64) ([48]byte, error) {
	// Create the tag hash (same as DeepHash for []byte)
	tag := append([]byte("blob"), []byte(fmt.Sprint(dataSize))...)
	tagHashed := sha512.Sum384(tag)

	// Stream the data through SHA512
	dataHasher := sha512.New384()
	_, err := io.Copy(dataHasher, reader)
	if err != nil {
		return [48]byte{}, err
	}
	dataHashed := dataHasher.Sum(nil)

	// Combine tag and data hashes (same as DeepHash)
	r := append(tagHashed[:], dataHashed[:]...)
	rHashed := sha512.Sum384(r)
	return rHashed, nil
}

// DeepHashMixed computes DeepHash for an array where one element is streamed
// This is specifically for DataItem signing where most fields are small but data can be huge
func DeepHashMixed(chunks [][]byte, streamReader io.Reader, streamSize int64) ([48]byte, error) {
	// Create list tag
	totalItems := len(chunks) + 1 // +1 for the streamed data
	tag := append([]byte("list"), []byte(fmt.Sprint(totalItems))...)
	acc := sha512.Sum384(tag)

	// Process each small chunk
	for _, chunk := range chunks {
		chunkHash := DeepHash(chunk)
		hashPair := append(acc[:], chunkHash[:]...)
		acc = sha512.Sum384(hashPair)
	}

	// Process the streamed data
	streamHash, err := DeepHashStream(streamReader, streamSize)
	if err != nil {
		return [48]byte{}, err
	}
	hashPair := append(acc[:], streamHash[:]...)
	finalHash := sha512.Sum384(hashPair)

	return finalHash, nil
}

func deepHashChunk(data []any, acc [48]byte) [48]byte {
	if len(data) < 1 {
		return acc
	}
	dHash := DeepHash(data[0])
	hashPair := append(acc[:], dHash[:]...)
	newAcc := sha512.Sum384(hashPair)
	return deepHashChunk(data[1:], newAcc)
}

func typeof(v any) string {
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
