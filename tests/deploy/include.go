package deploy

import "github.com/Zilliqa/gozilliqa-sdk/core"

//TYPES
type Pair struct {
	Argtypes    interface{} `json:"argtypes"`
	Arguments   []string    `json:"arguments"`
	Constructor string      `json:"constructor"`
}

type StateMap map[string]interface{}

type StateFieldTypes map[string]string

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

//CONSTANTS
const API_PROVIDER string = "http://zilliqa_server:5555"
const TX_CONFIRM_MAX_ATTEMPTS int = 5
const TX_CONFIRM_INTERVAL_SEC int = 0

//GLOBAL VARIABLES
//TODO: what is best practices for shared variable?
var TxIdLast string = ""
