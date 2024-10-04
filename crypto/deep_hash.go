package crypto

import (
	"crypto/sha512"
	"fmt"
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
