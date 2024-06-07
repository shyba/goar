package types

import "encoding/binary"

func encodeUint(x uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, x)
	return buf
}
