package tx

import (
	"errors"

	"github.com/liteseed/goar/crypto"
)

const (
	MAX_TAGS             = 128
	MAX_TAG_KEY_LENGTH   = 1024
	MAX_TAG_VALUE_LENGTH = 3072
)

func NewDataItem(rawData []byte, target string, anchor string, tags []Tag) (*DataItem, error) {
	return &DataItem{
		Target: target,
		Anchor: anchor,
		Tags:   tags,
		Data:   crypto.Base64Encode(rawData),
	}, nil
}

// Decode a DataItem from bytes
func DecodeDataItem(raw []byte) (*DataItem, error) {
	N := len(raw)
	if N < 2 {
		return nil, errors.New("binary too small")
	}

	signatureType, signatureLength, publicKeyLength, err := getSignatureMetadata(raw[:2])
	if err != nil {
		return nil, err
	}

	signatureStart := 2
	signatureEnd := signatureLength + signatureStart
	signature := crypto.Base64Encode(raw[signatureStart:signatureEnd])
	rawId, err := crypto.SHA256(raw[signatureStart:signatureEnd])
	if err != nil {
		return nil, err
	}
	id := crypto.Base64Encode(rawId)
	ownerStart := signatureEnd
	ownerEnd := ownerStart + publicKeyLength
	owner := crypto.Base64Encode(raw[ownerStart:ownerEnd])

	position := ownerEnd
	target, position := getTarget(&raw, position)
	anchor, position := getAnchor(&raw, position)
	tags, position, err := EncodeTags(raw, position)
	if err != nil {
		return nil, err
	}
	data := crypto.Base64Encode(raw[position:])

	return &DataItem{
		ID:            id,
		SignatureType: signatureType,
		Signature:     signature,
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          tags,
		Data:          data,
		Raw:           raw,
	}, nil
}

func VerifyDataItem(dataItem *DataItem) error {
	// Verify ID
	rawSignature, err := crypto.Base64Decode(dataItem.Signature)
	if err != nil {
		return err
	}
	rawId, err := crypto.SHA256(rawSignature)
	if err != nil {
		return err
	}
	id := crypto.Base64Encode(rawId)
	if id != dataItem.ID {
		return errors.New("invalid data item - signature and id don't match")
	}

	chunks, err := GetDataItemChunk(dataItem.Owner, dataItem.Target, dataItem.Anchor, dataItem.Tags, dataItem.Data)
	if err != nil {
		return err
	}

	err = crypto.Verify(chunks, rawSignature, dataItem.Owner)
	if err != nil {
		return err
	}

	// VERIFY TAGS
	if len(dataItem.Tags) > MAX_TAGS {
		return errors.New("invalid data item - tags cannot be more than 128")
	}

	for _, tag := range dataItem.Tags {
		if len([]byte(tag.Name)) == 0 || len([]byte(tag.Name)) > MAX_TAG_KEY_LENGTH {
			return errors.New("invalid data item - tag key too long")
		}
		if len([]byte(tag.Value)) == 0 || len([]byte(tag.Value)) > MAX_TAG_VALUE_LENGTH {
			return errors.New("invalid data item - tag value too long")
		}
	}

	if len([]byte(dataItem.Anchor)) > 32 {
		return errors.New("invalid data item - anchor should be 32 bytes")
	}
	return nil
}

func GetDataItemChunk(owner string, target string, anchor string, tags []Tag, data string) ([]byte, error) {
	rawOwner, err := crypto.Base64Decode(owner)
	if err != nil {
		return nil, err
	}

	rawTarget, err := crypto.Base64Decode(target)
	if err != nil {
		return nil, err
	}
	rawAnchor := []byte(anchor)

	rawTags, err := DecodeTags(tags)
	if err != nil {
		return nil, err
	}
	rawData, err := crypto.Base64Decode(data)
	if err != nil {
		return nil, err
	}

	chunks := []any{
		[]byte("dataitem"),
		[]byte("1"),
		[]byte("1"),
		rawOwner,
		rawTarget,
		rawAnchor,
		rawTags,
		rawData,
	}
	deepHashChunk := crypto.DeepHash(chunks)
	return deepHashChunk[:], nil
}
