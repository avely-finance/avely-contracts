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

type StubStakingContract struct {
	Contract
}

func (s *StubStakingContract) AddSSN(address string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			address,
		},
	}
	return s.Call("AddSSN", args, "0")
}

func (s *StubStakingContract) AssignStakeReward() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return s.Call("AssignStakeReward", args, "0")
}

func NewStubStakingContract(key string) (*StubStakingContract, error) {
	code, _ := ioutil.ReadFile("../contracts/stubStakingContract.scilla")

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
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

		return &StubStakingContract{Contract: contract}, nil
	} else {
		return nil, errors.New("deploy failed")
	}
}
