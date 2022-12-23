package contracts

import (
	"encoding/json"
	"errors"

	. "github.com/avely-finance/avely-contracts/sdk/core"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type HolderContract struct {
	Contract
}

func (b *HolderContract) ChangeOwner(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_owner",
			"ByStr20",
			new_addr,
		},
	}
	return b.Call("ChangeOwner", args, "0")
}

func (b *HolderContract) ClaimOwner() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("ClaimOwner", args, "0")
}

func (b *HolderContract) ChangeZimplAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			new_addr,
		},
	}
	return b.Call("ChangeZimplAddress", args, "0")
}

func (b *HolderContract) ChangeZproxyAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			new_addr,
		},
	}
	return b.Call("ChangeZproxyAddress", args, "0")
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

func (b *HolderContract) WithdrawStakeAmt(ssnaddr, amount string) (*transaction.Transaction, error) {
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
	return b.Call("WithdrawStakeAmt", args, "0")
}

func (b *HolderContract) CompleteWithdrawal() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("CompleteWithdrawal", args, "0")
}

func (b *HolderContract) ClaimRewards(ssnaddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
	}
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

func (b *HolderContract) DelegateStake(ssnaddr, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
	}
	return b.Call("DelegateStake", args, amount)
}

func NewHolderContract(sdk *AvelySDK, ownerAddr, stZilAddr, zproxyAddr, zimplAddr string, deployer *account.Wallet) (*HolderContract, error) {
	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "contract_owner",
			Type:  "ByStr20",
			Value: ownerAddr,
		}, {
			VName: "init_stzil_address",
			Type:  "ByStr20",
			Value: stZilAddr,
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

	contract := buildHolderContract(sdk, init)
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
		return &HolderContract{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed: " + string(data))
	}
}

func RestoreHolderContract(sdk *AvelySDK, contractAddress string) (*HolderContract, error) {
	contract := buildHolderContract(sdk, []core.ContractValue{})

	b32, err := bech32.ToBech32Address(contractAddress)
	if err != nil {
		return nil, errors.New("Config has invalid Holder address")
	}

	sdkContract := Contract{
		Sdk:      sdk,
		Provider: *contract.Provider,
		Addr:     contractAddress,
		Bech32:   b32,
		Wallet:   contract.Signer,
	}
	sdkContract.ErrorCodes = sdkContract.ParseErrorCodes(contract.Code)
	return &HolderContract{Contract: sdkContract}, nil
}

func buildHolderContract(sdk *AvelySDK, init []core.ContractValue) contract2.Contract {
	return Restore("holder", sdk.InitProvider(), init)
}
