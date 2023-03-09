package contracts

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

type FractionContract struct {
	Contract
}

func (a *FractionContract) Fraction(amount, x, y int) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			strconv.Itoa(amount),
		},
		{
			"x",
			"Uint128",
			strconv.Itoa(x),
		},
		{
			"y",
			"Uint128",
			strconv.Itoa(y),
		},
	}
	return a.Call("Fraction", args, "0")
}

func (a *FractionContract) FractionCeil(amount, x, y int) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			strconv.Itoa(amount),
		},
		{
			"x",
			"Uint128",
			strconv.Itoa(x),
		},
		{
			"y",
			"Uint128",
			strconv.Itoa(y),
		},
	}
	return a.Call("FractionCeil", args, "0")
}

func NewFractionContract(sdk *AvelySDK, deployer *account.Wallet) (*FractionContract, error) {
	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		},
	}

	contract := buildFractionContract(sdk, init)
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
		return &FractionContract{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func buildFractionContract(sdk *AvelySDK, init []core.ContractValue) contract2.Contract {
	return Restore("utils/fraction", sdk.InitProvider(), init)
}
