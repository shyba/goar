# Goar

A Go library for interacting with the Arweave blockchain. Goar provides a complete toolkit for creating, signing, and managing Arweave transactions, data items, and bundles.

## Features

- **Transaction Management**: Create, sign, and verify Arweave transactions
- **Data Items**: Support for ANS-104 data items and bundles
- **Wallet Operations**: Load wallets from JWK files and manage keys
- **Cryptographic Functions**: Sign/verify operations with RSA and Ed25519 keys
- **Tag Support**: Comprehensive tag handling for metadata
- **Upload Support**: Upload transactions and data to Arweave nodes
- **Merkle Proofs**: Generate and verify Merkle proofs for data integrity

## Install

```bash
go get github.com/liteseed/goar
```

## Quick Start

### Creating and Signing a Transaction

```go
package main

import (
    "fmt"
    "github.com/liteseed/goar/transaction"
    "github.com/liteseed/goar/wallet"
    "github.com/liteseed/goar/tag"
)

func main() {
    // Load wallet from JWK file
    w, err := wallet.LoadFromPath("wallet.json")
    if err != nil {
        panic(err)
    }

    // Create tags
    tags := []tag.Tag{
        {Name: "Content-Type", Value: "text/plain"},
        {Name: "App-Name", Value: "MyApp"},
    }

    // Create transaction
    data := []byte("Hello Arweave!")
    tx := transaction.New(data, "", "0", &tags)

    // Sign transaction
    signer := w.Signer()
    err = tx.Sign(signer)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Transaction ID: %s\n", tx.ID)
}
```

### Uploading Data

```go
package main

import (
    "github.com/liteseed/goar/client"
    "github.com/liteseed/goar/uploader"
)

func main() {
    // Create client
    c := client.New()

    // Create uploader
    up := uploader.New(c)

    // Upload transaction
    resp, err := up.SendTransaction(tx)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Upload response: %+v\n", resp)
}
```

## API Reference

### Core Packages

- **`transaction`**: Create and manage Arweave transactions
- **`wallet`**: Wallet loading and key management
- **`client`**: HTTP client for Arweave nodes
- **`uploader`**: Upload transactions and data items
- **`signer`**: Cryptographic signing operations
- **`tag`**: Tag creation and encoding
- **`crypto`**: Low-level cryptographic functions

### Transaction Package

The transaction package provides the core functionality for creating and managing Arweave transactions.

#### Key Functions

- `New(data []byte, target string, quantity string, tags *[]tag.Tag) *Transaction`: Creates a new transaction
- `(tx *Transaction) Sign(s *signer.Signer) error`: Signs a transaction
- `(tx *Transaction) Verify() error`: Verifies a transaction signature
- `(tx *Transaction) PrepareChunks(data []byte) error`: Prepares data chunks for large transactions

### Wallet Package

Handles wallet operations and key management.

#### Key Functions

- `LoadFromPath(path string) (*Wallet, error)`: Loads a wallet from a JWK file
- `(w *Wallet) Signer() *signer.Signer`: Creates a signer from the wallet

### Client Package

Provides HTTP client functionality for communicating with Arweave nodes.

#### Key Functions

- `New() *Client`: Creates a new client with default settings
- `NewWithURL(url string) *Client`: Creates a client with a custom node URL

## Examples

See the `examples/` directory for more detailed usage examples:

- `send_transaction.go`: Basic transaction creation and upload
- `send_data.go`: Data upload with tags
- `send_bundle.go`: Bundle creation and upload

## Testing

Run the test suite:

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
