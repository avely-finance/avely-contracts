package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/v3/core"
	"github.com/Zilliqa/gozilliqa-sdk/v3/transaction"
	sdk "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/ethereum/go-ethereum/core/types"
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

func AssertSuccessAny(tx interface{}, err error) (interface{}, error) {
	if err != nil {
		_, file, no, _ := runtime.Caller(1)
		GetLog().Error(tx)
		GetLog().Fatal("TRANSACTION FAILED, " + file + ":" + strconv.Itoa(no))
	}
	return tx, err
}

func AssertError(tx interface{}, code string) {
	var txError string
	var file string
	var no int
	var title string

	_, file, no, _ = runtime.Caller(1)

	if txn, ok := tx.(*transaction.Transaction); ok {
		title = "ASSERT_ERROR"
		if txn.Receipt.Success && txn.Status == core.Confirmed {
			GetLog().Error("ASSERT_ERROR FAILED. Tx does not have an issue, " + file + ":" + strconv.Itoa(no))
		}
		receipt, _ := json.Marshal(txn.Receipt)
		txError = string(receipt)
	} else if txEvm, ok := tx.(*types.Transaction); ok {
		title = "ASSERT_ERROR_EVM"
		receipt, err := _sdk.Evm.Client.TransactionReceipt(context.Background(), txEvm.Hash())
		if err != nil {
			log.Fatal(err)
		}
		for _, log := range receipt.Logs {
			res, found := _sdk.Evm.DecodeScillaEvent(log)
			if found {
				txError += res.Text + "\n"
			}
		}
		txError = strings.ReplaceAll(txError, `"`, `\"`)
	}
	errorMessage := fmt.Sprintf("Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 %s))])", code)
	AssertContainRaw(title, txError, errorMessage, file, no)
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

func AssertContainRaw(title, s1, s2, file string, no int) {
	if !strings.Contains(s1, s2) {
		GetLog().Error(title + " FAILED, " + file + ":" + strconv.Itoa(no))
		GetLog().Error("body: " + s1)
		GetLog().Error("pattern: " + s2)
		GetLog().Fatal("TESTS ARE FAILED")
	} else {
		GetLog().Info(title)
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
func AssertTransition(txAny interface{}, expectedTxn Transition) {
	if _, ok := txAny.(*types.Transaction); ok {
		// see https://github.com/Zilliqa/Zilliqa/issues/3924
		// [FEATURE REQUEST] add Scilla transitions to the EVM transaction receipt #3924
		GetLog().Debug("ASSERT_TRANSITION_EVM NOT SUPPORTED")
		return
	}

	txn, ok := txAny.(*transaction.Transaction)
	if !ok {
		GetLog().Fatal("Transaction type not supported")
	}

	if txn.Receipt.Transitions != nil {
		for _, txTransition := range txn.Receipt.Transitions {
			if txTransition.Addr == expectedTxn.Sender &&
				txTransition.Msg.Recipient == expectedTxn.Recipient &&
				txTransition.Msg.Tag == expectedTxn.Tag &&
				txTransition.Msg.Amount == expectedTxn.Amount &&
				compareParams(txTransition.Msg.Params, convertParams(expectedTxn.Params)) {
				GetLog().Info("ASSERT_TRANSITION SUCCESS")
				return
			}
		}
	}
	_, file, no, _ := runtime.Caller(1)
	GetLog().Error("ASSERT_TRANSITION FAILED, " + file + ":" + strconv.Itoa(no))
	actual, _ := json.MarshalIndent(txn, "", "     ")
	expected, _ := json.MarshalIndent(expectedTxn, "", "     ")
	GetLog().Error(fmt.Sprintf("Expected: %s", expected))
	GetLog().Error(fmt.Sprintf("Actual: %s", actual))
	GetLog().Fatal("TESTS ARE FAILED")
}

func AssertEvent(tx interface{}, expectedEvent Event) {
	title := ""

	var eventLogs []*EventLog

	if txn, ok := tx.(*transaction.Transaction); ok {
		title = "ASSERT_EVENT"
		if txn.Receipt.EventLogs != nil {
			for _, el := range txn.Receipt.EventLogs {
				txEvent := convertEventLog(el, GetLog())
				eventLogs = append(eventLogs, &txEvent)

				if txEvent.Address == expectedEvent.Sender &&
					txEvent.EventName == expectedEvent.EventName &&
					compareParams(txEvent.Params, convertParams(expectedEvent.Params)) {
					GetLog().Info(title + " SUCCESS")
					return
				}
			}
		}
	} else if txn, ok := tx.(*types.Transaction); ok {
		title = "ASSERT_EVENT_EVM"
		receipt, err := _sdk.Evm.Client.TransactionReceipt(context.Background(), txn.Hash())
		if err != nil {
			GetLog().Fatal(err)
		}
		for _, logEntry := range receipt.Logs {
			res, evtFound := _sdk.Evm.DecodeScillaEvent(logEntry)
			if evtFound && res.Kind == sdk.ForwardedScillaEvent {
				var txEvent EventLog
				err = json.Unmarshal([]byte(res.Text), &txEvent)
				if err != nil {
					GetLog().Fatal(err)
				}
				eventLogs = append(eventLogs, &txEvent)

				if txEvent.Address == expectedEvent.Sender &&
					txEvent.EventName == expectedEvent.EventName &&
					compareParams(txEvent.Params, convertParams(expectedEvent.Params)) {
					GetLog().Info(title + " SUCCESS")
					return
				}
			}
		}
	} else {
		GetLog().Fatal("Unknown transaction type")
	}

	_, file, no, _ := runtime.Caller(1)
	GetLog().Error(title + " FAILED, " + file + ":" + strconv.Itoa(no))
	expected, _ := json.Marshal(expectedEvent)
	GetLog().Error(fmt.Sprintf("EXPECTED: %s", expected))
	actual, _ := json.Marshal(eventLogs)
	GetLog().Error(fmt.Sprintf("ACTUAL: %s", actual))
	GetLog().Fatal("TESTS ARE FAILED")
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
