// Package client provides functionality to interact with the Arweave HTTP API.
//
// This package implements the complete set of HTTP endpoints documented at:
// https://docs.arweave.org/developers/server/http-api
//
// The client supports all major operations including:
// - Transaction retrieval and submission
// - Wallet balance and transaction history queries
// - Block and network information retrieval
// - Data uploading and chunk management
//
// Example usage:
//
//	client := client.New("https://arweave.net")
//
//	// Get transaction by ID
//	tx, err := client.GetTransactionByID("txid...")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Submit a new transaction
//	status, err := client.SubmitTransaction(myTransaction)
//	if err != nil {
//		log.Fatal(err)
//	}
package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/liteseed/goar/transaction"
)

// Client represents an HTTP client for communicating with Arweave nodes.
//
// The client maintains connection settings and provides methods for all
// Arweave HTTP API endpoints. It includes automatic timeout handling
// and error management for network operations.
type Client struct {
	Client  *http.Client // HTTP client with configured timeout
	Gateway string       // Base URL of the Arweave gateway
}

// New creates a new Arweave client with default settings.
//
// The client is configured with a 10-second timeout for all HTTP requests.
// This timeout applies to individual requests, not the overall operation time.
//
// Parameters:
//   - gateway: The base URL of the Arweave gateway (e.g., "https://arweave.net")
//
// Returns a configured Client instance ready for use.
//
// Example:
//
//	client := New("https://arweave.net")
//	// or use a custom gateway
//	client := New("https://my-arweave-node.com")
func New(gateway string) *Client {
	return &Client{
		Client:  &http.Client{Timeout: time.Second * 10},
		Gateway: gateway,
	}
}

// GetTransactionByID retrieves a complete transaction by its ID.
//
// This method fetches the full transaction data including all fields
// like data, tags, signature, and metadata. The transaction ID is the
// SHA256 hash of the transaction signature.
//
// Parameters:
//   - id: The transaction ID (base64url-encoded hash)
//
// Returns the complete Transaction struct or an error if the transaction
// is not found or cannot be retrieved.
//
// Example:
//
//	tx, err := client.GetTransactionByID("ABC123...")
//	if err != nil {
//		log.Printf("Transaction not found: %v", err)
//		return
//	}
//	fmt.Printf("Transaction from: %s\n", tx.Owner)
func (c *Client) GetTransactionByID(id string) (*transaction.Transaction, error) {
	body, err := c.get(fmt.Sprintf("tx/%s", id))
	if err != nil {
		return nil, err
	}
	t := &transaction.Transaction{}
	err = json.Unmarshal(body, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// GetTransactionStatus retrieves the confirmation status of a transaction.
//
// This method returns information about whether a transaction has been
// confirmed by the network, including block information if confirmed.
// Transactions typically take 2-10 minutes to be confirmed.
//
// Parameters:
//   - id: The transaction ID to check status for
//
// Returns TransactionStatus with confirmation details or an error if
// the transaction cannot be found.
//
// Example:
//
//	status, err := client.GetTransactionStatus("ABC123...")
//	if err != nil {
//		log.Printf("Failed to get status: %v", err)
//		return
//	}
//	if status.Confirmed {
//		fmt.Printf("Transaction confirmed in block %s\n", status.BlockIndepHash)
//	}
func (c *Client) GetTransactionStatus(id string) (*TransactionStatus, error) {
	body, err := c.get(fmt.Sprintf("tx/%s/status", id))
	if err != nil {
		return nil, err
	}

	t := &TransactionStatus{}
	err = json.Unmarshal(body, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// GetTransactionField retrieves a specific field from a transaction.
//
// This method allows fetching individual transaction fields without
// downloading the entire transaction. Useful for large transactions
// when only specific metadata is needed.
//
// Common fields include: "data", "tags", "target", "quantity", "signature"
//
// Parameters:
//   - id: The transaction ID
//   - field: The name of the field to retrieve
//
// Returns the field value as a string, or an error if the transaction
// or field is not found.
//
// Example:
//
//	tags, err := client.GetTransactionField("ABC123...", "tags")
//	if err != nil {
//		log.Printf("Failed to get tags: %v", err)
//		return
//	}
//	fmt.Printf("Transaction tags: %s\n", tags)
func (c *Client) GetTransactionField(id string, field string) (string, error) {
	body, err := c.get(fmt.Sprintf("tx/%s/%s", id, field))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// GetTransactionData retrieves the raw data from a transaction.
//
// This method downloads the actual data payload of a transaction.
// For large transactions, this may take some time and use significant
// bandwidth. The data is returned in its original format.
//
// Parameters:
//   - id: The transaction ID containing the data
//
// Returns the raw transaction data as bytes, or an error if the
// transaction is not found or data cannot be retrieved.
//
// Example:
//
//	data, err := client.GetTransactionData("ABC123...")
//	if err != nil {
//		log.Printf("Failed to get data: %v", err)
//		return
//	}
//	fmt.Printf("Downloaded %d bytes\n", len(data))
func (c *Client) GetTransactionData(id string) ([]byte, error) {
	body, err := c.get(id)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// GetTransactionPrice calculates the cost to store data of a given size.
//
// This method queries the network for the current transaction fee based
// on data size and optional target address. Prices are returned in Winston
// units (1 AR = 1,000,000,000,000 Winston).
//
// Parameters:
//   - size: The size of data in bytes
//   - target: Optional target address (use empty string if not applicable)
//
// Returns the transaction fee in Winston as a string, or an error if
// the price cannot be calculated.
//
// Example:
//
//	price, err := client.GetTransactionPrice(1024, "")
//	if err != nil {
//		log.Printf("Failed to get price: %v", err)
//		return
//	}
//	fmt.Printf("Cost for 1KB: %s Winston\n", price)
func (c *Client) GetTransactionPrice(size int, target string) (string, error) {
	url := fmt.Sprintf("price/%d/%s", size, target)
	body, err := c.get(url)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetTransactionAnchor retrieves the current transaction anchor.
//
// Transaction anchors are used to prevent replay attacks by referencing
// recent network state. Each transaction should use a recent anchor to
// be accepted by the network. Anchors are typically valid for about 50 blocks.
//
// Returns the current anchor as a base64url-encoded string, or an error
// if the anchor cannot be retrieved.
//
// Example:
//
//	anchor, err := client.GetTransactionAnchor()
//	if err != nil {
//		log.Printf("Failed to get anchor: %v", err)
//		return
//	}
//	fmt.Printf("Current anchor: %s\n", anchor)
func (c *Client) GetTransactionAnchor() (string, error) {
	body, err := c.get("tx_anchor")
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// SubmitTransaction submits a signed transaction to the network.
//
// This method uploads a complete, signed transaction to the Arweave network
// for inclusion in the next block. The transaction must be properly signed
// and include all required fields.
//
// Parameters:
//   - tx: The complete, signed transaction to submit
//
// Returns the HTTP status code from the submission, or an error if
// the submission fails. Status 200 indicates successful acceptance.
//
// Example:
//
//	status, err := client.SubmitTransaction(signedTx)
//	if err != nil {
//		log.Printf("Submission failed: %v", err)
//		return
//	}
//	if status == 200 {
//		fmt.Println("Transaction submitted successfully")
//	}
func (c *Client) SubmitTransaction(tx *transaction.Transaction) (int, error) {
	b, err := json.Marshal(tx)
	if err != nil {
		return -1, err
	}
	return c.post("tx", b)
}

// GetWalletBalance retrieves the current AR token balance for a wallet.
//
// This method queries the current confirmed balance for a given wallet
// address. The balance is returned in Winston units (1 AR = 1,000,000,000,000 Winston).
// Pending transactions are not included in the balance.
//
// Parameters:
//   - address: The wallet address to query (base64url-encoded public key hash)
//
// Returns the wallet balance in Winston as a string, or an error if
// the address is invalid or cannot be queried.
//
// Example:
//
//	balance, err := client.GetWalletBalance("1seRanklLU_1VTGkEk7P0xAwMJfA7owA1JHW5KyZKlY")
//	if err != nil {
//		log.Printf("Failed to get balance: %v", err)
//		return
//	}
//	fmt.Printf("Wallet balance: %s Winston\n", balance)
func (c *Client) GetWalletBalance(address string) (string, error) {
	body, err := c.get(fmt.Sprintf("wallet/%s/balance", address))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// GetLastTransactionID retrieves the last transaction ID for a wallet.
//
// This method returns the transaction ID of the most recent transaction
// sent from the specified wallet address. This is useful for building
// transaction chains and verifying wallet activity.
//
// Parameters:
//   - address: The wallet address to query
//
// Returns the last transaction ID as a string, or an error if the
// address has no transactions or cannot be queried.
//
// Example:
//
//	lastTx, err := client.GetLastTransactionID("1seRanklLU_1VTGkEk7P0xAwMJfA7owA1JHW5KyZKlY")
//	if err != nil {
//		log.Printf("Failed to get last tx: %v", err)
//		return
//	}
//	fmt.Printf("Last transaction: %s\n", lastTx)
func (c *Client) GetLastTransactionID(address string) (string, error) {
	body, err := c.get(fmt.Sprintf("wallet/%s/last_tx", address))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// GetBlockByID retrieves block information by block hash.
//
// This method fetches complete block data including all transactions,
// mining information, and block metadata. Blocks are identified by
// their unique hash (independent hash).
//
// Parameters:
//   - id: The block hash (independent hash)
//
// Returns the complete Block struct with all block data, or an error
// if the block is not found.
//
// Example:
//
//	block, err := client.GetBlockByID("ABC123...")
//	if err != nil {
//		log.Printf("Block not found: %v", err)
//		return
//	}
//	fmt.Printf("Block height: %d, TX count: %d\n", block.Height, len(block.Txs))
func (c *Client) GetBlockByID(id string) (*Block, error) {
	body, err := c.get(fmt.Sprintf("block/hash/%s", id))
	if err != nil {
		return nil, err
	}
	b := &Block{}
	err = json.Unmarshal(body, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetBlockByHeight retrieves block information by block height.
//
// This method fetches block data for a specific block height in the
// Arweave blockchain. Block heights start from 0 (genesis block) and
// increase sequentially.
//
// Parameters:
//   - height: The block height as a string
//
// Returns the complete Block struct for that height, or an error if
// the height is invalid or the block cannot be retrieved.
//
// Example:
//
//	block, err := client.GetBlockByHeight("1000000")
//	if err != nil {
//		log.Printf("Failed to get block: %v", err)
//		return
//	}
//	fmt.Printf("Block at height 1M: %s\n", block.IndepHash)
func (c *Client) GetBlockByHeight(height string) (*Block, error) {
	body, err := c.get(fmt.Sprintf("block/hash/%s", height))
	if err != nil {
		return nil, err
	}
	b := &Block{}
	err = json.Unmarshal(body, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetNetworkInfo retrieves current network information and statistics.
//
// This method provides information about the Arweave network including
// current block height, network hash rate, peer count, and other
// network-wide statistics.
//
// Returns NetworkInfo struct with current network data, or an error
// if the information cannot be retrieved.
//
// Example:
//
//	info, err := client.GetNetworkInfo()
//	if err != nil {
//		log.Printf("Failed to get network info: %v", err)
//		return
//	}
//	fmt.Printf("Network height: %d, Peers: %d\n", info.Height, info.Peers)
func (c *Client) GetNetworkInfo() (*NetworkInfo, error) {
	body, err := c.get("info")
	if err != nil {
		return nil, err
	}
	n := NetworkInfo{}
	err = json.Unmarshal(body, &n)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// UploadChunk uploads a data chunk with its Merkle proof.
//
// This method is used for uploading individual chunks of large transactions.
// Each chunk includes the data and a Merkle proof that verifies its
// position within the complete dataset. This is part of Arweave's
// chunked upload system for large files.
//
// Parameters:
//   - chunk: The chunk data with proof information
//
// Returns the HTTP status code from the upload, or an error if the
// upload fails. Status 200 indicates successful acceptance.
//
// Example:
//
//	status, err := client.UploadChunk(chunkWithProof)
//	if err != nil {
//		log.Printf("Chunk upload failed: %v", err)
//		return
//	}
//	if status == 200 {
//		fmt.Println("Chunk uploaded successfully")
//	}
func (c *Client) UploadChunk(chunk *transaction.GetChunkResult) (int, error) {
	b, err := json.Marshal(chunk)
	if err != nil {
		return -1, err
	}
	return c.post("chunk", b)
}
