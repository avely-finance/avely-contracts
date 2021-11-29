package deploy

import (
	"encoding/json"
	"errors"
	"log"
	"runtime"
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

const TxConfirmMaxAttempts = 5
const TxConfirmInterval = 0

type Contract struct {
	Code   string
	Init   []core.ContractValue
	Addr   string
	Bech32 string
	Wallet *account.Wallet
}

type MyEventParamsMap map[string]string
type MyEventLog struct {
	EventName string            `json:"_eventname"`
	Address   string            `json:"address"`
	Params    []MyContractValue `json:"params"`
}
type MyContractValue struct {
	VName string `json:"vname"`
	//Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

/* @param params String {"foo": "bar", "test" : "123"} */
func (c *Contract) Event(name string, params string) MyEventLog {

	//string to MyEventParamsMap
	var pmap MyEventParamsMap
	err := json.Unmarshal([]byte(params), &pmap)
	if err != nil {
		_, file, no, _ := runtime.Caller(1)
		log.Println("Can not parse json: " + params + " at " + file + ":" + strconv.Itoa(no))
		log.Fatal(err)
	}

	//transform MyEventParamsMap to array of ContractValue
	cvarr := []MyContractValue{}
	for key, val := range pmap {
		cvarr = append(cvarr, MyContractValue{
			Value: val,
			VName: key,
		})
	}

	return MyEventLog{
		EventName: name,
		Address:   "0x" + c.Addr,
		Params:    cvarr,
	}
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
	tx.Confirm(tx.ID, TxConfirmMaxAttempts, TxConfirmInterval, contract.Provider)
	if tx.Status != core.Confirmed {
		return tx, errors.New("transaction didn't get confirmed")
	}
	if !tx.Receipt.Success {
		return tx, errors.New("transaction failed")
	}
	return tx, nil
}
