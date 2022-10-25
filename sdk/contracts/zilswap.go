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

type ZilSwap struct {
	Contract
}

const blockShift = 3

func (z *ZilSwap) Initialize() (*transaction.Transaction, error) {
	args := []core.ContractValue{}

	return z.Call("Initialize", args, "0")
}

func (z *ZilSwap) AddLiquidity(tokenAddr, zilAmount, tokenAmount string, blockNum int) (*transaction.Transaction, error) {
	deadline := blockNum + blockShift

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
	return z.Call("AddLiquidity", args, zilAmount)
}

func (z *ZilSwap) SwapExactZILForTokens(tokenAddr, zilAmount, minTokenAmount, recipientAddress string, blockNum int) (*transaction.Transaction, error) {
	deadline := blockNum + blockShift

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
	return z.Call("SwapExactZILForTokens", args, zilAmount)
}

func NewZilSwap(sdk *AvelySDK) (*ZilSwap, error) {
	contract := buildZilSwapContract(sdk)

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
		return &ZilSwap{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreZilSwap(sdk *AvelySDK, contractAddress string) (*ZilSwap, error) {
	contract := buildZilSwapContract(sdk)

	b32, err := bech32.ToBech32Address(contractAddress)

	if err != nil {
		return nil, errors.New("Config has invalid ZilSwap address")
	}

	sdkContract := Contract{
		Sdk:      sdk,
		Provider: *contract.Provider,
		Addr:     contractAddress,
		Bech32:   b32,
		Wallet:   contract.Signer,
	}
	sdkContract.ErrorCodes = sdkContract.ParseErrorCodes(contract.Code)
	return &ZilSwap{Contract: sdkContract}, nil
}

func buildZilSwapContract(sdk *AvelySDK) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/zilswap/zilswap.scilla")
	key := sdk.Cfg.AdminKey

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "initial_owner",
			Type:  "ByStr20",
			Value: sdk.GetAddressFromPrivateKey(key),
		}, {
			VName: "initial_fee",
			Type:  "Uint256",
			Value: "30",
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
