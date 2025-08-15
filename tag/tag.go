// Package tag provides functionality for creating and managing Arweave transaction tags.
//
// This package handles tag encoding/decoding operations using Apache Avro format
// as specified in the Arweave protocol and ANS-104 standard. Tags are key-value
// pairs that attach metadata to transactions and data items.
//
// Example usage:
//
//	// Create tags
//	tags := []Tag{
//		{Name: "Content-Type", Value: "application/json"},
//		{Name: "App-Name", Value: "MyApp"},
//	}
//
//	// Serialize for transaction
//	serialized, err := Serialize(&tags)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Convert to base64 format
//	base64Tags := ConvertToBase64(&tags)
package tag

import (
	"encoding/binary"
	"errors"

	"github.com/linkedin/goavro/v2"
	"github.com/liteseed/goar/crypto"
)

// avroTagSchema defines the Apache Avro schema for encoding tags.
// This schema is used by both Arweave transactions and ANS-104 data items.
const avroTagSchema = `
{
	"type": "array",
	"items": {
		"type": "record",
		"name": "Tag",
		"fields": [
			{ "name": "name", "type": "bytes" },
			{ "name": "value", "type": "bytes" }
		]
	}
}`

// fromAvro converts Avro-encoded binary data to human-readable Tags.
//
// This internal function takes Avro-encoded tag data and converts it back
// to a slice of Tag structs. It uses the standard Avro schema defined
// in the Arweave protocol.
//
// Parameters:
//   - data: The Avro-encoded binary data
//
// Returns a slice of Tag structs or an error if decoding fails.
func fromAvro(data []byte) (*[]Tag, error) {
	codec, err := goavro.NewCodec(avroTagSchema)
	if err != nil {
		return nil, err
	}

	avroTags, _, err := codec.NativeFromBinary(data)
	if err != nil {
		return nil, err
	}

	var tags []Tag

	for _, v := range avroTags.([]any) {
		tag := v.(map[string]any)
		tags = append(tags, Tag{Name: string(tag["name"].([]byte)), Value: string(tag["value"].([]byte))})
	}
	return &tags, err
}

// toAvro converts human-readable Tags to Avro-encoded binary data.
//
// This internal function takes a slice of Tag structs and converts them
// to Avro-encoded binary data using the standard schema. The encoded
// data can be included in transactions or data items.
//
// Parameters:
//   - tags: A slice of Tag structs to encode
//
// Returns the Avro-encoded binary data or an error if encoding fails.
func toAvro(tags *[]Tag) ([]byte, error) {
	codec, err := goavro.NewCodec(avroTagSchema)
	if err != nil {
		return nil, err
	}

	var avroTags []map[string]any

	for _, tag := range *tags {
		m := map[string]any{"name": []byte(tag.Name), "value": []byte(tag.Value)}
		avroTags = append(avroTags, m)
	}
	data, err := codec.BinaryFromNative(nil, avroTags)
	if err != nil {
		return nil, err
	}
	return data, err
}

// Serialize converts readable Tag data into Avro-encoded bytes for Arweave transactions.
//
// This function takes a slice of tags and converts them to the binary format
// required by Arweave transactions. The encoding follows the ANS-104 standard
// for tag serialization using Apache Avro.
//
// Parameters:
//   - tags: A slice of Tag structs to serialize
//
// Returns the serialized tag data as bytes, or nil if there are no tags.
// Returns an error if serialization fails.
//
// Learn more: https://github.com/ArweaveTeam/arweave-standards/blob/master/ans/ANS-104.md
//
// Example:
//
//	tags := []Tag{
//		{Name: "Content-Type", Value: "application/json"},
//		{Name: "App-Name", Value: "MyApp"},
//	}
//	serialized, err := Serialize(&tags)
//	if err != nil {
//		log.Fatal(err)
//	}
func Serialize(tags *[]Tag) ([]byte, error) {
	if len(*tags) > 0 {
		data, err := toAvro(tags)
		if err != nil {
			return nil, err
		}

		return data, nil
	}
	return nil, nil
}

// Deserialize converts Avro-encoded byte data from an Arweave transaction into readable Tags.
//
// This function parses tag data from a binary stream, typically from a data item
// or transaction. It handles the binary format specified in ANS-104 which includes
// tag count and byte length headers followed by Avro-encoded tag data.
//
// Parameters:
//   - data: The binary data containing encoded tags
//   - startAt: The byte offset where tag data begins
//
// Returns the parsed tags, the ending offset, and any parsing error.
// The function enforces the ANS-104 limit of maximum 127 tags per item.
//
// Learn more: https://github.com/ArweaveTeam/arweave-standards/blob/master/ans/ANS-104.md
//
// Example:
//
//	tags, endOffset, err := Deserialize(binaryData, 0)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Parsed %d tags, data ends at offset %d\n", len(*tags), endOffset)
func Deserialize(data []byte, startAt int) (*[]Tag, int, error) {
	tags := &[]Tag{}
	tagsEnd := startAt + 8 + 8
	numberOfTags := int(data[startAt])
	numberOfTagBytesStart := startAt + 8
	numberOfTagBytesEnd := numberOfTagBytesStart + 8
	numberOfTagBytes := int(binary.LittleEndian.Uint16(data[numberOfTagBytesStart:numberOfTagBytesEnd]))
	if numberOfTags > 127 {
		return tags, tagsEnd, errors.New("invalid data item - max tags 127")
	}
	if numberOfTags > 0 && numberOfTagBytes > 0 {
		bytesDataStart := numberOfTagBytesEnd
		bytesDataEnd := numberOfTagBytesEnd + numberOfTagBytes
		bytesData := data[bytesDataStart:bytesDataEnd]

		tags, err := fromAvro(bytesData)
		if err != nil {
			return nil, tagsEnd, err
		}
		tagsEnd = bytesDataEnd
		return tags, tagsEnd, nil
	}
	return tags, tagsEnd, nil
}

// Decode converts base64url-encoded Tags to their raw byte representation.
//
// This function takes tags that have base64url-encoded names and values
// (as used in transactions) and decodes them back to their original byte
// form. This is useful for processing tags in their raw format.
//
// Parameters:
//   - tags: A slice of tags with base64url-encoded names and values
//
// Returns a 3D byte slice where each tag is represented as [name_bytes, value_bytes],
// or an error if any tag cannot be decoded.
//
// Example:
//
//	// Assuming tags have base64url-encoded values
//	decodedTags, err := Decode(&tags)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for i, tag := range decodedTags {
//		fmt.Printf("Tag %d: %s = %s\n", i, string(tag[0]), string(tag[1]))
//	}
func Decode(tags *[]Tag) ([][][]byte, error) {
	if len(*tags) == 0 {
		return nil, nil
	}
	data := make([][][]byte, 0)
	for _, tag := range *tags {
		name, err := crypto.Base64URLDecode(tag.Name)
		if err != nil {
			return nil, err
		}
		value, err := crypto.Base64URLDecode(tag.Value)
		if err != nil {
			return nil, err
		}
		data = append(data, [][]byte{name, value})
	}
	return data, nil
}

// ConvertToBase64 encodes all string values of Tags to base64url format.
//
// This function takes tags with plain string names and values and converts
// them to base64url-encoded format as required by Arweave transactions.
// This is typically used when preparing tags for inclusion in a transaction.
//
// Parameters:
//   - tags: A slice of tags with plain string names and values
//
// Returns a new slice of tags with base64url-encoded names and values.
//
// Example:
//
//	tags := []Tag{
//		{Name: "Content-Type", Value: "application/json"},
//		{Name: "App-Name", Value: "MyApp"},
//	}
//	encodedTags := ConvertToBase64(&tags)
//	// encodedTags now contains base64url-encoded names and values
func ConvertToBase64(tags *[]Tag) *[]Tag {
	var result []Tag
	for _, tag := range *tags {
		result = append(result, Tag{Name: crypto.Base64URLEncode([]byte(tag.Name)), Value: crypto.Base64URLEncode([]byte(tag.Value))})
	}
	return &result
}
