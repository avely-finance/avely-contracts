package helpers

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/v3/core"
	"github.com/Zilliqa/gozilliqa-sdk/v3/transaction"
	sdk "github.com/avely-finance/avely-contracts/sdk/core"
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

// replacement for core.EventLog, because of strange "undefined type" error
// we have https://github.com/Zilliqa/gozilliqa-sdk/v3/blob/master/core/types.go#L107
type EventLog struct {
	EventName string               `json:"_eventname"`
	Address   string               `json:"address"`
	Params    []core.ContractValue `json:"params"`
}

const walletErrorCode = "WalletError"

func Start(tag string) {
	GetLog().Infof("⚙️ === Start to test %s ===", tag)
}

func AssertContain(s1, s2 string) {
	_, file, no, _ := runtime.Caller(1)
	AssertContainRaw("ASSERT_CONTAIN", s1, s2, file, no)
}

func AssertEqual(s1, s2 string) {
	if s1 != s2 {
		_, file, no, _ := runtime.Caller(1)
		GetLog().Error("ASSERT_EQUAL FAILED, " + file + ":" + strconv.Itoa(no))
		GetLog().Error("EXPECTED: " + s2)
		GetLog().Error("ACTUAL: " + s1)
		GetLog().Fatal("TESTS ARE FAILED")
	} else {
		GetLog().Info("ASSERT_EQUAL SUCCESS")
	}
}

func AssertNotEqual(s1, s2 string) {
	if s1 == s2 {
		_, file, no, _ := runtime.Caller(1)
		GetLog().Error("ASSERT_NOT_EQUAL FAILED, " + file + ":" + strconv.Itoa(no))
		GetLog().Error("EXPECTED: " + s2)
		GetLog().Error("ACTUAL: " + s1)
		GetLog().Fatal("TESTS ARE FAILED")
	} else {
		GetLog().Info("ASSERT_NOT_EQUAL SUCCESS")
	}
}

func AssertSuccess(tx *transaction.Transaction, err error) (*transaction.Transaction, error) {
	if err != nil {
		_, file, no, _ := runtime.Caller(1)
		GetLog().Error(tx)
		GetLog().Fatal("TRANSACTION FAILED, " + file + ":" + strconv.Itoa(no))
	}
	return tx, err
}

func AssertError(txn *transaction.Transaction, code string) {
	_, file, no, _ := runtime.Caller(1)

	if txn.Receipt.Success && txn.Status == core.Confirmed {
		GetLog().Error("ASSERT_ERROR FAILED. Tx does not have an issue, " + file + ":" + strconv.Itoa(no))
	}

	receipt, _ := json.Marshal(txn.Receipt)
	txError := string(receipt)
	errorMessage := fmt.Sprintf("Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 %s))])", code)
	AssertContainRaw("ASSERT_ERROR", txError, errorMessage, file, no)
}

func AssertMultisigSuccess(txn *transaction.Transaction, err error) (*transaction.Transaction, error) {
	receipt, _ := json.Marshal(txn.Receipt)
	txnData := string(receipt)

	if err != nil || strings.Contains(txnData, walletErrorCode) {
		_, file, no, _ := runtime.Caller(1)
		GetLog().Error(txn)
		GetLog().Error(err)
		GetLog().Fatal("ASSERT_MULTISIG_SUCCESS FAILED, " + file + ":" + strconv.Itoa(no))
	}
	return txn, err
}

func AssertMultisigError(txn *transaction.Transaction, code string) {
	_, file, no, _ := runtime.Caller(1)

	receipt, _ := json.Marshal(txn.Receipt)
	txError := string(receipt)

	AssertContainRaw("ASSERT_ERROR", txError, walletErrorCode, file, no)
	AssertContainRaw("ASSERT_ERROR", txError, code, file, no)
}

func AssertZimplError(txn *transaction.Transaction, code string) {
	_, file, no, _ := runtime.Caller(1)

	if txn.Receipt.Success && txn.Status == core.Confirmed {
		GetLog().Error("ASSERT_SSNLIST_ERROR FAILED. Tx does not have an issue, " + file + ":" + strconv.Itoa(no))
	}

	receipt, _ := json.Marshal(txn.Receipt)
	txError := string(receipt)
	errorMessage := fmt.Sprintf("Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 %s))])", code)
	AssertContainRaw("ASSERT_SSNLIST_ERROR", txError, errorMessage, file, no)
}

func AssertASwapError(txn *transaction.Transaction, code string) {
	_, file, no, _ := runtime.Caller(1)

	if txn.Receipt.Success && txn.Status == core.Confirmed {
		GetLog().Error("ASSERT_ASWAP_ERROR FAILED. Tx does not have an issue, " + file + ":" + strconv.Itoa(no))
	}

	receipt, _ := json.Marshal(txn.Receipt)
	txError := string(receipt)
	errorMessage := fmt.Sprintf("(code : (Int32 %s))", code)
	AssertContainRaw("ASSERT_ASWAP_ERROR", txError, errorMessage, file, no)
}

func AssertContainRaw(code, s1, s2, file string, no int) {
	if !strings.Contains(s1, s2) {
		GetLog().Error(code + " FAILED, " + file + ":" + strconv.Itoa(no))
		GetLog().Error(s1)
		GetLog().Error(s2)
		GetLog().Fatal("TESTS ARE FAILED")
	} else {
		GetLog().Info(code)
	}
}

/*
https://github.com/Zilliqa/gozilliqa-sdk/v3/blob/master/core/types.go#L129

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
func AssertTransition(txn *transaction.Transaction, expectedTxn Transition) {
	found := false
	if txn.Receipt.Transitions != nil {
		for _, txTransition := range txn.Receipt.Transitions {
			if txTransition.Addr == expectedTxn.Sender &&
				txTransition.Msg.Recipient == expectedTxn.Recipient &&
				txTransition.Msg.Tag == expectedTxn.Tag &&
				txTransition.Msg.Amount == expectedTxn.Amount &&
				compareParams(txTransition.Msg.Params, convertParams(expectedTxn.Params)) {
				found = true
				break
			}
		}
	}
	if found {
		GetLog().Info("ASSERT_TRANSITION SUCCESS")
	} else {
		_, file, no, _ := runtime.Caller(1)
		GetLog().Error("ASSERT_TRANSITION FAILED, " + file + ":" + strconv.Itoa(no))
		actual, _ := json.MarshalIndent(txn, "", "     ")
		expected, _ := json.MarshalIndent(expectedTxn, "", "     ")
		GetLog().Error(fmt.Sprintf("Expected: %s", expected))
		GetLog().Error(fmt.Sprintf("Actual: %s", actual))
		GetLog().Fatal("TESTS ARE FAILED")
	}
}

func AssertEvent(txn *transaction.Transaction, expectedEvent Event) {
	found := false
	if txn.Receipt.EventLogs != nil {
		for _, el := range txn.Receipt.EventLogs {
			txEvent := convertEventLog(el, GetLog())
			if txEvent.Address == expectedEvent.Sender &&
				txEvent.EventName == expectedEvent.EventName &&
				compareParams(txEvent.Params, convertParams(expectedEvent.Params)) {
				found = true
				break
			}
		}
	}

	if found {
		GetLog().Info("ASSERT_EVENT SUCCESS")
	} else {
		_, file, no, _ := runtime.Caller(1)
		GetLog().Error("ASSERT_EVENT FAILED, " + file + ":" + strconv.Itoa(no))
		expected, _ := json.Marshal(expectedEvent)
		GetLog().Error(fmt.Sprintf("EXPECTED: %s", expected))
		actual, _ := json.Marshal(txn.Receipt.EventLogs)
		GetLog().Error(fmt.Sprintf("ACTUAL: %s", actual))
		GetLog().Fatal("TESTS ARE FAILED")
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

func convertEventLog(el interface{}, log *sdk.Log) EventLog {
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
