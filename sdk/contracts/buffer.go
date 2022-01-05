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
			"0x" + ssnaddr,
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
			"0x" + ssnaddr,
		},
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return b.Call("DelegateStakeSuccessCallBack", args, "0")
}

func (b *BufferContract) ChangeZproxyAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return b.Call("ChangeZproxyAddress", args, "0")
}

func (b *BufferContract) ChangeZimplAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return b.Call("ChangeZimplAddress", args, "0")
}

func (b *BufferContract) ChangeAzilSSNAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return b.Call("ChangeAzilSSNAddress", args, "0")
}

func (b *BufferContract) ChangeAimplAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return b.Call("ChangeAimplAddress", args, "0")
}

func (b *BufferContract) DelegateStake() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("DelegateStake", args, "0")
}

func (b *BufferContract) ClaimRewards() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("ClaimRewards", args, "0")
}

func (b *BufferContract) RequestDelegatorSwap(new_deleg_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_deleg_addr",
			"ByStr20",
			"0x" + new_deleg_addr,
		},
	}
	return b.Call("RequestDelegatorSwap", args, "0")
}

func NewBufferContract(sdk *AvelySDK, aimplAddr, zproxyAddr, zimplAddr string) (*BufferContract, error) {
	contract := buildBufferContract(sdk, aimplAddr, zproxyAddr, zimplAddr)

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
			Addr:            tx.ContractAddress,
			Bech32:          b32,
			Wallet:          contract.Signer,
			StateFieldTypes: stateFieldTypes,
		}
		return &BufferContract{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed: " + string(data))
	}
}

func RestoreBufferContract(sdk *AvelySDK, contractAddress, aimplAddr, zproxyAddr, zimplAddr string) (*BufferContract, error) {
	contract := buildBufferContract(sdk, aimplAddr, zproxyAddr, zimplAddr)

	b32, err := bech32.ToBech32Address(contractAddress)

	if err != nil {
		return nil, errors.New("Config has invalid Buffer address")
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

	return &BufferContract{Contract: sdkContract}, nil
}

func buildBufferContract(sdk *AvelySDK, aimplAddr, zproxyAddr, zimplAddr string) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/buffer.scilla")
	key := sdk.Cfg.AdminKey
	aZilSSNAddress := sdk.Cfg.AzilSsnAddress

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_admin_address",
			Type:  "ByStr20",
			Value: "0x" + sdk.GetAddressFromPrivateKey(key),
		}, {
			VName: "init_aimpl_address",
			Type:  "ByStr20",
			Value: "0x" + aimplAddr,
		}, {
			VName: "init_azil_ssn_address",
			Type:  "ByStr20",
			Value: "0x" + aZilSSNAddress,
		}, {
			VName: "init_zproxy_address",
			Type:  "ByStr20",
			Value: "0x" + zproxyAddr,
		}, {
			VName: "init_zimpl_address",
			Type:  "ByStr20",
			Value: "0x" + zimplAddr,
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
