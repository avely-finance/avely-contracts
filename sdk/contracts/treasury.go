package contracts

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

type TreasuryContract struct {
	Contract
}

func (a *TreasuryContract) WithUser(key string) *TreasuryContract {
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)
	a.Contract.Wallet = wallet
	return a
}

func (a *TreasuryContract) ChangeOwner(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_owner",
			"ByStr20",
			new_addr,
		},
	}
	return a.Call("ChangeOwner", args, "0")
}

func (a *TreasuryContract) Withdraw(recipient, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"recipient",
			"ByStr20",
			recipient,
		},
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return a.Call("Withdraw", args, "0")
}

func (a *TreasuryContract) AddFunds(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("AddFunds", args, amount)
}

func (a *TreasuryContract) ClaimOwner() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("ClaimOwner", args, "0")
}

func NewTreasuryContract(sdk *AvelySDK, init_owner string) (*TreasuryContract, error) {
	contract := buildTreasuryContract(sdk, init_owner)

	tx, err := sdk.DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, sdk.Cfg.TxConfrimMaxAttempts, sdk.Cfg.TxConfirmIntervalSec, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)

		contract := Contract{
			Sdk:      sdk,
			Provider: *contract.Provider,
			Addr:     "0x" + tx.ContractAddress,
			Bech32:   b32,
			Wallet:   contract.Signer,
		}

		return &TreasuryContract{Contract: contract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreTreasuryContract(sdk *AvelySDK, contractAddress string, init_owner string) (*TreasuryContract, error) {
	contract := buildTreasuryContract(sdk, init_owner)

	b32, err := bech32.ToBech32Address(contractAddress)

	if err != nil {
		return nil, errors.New("Config has invalid Treasury address")
	}

	sdkContract := Contract{
		Sdk:      sdk,
		Provider: *contract.Provider,
		Addr:     contractAddress,
		Bech32:   b32,
		Wallet:   contract.Signer,
	}

	return &TreasuryContract{Contract: sdkContract}, nil
}

func buildTreasuryContract(sdk *AvelySDK, init_owner string) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/treasury.scilla")
	key := sdk.Cfg.AdminKey

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_owner",
			Type:  "ByStr20",
			Value: init_owner,
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