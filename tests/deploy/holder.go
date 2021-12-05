package deploy

import (
	"errors"
	"io/ioutil"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type HolderContract struct {
	Contract
}

func (b *HolderContract) ChangeProxyStakingContractAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return b.Call("ChangeProxyStakingContractAddress", args, "0")
}

func (b *HolderContract) CompleteWithdrawal() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("CompleteWithdrawal", args, "0")
}

func NewHolderContract(key string, aimplAddress string, aZilSSNAddress string, stubStakingAddr string) (*HolderContract, error) {
	code, _ := ioutil.ReadFile("../contracts/holder.scilla")

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_aimpl_address",
			Type:  "ByStr20",
			Value: "0x" + aimplAddress,
		}, {
			VName: "init_azil_ssn_address",
			Type:  "ByStr20",
			Value: aZilSSNAddress,
		}, {
			VName: "init_proxy_staking_contract_address",
			Type:  "ByStr20",
			Value: "0x" + stubStakingAddr,
		},
	}

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)

	contract := contract2.Contract{
		Code:   string(code),
		Init:   init,
		Signer: wallet,
	}

	tx, err := DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, TxConfirmMaxAttempts, TxConfirmInterval, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)
		contract := Contract{
			Code:   string(code),
			Init:   init,
			Addr:   tx.ContractAddress,
			Bech32: b32,
			Wallet: wallet,
		}
		return &HolderContract{Contract: contract}, nil
	} else {
		return nil, errors.New("deploy failed")
	}
}
