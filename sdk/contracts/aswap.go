package contracts

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"regexp"
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

func (a *ASwap) SetTreasuryAddress(new_address string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "new_address",
			Type:  "ByStr20",
			Value: new_address,
		},
	}
	return a.Call("SetTreasuryAddress", args, "0")
}

func (a *ASwap) SetTreasuryFee(new_fee string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "new_fee",
			Type:  "Uint128",
			Value: new_fee,
		},
	}
	return a.Call("SetTreasuryFee", args, "0")
}

func (a *ASwap) SetLiquidityFee(new_fee string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "new_fee",
			Type:  "Uint256",
			Value: new_fee,
		},
	}
	return a.Call("SetLiquidityFee", args, "0")
}

func (a *ASwap) AddLiquidity(_amount, tokenAddr, minContributionAmount, tokenAmount string, blockNum int) (*transaction.Transaction, error) {
	deadline := blockNum + ASwapBlockShift

	args := []core.ContractValue{
		{
			VName: "token_address",
			Type:  "ByStr20",
			Value: tokenAddr,
		}, {
			VName: "min_contribution_amount",
			Type:  "Uint128",
			Value: minContributionAmount,
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
	return a.Call("AddLiquidity", args, _amount)
}

func (a *ASwap) RemoveLiquidity(tokenAddress, contributionAmount, minZilAmount, minTokenAmount string, blockNum int) (*transaction.Transaction, error) {
	deadline := blockNum + ASwapBlockShift

	args := []core.ContractValue{
		{
			VName: "token_address",
			Type:  "ByStr20",
			Value: tokenAddress,
		}, {
			VName: "contribution_amount",
			Type:  "Uint128",
			Value: contributionAmount,
		}, {
			VName: "min_zil_amount",
			Type:  "Uint128",
			Value: minZilAmount,
		}, {
			VName: "min_token_amount",
			Type:  "Uint128",
			Value: minTokenAmount,
		}, {
			VName: "deadline_block",
			Type:  "BNum",
			Value: strconv.Itoa(deadline),
		},
	}
	return a.Call("RemoveLiquidity", args, "0")
}

func (a *ASwap) SwapExactZILForTokens(_amount, tokenAddr, minTokenAmount, recipientAddress string, blockNum int) (*transaction.Transaction, error) {
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
	return a.Call("SwapExactZILForTokens", args, _amount)
}

func (a *ASwap) SwapExactTokensForZIL(tokenAddress, tokenAmount, minZilAmount, recipientAddress string, blockNum int) (*transaction.Transaction, error) {
	deadline := blockNum + ASwapBlockShift

	args := []core.ContractValue{
		{
			VName: "token_address",
			Type:  "ByStr20",
			Value: tokenAddress,
		}, {
			VName: "token_amount",
			Type:  "Uint128",
			Value: tokenAmount,
		}, {
			VName: "min_zil_amount",
			Type:  "Uint128",
			Value: minZilAmount,
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
	return a.Call("SwapExactTokensForZIL", args, "0")
}

func (a *ASwap) ChangeOwner(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_owner",
			"ByStr20",
			new_addr,
		},
	}
	return a.Call("ChangeOwner", args, "0")
}

func (a *ASwap) ClaimOwner() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("ClaimOwner", args, "0")
}

func NewASwap(sdk *AvelySDK, init_owner string) (*ASwap, error) {
	contract := buildASwapContract(sdk, init_owner)

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

		aswap := &ASwap{Contract: sdkContract}
		aswap.ErrorCodes = aswap.ParseErrorCodes(contract.Code)

		return aswap, nil

	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreASwap(sdk *AvelySDK, contractAddress string, init_owner string) (*ASwap, error) {
	contract := buildASwapContract(sdk, init_owner)

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

	aswap := &ASwap{Contract: sdkContract}
	aswap.ErrorCodes = aswap.ParseErrorCodes(contract.Code)

	return aswap, nil
}

func buildASwapContract(sdk *AvelySDK, init_owner string) contract2.Contract {
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

func (a *ASwap) ParseErrorCodes(contractCode string) ContractErrorCodes {
	codes := make(ContractErrorCodes)
	re := regexp.MustCompile(`(?s)\| *?([A-Za-z0-9]+) *?=>.*?code *?: *?Int32 *?([0-9-]+)`)
	results := re.FindAllStringSubmatch(contractCode, -1)
	for _, row := range results {
		codes[row[1]] = row[2]
	}
	return codes
}
