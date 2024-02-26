package contracts

import (
	"encoding/json"
	"errors"

	"github.com/Zilliqa/gozilliqa-sdk/v3/account"
	"github.com/Zilliqa/gozilliqa-sdk/v3/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/v3/contract"
	"github.com/Zilliqa/gozilliqa-sdk/v3/core"
	"github.com/Zilliqa/gozilliqa-sdk/v3/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

type TreasuryContract struct {
	Contract
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

func NewTreasuryContract(sdk *AvelySDK, init_owner string, deployer *account.Wallet) (*TreasuryContract, error) {
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

	contract := buildTreasuryContract(sdk, init)
	contract.Signer = deployer

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
		sdkContract.ErrorCodes = sdkContract.ParseErrorCodes(contract.Code)
		return &TreasuryContract{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreTreasuryContract(sdk *AvelySDK, contractAddress string) (*TreasuryContract, error) {
	contract := buildTreasuryContract(sdk, []core.ContractValue{})

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
	sdkContract.ErrorCodes = sdkContract.ParseErrorCodes(contract.Code)
	return &TreasuryContract{Contract: sdkContract}, nil
}

func buildTreasuryContract(sdk *AvelySDK, init []core.ContractValue) contract2.Contract {
	return Restore("treasury", sdk.InitProvider(), init)
}
