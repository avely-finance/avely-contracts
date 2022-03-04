package contracts

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"

	. "github.com/avely-finance/avely-contracts/sdk/core"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type MultisigWallet struct {
	Contract
}

func (a *MultisigWallet) WithUser(key string) *MultisigWallet {
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)
	a.Contract.Wallet = wallet
	return a
}

// func (s *AZil) BalanceOf(addr string) *big.Int {
// 	rawState := s.Contract.SubState("balances", []string{addr})
// 	state := NewState(rawState)

// 	return state.Dig("result.balances." + addr).BigInt()
// }

func (s *MultisigWallet) SignTransaction(transactionId int) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "transactionId",
			Type:  "Uint32",
			Value: strconv.Itoa(transactionId),
		},
	}

	return s.Call("SignTransaction", args, "0")
}

func (s *MultisigWallet) ExecuteTransaction(transactionId int) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "transactionId",
			Type:  "Uint32",
			Value: strconv.Itoa(transactionId),
		},
	}

	return s.Call("ExecuteTransaction", args, "0")
}

func (s *MultisigWallet) SubmitChangeAzilSSNAddressTransaction(azilAddr, addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
		{
			VName: "address",
			Type:  "ByStr20",
			Value: addr,
		},
	}

	return s.Call("SubmitChangeAzilSSNAddressTransaction", args, "0")
}

func NewMultisigContract(sdk *AvelySDK, owners []string, requiredSignaturesCount int) (*MultisigWallet, error) {
	// TOOD: add requiredSignaturesCount validation
	contract := buildMultisigContract(sdk, owners, strconv.Itoa(requiredSignaturesCount))

	tx, err := sdk.DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, sdk.Cfg.TxConfrimMaxAttempts, sdk.Cfg.TxConfirmIntervalSec, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)

		sdkContract := Contract{
			Sdk:      sdk,
			Provider: *contract.Provider,
			Addr:     "0x" + tx.ContractAddress,
			Bech32:   b32,
			Wallet:   contract.Signer,
		}
		return &MultisigWallet{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreMultisigContract(sdk *AvelySDK, contractAddress string, owners []string, requiredSignaturesCount int) (*MultisigWallet, error) {
	contract := buildMultisigContract(sdk, owners, strconv.Itoa(requiredSignaturesCount))

	b32, err := bech32.ToBech32Address(contractAddress)

	if err != nil {
		return nil, errors.New("Config has invalid MultisigWallet address")
	}

	sdkContract := Contract{
		Sdk:      sdk,
		Provider: *contract.Provider,
		Addr:     contractAddress,
		Bech32:   b32,
		Wallet:   contract.Signer,
	}
	return &MultisigWallet{Contract: sdkContract}, nil
}

func buildMultisigContract(sdk *AvelySDK, owners []string, requiredSignaturesCount string) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/multisig_wallet.scilla")
	key := sdk.Cfg.AdminKey

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "owners_list",
			Type:  "List ByStr20",
			Value: owners,
		}, {
			VName: "required_signatures",
			Type:  "Uint32",
			Value: requiredSignaturesCount,
		},
	}

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)

	return contract2.Contract{
		Provider: sdk.InitProvider(),
		Code:     string(code),
		Init:     init,
		Signer:   wallet,
	}
}
