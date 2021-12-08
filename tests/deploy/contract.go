package deploy

import (
	"encoding/json"
	"errors"
	"log"
	//"runtime"
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

const TxConfirmMaxAttempts = 5
const TxConfirmInterval = 0

type StateMap map[string]interface{}

type Contract struct {
	Code            string
	Init            []core.ContractValue
	Addr            string
	Bech32          string
	Wallet          *account.Wallet
	TxIdLast        string
	TxIdStateParsed string
	StateMap        StateMap
}

type ParamsMap map[string]string
type Transition struct {
	Sender    string
	Tag       string
	Recipient string
	Amount    string
	Params    ParamsMap
}
type Event struct {
	Sender    string
	EventName string
	Params    ParamsMap
}

//replacement for core.EventLog, because of strange "undefined type" error
//we have https://github.com/Zilliqa/gozilliqa-sdk/blob/master/core/types.go#L107
type EventLog struct {
	EventName string               `json:"_eventname"`
	Address   string               `json:"address"`
	Params    []core.ContractValue `json:"params"`
}

func (c *Contract) LogContractStateJson() string {
	provider := provider2.NewProvider("http://zilliqa_server:5555")
	rsp, _ := provider.GetSmartContractState(c.Addr)
	j, _ := json.Marshal(rsp)
	//c.LogPrettyStateJson(rsp)
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
	c.TxIdLast = tx.ID
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

func (c *Contract) StateFieldArray(v interface{}) map[string]interface{} {
	tmp, _ := json.Marshal(v)
	var field []string
	json.Unmarshal([]byte(tmp), &field)
	res := make(map[string]interface{})
	for i, w := range field {
		res[strconv.Itoa(i)] = w
	}
	return res
}

func (c *Contract) StateFieldMap(v interface{}) map[string]interface{} {
	tmp, _ := json.Marshal(v)
	var field map[string]interface{}
	json.Unmarshal([]byte(tmp), &field)
	return field
}

func (c *Contract) StateFieldMapWithdrawal(v interface{}) map[string]interface{} {
	tmp, _ := json.Marshal(v)
	var field map[string]Withdrawal
	json.Unmarshal([]byte(tmp), &field)
	res := make(map[string]interface{})
	for i, w := range field {
		res[string(i)] = w.Arguments[1] //0=>token,1=>stake
	}
	return res
}

func (c *Contract) StateFieldMapMapWithdrawal(v interface{}) map[string]interface{} {
	tmp, _ := json.Marshal(v)
	var field map[string](map[string]Withdrawal)
	json.Unmarshal([]byte(tmp), &field)
	res := make(map[string]interface{})
	for i, w := range field {
		tmpmap := make(map[string]interface{})
		for ii, ww := range w {
			tmpmap[string(ii)] = ww.Arguments[1] //0=>token,1=>stake
		}
		res[string(i)] = tmpmap
	}
	return res
}
