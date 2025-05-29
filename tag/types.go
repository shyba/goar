package tag

// Tag represents a key-value pair of metadata for Arweave transactions.
//
// Tags are used to attach metadata to transactions and data items in the
// Arweave protocol. They are encoded using Apache Avro format and can be
// used to store information like content type, application name, or any
// other custom metadata.
//
// Example:
//
//	tag := Tag{
//		Name:  "Content-Type",
//		Value: "application/json",
//	}
type Tag struct {
	Name  string `json:"name"`  // The tag name/key
	Value string `json:"value"` // The tag value
}
