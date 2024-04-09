package client

type TransactionOffset struct {
	Size   int64 `json:"size"`
	Offset int64 `json:"offset"`
}
type TransactionChunk struct {
	Chunk    string `json:"chunk"`
	DataPath string `json:"data_path"`
	TxPath   string `json:"tx_path"`
}

type TransactionStatus struct {
	BlockHeight           int    `json:"block_height"`
	BlockIndepHash        string `json:"block_indep_hash"`
	NumberOfConfirmations int    `json:"number_of_confirmations"`
}
