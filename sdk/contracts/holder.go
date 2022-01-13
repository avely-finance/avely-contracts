package contracts

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	. "github.com/avely-finance/avely-contracts/sdk/core"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type HolderContract struct {
	Contract
}

func (b *HolderContract) AddFunds(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("AddFunds", args, amount)
}

func (b *HolderContract) CompleteWithdrawalNoUnbondedStakeCallBack(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return b.Call("CompleteWithdrawalNoUnbondedStakeCallBack", args, "0")
}

func (b *HolderContract) CompleteWithdrawalSuccessCallBack(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return b.Call("CompleteWithdrawalSuccessCallBack", args, "0")
}

func (b *HolderContract) DelegateStakeSuccessCallBack(ssnaddr, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return b.Call("DelegateStakeSuccessCallBack", args, "0")
}

func (b *HolderContract) WithdrawStakeRewardsSuccessCallBack(ssnaddr, rewards string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
		{
			"rewards",
			"Uint128",
			rewards,
		},
	}
	return b.Call("WithdrawStakeRewardsSuccessCallBack", args, "0")
}

func (b *HolderContract) WithdrawStakeAmtSuccessCallBack(ssnaddr, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return b.Call("WithdrawStakeAmtSuccessCallBack", args, "0")
}

func (b *HolderContract) WithdrawStakeAmt(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return b.Call("WithdrawStakeAmt", args, "0")
}

func (b *HolderContract) CompleteWithdrawal() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("CompleteWithdrawal", args, "0")
}

func (b *HolderContract) ClaimRewards() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("ClaimRewards", args, "0")
}

func (b *HolderContract) ConfirmDelegatorSwap(requestor string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"requestor",
			"ByStr20",
			requestor,
		},
	}
	return b.Call("ConfirmDelegatorSwap", args, "0")
}

func (b *HolderContract) DelegateStake(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("DelegateStake", args, amount)
}

func (b *HolderContract) ReDelegateStake(ssnaddr, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return b.Call("ReDelegateStake", args, "0")
}

func (b *HolderContract) ReDelegateStakeSuccessCallBack(ssnaddr, to_ssn, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
		{
			"to_ssn",
			"ByStr20",
			to_ssn,
		},
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return b.Call("ReDelegateStakeSuccessCallBack", args, "0")
}

func NewHolderContract(sdk *AvelySDK, aimplAddr, zproxyAddr, zimplAddr string) (*HolderContract, error) {
	contract := buildHolderContract(sdk, aimplAddr, zproxyAddr, zimplAddr)

	tx, err := sdk.DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, sdk.Cfg.TxConfrimMaxAttempts, sdk.Cfg.TxConfirmIntervalSec, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)
		stateFieldTypes := make(StateFieldTypes)
		sdkContract := Contract{
			Sdk:             sdk,
			Provider:        *contract.Provider,
			Addr:            "0x" + tx.ContractAddress,
			Bech32:          b32,
			Wallet:          contract.Signer,
			StateFieldTypes: stateFieldTypes,
		}
		return &HolderContract{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed: " + string(data))
	}
}

func RestoreHolderContract(sdk *AvelySDK, contractAddress, aimplAddr, zproxyAddr, zimplAddr string) (*HolderContract, error) {
	contract := buildHolderContract(sdk, aimplAddr, zproxyAddr, zimplAddr)

	b32, err := bech32.ToBech32Address(contractAddress)
	if err != nil {
		return nil, errors.New("Config has invalid Holder address")
	}

	stateFieldTypes := make(StateFieldTypes)
	sdkContract := Contract{
		Sdk:             sdk,
		Provider:        *contract.Provider,
		Addr:            contractAddress,
		Bech32:          b32,
		Wallet:          contract.Signer,
		StateFieldTypes: stateFieldTypes,
	}
	return &HolderContract{Contract: sdkContract}, nil
}

func buildHolderContract(sdk *AvelySDK, aimplAddr, zproxyAddr, zimplAddr string) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/holder.scilla")
	key := sdk.Cfg.AdminKey
	aZilSSNAddress := sdk.Cfg.AzilSsnAddress

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_aimpl_address",
			Type:  "ByStr20",
			Value: aimplAddr,
		}, {
			VName: "init_azil_ssn_address",
			Type:  "ByStr20",
			Value: aZilSSNAddress,
		}, {
			VName: "init_zproxy_address",
			Type:  "ByStr20",
			Value: zproxyAddr,
		}, {
			VName: "init_zimpl_address",
			Type:  "ByStr20",
			Value: zimplAddr,
		},
	}

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)

	return contract2.Contract{
		Provider: provider2.NewProvider(sdk.Cfg.ApiUrl),
		Code:     string(code),
		Init:     init,
		Signer:   wallet,
	}
}
