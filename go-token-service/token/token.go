package token

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"ERC20Whitelistable/go-token-service/contracts"
)

var (
	InvalidAddressError = errors.New("Invalid Address")
)

type WhitelistableToken struct {
	EthClient    *ethclient.Client  // infura client
	TransactOpts *bind.TransactOpts // transaction options
	CallerAddres *common.Address    // address of the contract's owner
	Token        *token.Token       // contract instance

	WhitelistedRole [32]byte // simple can do keccak256("WHITELISTED_ROLE")
	MinterRole      [32]byte // simple can do keccak256("MINTER_ROLE") but taking it from contract is safer

	*sync.Mutex // used to protect TransactOpts.Nonce
}

// GetWhitelistableToken
// Generates WhitelistablToken's context needed for contract's method calls
func GetWhitelistableToken() (*WhitelistableToken, error) {
	// reading all specific and sensitive data from config file
	cfg := GetConfig()

	// set up client
	client, err := ethclient.Dial(fmt.Sprintf("https://%s.infura.io/v3/%s", cfg.Network, cfg.InfuraKey))
	if err != nil {
		return nil, err
	}

	// set up keys
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("Error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// set up TransactOpts
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	trOpts := bind.NewKeyedTransactor(privateKey)
	// start nonce with one down and increase it just before transaction
	trOpts.Nonce = big.NewInt(int64(nonce - 1))
	trOpts.Value = big.NewInt(0)     // in wei
	trOpts.GasLimit = uint64(300000) // in units
	trOpts.GasPrice = gasPrice

	// contract instance
	address := common.HexToAddress(cfg.ContractAddress)
	instance, err := token.NewToken(address, client)
	if err != nil {
		return nil, err
	}

	// WhitelistedRole
	whitelistedRole, err := instance.WHITELISTEDROLE(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	// MinterRole
	minterRole, err := instance.MINTERROLE(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	obj := &WhitelistableToken{
		client,
		trOpts,
		&fromAddress,
		instance,
		whitelistedRole,
		minterRole,
		&sync.Mutex{},
	}

	return obj, nil
}

func (wlt *WhitelistableToken) WhitelistAddress(i *WhitelistInput) (*TxOutput, error) {
	txo := &TxOutput{ i.Address, "", false }

	if ok := IsValidAddress(i.Address); !ok {
		return txo, InvalidAddressError
	}

	// increment nonce at start to prevent "Error: Known Transaction"
	wlt.incrementNonce()

	tx, err := wlt.Token.GrantRole(
		wlt.TransactOpts,
		wlt.WhitelistedRole,
		common.HexToAddress(i.Address),
	)
	if err != nil {
		return txo, err
	}

	txo.OK = wlt.getStatusOfTX(tx)
	txo.TransactionHash = tx.Hash().Hex()
	return txo, nil
}

func (wlt *WhitelistableToken) Mint(i *MintInput) (*TxOutput, error) {
	txo := &TxOutput{ i.Address, "", false }

	if ok := IsValidAddress(i.Address); !ok {
		return txo, InvalidAddressError
	}

	// increment nonce at start to prevent "Error: Known Transaction"
	wlt.incrementNonce()

	tx, err := wlt.Token.Mint(
		wlt.TransactOpts,
		common.HexToAddress(i.Address),
		big.NewInt(int64(i.Amount)),
	)
	if err != nil {
		return txo, err
	}

	txo.OK = wlt.getStatusOfTX(tx)
	txo.TransactionHash = tx.Hash().Hex()
	return txo, nil
}

// getStatusOfTX checks transaction Status - returns error on Status != 0
func (wlt *WhitelistableToken) getStatusOfTX(tx *types.Transaction) bool {
	receipt, err := bind.WaitMined(context.Background(), wlt.EthClient, tx)
	if err != nil {
		return false
	}

	// 0 - on revert or failure and 1 - on success
	// https://ethereum.stackexchange.com/questions/28889/what-is-the-exact-meaning-of-a-transactions-new-receipt-status-field
	if receipt.Status != 1 {
		return false
	}

	return true
}

// incrementNonce manages nonce
func (wlt *WhitelistableToken) incrementNonce() {
	wlt.Lock()
	wlt.TransactOpts.Nonce.Add(wlt.TransactOpts.Nonce, big.NewInt(int64(1)))
	wlt.Unlock()
}
