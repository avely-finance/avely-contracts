package contracts

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"

	. "github.com/avely-finance/avely-contracts/sdk/core"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type Pair struct {
	Argtypes    interface{} `json:"argtypes"`
	Arguments   []string    `json:"arguments"`
	Constructor string      `json:"constructor"`
}

type StateMap map[string]interface{}

type StateFieldTypes map[string]string

type Contract struct {
	Sdk             *AvelySDK
	Provider        provider2.Provider
	Addr            string
	Bech32          string
	Wallet          *account.Wallet
	TxIdStateParsed string
	StateMap        StateMap
	StateFieldTypes StateFieldTypes
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

	tx, err := c.Sdk.CallFor(&contract, transition, params, false, amount)
	if err != nil {
		return tx, err
	}
	tx.Confirm(tx.ID, c.Sdk.Cfg.TxConfrimMaxAttempts, c.Sdk.Cfg.TxConfirmIntervalSec, contract.Provider)
	if tx.Status != core.Confirmed {
		return tx, errors.New("transaction didn't get confirmed")
	}
	if !tx.Receipt.Success {
		return tx, errors.New("transaction failed")
	}
	return tx, nil
}

func (c *Contract) Field(key ...string) string {
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

func (c *Contract) State() string {
	rsp, _ := c.Provider.GetSmartContractState(c.Addr)
	result, _ := json.MarshalIndent(rsp.Result, "", "     ")
	state := string(result)
	return state
}

func (c *Contract) stateParse() {
	if c.TxIdStateParsed == c.Sdk.TxIdLast {
		return
	}
	state := c.State()

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
			case "StateFieldOption":
				statemap[k] = stateFieldOption(v)
				break
			case "StateFieldMap":
				statemap[k] = stateFieldMap(v)
				break
			case "StateFieldMapMap":
				statemap[k] = stateFieldMapMap(v)
				break
			case "StateFieldMapMapMap":
				statemap[k] = stateFieldMapMapMap(v)
				break
			case "StateFieldArray":
				statemap[k] = stateFieldArray(v)
				break
			case "StateFieldMapMapPair":
				statemap[k] = stateFieldMapMapPair(v)
				break
			case "StateFieldMapPair":
				statemap[k] = stateFieldMapPair(v)
				break
			default:
				panic("State field type not found: " + typ)
			}
		}
	}
	c.StateMap = statemap
	c.TxIdStateParsed = c.Sdk.TxIdLast
}

func stateFieldOption(v interface{}) string {
	tmp, _ := json.Marshal(v)
	var field map[string]interface{}
	json.Unmarshal([]byte(tmp), &field)

	args := field["arguments"].([]interface{})
	if 0 == len(args) {
		return ""
	}
	return args[0].(string)
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

func stateFieldMapPair(v interface{}) map[string]interface{} {
	tmp, _ := json.Marshal(v)
	var field map[string]Pair
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

func stateFieldMapMapPair(v interface{}) map[string]interface{} {
	tmp, _ := json.Marshal(v)
	var field map[string](map[string]Pair)
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

func stateFieldMap(v interface{}) map[string]interface{} {
	tmp, _ := json.Marshal(v)
	var field map[string]interface{}
	json.Unmarshal([]byte(tmp), &field)
	return field
}

func stateFieldMapMap(v interface{}) map[string]interface{} {
	tmp, _ := json.Marshal(v)
	var field map[string](map[string](map[string]interface{}))
	json.Unmarshal([]byte(tmp), &field)
	res := make(map[string]interface{})
	for i, w := range field {
		mapi := make(map[string]interface{})
		for ii, ww := range w {
			mapi[string(ii)] = ww
		}
		res[string(i)] = mapi
	}
	return res
}

func stateFieldMapMapMap(v interface{}) map[string]interface{} {
	tmp, _ := json.Marshal(v)
	var field map[string](map[string](map[string]interface{}))
	json.Unmarshal([]byte(tmp), &field)
	res := make(map[string]interface{})
	for i, w := range field {
		mapi := make(map[string]interface{})
		for ii, ww := range w {
			mapii := make(map[string]interface{})
			for iii, www := range ww {
				mapii[string(iii)] = www
			}
			mapi[string(ii)] = ww
		}
		res[string(i)] = mapi
	}
	return res
}
