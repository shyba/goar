package client

import "github.com/liteseed/goar/tag"

type Block struct {
	Nonce          string    `json:"nonce"`
	PreviousBlock  string    `json:"previous_block"`
	Timestamp      uint64    `json:"timestamp"`
	LastRetarget   uint64    `json:"last_retarget"`
	Diff           string    `json:"diff"`
	Height         uint64    `json:"height"`
	Hash           string    `json:"hash"`
	IndepHash      string    `json:"indep_hash"`
	Txs            []string  `json:"txs"`
	TxRoot         string    `json:"tx_root"`
	WalletList     string    `json:"wallet_list"`
	RewardAddr     string    `json:"reward_addr"`
	Tags           []tag.Tag `json:"tags"`
	RewardPool     uint64    `json:"reward_pool"`
	WeaveSize      uint64    `json:"weave_size"`
	BlockSize      uint64    `json:"block_size"`
	CumulativeDiff string    `json:"cumulative_diff"`
	HashListMerkle string    `json:"hash_list_merkle"`
}

type NetworkInfo struct {
	Network          string `json:"network"`
	Version          int64  `json:"version"`
	Release          int64  `json:"release"`
	Height           int64  `json:"height"`
	Current          string `json:"current"`
	Blocks           int64  `json:"blocks"`
	Peers            int64  `json:"peers"`
	QueueLength      int64  `json:"queue_length"`
	NodeStateLatency int64  `json:"node_state_latency"`
}

type TransactionStatus struct {
	BlockHeight           int    `json:"block_height"`
	BlockIndepHash        string `json:"block_indep_hash"`
	NumberOfConfirmations int    `json:"number_of_confirmations"`
}
