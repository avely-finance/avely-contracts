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
	// "github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
)

type StubStakingContract struct {
	Code string
	Init []core.ContractValue
	Addr string
	Bech32 string
	Wallet *account.Wallet
}

func (s *StubStakingContract) LogContractStateJson() string {
	provider := provider2.NewProvider("http://zilliqa_server:5555")
	rsp, _ := provider.GetSmartContractState(s.Addr)
	j, _ := json.Marshal(rsp)
	s.LogPrettyStateJson(rsp)
	return string(j)
}

func (s *StubStakingContract) LogPrettyStateJson(data interface{}) {
	j, _ := json.MarshalIndent(data, "", "   ")
	log.Println(string(j))
}

func (s *StubStakingContract) GetBalance() string {
	provider := provider2.NewProvider("http://zilliqa_server:5555")
	balAndNonce, _ := provider.GetBalance(s.Addr)
	return balAndNonce.Balance
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

func (s *StubStakingContract) Call(transition string, params []core.ContractValue, amount string) (*transaction.Transaction, error) {
	contract := contract2.Contract{
		Address: s.Bech32,
		Signer:  s.Wallet,
	}

	tx, err := CallFor(&contract, transition, params, false, amount)
	if err != nil {
		return tx, err
	}
	tx.Confirm(tx.ID, 1, 1, contract.Provider)
	if tx.Status != core.Confirmed {
		return tx, errors.New("transaction didn't get confirmed")
	}
	if !tx.Receipt.Success {
		return tx, errors.New("transaction failed")
	}
	return tx, nil
}

func NewStubStakingContract(key string) (*StubStakingContract, error) {
	code, _ := ioutil.ReadFile("../contracts/stubStakingContract.scilla")
	// adminAddr := keytools.GetAddressFromPrivateKey(util.DecodeHex(key))

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

		return &StubStakingContract{
			Code: string(code),
			Init: init,
			Addr: tx.ContractAddress,
			Bech32: b32,
			Wallet: wallet,
		}, nil
	} else {
		return nil, errors.New("deploy failed")
	}
}
