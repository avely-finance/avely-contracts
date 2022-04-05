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
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type BufferContract struct {
	Contract
}

func (b *BufferContract) AddFunds(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("AddFunds", args, amount)
}

func (b *BufferContract) WithdrawStakeRewardsSuccessCallBack(ssnaddr, rewards string) (*transaction.Transaction, error) {
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

func (b *BufferContract) DelegateStakeSuccessCallBack(ssnaddr, amount string) (*transaction.Transaction, error) {
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

func (b *BufferContract) DelegateStake(ssnaddr, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
	}
	return b.Call("DelegateStake", args, amount)
}

func (b *BufferContract) ClaimRewards(ssnaddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
	}
	return b.Call("ClaimRewards", args, "0")
}

func (b *BufferContract) ConfirmDelegatorSwap(requestor string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"requestor",
			"ByStr20",
			requestor,
		},
	}
	return b.Call("ConfirmDelegatorSwap", args, "0")
}

func (b *BufferContract) RejectDelegatorSwap(requestor string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"requestor",
			"ByStr20",
			requestor,
		},
	}
	return b.Call("RejectDelegatorSwap", args, "0")
}

func (b *BufferContract) RequestDelegatorSwap(new_deleg_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_deleg_addr",
			"ByStr20",
			new_deleg_addr,
		},
	}
	return b.Call("RequestDelegatorSwap", args, "0")
}

func (b *BufferContract) ReDelegateStake(ssnaddr, to_ssn, amount string) (*transaction.Transaction, error) {
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
	return b.Call("ReDelegateStake", args, "0")
}

func (b *BufferContract) ReDelegateStakeSuccessCallBack(ssnaddr, to_ssn, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
		{
			"tossn",
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

func NewBufferContract(sdk *AvelySDK, azilAddr, zproxyAddr string) (*BufferContract, error) {
	contract := buildBufferContract(sdk, azilAddr, zproxyAddr)

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
		return &BufferContract{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed: " + string(data))
	}
}

func RestoreBufferContract(sdk *AvelySDK, contractAddress, azilAddr, zproxyAddr string) (*BufferContract, error) {
	contract := buildBufferContract(sdk, azilAddr, zproxyAddr)

	b32, err := bech32.ToBech32Address(contractAddress)

	if err != nil {
		return nil, errors.New("Config has invalid Buffer address")
	}

	sdkContract := Contract{
		Sdk:      sdk,
		Provider: *contract.Provider,
		Addr:     contractAddress,
		Bech32:   b32,
		Wallet:   contract.Signer,
	}

	return &BufferContract{Contract: sdkContract}, nil
}

func buildBufferContract(sdk *AvelySDK, azilAddr, zproxyAddr string) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/buffer.scilla")
	key := sdk.Cfg.AdminKey
	aZilSSNAddress := sdk.Cfg.AzilSsnAddress

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_azil_address",
			Type:  "ByStr20",
			Value: azilAddr,
		}, {
			VName: "init_azil_ssn_address",
			Type:  "ByStr20",
			Value: aZilSSNAddress,
		}, {
			VName: "init_zproxy_address",
			Type:  "ByStr20",
			Value: zproxyAddr,
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
