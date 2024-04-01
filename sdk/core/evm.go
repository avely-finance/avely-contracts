package core

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Evm struct {
	Cfg    EvmConfig
	Client *ethclient.Client
	// Store accounts directly within Evm struct
	// I had used KeyStore, but it stores data on filesystem, so I did simple wrapper for that instead
	Accounts []*ManagedAccount
}

type ManagedAccount struct {
	Account    accounts.Account
	PrivateKey *ecdsa.PrivateKey
}

// ScillaKind represents the kind of Scilla message.
type ScillaKind int

// Enumeration of Scilla message kinds.
const (
	Unknown ScillaKind = iota // Use of iota for automatic increment.
	ForwardedScillaError
	ForwardedScillaException
)

// ScillaDecoded represents the decoded structure for Scilla logs.
type ScillaDecoded struct {
	Kind        ScillaKind
	Description string
}

func NewEvm(config Config) *Evm {
	client, err := ethclient.Dial(config.Api.HttpUrl)
	if err != nil {
		panic(err)
	}

	return &Evm{
		Cfg:      config.Evm,
		Client:   client,
		Accounts: []*ManagedAccount{},
	}
}

func (evm *Evm) AddAccountByPrivateKey(hex string) (*accounts.Account, error) {
	privateKey, err := crypto.HexToECDSA(hex)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	account := accounts.Account{Address: fromAddress}
	managedAccount := &ManagedAccount{
		Account:    account,
		PrivateKey: privateKey,
	}
	evm.Accounts = append(evm.Accounts, managedAccount)

	return &account, nil
}

func (evm *Evm) findAccount(address common.Address) (*ManagedAccount, error) {
	for _, ma := range evm.Accounts {
		if ma.Account.Address == address {
			return ma, nil
		}
	}
	return nil, fmt.Errorf("account not found")
}

func (evm *Evm) GetTransactOpts(fromAccount *accounts.Account) (*bind.TransactOpts, error) {
	ma, err := evm.findAccount(fromAccount.Address)
	if err != nil {
		return nil, err
	}

	nonce, err := evm.Client.PendingNonceAt(context.Background(), fromAccount.Address)
	if err != nil {
		return nil, err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(ma.PrivateKey, big.NewInt(int64(evm.Cfg.ChainId)))

	if err != nil {
		return nil, err
	}

	opts.Nonce = big.NewInt(int64(nonce))
	opts.Value = big.NewInt(0) // in wei
	// opts.GasLimit = uint64(300000)         // in units
	opts.GasLimit = uint64(750000)         // in units
	opts.GasPrice = big.NewInt(1000000000) // in wei

	return opts, nil
}

func (evm *Evm) GetTransactOptsOrPanic(fromAccount *accounts.Account) *bind.TransactOpts {
	opts, err := evm.GetTransactOpts(fromAccount)
	if err != nil {
		panic(err)
	}
	return opts
}

// decodeScillaErrorOrException checks if the log is a Scilla error or exception and decodes the message.
func (evm *Evm) DecodeScillaErrorOrException(log *types.Log) (*ScillaDecoded, bool) {
	scillaErrorTopic := crypto.Keccak256Hash([]byte("ScillaError(string)")).String()
	scillaExceptionTopic := crypto.Keccak256Hash([]byte("ScillaException(string)")).String()

	var kind ScillaKind

	// Compare the first topic of the log to identify the kind of Scilla message.
	switch log.Topics[0].Hex() {
	case scillaErrorTopic:
		kind = ForwardedScillaError
	case scillaExceptionTopic:
		kind = ForwardedScillaException
	default:
		return nil, false // Log does not match known Scilla topics.
	}

	// Extracts the string length: The first 32 bytes are the offset, followed by 32 bytes for the string's length.
	stringLength := big.NewInt(0).SetBytes(log.Data[32:64]).Int64()

	// Extract the string, starting at byte 64 for the length of stringLength
	description := string(log.Data[64 : 64+stringLength])

	// for scilla event we could use this
	/*var eventData map[string]interface{}
	if err := json.Unmarshal([]byte(description), &eventData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	if _, ok := eventData["_eventname"]; !ok {
		return nil, errors.New("_eventname is required but not found")
		}*/

	return &ScillaDecoded{
		Kind:        kind,
		Description: description,
	}, true
}
