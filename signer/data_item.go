package signer

import (
	"encoding/binary"

	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/tx"
	"github.com/liteseed/goar/types"
)

func (s *Signer) SignDataItem(dataItem *types.DataItem) error {
	deepHashChunk, err := tx.GetDataItemChunk(s.Owner(), dataItem.Target, dataItem.Anchor, dataItem.Tags, dataItem.Data)
	if err != nil {
		return err
	}

	rawSignature, err := s.Sign(deepHashChunk)
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

	rawTags, err := tx.DecodeTags(dataItem.Tags)
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
