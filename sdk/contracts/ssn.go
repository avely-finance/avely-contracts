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

type SsnContract struct {
	Contract
}

func (a *SsnContract) WithUser(key string) *SsnContract {
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)
	a.Contract.Wallet = wallet
	return a
}

func (a *SsnContract) ChangeOwner(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_owner",
			"ByStr20",
			new_addr,
		},
	}
	return a.Call("ChangeOwner", args, "0")
}

func (a *SsnContract) ClaimOwner() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("ClaimOwner", args, "0")
}

func (a *SsnContract) ChangeZproxy(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_address",
			"ByStr20",
			new_addr,
		},
	}
	return a.Call("ChangeZproxy", args, "0")
}

func (a *SsnContract) UpdateReceivingAddr(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_addr",
			"ByStr20",
			new_addr,
		},
	}
	return a.Call("UpdateReceivingAddr", args, "0")
}

func (a *SsnContract) UpdateComm(new_rate int) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_rate",
			"Uint128",
			strconv.Itoa(new_rate),
		},
	}
	return a.Call("UpdateComm", args, "0")
}

func (a *SsnContract) WithdrawComm() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("WithdrawComm", args, "0")
}

func NewSsnContract(sdk *AvelySDK, init_owner, init_zproxy string) (*SsnContract, error) {
	contract := buildSsnContract(sdk, init_owner, init_zproxy)

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

		return &SsnContract{Contract: contract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreSsnContract(sdk *AvelySDK, contractAddress, init_owner, init_zproxy string) (*SsnContract, error) {
	contract := buildSsnContract(sdk, init_owner, init_zproxy)

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

	return &SsnContract{Contract: sdkContract}, nil
}

func buildSsnContract(sdk *AvelySDK, init_owner, init_zproxy string) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/ssn.scilla")
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
			VName: "init_zproxy",
			Type:  "ByStr20",
			Value: init_zproxy,
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