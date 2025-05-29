// Package wallet provides high-level functionality for managing Arweave wallets and transactions.
//
// This package combines the signer, client, and transaction functionality into a convenient
// wallet interface that handles the common workflow of creating, signing, and sending
// transactions to the Arweave network.
//
// Example usage:
//
//	// Create wallet from JWK file
//	wallet, err := wallet.FromPath("wallet.json", "https://arweave.net")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Create and send a transaction
//	tx := wallet.CreateTransaction([]byte("Hello Arweave!"), "", "0", nil)
//	signedTx, err := wallet.SignTransaction(tx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = wallet.SendTransaction(signedTx)
//	if err != nil {
//		log.Fatal(err)
//	}
package wallet

import (
	"errors"
	"os"

	"github.com/liteseed/goar/client"
	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/tag"
	"github.com/liteseed/goar/transaction"
	"github.com/liteseed/goar/transaction/bundle"
	"github.com/liteseed/goar/transaction/data_item"
	"github.com/liteseed/goar/uploader"
)

// Wallet represents an Arweave wallet with signing and network capabilities.
//
// A Wallet combines a cryptographic signer for creating transaction signatures
// and a client for communicating with Arweave nodes. It provides a high-level
// interface for common Arweave operations like creating transactions, data items,
// and bundles.
type Wallet struct {
	Client *client.Client // HTTP client for communicating with Arweave nodes
	Signer *signer.Signer // Cryptographic signer for transaction signing
}

// New creates a new wallet with a randomly generated private key.
//
// This function creates a new wallet with:
// - A freshly generated RSA private key for signing
// - A client configured to use the specified gateway URL
//
// Parameters:
//   - gateway: The URL of the Arweave gateway to use (e.g., "https://arweave.net")
//
// Returns a new Wallet instance or an error if key generation fails.
//
// Example:
//
//	wallet, err := New("https://arweave.net")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Created wallet with address: %s\n", wallet.Signer.Address())
func New(gateway string) (w *Wallet, err error) {
	s, err := signer.New()
	if err != nil {
		return nil, err
	}
	return &Wallet{
		Client: client.New(gateway),
		Signer: s,
	}, nil
}

// FromPath creates a wallet from a JWK file on disk.
//
// This function loads a wallet from a JSON Web Key (JWK) file, which is the
// standard format for storing Arweave private keys. The file should contain
// an RSA private key in JWK format.
//
// Parameters:
//   - path: The file system path to the JWK file
//   - gateway: The URL of the Arweave gateway to use
//
// Returns a Wallet instance loaded with the key from the file, or an error
// if the file cannot be read or the key format is invalid.
//
// Example:
//
//	wallet, err := FromPath("./wallet.json", "https://arweave.net")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Loaded wallet with address: %s\n", wallet.Signer.Address())
func FromPath(path string, gateway string) (*Wallet, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return FromJWK(b, gateway)
}

// FromJWK creates a wallet from JWK data in memory.
//
// This function creates a wallet from JSON Web Key (JWK) data provided as
// a byte slice. This is useful when the JWK data is stored in memory or
// retrieved from a source other than the file system.
//
// Parameters:
//   - jwk: The JWK data as bytes (should be valid JSON)
//   - gateway: The URL of the Arweave gateway to use
//
// Returns a Wallet instance with the loaded key, or an error if the JWK
// format is invalid or cannot be parsed.
//
// Example:
//
//	jwkData := []byte(`{"kty":"RSA","n":"...","e":"AQAB",...}`)
//	wallet, err := FromJWK(jwkData, "https://arweave.net")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Created wallet from JWK\n")
func FromJWK(jwk []byte, gateway string) (*Wallet, error) {
	s, err := signer.FromJWK(jwk)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		Client: client.New(gateway),
		Signer: s,
	}, nil
}

// CreateTransaction creates a new Arweave transaction.
//
// This method creates a transaction with the provided data and metadata.
// The transaction is not signed or sent - use SignTransaction and SendTransaction
// for those operations.
//
// Parameters:
//   - data: The data to include in the transaction (can be nil for AR transfers)
//   - target: The target wallet address for AR transfers (empty string for data-only)
//   - quantity: The amount of AR to transfer in Winston units ("0" for data-only)
//   - tags: Optional metadata tags (can be nil)
//
// Returns a new Transaction instance ready for signing.
//
// Example:
//
//	// Data transaction
//	tags := []tag.Tag{{Name: "Content-Type", Value: "text/plain"}}
//	tx := wallet.CreateTransaction([]byte("Hello!"), "", "0", &tags)
//
//	// AR transfer
//	tx := wallet.CreateTransaction(nil, targetAddr, "1000000000000", nil)
func (w *Wallet) CreateTransaction(data []byte, target string, quantity string, tags *[]tag.Tag) *transaction.Transaction {
	return transaction.New(data, target, quantity, tags)
}

// SignTransaction signs a transaction and fills in required network fields.
//
// This method performs several operations:
// 1. Sets the transaction owner to this wallet's public key
// 2. Gets the current transaction anchor from the network
// 3. Calculates the required transaction fee
// 4. Signs the transaction with this wallet's private key
//
// Parameters:
//   - tx: The transaction to sign (created with CreateTransaction)
//
// Returns the signed transaction with all fields populated, or an error if
// any network calls fail or signing fails.
//
// Example:
//
//	tx := wallet.CreateTransaction(data, "", "0", nil)
//	signedTx, err := wallet.SignTransaction(tx)
//	if err != nil {
//		log.Printf("Failed to sign transaction: %v", err)
//		return err
//	}
//	fmt.Printf("Transaction signed with ID: %s\n", signedTx.ID)
func (w *Wallet) SignTransaction(tx *transaction.Transaction) (*transaction.Transaction, error) {
	tx.Owner = w.Signer.Owner()

	anchor, err := w.Client.GetTransactionAnchor()
	if err != nil {
		return nil, err
	}
	tx.LastTx = anchor

	reward, err := w.Client.GetTransactionPrice(len(tx.Data), "")
	if err != nil {
		return nil, err
	}
	tx.Reward = reward

	if err = tx.Sign(w.Signer); err != nil {
		return nil, err
	}
	return tx, nil
}

// SendTransaction sends a signed transaction to the Arweave network.
//
// This method uploads the transaction to the configured Arweave gateway.
// The transaction must be signed before calling this method.
//
// Parameters:
//   - tx: The signed transaction to send
//
// Returns an error if the transaction is not signed or if the upload fails.
//
// Example:
//
//	err := wallet.SendTransaction(signedTx)
//	if err != nil {
//		log.Printf("Failed to send transaction: %v", err)
//		return err
//	}
//	fmt.Printf("Transaction sent successfully: %s\n", signedTx.ID)
func (w *Wallet) SendTransaction(tx *transaction.Transaction) error {
	if tx.ID == "" || tx.Signature == "" {
		return errors.New("transaction not signed")
	}
	tu, err := uploader.New(w.Client, tx)
	if err != nil {
		return err
	}
	if err = tu.PostTransaction(); err != nil {
		return err
	}
	return nil
}

// CreateDataItem creates a new ANS-104 data item.
//
// Data items are a more efficient way to upload data to Arweave when using
// bundling services. They follow the ANS-104 specification and can be
// aggregated into bundles for cost-effective uploads.
//
// Parameters:
//   - data: The data to include in the data item
//   - target: Optional target address for the data item
//   - anchor: Optional anchor value for the data item
//   - tags: Optional metadata tags
//
// Returns a new DataItem instance ready for signing.
//
// Example:
//
//	tags := []tag.Tag{{Name: "Content-Type", Value: "image/jpeg"}}
//	dataItem := wallet.CreateDataItem(imageData, "", "", &tags)
func (w *Wallet) CreateDataItem(data []byte, target string, anchor string, tags *[]tag.Tag) *data_item.DataItem {
	return data_item.New(data, target, anchor, tags)
}

// SignDataItem signs a data item with this wallet's private key.
//
// This method signs the data item using the wallet's signer, making it
// ready for inclusion in a bundle or direct upload.
//
// Parameters:
//   - di: The data item to sign
//
// Returns the signed data item, or an error if signing fails.
//
// Example:
//
//	dataItem := wallet.CreateDataItem(data, "", "", nil)
//	signedItem, err := wallet.SignDataItem(dataItem)
//	if err != nil {
//		log.Printf("Failed to sign data item: %v", err)
//		return err
//	}
//	fmt.Printf("Data item signed with ID: %s\n", signedItem.ID)
func (w *Wallet) SignDataItem(di *data_item.DataItem) (*data_item.DataItem, error) {
	if err := di.Sign(w.Signer); err != nil {
		return nil, err
	}
	return di, nil
}

// CreateBundle creates a new ANS-104 bundle from multiple data items.
//
// Bundles allow multiple data items to be uploaded together in a single
// transaction, reducing costs compared to individual uploads. This follows
// the ANS-104 specification for data bundles.
//
// Parameters:
//   - dataItems: A slice of data items to include in the bundle
//
// Returns a new Bundle instance, or an error if bundle creation fails.
//
// Example:
//
//	dataItems := []data_item.DataItem{item1, item2, item3}
//	bundle, err := wallet.CreateBundle(&dataItems)
//	if err != nil {
//		log.Printf("Failed to create bundle: %v", err)
//		return err
//	}
//	fmt.Printf("Bundle created with %d items\n", len(dataItems))
func (w *Wallet) CreateBundle(dataItems *[]data_item.DataItem) (*bundle.Bundle, error) {
	return bundle.New(dataItems)
}
