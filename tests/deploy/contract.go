package deploy

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

const TxConfirmMaxAttempts = 5
const TxConfirmInterval = 0

var TxIdLast = ""

type StateMap map[string]interface{}

type StateFieldTypes map[string]string

type Contract struct {
	Code            string
	Init            []core.ContractValue
	Addr            string
	Bech32          string
	Wallet          *account.Wallet
	TxIdStateParsed string
	StateMap        StateMap
	StateFieldTypes StateFieldTypes
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
	TxIdLast = tx.ID
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

func (c *Contract) StateField(key ...string) string {
	c.stateParse()
	src := c.StateMap
	for _, v := range key {
		val, ok := src[v]
		if !ok {
			//key not found in map
			return ""
		} else if reflect.String == reflect.ValueOf(val).Kind() {
			return val.(string)
		} else if reflect.Map == reflect.ValueOf(val).Kind() && 0 == len(val.(map[string]interface{})) {
			//empty map
			return "empty"
		}
		src = val.(map[string]interface{})
	}
	return "map"
}

func (c *Contract) stateParse() {
	if c.TxIdStateParsed == TxIdLast {
		return
	}
	provider := provider2.NewProvider("http://zilliqa_server:5555")
	rsp, _ := provider.GetSmartContractState(c.Addr)
	result, _ := json.Marshal(rsp.Result)
	state := string(result)

	var statemap StateMap
	json.Unmarshal([]byte(state), &statemap)
	for k, v := range statemap {
		typ, ok := c.StateFieldTypes[k]
		if !ok {
			if reflect.String == reflect.ValueOf(v).Kind() {
				statemap[k] = v.(string)
			} else {
				statemap[k] = "not_parsed"
			}
		} else {
			switch typ {
			case "StateFieldMap":
				statemap[k] = stateFieldMap(v)
				break
			case "StateFieldArray":
				statemap[k] = stateFieldArray(v)
				break
			case "StateFieldMapMapWithdrawal":
				statemap[k] = stateFieldMapMapWithdrawal(v)
				break
			case "StateFieldMapWithdrawal":
				statemap[k] = stateFieldMapWithdrawal(v)
				break
			default:
				panic("State field type not found: " + typ)
			}
		}
	}
	log.Println("State parsed after txid=" + TxIdLast)
	c.StateMap = statemap
	c.TxIdStateParsed = TxIdLast
}

func stateFieldArray(v interface{}) map[string]interface{} {
	tmp, _ := json.Marshal(v)
	var field []string
	json.Unmarshal([]byte(tmp), &field)
	res := make(map[string]interface{})
	for i, w := range field {
		res[strconv.Itoa(i)] = w
	}
	return res
}

func stateFieldMap(v interface{}) map[string]interface{} {
	tmp, _ := json.Marshal(v)
	var field map[string]interface{}
	json.Unmarshal([]byte(tmp), &field)
	return field
}

func stateFieldMapWithdrawal(v interface{}) map[string]interface{} {
	tmp, _ := json.Marshal(v)
	var field map[string]Withdrawal
	json.Unmarshal([]byte(tmp), &field)
	res := make(map[string]interface{})
	for i, w := range field {
		inner := make(map[string]interface{})
		inner["0"] = w.Arguments[0] //token
		inner["1"] = w.Arguments[1] //stake
		res[string(i)] = inner
	}
	return res
}

func stateFieldMapMapWithdrawal(v interface{}) map[string]interface{} {
	tmp, _ := json.Marshal(v)
	var field map[string](map[string]Withdrawal)
	json.Unmarshal([]byte(tmp), &field)
	res := make(map[string]interface{})
	for i, w := range field {
		tmpmap := make(map[string]interface{})
		for ii, ww := range w {
			inner := make(map[string]interface{})
			inner["0"] = ww.Arguments[0] //token
			inner["1"] = ww.Arguments[1] //stake
			tmpmap[string(ii)] = inner
		}
		res[string(i)] = tmpmap
	}
	return res
}
