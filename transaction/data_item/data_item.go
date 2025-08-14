package data_item

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

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

// NewFromReader Create a new DataItem from a seekable reader for streaming large data
// This avoids loading the entire data into memory. The reader must be seekable (implement io.ReadSeeker)
// for multiple passes during signing and verification.
func NewFromReader(dataReader io.ReadSeeker, dataSize int64, target string, anchor string, tags *[]tag.Tag) *DataItem {
	if tags == nil {
		tags = &[]tag.Tag{}
	}
	return &DataItem{
		Target:     target,
		Anchor:     anchor,
		Tags:       tags,
		DataReader: dataReader,
		DataSize:   dataSize,
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
	// For streaming data, we now handle large files without loading into memory
	// The Raw field construction will handle streaming separately
	var rawData []byte
	var isStreaming = d.DataReader != nil && d.DataSize > 0

	if !isStreaming {
		// Handle small/in-memory data
		if d.Data != "" {
			var err error
			rawData, err = crypto.Base64URLDecode(d.Data)
			if err != nil {
				return err
			}
		} else {
			rawData = []byte{}
		}
	}

	if isStreaming {
		// For streaming data, we'll defer Raw construction
		// The Raw field will be built on-demand when GetRawWithData() is called
		// Here we just create a partial Raw without the data portion
		raw := d.buildHeaderOnly(rawSignature, rawOwner, rawTarget, rawAnchor, rawTags)
		rawID := crypto.SHA256(rawSignature)

		d.Owner = s.Owner()
		d.Signature = crypto.Base64URLEncode(rawSignature)
		d.ID = crypto.Base64URLEncode(rawID)
		d.Raw = raw // Contains only header, data streamed later
		return nil
	}

	// Build Raw for small/in-memory data
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

// buildHeaderOnly creates the header portion of Raw data without the data payload
func (d *DataItem) buildHeaderOnly(rawSignature, rawOwner, rawTarget, rawAnchor, rawTags []byte) []byte {
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

	return raw
}

// GetRawWithData returns the complete raw data including the data payload
// This is needed for bundle creation where the full DataItem binary is required
func (d *DataItem) GetRawWithData() ([]byte, error) {
	if d.DataReader != nil && d.DataSize > 0 {
		// For streaming data, combine header (in Raw) with streamed data
		reader, err := d.getDataReader()
		if err != nil {
			return nil, err
		}

		_, err = reader.Seek(0, io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("failed to seek to beginning: %v", err)
		}

		return d.combineHeaderWithStreamedData(reader)
	}

	return d.Raw, nil
}

// combineHeaderWithStreamedData combines the header (stored in Raw) with streamed data
// WARNING: This method reads the entire data stream into memory for bundle compatibility
// This defeats the purpose of streaming for very large files, but is required for ANS-104 compatibility
func (d *DataItem) combineHeaderWithStreamedData(reader io.ReadSeeker) ([]byte, error) {
	// Allocate buffer for the complete raw data
	totalSize := int64(len(d.Raw)) + d.DataSize
	result := make([]byte, 0, totalSize)

	// Add the header portion (already in d.Raw)
	result = append(result, d.Raw...)

	// Stream the data in chunks to avoid huge single allocations
	const chunkSize = 32768 // 32KB chunks
	buffer := make([]byte, chunkSize)

	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			result = append(result, buffer[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading data stream: %v", err)
		}
	}

	return result, nil
}

// getDataReader returns the provided data reader
func (d *DataItem) getDataReader() (io.ReadSeeker, error) {
	if d.DataReader != nil {
		return d.DataReader, nil
	}

	return nil, fmt.Errorf("no data reader available")
}

// WriteRawFile streams the complete DataItem raw bytes to a file without loading everything into memory.
// For streaming DataItems, this writes the header followed by the data stream directly to the file.
// For non-streaming DataItems, this writes the existing Raw data to the file.
// This method is memory-efficient for large files as it avoids the memory allocation required by GetRawWithData().
func (d *DataItem) WriteRawFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	return d.WriteRawTo(file)
}

// WriteRawTo streams the complete DataItem raw bytes to an io.Writer without loading everything into memory.
// For streaming DataItems, this writes the header followed by the data stream directly to the writer.
// For non-streaming DataItems, this writes the existing Raw data to the writer.
// This method is memory-efficient for large files as it avoids the memory allocation required by GetRawWithData().
func (d *DataItem) WriteRawTo(writer io.Writer) error {
	// Check if this is streaming data
	if d.DataReader != nil && d.DataSize > 0 {
		// Stream the header first (already in d.Raw)
		_, err := writer.Write(d.Raw)
		if err != nil {
			return fmt.Errorf("failed to write header: %v", err)
		}

		// Get the data reader and seek to beginning
		reader, err := d.getDataReader()
		if err != nil {
			return fmt.Errorf("failed to get data reader: %v", err)
		}

		_, err = reader.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("failed to seek to beginning of data: %v", err)
		}

		// Stream the data directly to the writer
		_, err = io.Copy(writer, reader)
		if err != nil {
			return fmt.Errorf("failed to stream data: %v", err)
		}

		return nil
	}

	// For non-streaming data, just write the Raw bytes
	if d.Raw != nil {
		_, err := writer.Write(d.Raw)
		if err != nil {
			return fmt.Errorf("failed to write raw data: %v", err)
		}
	}

	return nil
}

// GetDataSize returns the size of the data payload
func (d *DataItem) GetDataSize() int64 {
	if d.DataSize > 0 {
		return d.DataSize
	}
	// For base64 encoded data, decode to get actual size
	rawData, err := crypto.Base64URLDecode(d.Data)
	if err != nil {
		return 0
	}
	return int64(len(rawData))
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

	// For verification, we need to compute the DeepHash
	// This requires reading the data, which we'll do temporarily
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
