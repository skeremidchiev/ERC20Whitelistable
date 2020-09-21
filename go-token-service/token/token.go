package token

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"strconv"

	ethereum "github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"golang.org/x/crypto/sha3"

	"ERC20Whitelistable/go-token-service/contracts"
)

var (
	InvalidAddressError = errors.New("Invalid Address")
)

type WhitelistableToken struct {
	EthClient       *ethclient.Client  // infura client
	TransactOpts    *bind.TransactOpts // transaction options
	CallerAddres    *common.Address    // address of the contract's owner
	ContractAddress *common.Address    // address of the contract's owner
	Token           *token.Token       // contract instance

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
		&address,
		instance,
		whitelistedRole,
		minterRole,
		&sync.Mutex{},
	}

	return obj, nil
}

func (wlt *WhitelistableToken) WhitelistAddress(i *WhitelistInput) (*TxOutput, error) {
	txo := &TxOutput{i.Address, "", false}

	// check if address is valid
	if ok := IsValidAddress(i.Address); !ok {
		return txo, InvalidAddressError
	}

	// check estimateGas
	if _, err := wlt.egWhitelisting(i.Address); err != nil {
		return txo, err
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

	// can check if transaction is Mined and OK but it will slow down response
	txo.OK = true // wlt.getStatusOfTX(tx)
	txo.TransactionHash = tx.Hash().Hex()
	return txo, nil
}

func (wlt *WhitelistableToken) Mint(i *MintInput) (*TxOutput, error) {
	txo := &TxOutput{i.Address, "", false}

	// check if address is valid
	if ok := IsValidAddress(i.Address); !ok {
		return txo, InvalidAddressError
	}

	// check estimateGas
	if _, err := wlt.egMint(i.Address, i.Amount); err != nil {
		return txo, err
	}

	// increment nonce at start to prevent "Error: Known Transaction"
	wlt.incrementNonce()

	amountInt64, _ := strconv.ParseInt(i.Amount, 10, 64)
	tx, err := wlt.Token.Mint(
		wlt.TransactOpts,
		common.HexToAddress(i.Address),
		big.NewInt(amountInt64),
	)
	if err != nil {
		return txo, err
	}

	txo.OK = true
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

// egGrantRole Estimate Gas for Whitelisting
func (wlt *WhitelistableToken) egWhitelisting(address string) (uint64, error){
	// method
	transferFnSignature := []byte("grantRole(bytes32,address)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	// address
	addr := common.HexToAddress(address)
	paddedAddress := common.LeftPadBytes(addr.Bytes(), 32)
	// role
	paddedRole := common.LeftPadBytes(wlt.WhitelistedRole[:], 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedRole...)
	data = append(data, paddedAddress...)

	gasLimit, err := wlt.EthClient.EstimateGas(context.Background(), ethereum.CallMsg{
		From: *wlt.CallerAddres,
		To:   wlt.ContractAddress,
		Data: data,
	})
	if err != nil {
		return 0, err
	}
	return gasLimit, nil
}

// egMint Estimate Gas for minting
func (wlt *WhitelistableToken) egMint(address, amount string) (uint64, error) {
	// method
	transferFnSignature := []byte("mint(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	// address
	addr := common.HexToAddress(address)
	paddedAddress := common.LeftPadBytes(addr.Bytes(), 32)

	// amount
	amountBN := new(big.Int)
	amountBN.SetString(amount, 10)
	paddedAmount := common.LeftPadBytes(amountBN.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := wlt.EthClient.EstimateGas(context.Background(), ethereum.CallMsg{
		From: *wlt.CallerAddres,
		To:   wlt.ContractAddress,
		Data: data,
	})
	if err != nil {
		return 0, err
	}
	return gasLimit, nil
}