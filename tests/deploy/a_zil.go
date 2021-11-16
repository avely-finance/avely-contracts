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

type AZil struct {
	Contract
}

func (a *AZil) ChangeBufferAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return a.Contract.Call("ChangeBufferAddress", args, "0")
}

func (a *AZil) ChangeHolderAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return a.Contract.Call("ChangeHolderAddress", args, "0")
}

func (a *AZil) DelegateStake(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{}

	return a.Call("DelegateStake", args, amount)
}

func (a *AZil) WithdrawStakeAmt(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return a.Call("WithdrawStakeAmt", args, "0")
}

func (a *AZil) ZilBalanceOf (addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + addr,
		},
	}
	return a.Contract.Call("ZilBalanceOf", args, "0")
}


func NewAZilContract(key string, aZilSSNAddress string, stubStakingAddr string) (*AZil, error) {
	code, _ := ioutil.ReadFile("../contracts/aZil.scilla")

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_azil_ssn_address",
			Type:  "ByStr20",
			Value: aZilSSNAddress,
		}, {
			VName: "init_proxy_staking_contract_address",
			Type:  "ByStr20",
			Value: "0x" + stubStakingAddr,
		}, {
			VName: "init_buffer_address",
			Type:  "ByStr20",
			Value: "0xb2e2c996e6068f4ae11c4cc2c6a189b774819f79",
		}, {
			VName: "init_holder_address",
			Type:  "ByStr20",
			Value: "0xb2e2c996e6068f4ae11c4cc2c6a189b774819f79",
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
	tx.Confirm(tx.ID, 1, 1, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)
		contract := Contract{
			Code:   string(code),
			Init:   init,
			Addr:   tx.ContractAddress,
			Bech32: b32,
			Wallet: wallet,
		}

		return &AZil{Contract: contract}, nil
	} else {
		return nil, errors.New("deploy failed")
	}
}
