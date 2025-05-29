package data_item

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/tag"
	"github.com/stretchr/testify/assert"
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
