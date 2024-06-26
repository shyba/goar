package bundle

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	data, err := os.ReadFile("../../test/signed-bundle")
	assert.NoError(t, err)

	b, err := Decode(data)
	assert.NoError(t, err)

	assert.NotNil(t, b)

}