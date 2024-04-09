package crypto

import (
	"crypto/sha512"
	"fmt"
)

func DeepHash(data any) [48]byte {
	if typeof(data) == "[]uint8" {
		tag := append([]byte("blob"), []byte(fmt.Sprintf("%d", len(data.([]byte))))...)
		tagHashed := sha512.Sum384(tag)
		dataHashed := sha512.Sum384(data.([]byte))
		r := append(tagHashed[:], dataHashed[:]...)
		rHashed := sha512.Sum384(r)
		return rHashed
	} else {
		_data := unpackArray(data)
		tag := append([]byte("list"), []byte(fmt.Sprintf("%d", len(_data)))...)
		tagHashed := sha512.Sum384(tag)
		return deepHashChunk(_data, tagHashed)
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
