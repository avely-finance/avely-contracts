package deploy

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type Contract struct {
	Code   string
	Init   []core.ContractValue
	Addr   string
	Bech32 string
	Wallet *account.Wallet
}

func (c *Contract) LogContractStateJson() string {
	provider := provider2.NewProvider("http://zilliqa_server:5555")
	rsp, _ := provider.GetSmartContractState(c.Addr)
	j, _ := json.Marshal(rsp)
	c.LogPrettyStateJson(rsp)
	return string(j)
}

func (c *Contract) LogPrettyStateJson(data interface{}) {
	j, _ := json.MarshalIndent(data, "", "   ")
	log.Println(string(j))
}

func (c *Contract) GetBalance() string {
	provider := provider2.NewProvider("http://zilliqa_server:5555")
	balAndNonce, _ := provider.GetBalance(c.Addr)
	return balAndNonce.Balance
}

func (c *Contract) UpdateWallet(newKey string) {
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(newKey)
	c.Wallet = wallet
}

func (c *Contract) Call(transition string, params []core.ContractValue, amount string) (*transaction.Transaction, error) {
	contract := contract2.Contract{
		Address: c.Bech32,
		Signer:  c.Wallet,
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
