package data_item

import (
	"encoding/binary"
	"errors"

	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/tag"
)

const (
	MAX_TAGS             = 128
	MAX_TAG_KEY_LENGTH   = 1024
	MAX_TAG_VALUE_LENGTH = 3072
)

func New(rawData []byte, target string, anchor string, tags []tag.Tag) *DataItem {
	return &DataItem{
		Target: target,
		Anchor: anchor,
		Tags:   tags,
		Data:   crypto.Base64Encode(rawData),
	}
}

// Decode a DataItem from bytes
func Decode(raw []byte) (*DataItem, error) {
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
	tags, position, err := tag.Encode(raw, position)
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

func Verify(dataItem *DataItem) error {
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

	publicKey, err := crypto.GetPublicKeyFromOwner(dataItem.Owner)
	if err != nil {
		return err
	}
	err = crypto.Verify(chunks, rawSignature, publicKey)
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

func GetDataItemChunk(owner string, target string, anchor string, tags []tag.Tag, data string) ([]byte, error) {
	rawOwner, err := crypto.Base64Decode(owner)
	if err != nil {
		return nil, err
	}

	rawTarget, err := crypto.Base64Decode(target)
	if err != nil {
		return nil, err
	}
	rawAnchor := []byte(anchor)

	rawTags, err := tag.Decode(tags)
	if err != nil {
		return nil, err
	}
	rawData, err := crypto.Base64Decode(data)
	if err != nil {
		return nil, err
	}

	chunks := [][]byte{
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

func (dataItem *DataItem) Sign(s *signer.Signer) error {
	deepHashChunk, err := GetDataItemChunk(s.Owner(), dataItem.Target, dataItem.Anchor, dataItem.Tags, dataItem.Data)
	if err != nil {
		return err
	}

	rawSignature, err := crypto.Sign(deepHashChunk, s.PrivateKey)
	if err != nil {
		return err
	}

	rawOwner, err := crypto.Base64Decode(s.Owner())
	if err != nil {
		return err
	}

	rawTarget, err := crypto.Base64Decode(dataItem.Target)
	if err != nil {
		return err
	}
	rawAnchor := []byte(dataItem.Anchor)

	rawTags, err := tag.Decode(dataItem.Tags)
	if err != nil {
		return err
	}
	rawData, err := crypto.Base64Decode(dataItem.Data)
	if err != nil {
		return err
	}

	raw := make([]byte, 0)
	raw = binary.LittleEndian.AppendUint16(raw, uint16(1))
	raw = append(raw, rawSignature...)
	raw = append(raw, rawOwner...)

	if dataItem.Target == "" {
		raw = append(raw, 0)
	} else {
		raw = append(raw, 1)
	}
	raw = append(raw, rawTarget...)

	if dataItem.Anchor == "" {
		raw = append(raw, 0)
	} else {
		raw = append(raw, 1)
	}
	raw = append(raw, rawAnchor...)
	numberOfTags := make([]byte, 8)
	binary.LittleEndian.PutUint16(numberOfTags, uint16(len(dataItem.Tags)))
	raw = append(raw, numberOfTags...)

	tagsLength := make([]byte, 8)
	binary.LittleEndian.PutUint16(tagsLength, uint16(len(rawTags)))
	raw = append(raw, tagsLength...)
	raw = append(raw, rawTags...)
	raw = append(raw, rawData...)
	rawID, err := crypto.SHA256(rawSignature)
	if err != nil {
		return err
	}

	dataItem.Owner = s.Owner()
	dataItem.Signature = crypto.Base64Encode(rawSignature)
	dataItem.ID = crypto.Base64Encode(rawID)
	dataItem.Raw = raw
	return nil
}
