package bundle

import (
	"log"
	"os"
	"testing"

	"github.com/liteseed/goar/transaction/data_item"
	"github.com/stretchr/testify/assert"
)

func TestDecodeBundleHeader(t *testing.T) {
	data, err := os.ReadFile("../../test//bundleHeader")
	if err != nil {
		log.Fatal(err)
	}
	headers, N := decodeBundleHeader(data)
	assert.Equal(t, N, 1)
	assert.Equal(t, 1115, (*headers)[0].Size)
	assert.Equal(t, 4617428110304385346, (*headers)[0].ID)
}

func TestGenerateBundleHeader(t *testing.T) {
	data, err := os.ReadFile("../../test/1115BDataItem")
	assert.NoError(t, err)

	dataItem, err := data_item.Decode(data)
	assert.NoError(t, err)
	headers, err := generateBundleHeader(&[]data_item.DataItem{*dataItem})

	assert.NoError(t, err)
	assert.Equal(t, 1115, (*headers)[0].Size)
	assert.Equal(t, 4617428110304385346, (*headers)[0].ID)
}
