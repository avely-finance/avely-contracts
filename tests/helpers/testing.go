package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"runtime"
	"strconv"
	"strings"
)

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

type Testing struct {
}

func NewTesting() *Testing {
	return &Testing{}
}

func (t *Testing) AssertContain(s1, s2 string) {
	_, file, no, _ := runtime.Caller(1)
	t.AssertContainRaw("ASSERT_CONTAIN", s1, s2, file, no)
}

func (t *Testing) AssertEqual(s1, s2 string) {
	if s1 != s2 {
		_, file, no, _ := runtime.Caller(1)
		log.Error("ASSERT_EQUAL FAILED, " + file + ":" + strconv.Itoa(no))
		log.Error("EXPECTED: " + s2)
		log.Error("ACTUAL: " + s1)
		log.Fatal("TESTS ARE FAILED")
	} else {
		log.Success("ASSERT_EQUAL SUCCESS")
	}
}

func (t *Testing) AssertSuccess(tx *transaction.Transaction, err error) (*transaction.Transaction, error) {
	if err != nil {
		_, file, no, _ := runtime.Caller(1)
		log.Error(tx)
		log.Fatal("TRANSACTION FAILED, " + file + ":" + strconv.Itoa(no))
	}
	return tx, err
}

func (t *Testing) AssertError(txn *transaction.Transaction, err error, code int) {
	_, file, no, _ := runtime.Caller(1)

	if err == nil {
		log.Error("ASSERT_ERROR FAILED. Tx does not have an issue, " + file + ":" + strconv.Itoa(no))
	}

	receipt, _ := json.Marshal(txn.Receipt)
	txError := string(receipt)
	errorMessage := fmt.Sprintf("Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 %d))])", code)

	t.AssertContainRaw("ASSERT_ERROR", txError, errorMessage, file, no)
}

func (t *Testing) AssertContainRaw(code, s1, s2, file string, no int) {
	if !strings.Contains(s1, s2) {
		log.Error(code + " FAILED, " + file + ":" + strconv.Itoa(no))
		log.Error(s1)
		log.Error(s2)
		log.Fatal("TESTS ARE FAILED")
	} else {
		log.Success(code)
	}
}

/*
https://github.com/Zilliqa/gozilliqa-sdk/blob/master/core/types.go#L129
type Transition struct {
	Accept bool               `json:"accept"`
	Addr   string             `json:"addr"`
	Depth  int                `json:"depth"`
	Msg    TransactionMessage `json:"msg"`
}

type TransactionMessage struct {
	Amount    string          `json:"_amount"`
	Recipient string          `json:"_recipient"`
	Tag       string          `json:"_tag"`
	Params    []ContractValue `json:"params"`
}
*/
func (t *Testing) AssertTransition(txn *transaction.Transaction, expectedTxn Transition) {
	found := false
	if txn.Receipt.Transitions != nil {
		for _, txTransition := range txn.Receipt.Transitions {
			if txTransition.Addr == "0x"+expectedTxn.Sender &&
				txTransition.Msg.Recipient == "0x"+expectedTxn.Recipient &&
				txTransition.Msg.Tag == expectedTxn.Tag &&
				txTransition.Msg.Amount == expectedTxn.Amount &&
				compareParams(txTransition.Msg.Params, convertParams(expectedTxn.Params)) {
				found = true
				break
			}
		}
	}
	if found {
		log.Success("ASSERT_TRANSITION SUCCESS")
	} else {
		_, file, no, _ := runtime.Caller(1)
		log.Error("ASSERT_TRANSITION FAILED, " + file + ":" + strconv.Itoa(no))
		actual, _ := json.MarshalIndent(txn, "", "     ")
		expected, _ := json.MarshalIndent(expectedTxn, "", "     ")
		log.Error(fmt.Sprintf("Expected: %s", expected))
		log.Error(fmt.Sprintf("Actual: %s", actual))
		log.Fatal("TESTS ARE FAILED")
	}
}

func (t *Testing) AssertEvent(txn *transaction.Transaction, expectedEvent Event) {
	found := false
	if txn.Receipt.EventLogs != nil {
		for _, el := range txn.Receipt.EventLogs {
			txEvent := convertEventLog(el)
			if txEvent.Address == "0x"+expectedEvent.Sender &&
				txEvent.EventName == expectedEvent.EventName &&
				compareParams(txEvent.Params, convertParams(expectedEvent.Params)) {
				found = true
				break
			}
		}
	}

	if found {
		log.Success("ASSERT_EVENT SUCCESS")
	} else {
		_, file, no, _ := runtime.Caller(1)
		log.Error("ASSERT_EVENT FAILED, " + file + ":" + strconv.Itoa(no))
		expected, _ := json.Marshal(expectedEvent)
		log.Error(fmt.Sprintf("EXPECTED: %s", expected))
		actual, _ := json.Marshal(txn.Receipt.EventLogs)
		log.Error(fmt.Sprintf("ACTUAL: %s", actual))
		log.Fatal("TESTS ARE FAILED")
	}
}

func convertParams(pmap ParamsMap) []core.ContractValue {
	cvarr := []core.ContractValue{}
	for key, val := range pmap {
		cvarr = append(cvarr, core.ContractValue{
			Value: val,
			Type:  "foo",
			VName: key,
		})
	}
	return cvarr
}

func convertEventLog(el interface{}) EventLog {
	//TODO: correct way to get txEvent EventLog structure
	b, err := json.Marshal(el)
	if err != nil {
		log.Fatal(err)
	}
	var txEvent EventLog // it's strange, but undefined: core.EventLog
	err = json.Unmarshal([]byte(b), &txEvent)
	if err != nil {
		log.Fatal(err)
	}
	//--
	return txEvent
}

func compareParams(all, wanted []core.ContractValue) bool {
	makeKey := func(cv core.ContractValue) string {
		return cv.VName + "=====" + fmt.Sprintf("%v", cv.Value)
	}
	allMap := make(map[string]bool)
	for _, _map := range all {
		allMap[makeKey(_map)] = true
	}

	//all test event parameters should be present, else events are not matching
	for _, _tmap := range wanted {
		if !allMap[makeKey(_tmap)] {
			return false
		}
	}
	return true
}
