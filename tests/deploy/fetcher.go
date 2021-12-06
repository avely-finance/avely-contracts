package deploy

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type FetcherContract struct {
	Contract
}

func (b *FetcherContract) AimplState() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("AimplState", args, "0")
}

func (b *FetcherContract) ZimplState() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("ZimplState", args, "0")
}

func (b *FetcherContract) AimplWithdrawalPending(bnum, delegator string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"bnum",
			"BNum",
			bnum,
		}, {
			"delegator",
			"ByStr20",
			delegator,
		},
	}
	return b.Call("AimplWithdrawalPending", args, "0")
}

func NewFetcherContract(key string, azilUtilsAddress string, aimplAddress string, stubStakingAddr string) (*FetcherContract, error) {
	code, _ := ioutil.ReadFile("../contracts/fetcher.scilla")
	type Constructor struct {
		Constructor string   `json:"constructor"`
		ArgTypes    []string `json:"argtypes"`
		Arguments   []string `json:"arguments"`
	}
	ats := []string{
		"String",
		"ByStr20",
	}
	ars := []string{
		"AzilUtils",
		"0x" + azilUtilsAddress,
	}
	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "_extlibs",
			Type:  "List(Pair String ByStr20)",
			Value: []Constructor{
				{
					Constructor: "Pair",
					ArgTypes:    ats,
					Arguments:   ars,
				},
			},
		}, {
			VName: "init_aimpl_address",
			Type:  "ByStr20",
			Value: "0x" + aimplAddress,
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
		return &FetcherContract{Contract: contract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		log.Println(string(data))
		return nil, errors.New("deploy failed")
	}
}
