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

// New Create a new DataItem
// Learn more: https://github.com/ArweaveTeam/arweave-standards/blob/master/ans/ANS-104.md
func New(rawData []byte, target string, anchor string, tags *[]tag.Tag) *DataItem {
	if tags == nil {
		tags = &[]tag.Tag{}
	}
	return &DataItem{
		Target: target,
		Anchor: anchor,
		Tags:   tags,
		Data:   crypto.Base64URLEncode(rawData),
	}
}

// Decode a [DataItem] from bytes
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

	signature := crypto.Base64URLEncode(raw[signatureStart:signatureEnd])
	rawId := crypto.SHA256(raw[signatureStart:signatureEnd])
	id := crypto.Base64URLEncode(rawId)
	ownerStart := signatureEnd
	ownerEnd := ownerStart + publicKeyLength
	owner := crypto.Base64URLEncode(raw[ownerStart:ownerEnd])

	position := ownerEnd
	target, position := getTarget(&raw, position)
	anchor, position := getAnchor(&raw, position)
	tags, position, err := tag.Deserialize(raw, position)
	if err != nil {
		return nil, err
	}
	data := crypto.Base64URLEncode(raw[position:])

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

func (d *DataItem) Sign(s *signer.Signer) error {
	d.Owner = s.Owner()
	deepHashChunk, err := d.getDataItemChunk()
	if err != nil {
		return err
	}

	rawSignature, err := crypto.Sign(deepHashChunk, s.PrivateKey)
	if err != nil {
		return err
	}

	rawOwner, err := crypto.Base64URLDecode(s.Owner())
	if err != nil {
		return err
	}

	rawTarget, err := crypto.Base64URLDecode(d.Target)
	if err != nil {
		return err
	}
	rawAnchor := []byte(d.Anchor)

	rawTags, err := tag.Serialize(d.Tags)
	if err != nil {
		return err
	}
	rawData, err := crypto.Base64URLDecode(d.Data)
	if err != nil {
		return err
	}

	raw := make([]byte, 0)
	raw = binary.LittleEndian.AppendUint16(raw, uint16(1))
	raw = append(raw, rawSignature...)
	raw = append(raw, rawOwner...)

	if d.Target == "" {
		raw = append(raw, 0)
	} else {
		raw = append(raw, 1)
	}
	raw = append(raw, rawTarget...)

	if d.Anchor == "" {
		raw = append(raw, 0)
	} else {
		raw = append(raw, 1)
	}
	raw = append(raw, rawAnchor...)
	numberOfTags := make([]byte, 8)
	binary.LittleEndian.PutUint16(numberOfTags, uint16(len(*d.Tags)))
	raw = append(raw, numberOfTags...)

	tagsLength := make([]byte, 8)
	binary.LittleEndian.PutUint16(tagsLength, uint16(len(rawTags)))
	raw = append(raw, tagsLength...)
	raw = append(raw, rawTags...)
	raw = append(raw, rawData...)
	rawID := crypto.SHA256(rawSignature)

	d.Owner = s.Owner()
	d.Signature = crypto.Base64URLEncode(rawSignature)
	d.ID = crypto.Base64URLEncode(rawID)
	d.Raw = raw
	return nil
}

func (d *DataItem) Verify() error {
	// Verify ID
	rawSignature, err := crypto.Base64URLDecode(d.Signature)
	if err != nil {
		return err
	}

	rawId := crypto.SHA256(rawSignature)
	id := crypto.Base64URLEncode(rawId)

	if id != d.ID {
		return errors.New("invalid data item - signature and id don't match")
	}

	chunks, err := d.getDataItemChunk()
	if err != nil {
		return err
	}

	publicKey, err := crypto.GetPublicKeyFromOwner(d.Owner)
	if err != nil {
		return err
	}
	err = crypto.Verify(chunks, rawSignature, publicKey)
	if err != nil {
		return err
	}

	// VERIFY TAGS
	if len(*d.Tags) > MAX_TAGS {
		return errors.New("invalid data item - tags cannot be more than 128")
	}

	for _, t := range *d.Tags {
		if len([]byte(t.Name)) == 0 || len([]byte(t.Name)) > MAX_TAG_KEY_LENGTH {
			return errors.New("invalid data item - tag key too long")
		}
		if len([]byte(t.Value)) == 0 || len([]byte(t.Value)) > MAX_TAG_VALUE_LENGTH {
			return errors.New("invalid data item - tag value too long")
		}
	}

	if len([]byte(d.Anchor)) > 32 {
		return errors.New("invalid data item - anchor should be 32 bytes")
	}
	return nil
}
