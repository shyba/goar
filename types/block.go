package types

type Block struct {
	Nonce          string   `json:"nonce"`
	PreviousBlock  string   `json:"previous_block"`
	Timestamp      uint64   `json:"timestamp"`
	LastRetarget   uint64   `json:"last_retarget"`
	Diff           string   `json:"diff"`
	Height         uint64   `json:"height"`
	Hash           string   `json:"hash"`
	IndepHash      string   `json:"indep_hash"`
	Txs            []string `json:"txs"`
	TxRoot         string   `json:"tx_root"`
	WalletList     string   `json:"wallet_list"`
	RewardAddr     string   `json:"reward_addr"`
	Tags           []Tag    `json:"tags"`
	RewardPool     uint64   `json:"reward_pool"`
	WeaveSize      uint64   `json:"weave_size"`
	BlockSize      uint64   `json:"block_size"`
	CumulativeDiff string   `json:"cumulative_diff"`
	HashListMerkle string   `json:"hash_list_merkle"`
}
