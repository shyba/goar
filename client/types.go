package client

import "github.com/liteseed/goar/tag"

// Block represents a block in the Arweave blockchain.
//
// This struct contains all the information about a block including
// its hash, height, transactions, and mining-related data. Blocks
// are the fundamental units of the Arweave blockchain that contain
// batches of transactions.
type Block struct {
	Nonce          string    `json:"nonce"`            // Mining nonce used to find the block
	PreviousBlock  string    `json:"previous_block"`   // Hash of the previous block
	Timestamp      uint64    `json:"timestamp"`        // Unix timestamp when block was mined
	LastRetarget   uint64    `json:"last_retarget"`    // Timestamp of last difficulty retarget
	Diff           string    `json:"diff"`             // Current mining difficulty
	Height         uint64    `json:"height"`           // Block height (number of blocks since genesis)
	Hash           string    `json:"hash"`             // Block hash (dependent on transaction order)
	IndepHash      string    `json:"indep_hash"`       // Independent hash (does not depend on transaction order)
	Txs            []string  `json:"txs"`              // List of transaction IDs in this block
	TxRoot         string    `json:"tx_root"`          // Merkle root of transaction tree
	WalletList     string    `json:"wallet_list"`      // Hash of wallet list at this block
	RewardAddr     string    `json:"reward_addr"`      // Address that will receive mining reward
	Tags           []tag.Tag `json:"tags"`             // Optional tags attached to the block
	RewardPool     uint64    `json:"reward_pool"`      // Current size of mining reward pool
	WeaveSize      uint64    `json:"weave_size"`       // Total size of data stored in Arweave
	BlockSize      uint64    `json:"block_size"`       // Size of this block in bytes
	CumulativeDiff string    `json:"cumulative_diff"`  // Cumulative difficulty since genesis
	HashListMerkle string    `json:"hash_list_merkle"` // Merkle root of block hash list
}

// NetworkInfo represents current information about the Arweave network.
//
// This struct contains statistics and metadata about the overall state
// of the Arweave network, including version information, network size,
// and performance metrics.
type NetworkInfo struct {
	Network          string `json:"network"`            // Network identifier (usually "arweave.N.1")
	Version          int64  `json:"version"`            // Network protocol version
	Release          int64  `json:"release"`            // Node software release number
	Height           int64  `json:"height"`             // Current block height
	Current          string `json:"current"`            // Hash of current block
	Blocks           int64  `json:"blocks"`             // Total number of blocks
	Peers            int64  `json:"peers"`              // Number of connected peers
	QueueLength      int64  `json:"queue_length"`       // Number of transactions in mempool
	NodeStateLatency int64  `json:"node_state_latency"` // Node state synchronization latency
}

// TransactionStatus represents the confirmation status of a transaction.
//
// This struct provides information about whether a transaction has been
// confirmed by the network and included in a block. It includes the
// block information if the transaction is confirmed.
type TransactionStatus struct {
	BlockHeight           int    `json:"block_height"`            // Height of block containing this transaction (0 if unconfirmed)
	BlockIndepHash        string `json:"block_indep_hash"`        // Independent hash of block containing this transaction
	NumberOfConfirmations int    `json:"number_of_confirmations"` // Number of confirmations (blocks since inclusion)
	Confirmed             bool   `json:"-"`                       // Whether the transaction is confirmed (derived field)
}
