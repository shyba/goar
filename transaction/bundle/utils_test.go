package bundle

import (
	"log"
	"os"
	"testing"

	"github.com/liteseed/goar/transaction/data_item"
	"github.com/stretchr/testify/assert"
)

func TestDecodeBundleHeader(t *testing.T) {
	data, err := os.ReadFile("../../test/signed-bundle")
	if err != nil {
		log.Fatal(err)
	}
	headers, N := decodeBundleHeader(data)
	assert.Equal(t, N, 1)
	assert.Equal(t, 1063, headers[0].Size)
	assert.Equal(t, "Rh71hbi1SjdweiLSgJQioZ4VLlsnN0PM1Zzkzo_S3w0", headers[0].ID)
}

func TestGenerateBundleHeader(t *testing.T) {
	data, err := os.ReadFile("../../test/1115BDataItem")
	assert.NoError(t, err)

	dataItem, err := data_item.Decode(data)
	assert.NoError(t, err)
	headers, err := generateBundleHeader(&[]data_item.DataItem{*dataItem})

	assert.NoError(t, err)
	assert.Equal(t, 1115, (*headers)[0].Size)
	assert.Equal(t, "QpmY8mZmFEC8RxNsgbxSV6e36OF6quIYaPRKzvUco0o", (*headers)[0].ID)
}

func TestByteArrayToLong(t *testing.T) {
	v0Int := 281474976710655
	v0Bytes := []byte{255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	res0 := byteArrayToLong(v0Bytes)

	assert.Equal(t, v0Int, res0)

	v1Int := 34566888345923
	v1Bytes := []byte{67, 209, 25, 59, 112, 31, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	res1 := byteArrayToLong(v1Bytes)

	assert.Equal(t, v1Int, res1)
}

func TestLongToByteArray(t *testing.T) {
	v0Int := 281474976710655
	v0Bytes := []byte{255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	res0 := longTo32ByteArray(v0Int)

	assert.Equal(t, v0Bytes, res0)

	v1Int := 34566888345923
	v1Bytes := []byte{67, 209, 25, 59, 112, 31, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	res1 := longTo32ByteArray(v1Int)

	assert.Equal(t, v1Bytes, res1)
}
