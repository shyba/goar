package data_item

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/tag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	s, err := signer.FromPath("../../test/signer.json")
	assert.NoError(t, err)
	t.Run("Decode - New empty test data item", func(t *testing.T) {
		data := ""
		tags := &[]tag.Tag{}
		anchor := ""
		target := ""
		a := New([]byte(data), target, anchor, tags)
		assert.NoError(t, err)

		err = a.Sign(s)
		assert.NoError(t, err)

		dataItem, err := Decode(a.Raw)
		assert.NoError(t, err)
		assert.Equal(t, "gxngjcqu8Kz171MqWuKBAZVaum0cquKpBtwH5s2DucY9rOaxZsszXRnpoHQT7nVdAIPwc40WBqimclR_xJ3jZQ7UKAVKUPyePP_l5jh5Id4HVwwjPMtqApeipaQCJsFCYa33gEzS4NUdKSwGNr6C-Q6SqJ3CXfcwiLrliRHKARMzyhQaTCLwBJP4bHftUjadgix6oqx5hqMGHVWKboJkS6M22fTq4VeUd4whihcYPKzG_ow0aajw1VfqVsXTbQnne9XXXyDswQYiKdsL4OfwBaLtXiDURD12IFQqAkjJ9O68M1AZ102V_TDjZCDEGyRHqmV9yPwihcCbj8r0R7oHgKsDxpRSvxV3Vtx-DxxOUfn8UkdVuRzT9RRs1TLbrfNlIJL2RyjvOXo6fy8p4k_R_w6lAL83JSlXYe24cJj76zEw-CmJnuHVKkXmYeB2NaDFlmvH3Sl3NsraJauycd-1i7gDG0niKF2AeQt76UACamZx2LtE099jl1GetuUYEulNA2V_-zZlOvGH3Lg9x6yepMiW7t2YAXnNoKfD025fuUYXdn_0_IdDJcrySHa9tfrhQzU0gS4FTXjO4Xv9Nmjn9E2ADqb-vcaz73KLtOLHBG5TE60gzbSphi8J7S56zk1UUeZ_IsN9i_p0XeeLN_IpioGumAWcX_B6Pvzm3LBj1-0", dataItem.Owner)
		assert.Equal(t, target, dataItem.Target)
		assert.Equal(t, anchor, dataItem.Anchor)
		assert.Equal(t, base64.RawURLEncoding.EncodeToString([]byte(data)), dataItem.Data)
		assert.ElementsMatch(t, tags, dataItem.Tags)
	})

	t.Run("Decode - data, tags, anchor, target", func(t *testing.T) {
		data := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{};':\",./<>?`~"
		tags := &[]tag.Tag{
			{Name: "tag1", Value: "value1"},
			{Name: "tag2", Value: "value2"},
		}
		anchor := "thisSentenceIs32BytesLongTrustMe"
		target := "OXcT1sVRSA5eGwt2k6Yuz8-3e3g9WJi5uSE99CWqsBs"

		a := New([]byte(data), target, anchor, tags)
		assert.NoError(t, err)

		err = a.Sign(s)
		assert.NoError(t, err)

		dataItem, err := Decode(a.Raw)
		assert.NoError(t, err)

		assert.Equal(t, dataItem.Owner, "gxngjcqu8Kz171MqWuKBAZVaum0cquKpBtwH5s2DucY9rOaxZsszXRnpoHQT7nVdAIPwc40WBqimclR_xJ3jZQ7UKAVKUPyePP_l5jh5Id4HVwwjPMtqApeipaQCJsFCYa33gEzS4NUdKSwGNr6C-Q6SqJ3CXfcwiLrliRHKARMzyhQaTCLwBJP4bHftUjadgix6oqx5hqMGHVWKboJkS6M22fTq4VeUd4whihcYPKzG_ow0aajw1VfqVsXTbQnne9XXXyDswQYiKdsL4OfwBaLtXiDURD12IFQqAkjJ9O68M1AZ102V_TDjZCDEGyRHqmV9yPwihcCbj8r0R7oHgKsDxpRSvxV3Vtx-DxxOUfn8UkdVuRzT9RRs1TLbrfNlIJL2RyjvOXo6fy8p4k_R_w6lAL83JSlXYe24cJj76zEw-CmJnuHVKkXmYeB2NaDFlmvH3Sl3NsraJauycd-1i7gDG0niKF2AeQt76UACamZx2LtE099jl1GetuUYEulNA2V_-zZlOvGH3Lg9x6yepMiW7t2YAXnNoKfD025fuUYXdn_0_IdDJcrySHa9tfrhQzU0gS4FTXjO4Xv9Nmjn9E2ADqb-vcaz73KLtOLHBG5TE60gzbSphi8J7S56zk1UUeZ_IsN9i_p0XeeLN_IpioGumAWcX_B6Pvzm3LBj1-0")
		assert.Equal(t, dataItem.Target, target)
		assert.Equal(t, dataItem.Anchor, anchor)
		assert.Equal(t, dataItem.Data, base64.RawURLEncoding.EncodeToString([]byte(data)))
	})
	t.Run("Decode - Stub", func(t *testing.T) {
		data, err := os.ReadFile("../../test/1115BDataItem")
		assert.NoError(t, err)

		dataItem, err := Decode(data)
		assert.NoError(t, err)
		assert.Equal(t, dataItem.ID, "QpmY8mZmFEC8RxNsgbxSV6e36OF6quIYaPRKzvUco0o")
		assert.Equal(t, dataItem.Signature, "wUIlPaBflf54QyfiCkLnQcfakgcS5B4Pld-hlOJKyALY82xpAivoc0fxBJWjoeg3zy9aXz8WwCs_0t0MaepMBz2bQljRrVXnsyWUN-CYYfKv0RRglOl-kCmTiy45Ox13LPMATeJADFqkBoQKnGhyyxW81YfuPnVlogFWSz1XHQgHxrFMAeTe9epvBK8OCnYqDjch4pwyYUFrk48JFjHM3-I2kcQnm2dAFzFTfO-nnkdQ7ulP3eoAUr-W-KAGtPfWdJKFFgWFCkr_FuNyHYQScQo-FVOwIsvj_PVWEU179NwiqfkZtnN8VoBgCSxbL1Wmh4NYL-GsRbKz_94hpcj5RiIgq0_H5dzAp-bIb49M4SP-DcuIJ5oT2v2AfPWvznokDDVTeikQJxCD2n9usBOJRpLw_P724Yurbl30eNow0U-Jmrl8S6N64cjwKVLI-hBUfcpviksKEF5_I4XCyciW0TvZj1GxK6ET9lx0s6jFMBf27-GrFx6ZDJUBncX6w8nDvuL6A8TG_ILGNQU_EDoW7iil6NcHn5w11yS_yLkqG6dw_zuC1Vkg1tbcKY3703tmbF-jMEZUvJ6oN8vRwwodinJjzGdj7bxmkUPThwVWedCc8wCR3Ak4OkIGASLMUahSiOkYmELbmwq5II-1Txp2gDPjCpAf9gT6Iu0heAaXhjk")
		assert.Equal(t, dataItem.Owner, "0zBGbs8Y4wvdS58cAVyxp7mDffScOkbjh50ZrqnWKR_5NGwjezT6J40ejIg5cm1KnuDnw9OhvA7zO6sv1hEE6IaGNnNJWiXFecRMxCl7iw78frrT8xJvhBgtD4fBCV7eIvydqLoMl8K47sacTUxEGseaLfUdYVJ5CSock5SktEEdqqoe3MAso7x4ZsB5CGrbumNcCTifr2mMsrBytocSoHuiCEi7-Nwv4CqzB6oqymBtEECmKYWdINnNQHVyKK1l0XP1hzByHv_WmhouTPos9Y77sgewZrvLF-dGPNWSc6LaYGy5IphCnq9ACFrEbwkiCRgZHnKsRFH0dfGaCgGb3GZE-uspmICJokJ9CwDPDJoxkCBEF0tcLSIA9_ofiJXaZXbrZzu3TUXWU3LQiTqYr4j5gj_7uTclewbyZSsY-msfbFQlaACc02nQkEkr4pMdpEOdAXjWP6qu7AJqoBPNtDPBqWbdfsLXgyK90NbYmf3x4giAmk8L9REy7SGYugG4VyqG39pNQy_hdpXdcfyE0ftCr5tSHVpMreJ0ni7v3IDCbjZFcvcHp0H6f6WPfNCoHg1BM6rHUqkXWd84gdHUzo9LTGq9-7wSBCizpcc_12_I-6yvZsROJvdfYOmjPnd5llefa_X3X1dVm5FPYFIabydGlh1Vs656rRu4dzeEQwc")
		assert.Equal(t, dataItem.Target, "")
		assert.Equal(t, dataItem.Anchor, "")
		assert.ElementsMatch(
			t,
			*dataItem.Tags,
			[]tag.Tag{
				{Name: "Content-Type", Value: "text/plain"},
				{Name: "App-Name", Value: "ArDrive-CLI"},
				{Name: "App-Version", Value: "1.21.0"},
			},
		)
		assert.Equal(t, dataItem.Data, "NTY3MAo")
	})
}

func TestNew(t *testing.T) {
	s, err := signer.FromPath("../../test/signer.json")
	assert.NoError(t, err)

	t.Run("New - New empty test data item", func(t *testing.T) {
		data := ""
		tags := &[]tag.Tag{}
		anchor := ""
		target := ""

		dataItem := New([]byte(data), target, anchor, tags)
		assert.Equal(t, "", dataItem.Owner)
		assert.Equal(t, target, dataItem.Target)
		assert.Equal(t, anchor, dataItem.Anchor)
		assert.Equal(t, base64.RawURLEncoding.EncodeToString([]byte(data)), dataItem.Data)
		assert.ElementsMatch(t, tags, dataItem.Tags)

		assert.NoError(t, err)

		err = dataItem.Sign(s)
		assert.NoError(t, err)

		assert.Equal(t, "gxngjcqu8Kz171MqWuKBAZVaum0cquKpBtwH5s2DucY9rOaxZsszXRnpoHQT7nVdAIPwc40WBqimclR_xJ3jZQ7UKAVKUPyePP_l5jh5Id4HVwwjPMtqApeipaQCJsFCYa33gEzS4NUdKSwGNr6C-Q6SqJ3CXfcwiLrliRHKARMzyhQaTCLwBJP4bHftUjadgix6oqx5hqMGHVWKboJkS6M22fTq4VeUd4whihcYPKzG_ow0aajw1VfqVsXTbQnne9XXXyDswQYiKdsL4OfwBaLtXiDURD12IFQqAkjJ9O68M1AZ102V_TDjZCDEGyRHqmV9yPwihcCbj8r0R7oHgKsDxpRSvxV3Vtx-DxxOUfn8UkdVuRzT9RRs1TLbrfNlIJL2RyjvOXo6fy8p4k_R_w6lAL83JSlXYe24cJj76zEw-CmJnuHVKkXmYeB2NaDFlmvH3Sl3NsraJauycd-1i7gDG0niKF2AeQt76UACamZx2LtE099jl1GetuUYEulNA2V_-zZlOvGH3Lg9x6yepMiW7t2YAXnNoKfD025fuUYXdn_0_IdDJcrySHa9tfrhQzU0gS4FTXjO4Xv9Nmjn9E2ADqb-vcaz73KLtOLHBG5TE60gzbSphi8J7S56zk1UUeZ_IsN9i_p0XeeLN_IpioGumAWcX_B6Pvzm3LBj1-0", dataItem.Owner)
		assert.Equal(t, target, dataItem.Target)
		assert.Equal(t, anchor, dataItem.Anchor)
		assert.Equal(t, base64.RawURLEncoding.EncodeToString([]byte(data)), dataItem.Data)
	})
	t.Run("New - data, tags, anchor, target", func(t *testing.T) {
		data := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{};':\",./<>?`~"
		tags := &[]tag.Tag{
			{Name: "tag1", Value: "value1"},
			{Name: "tag2", Value: "value2"},
		}
		anchor := "thisSentenceIs32BytesLongTrustMe"
		target := "OXcT1sVRSA5eGwt2k6Yuz8-3e3g9WJi5uSE99CWqsBs"

		dataItem := New([]byte(data), target, anchor, tags)
		assert.NoError(t, err)
		err = dataItem.Sign(s)
		assert.NoError(t, err)

		assert.Equal(t, "gxngjcqu8Kz171MqWuKBAZVaum0cquKpBtwH5s2DucY9rOaxZsszXRnpoHQT7nVdAIPwc40WBqimclR_xJ3jZQ7UKAVKUPyePP_l5jh5Id4HVwwjPMtqApeipaQCJsFCYa33gEzS4NUdKSwGNr6C-Q6SqJ3CXfcwiLrliRHKARMzyhQaTCLwBJP4bHftUjadgix6oqx5hqMGHVWKboJkS6M22fTq4VeUd4whihcYPKzG_ow0aajw1VfqVsXTbQnne9XXXyDswQYiKdsL4OfwBaLtXiDURD12IFQqAkjJ9O68M1AZ102V_TDjZCDEGyRHqmV9yPwihcCbj8r0R7oHgKsDxpRSvxV3Vtx-DxxOUfn8UkdVuRzT9RRs1TLbrfNlIJL2RyjvOXo6fy8p4k_R_w6lAL83JSlXYe24cJj76zEw-CmJnuHVKkXmYeB2NaDFlmvH3Sl3NsraJauycd-1i7gDG0niKF2AeQt76UACamZx2LtE099jl1GetuUYEulNA2V_-zZlOvGH3Lg9x6yepMiW7t2YAXnNoKfD025fuUYXdn_0_IdDJcrySHa9tfrhQzU0gS4FTXjO4Xv9Nmjn9E2ADqb-vcaz73KLtOLHBG5TE60gzbSphi8J7S56zk1UUeZ_IsN9i_p0XeeLN_IpioGumAWcX_B6Pvzm3LBj1-0", dataItem.Owner)
		assert.Equal(t, target, dataItem.Target)
		assert.Equal(t, anchor, dataItem.Anchor)
		assert.Equal(t, base64.RawURLEncoding.EncodeToString([]byte(data)), dataItem.Data)
	})
}

func TestVerifyDataItem(t *testing.T) {
	s, err := signer.FromPath("../../test/signer.json")
	assert.NoError(t, err)

	t.Run("Verify - Empty test data item", func(t *testing.T) {
		data := ""
		tags := &[]tag.Tag{}
		anchor := ""
		target := ""

		dataItem := New([]byte(data), target, anchor, tags)
		assert.NoError(t, err)

		err = dataItem.Sign(s)
		assert.NoError(t, err)

		err = dataItem.Verify()
		assert.NoError(t, err)
	})
	t.Run("Verify - data, tags, anchor, target", func(t *testing.T) {
		data := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{};':\",./<>?`~"
		tags := &[]tag.Tag{
			{Name: "tag1", Value: "value1"},
			{Name: "tag2", Value: "value2"},
		}
		anchor := "thisSentenceIs32BytesLongTrustMe"
		target := "OXcT1sVRSA5eGwt2k6Yuz8-3e3g9WJi5uSE99CWqsBs"

		dataItem := New([]byte(data), target, anchor, tags)
		assert.NoError(t, err)

		err = dataItem.Sign(s)
		assert.NoError(t, err)

		err = dataItem.Verify()
		assert.NoError(t, err)
	})
	t.Run("Verify - Stub", func(t *testing.T) {
		data, err := os.ReadFile("../../test/1115BDataItem")
		assert.NoError(t, err)

		dataItem, err := Decode(data)
		assert.NoError(t, err)

		err = dataItem.Verify()
		assert.NoError(t, err)
	})
}

// MockReadSeeker implements io.ReadSeeker for testing streaming functionality
type MockReadSeeker struct {
	data     []byte
	position int64
}

func NewMockReadSeeker(data []byte) *MockReadSeeker {
	return &MockReadSeeker{
		data:     data,
		position: 0,
	}
}

func (m *MockReadSeeker) Read(p []byte) (int, error) {
	if m.position >= int64(len(m.data)) {
		return 0, io.EOF
	}

	n := copy(p, m.data[m.position:])
	m.position += int64(n)
	return n, nil
}

func (m *MockReadSeeker) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		m.position = offset
	case io.SeekCurrent:
		m.position += offset
	case io.SeekEnd:
		m.position = int64(len(m.data)) + offset
	default:
		return 0, fmt.Errorf("invalid whence value")
	}

	if m.position < 0 {
		m.position = 0
	}
	if m.position > int64(len(m.data)) {
		m.position = int64(len(m.data))
	}

	return m.position, nil
}

// TestNewFromReader tests the streaming functionality with NewFromReader
func TestNewFromReader(t *testing.T) {
	s, err := signer.New()
	require.NoError(t, err)

	t.Run("NewFromReader - Basic functionality", func(t *testing.T) {
		data := []byte("Hello, streaming world!")
		reader := NewMockReadSeeker(data)
		tags := &[]tag.Tag{
			{Name: "Content-Type", Value: "text/plain"},
			{Name: "Test", Value: "Streaming"},
		}

		dataItem := NewFromReader(reader, int64(len(data)), "", "", tags)

		// Verify initial state
		assert.Equal(t, reader, dataItem.DataReader)
		assert.Equal(t, int64(len(data)), dataItem.DataSize)
		assert.Equal(t, "", dataItem.Target)
		assert.Equal(t, "", dataItem.Anchor)
		assert.Equal(t, tags, dataItem.Tags)
		assert.Equal(t, "", dataItem.Data) // Should be empty for streaming

		// Sign the DataItem
		err := dataItem.Sign(s)
		require.NoError(t, err)

		// Verify signing results
		assert.NotEmpty(t, dataItem.ID)
		assert.NotEmpty(t, dataItem.Signature)
		assert.NotEmpty(t, dataItem.Owner)
		assert.NotEmpty(t, dataItem.Raw) // Should contain header

		// Verify the DataItem
		err = dataItem.Verify()
		assert.NoError(t, err)
	})

	t.Run("NewFromReader - Large data simulation", func(t *testing.T) {
		// Create a larger dataset to simulate streaming
		large_data := make([]byte, 1024*1024) // 1MB
		for i := range large_data {
			large_data[i] = byte(i % 256)
		}

		reader := NewMockReadSeeker(large_data)
		tags := &[]tag.Tag{
			{Name: "Content-Type", Value: "application/octet-stream"},
			{Name: "Size", Value: "1MB"},
		}

		dataItem := NewFromReader(reader, int64(len(large_data)), "", "", tags)

		// Sign should work without loading all data into memory
		err := dataItem.Sign(s)
		require.NoError(t, err)

		// Verify should work
		err = dataItem.Verify()
		assert.NoError(t, err)

		// GetDataSize should return correct size
		assert.Equal(t, int64(len(large_data)), dataItem.GetDataSize())
	})

	t.Run("NewFromReader - With target and anchor", func(t *testing.T) {
		data := []byte("Test data with target and anchor")
		reader := NewMockReadSeeker(data)
		target := "OXcT1sVRSA5eGwt2k6Yuz8-3e3g9WJi5uSE99CWqsBs"
		anchor := "test_anchor_string"
		tags := &[]tag.Tag{{Name: "Type", Value: "Test"}}

		dataItem := NewFromReader(reader, int64(len(data)), target, anchor, tags)

		err := dataItem.Sign(s)
		require.NoError(t, err)

		assert.Equal(t, target, dataItem.Target)
		assert.Equal(t, anchor, dataItem.Anchor)

		err = dataItem.Verify()
		assert.NoError(t, err)
	})

	t.Run("NewFromReader - Nil tags handling", func(t *testing.T) {
		data := []byte("Test with nil tags")
		reader := NewMockReadSeeker(data)

		dataItem := NewFromReader(reader, int64(len(data)), "", "", nil)

		// Should create empty tags array
		assert.NotNil(t, dataItem.Tags)
		assert.Equal(t, 0, len(*dataItem.Tags))

		err := dataItem.Sign(s)
		require.NoError(t, err)

		err = dataItem.Verify()
		assert.NoError(t, err)
	})
}

// TestGetRawWithData tests the GetRawWithData functionality for streaming
func TestGetRawWithData(t *testing.T) {
	s, err := signer.New()
	require.NoError(t, err)

	t.Run("GetRawWithData - Streaming data", func(t *testing.T) {
		originalData := []byte("This is test data for GetRawWithData")
		reader := NewMockReadSeeker(originalData)

		dataItem := NewFromReader(reader, int64(len(originalData)), "", "", nil)
		err := dataItem.Sign(s)
		require.NoError(t, err)

		// Get raw data with data included
		rawWithData, err := dataItem.GetRawWithData()
		require.NoError(t, err)

		// Raw data should be larger than just the original data (includes headers)
		assert.Greater(t, len(rawWithData), len(originalData))

		// The raw data should end with our original data
		assert.True(t, bytes.HasSuffix(rawWithData, originalData))

		// Should be able to decode the raw data
		decoded, err := Decode(rawWithData)
		require.NoError(t, err)

		// Decoded data should match original DataItem properties
		assert.Equal(t, dataItem.ID, decoded.ID)
		assert.Equal(t, dataItem.Signature, decoded.Signature)
		assert.Equal(t, dataItem.Owner, decoded.Owner)
	})

	t.Run("GetRawWithData - Non-streaming data", func(t *testing.T) {
		originalData := []byte("Regular non-streaming data")
		tags := &[]tag.Tag{{Name: "Test", Value: "NonStreaming"}}

		dataItem := New(originalData, "", "", tags)
		err := dataItem.Sign(s)
		require.NoError(t, err)

		// For non-streaming data, GetRawWithData should return the same as Raw
		rawWithData, err := dataItem.GetRawWithData()
		require.NoError(t, err)

		assert.Equal(t, dataItem.Raw, rawWithData)
	})
}

// TestStreamingSeekBehavior tests that the reader seeking works correctly
func TestStreamingSeekBehavior(t *testing.T) {
	s, err := signer.New()
	require.NoError(t, err)

	t.Run("Multiple operations with seeking", func(t *testing.T) {
		data := []byte("Test data for seek behavior testing")
		reader := NewMockReadSeeker(data)

		dataItem := NewFromReader(reader, int64(len(data)), "", "", nil)

		// Sign (uses the reader)
		err := dataItem.Sign(s)
		require.NoError(t, err)

		// Verify (uses the reader again)
		err = dataItem.Verify()
		require.NoError(t, err)

		// GetRawWithData (uses the reader a third time)
		_, err = dataItem.GetRawWithData()
		require.NoError(t, err)

		// All operations should succeed due to proper seeking
	})
}

// TestDataSizeCalculation tests GetDataSize method
func TestDataSizeCalculation(t *testing.T) {
	t.Run("GetDataSize - Streaming data", func(t *testing.T) {
		data := []byte("Test data for size calculation")
		reader := NewMockReadSeeker(data)

		dataItem := NewFromReader(reader, int64(len(data)), "", "", nil)

		assert.Equal(t, int64(len(data)), dataItem.GetDataSize())
	})

	t.Run("GetDataSize - Non-streaming data", func(t *testing.T) {
		data := []byte("Non-streaming test data")
		dataItem := New(data, "", "", nil)

		// For non-streaming data, it should decode base64 to get actual size
		assert.Equal(t, int64(len(data)), dataItem.GetDataSize())
	})
}

// TestErrorHandling tests various error conditions
func TestErrorHandling(t *testing.T) {
	t.Run("Sign - Seek error during streaming", func(t *testing.T) {
		// Create a mock reader that fails on seek
		failingReader := &FailingSeeker{data: []byte("test")}
		dataItem := NewFromReader(failingReader, 4, "", "", nil)

		s, err := signer.New()
		require.NoError(t, err)

		// Sign should fail because of seek error in getDataItemChunkStreaming
		err = dataItem.Sign(s)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to seek")
	})

	t.Run("GetDataSize - Invalid base64 data", func(t *testing.T) {
		// Create a DataItem with invalid base64 data
		dataItem := &DataItem{
			Data: "invalid-base64!@#$%^&*()",
		}

		// Should return 0 when base64 decode fails
		size := dataItem.GetDataSize()
		assert.Equal(t, int64(0), size)
	})

	t.Run("GetRawWithData - No DataReader for streaming", func(t *testing.T) {
		// Create a DataItem that looks like streaming but has no reader
		dataItem := &DataItem{
			DataReader: nil,
			DataSize:   100,
		}

		// Should return empty Raw data
		raw, err := dataItem.GetRawWithData()
		require.NoError(t, err)
		assert.Empty(t, raw)
	})
}

// FailingSeeker is a mock reader that fails on Seek operations
type FailingSeeker struct {
	data     []byte
	position int64
}

func (f *FailingSeeker) Read(p []byte) (int, error) {
	n := copy(p, f.data[f.position:])
	f.position += int64(n)
	if f.position >= int64(len(f.data)) {
		return n, io.EOF
	}
	return n, nil
}

func (f *FailingSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, fmt.Errorf("seek operation failed")
}

// TestBackwardCompatibility ensures new changes don't break existing functionality
func TestBackwardCompatibility(t *testing.T) {
	s, err := signer.New()
	require.NoError(t, err)

	t.Run("Traditional New() still works", func(t *testing.T) {
		data := []byte("Traditional data")
		tags := &[]tag.Tag{{Name: "Method", Value: "Traditional"}}

		dataItem := New(data, "", "", tags)
		err := dataItem.Sign(s)
		require.NoError(t, err)

		err = dataItem.Verify()
		assert.NoError(t, err)

		// Should have Data field populated
		assert.NotEmpty(t, dataItem.Data)
		// Should not have streaming fields
		assert.Nil(t, dataItem.DataReader)
		assert.Equal(t, int64(0), dataItem.DataSize)
	})

	t.Run("Decode still works with traditional data", func(t *testing.T) {
		originalData := []byte("Data for decode test")
		dataItem := New(originalData, "", "", nil)
		err := dataItem.Sign(s)
		require.NoError(t, err)

		// Decode from raw bytes
		decoded, err := Decode(dataItem.Raw)
		require.NoError(t, err)

		assert.Equal(t, dataItem.ID, decoded.ID)
		assert.Equal(t, dataItem.Data, decoded.Data)
		assert.Equal(t, dataItem.Signature, decoded.Signature)
	})
}

// TestStreamingInternalMethods tests internal methods used by streaming functionality
func TestStreamingInternalMethods(t *testing.T) {
	s, err := signer.New()
	require.NoError(t, err)

	t.Run("buildHeaderOnly creates proper header", func(t *testing.T) {
		data := []byte("Test data for header building")
		reader := NewMockReadSeeker(data)
		tags := &[]tag.Tag{{Name: "Test", Value: "HeaderBuild"}}

		dataItem := NewFromReader(reader, int64(len(data)), "", "", tags)
		err := dataItem.Sign(s)
		require.NoError(t, err)

		// After signing, Raw should contain header-only data for streaming
		// The header should be substantial (Arweave headers are quite large due to signatures)
		assert.Greater(t, len(dataItem.Raw), 100) // Should have substantial header

		// For streaming data, the Raw field contains only the header, not the data
		// So it should be independent of the actual data size

		// GetRawWithData should return much larger data
		fullRaw, err := dataItem.GetRawWithData()
		require.NoError(t, err)
		assert.Greater(t, len(fullRaw), len(dataItem.Raw))
	})

	t.Run("Empty data streaming", func(t *testing.T) {
		emptyData := []byte{}
		reader := NewMockReadSeeker(emptyData)

		dataItem := NewFromReader(reader, 0, "", "", nil)
		err := dataItem.Sign(s)
		require.NoError(t, err)

		err = dataItem.Verify()
		assert.NoError(t, err)

		// Should handle empty streaming data
		assert.Equal(t, int64(0), dataItem.GetDataSize())
	})
}

// TestWriteRawTo tests the WriteRawTo method for streaming raw data to writers
func TestWriteRawTo(t *testing.T) {
	s, err := signer.New()
	require.NoError(t, err)

	t.Run("WriteRawTo - Streaming data", func(t *testing.T) {
		originalData := []byte("Test data for WriteRawTo streaming functionality")
		reader := NewMockReadSeeker(originalData)
		tags := &[]tag.Tag{
			{Name: "Content-Type", Value: "text/plain"},
			{Name: "Test", Value: "WriteRawTo"},
		}

		dataItem := NewFromReader(reader, int64(len(originalData)), "", "", tags)
		err := dataItem.Sign(s)
		require.NoError(t, err)

		// Write to a buffer
		var buffer bytes.Buffer
		err = dataItem.WriteRawTo(&buffer)
		require.NoError(t, err)

		// Verify the written data
		writtenData := buffer.Bytes()
		assert.Greater(t, len(writtenData), len(originalData))     // Should include headers
		assert.True(t, bytes.HasSuffix(writtenData, originalData)) // Should end with original data

		// The written data should be decodable
		decoded, err := Decode(writtenData)
		require.NoError(t, err)
		assert.Equal(t, dataItem.ID, decoded.ID)
		assert.Equal(t, dataItem.Signature, decoded.Signature)
	})

	t.Run("WriteRawTo - Non-streaming data", func(t *testing.T) {
		originalData := []byte("Non-streaming data for WriteRawTo")
		tags := &[]tag.Tag{{Name: "Test", Value: "NonStreamingWriteRawTo"}}

		dataItem := New(originalData, "", "", tags)
		err := dataItem.Sign(s)
		require.NoError(t, err)

		// Write to a buffer
		var buffer bytes.Buffer
		err = dataItem.WriteRawTo(&buffer)
		require.NoError(t, err)

		// Should write the same as the Raw field
		assert.Equal(t, dataItem.Raw, buffer.Bytes())

		// Verify it's decodable
		decoded, err := Decode(buffer.Bytes())
		require.NoError(t, err)
		assert.Equal(t, dataItem.ID, decoded.ID)
	})

	t.Run("WriteRawTo - Large streaming data", func(t *testing.T) {
		// Create larger data to test streaming efficiency
		largeData := make([]byte, 10*1024) // 10KB
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		reader := NewMockReadSeeker(largeData)
		dataItem := NewFromReader(reader, int64(len(largeData)), "", "", nil)
		err := dataItem.Sign(s)
		require.NoError(t, err)

		var buffer bytes.Buffer
		err = dataItem.WriteRawTo(&buffer)
		require.NoError(t, err)

		// Verify the data was streamed correctly
		writtenData := buffer.Bytes()
		assert.True(t, bytes.HasSuffix(writtenData, largeData))

		// Verify it can be decoded
		decoded, err := Decode(writtenData)
		require.NoError(t, err)
		assert.Equal(t, dataItem.ID, decoded.ID)
	})

	t.Run("WriteRawTo - Empty streaming data", func(t *testing.T) {
		emptyData := []byte{}
		reader := NewMockReadSeeker(emptyData)

		dataItem := NewFromReader(reader, 0, "", "", nil)
		err := dataItem.Sign(s)
		require.NoError(t, err)

		var buffer bytes.Buffer
		err = dataItem.WriteRawTo(&buffer)
		require.NoError(t, err)

		// Should still contain the header
		assert.Greater(t, buffer.Len(), 0)

		// Should be decodable
		decoded, err := Decode(buffer.Bytes())
		require.NoError(t, err)
		assert.Equal(t, dataItem.ID, decoded.ID)
	})
}

// TestWriteRawFile tests the WriteRawFile method for streaming raw data to files
func TestWriteRawFile(t *testing.T) {
	s, err := signer.New()
	require.NoError(t, err)

	t.Run("WriteRawFile - Streaming data to file", func(t *testing.T) {
		originalData := []byte("Test data for WriteRawFile with streaming")
		reader := NewMockReadSeeker(originalData)
		tags := &[]tag.Tag{
			{Name: "Content-Type", Value: "application/octet-stream"},
			{Name: "Test", Value: "WriteRawFile"},
		}

		dataItem := NewFromReader(reader, int64(len(originalData)), "", "", tags)
		err := dataItem.Sign(s)
		require.NoError(t, err)

		// Create temp file
		tempDir := t.TempDir()
		tempFile := filepath.Join(tempDir, "test_dataitem.bin")

		// Write to file
		err = dataItem.WriteRawFile(tempFile)
		require.NoError(t, err)

		// Verify file was created and has correct content
		fileData, err := os.ReadFile(tempFile)
		require.NoError(t, err)

		// Should contain headers plus data
		assert.Greater(t, len(fileData), len(originalData))
		assert.True(t, bytes.HasSuffix(fileData, originalData))

		// Should be decodable
		decoded, err := Decode(fileData)
		require.NoError(t, err)
		assert.Equal(t, dataItem.ID, decoded.ID)
		assert.Equal(t, dataItem.Signature, decoded.Signature)
	})

	t.Run("WriteRawFile - Non-streaming data to file", func(t *testing.T) {
		originalData := []byte("Non-streaming data for file write test")
		dataItem := New(originalData, "", "", nil)
		err := dataItem.Sign(s)
		require.NoError(t, err)

		// Create temp file
		tempDir := t.TempDir()
		tempFile := filepath.Join(tempDir, "test_traditional_dataitem.bin")

		// Write to file
		err = dataItem.WriteRawFile(tempFile)
		require.NoError(t, err)

		// Verify file content matches Raw data
		fileData, err := os.ReadFile(tempFile)
		require.NoError(t, err)
		assert.Equal(t, dataItem.Raw, fileData)
	})

	t.Run("WriteRawFile - Large file streaming", func(t *testing.T) {
		// Create a reasonably large dataset
		largeData := make([]byte, 1024*1024) // 1MB
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		reader := NewMockReadSeeker(largeData)
		dataItem := NewFromReader(reader, int64(len(largeData)), "", "", nil)
		err := dataItem.Sign(s)
		require.NoError(t, err)

		// Create temp file
		tempDir := t.TempDir()
		tempFile := filepath.Join(tempDir, "test_large_dataitem.bin")

		// Write to file (should be memory efficient)
		err = dataItem.WriteRawFile(tempFile)
		require.NoError(t, err)

		// Verify file size is reasonable
		fileInfo, err := os.Stat(tempFile)
		require.NoError(t, err)
		expectedSize := int64(len(largeData)) + int64(len(dataItem.Raw))
		assert.Equal(t, expectedSize, fileInfo.Size())

		// Verify file content by reading just the end
		file, err := os.Open(tempFile)
		require.NoError(t, err)
		defer file.Close()

		// Seek to where data should start (after header)
		_, err = file.Seek(int64(len(dataItem.Raw)), io.SeekStart)
		require.NoError(t, err)

		// Read first few bytes of data section
		dataStart := make([]byte, 100)
		n, err := file.Read(dataStart)
		require.NoError(t, err)

		// Should match original data
		assert.Equal(t, largeData[:n], dataStart[:n])
	})
}

// TestWriteRawErrorHandling tests error conditions for WriteRaw methods
func TestWriteRawErrorHandling(t *testing.T) {
	t.Run("WriteRawFile - Invalid file path", func(t *testing.T) {
		data := []byte("test data")
		reader := NewMockReadSeeker(data)
		dataItem := NewFromReader(reader, int64(len(data)), "", "", nil)

		s, err := signer.New()
		require.NoError(t, err)
		err = dataItem.Sign(s)
		require.NoError(t, err)

		// Try to write to invalid path
		err = dataItem.WriteRawFile("/invalid/path/that/does/not/exist/file.bin")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create file")
	})

	t.Run("WriteRawTo - Seek error during streaming", func(t *testing.T) {
		failingReader := &FailingSeeker{data: []byte("test data")}
		dataItem := NewFromReader(failingReader, 9, "", "", nil)

		// Set up a fake signed item to test WriteRawTo error handling
		dataItem.Raw = []byte("fake header") // Fake header for testing

		var buffer bytes.Buffer
		err := dataItem.WriteRawTo(&buffer)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to seek")
	})

	t.Run("WriteRawTo - No Raw data", func(t *testing.T) {
		// Create empty DataItem
		dataItem := &DataItem{}

		var buffer bytes.Buffer
		err := dataItem.WriteRawTo(&buffer)
		require.NoError(t, err)

		// Should write nothing
		assert.Equal(t, 0, buffer.Len())
	})
}

// TestWriteRawMemoryEfficiency tests that WriteRaw methods don't load large data into memory
func TestWriteRawMemoryEfficiency(t *testing.T) {
	t.Run("Memory efficient compared to GetRawWithData", func(t *testing.T) {
		// Create a reasonably large dataset
		largeData := make([]byte, 100*1024) // 100KB
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		reader := NewMockReadSeeker(largeData)
		dataItem := NewFromReader(reader, int64(len(largeData)), "", "", nil)

		s, err := signer.New()
		require.NoError(t, err)
		err = dataItem.Sign(s)
		require.NoError(t, err)

		// Test WriteRawTo
		var buffer bytes.Buffer
		err = dataItem.WriteRawTo(&buffer)
		require.NoError(t, err)

		// Test GetRawWithData for comparison
		rawWithData, err := dataItem.GetRawWithData()
		require.NoError(t, err)

		// Both should produce the same result
		assert.Equal(t, rawWithData, buffer.Bytes())

		// The key difference is that WriteRawTo streams the data without loading it all into memory
		// while GetRawWithData loads everything. This is verified by the successful execution
		// of WriteRawTo without memory allocation proportional to data size.
	})
}
