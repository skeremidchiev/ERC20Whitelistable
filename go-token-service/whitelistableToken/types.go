package whitelistableToken

// WhitelistInput simple wrapper for WhitelistAddress() inputs
type WhitelistInput struct {
	Address string `json:"address"`
}

// WhitelistInput simple wrapper for Mint() inputs
type MintInput struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

// TxOutput simple wrapper for WhitelistAddress() and Mint() outputs
type TxOutput struct {
	TransactionHash string `json:"txHash"`
}
