package contracts

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

type ASwap struct {
	Contract
}

const ASwapBlockShift = 3

func (a *ASwap) WithUser(key string) *ASwap {
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)
	a.Contract.Wallet = wallet
	return a
}

func (a *ASwap) TogglePause() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("TogglePause", args, "0")
}

func (a *ASwap) SetTreasuryFee(new_fee string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "new_fee",
			Type:  "Uint128",
			Value: new_fee,
		},
	}
	return a.Call("SetTreasuryFee", args, new_fee)
}

func (a *ASwap) SetLiquidityFee(new_fee string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "new_fee",
			Type:  "Uint256",
			Value: new_fee,
		},
	}
	return a.Call("SetLiquidityFee", args, new_fee)
}

func (a *ASwap) AddLiquidity(tokenAddr, zilAmount, tokenAmount string, blockNum int) (*transaction.Transaction, error) {
	deadline := blockNum + ASwapBlockShift

	args := []core.ContractValue{
		{
			VName: "token_address",
			Type:  "ByStr20",
			Value: tokenAddr,
		}, {
			VName: "min_contribution_amount",
			Type:  "Uint128",
			Value: "0",
		}, {
			VName: "max_token_amount",
			Type:  "Uint128",
			Value: tokenAmount,
		}, {
			VName: "deadline_block",
			Type:  "BNum",
			Value: strconv.Itoa(deadline),
		},
	}
	return a.Call("AddLiquidity", args, zilAmount)
}

func (a *ASwap) SwapExactZILForTokens(tokenAddr, zilAmount, minTokenAmount, recipientAddress string, blockNum int) (*transaction.Transaction, error) {
	deadline := blockNum + ASwapBlockShift

	args := []core.ContractValue{
		{
			VName: "token_address",
			Type:  "ByStr20",
			Value: tokenAddr,
		}, {
			VName: "min_token_amount",
			Type:  "Uint128",
			Value: minTokenAmount,
		}, {
			VName: "deadline_block",
			Type:  "BNum",
			Value: strconv.Itoa(deadline),
		}, {
			VName: "recipient_address",
			Type:  "ByStr20",
			Value: recipientAddress,
		},
	}
	return a.Call("SwapExactZILForTokens", args, zilAmount)
}

func NewASwap(sdk *AvelySDK, init_owner string, operators []string) (*ASwap, error) {
	contract := buildASwapContract(sdk, init_owner, operators)

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

		return &ASwap{Contract: contract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreASwap(sdk *AvelySDK, contractAddress string, init_owner string, operators []string) (*ASwap, error) {
	contract := buildASwapContract(sdk, init_owner, operators)

	b32, err := bech32.ToBech32Address(contractAddress)

	if err != nil {
		return nil, errors.New("Config has invalid ASwap address")
	}

	sdkContract := Contract{
		Sdk:      sdk,
		Provider: *contract.Provider,
		Addr:     contractAddress,
		Bech32:   b32,
		Wallet:   contract.Signer,
	}

	return &ASwap{Contract: sdkContract}, nil
}

func buildASwapContract(sdk *AvelySDK, init_owner string, operators []string) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/aswap.scilla")
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
		}, {
			VName: "operators",
			Type:  "List ByStr20",
			Value: operators,
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
