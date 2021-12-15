package transitions

import (
	"Azil/test/deploy"
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"log"
	"runtime"
	"strconv"
	"strings"
)

const AZIL_SSN_ADDRESS = "0x166862bdd5d76b3a4775d2494820179d582acac5"
const AZIL_SSN_REWARD_SHARE_PERCENT = "50"

const HOLDER_INITIAL_DELEGATE_ZIL = 1000

// owner
const adminKey = "d96e9eb5b782a80ea153c937fa83e5948485fbfc8b7e7c069d7b914dbc350aba"
const admin = "381f4008505e940ad7681ec3468a719060caf796"

const verifierKey = "5430365143ce0154b682301d0ab731897221906a7054bbf5bd83c7663a6cbc40"
const verifier = "10200e3da08ee88729469d6eabc055cb225821e7"

const key1 = "1080d2cca18ace8225354ac021f9977404cee46f1d12e9981af8c36322eac1a4"
const addr1 = "ac941274c3b6a50203cc5e7939b7dad9f32a0c12"
const key2 = "254d9924fc1dcdca44ce92d80255c6a0bb690f867abde80e626fbfef4d357004"
const addr2 = "ec902fe17d90203d0bddd943d97b29576ece3177"
const key3 = "b8fc4e270594d87d3f728d0873a38fb0896ea83bd6f96b4f3c9ff0a29122efe4"
const addr3 = "c2035715831ab100ec42e562ce341b834bed1f4c"
const key4 = "b87f4ba7dcd6e60f2cca8352c89904e3993c5b2b0b608d255002edcda6374de4"
const addr4 = "6cd3667ba79310837e33f0aecbe13688a6cbca32"

const qa = "000000000000"

func zil(amount int) string {
	if amount == 0 {
		return "0"
	}
	return fmt.Sprintf("%d%s", amount, qa)
}

func azil(amount int) string {
	if amount == 0 {
		return "0"
	}
	return fmt.Sprintf("%d%s", amount, qa)
}

type Testing struct {
	debug map[string]string
}

func NewTesting() *Testing {
	return &Testing{
		debug: make(map[string]string),
	}
}

func (t *Testing) LogStart(tag string) {
	log.Printf("⚙️  === Start to test %s === \n", tag)
}

func (t *Testing) LogEnd() {
	log.Println("🏁 TESTS PASSED SUCCESSFULLY")
}

func (t *Testing) LogError(tag string, err error) {
	log.Fatalf("🔴 Failed at %s, err = %s\n", tag, err.Error())
}

func (t *Testing) GetReceiptString(txn *transaction.Transaction) string {
	receipt, _ := json.Marshal(txn.Receipt)
	return string(receipt)
}

func (t *Testing) LogPrettyReceipt(txn *transaction.Transaction) {
	data, _ := json.MarshalIndent(txn.Receipt, "", "     ")
	log.Println(string(data))
}

func (t *Testing) AssertContain(s1, s2 string) {
	_, file, no, _ := runtime.Caller(1)
	t.AssertContainRaw("ASSERT_CONTAIN", s1, s2, file, no)
}

func (t *Testing) AssertEqual(s1, s2 string) {
	if s1 != s2 {
		_, file, no, _ := runtime.Caller(1)
		log.Println("🔴 ASSERT_EQUAL FAILED, " + file + ":" + strconv.Itoa(no))
		log.Println("🔴 EXPECTED: " + s2)
		log.Println("🔴 ACTUAL: " + s1)
		log.Fatalf("💔 TESTS ARE FAILED")
	} else {
		log.Println("🟢 ASSERT_EQUAL SUCCESS")
	}
}

func (t *Testing) AssertError(txn *transaction.Transaction, err error, code int) {
	_, file, no, _ := runtime.Caller(1)

	if err == nil {
		log.Println("🔴 ASSERT_ERROR FAILED. Tx does not have an issue, " + file + ":" + strconv.Itoa(no))
	}

	tx := t.GetReceiptString(txn)
	errorMessage := fmt.Sprintf("Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 %d))])", code)

	t.AssertContainRaw("ASSERT_ERROR", tx, errorMessage, file, no)
}

func (t *Testing) AssertContainRaw(code, s1, s2, file string, no int) {
	if !strings.Contains(s1, s2) {
		log.Println("🔴 " + code + " FAILED, " + file + ":" + strconv.Itoa(no))
		log.Println(s1)
		log.Println(s2)
		log.Fatalf("💔 TESTS ARE FAILED")
	} else {
		log.Println("🟢 " + code + " SUCCESS")
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
func (t *Testing) AssertTransition(txn *transaction.Transaction, expectedTxn deploy.Transition) {
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
		log.Println("🟢 ASSERT_TRANSITION SUCCESS")
	} else {
		_, file, no, _ := runtime.Caller(1)
		log.Println("🔴 ASSERT_TRANSITION FAILED, " + file + ":" + strconv.Itoa(no))
		actual, _ := json.MarshalIndent(txn, "", "     ")
		expected, _ := json.MarshalIndent(expectedTxn, "", "     ")
		log.Println(fmt.Sprintf("Expected: %s", expected))
		log.Println(fmt.Sprintf("Actual: %s", actual))
		t.LogDebug()
		log.Fatalf("💔 TESTS ARE FAILED")
	}
}

func (t *Testing) AssertEvent(txn *transaction.Transaction, wantedEvent deploy.Event) {
	found := false
	if txn.Receipt.EventLogs != nil {
		for _, el := range txn.Receipt.EventLogs {
			txEvent := convertEventLog(el)
			if txEvent.Address == "0x"+wantedEvent.Sender &&
				txEvent.EventName == wantedEvent.EventName &&
				compareParams(txEvent.Params, convertParams(wantedEvent.Params)) {
				found = true
				break
			}
		}
	}

	if found {
		log.Println("🟢 ASSERT_EVENT SUCCESS")
	} else {
		_, file, no, _ := runtime.Caller(1)
		log.Println("🔴 ASSERT_EVENT FAILED, " + file + ":" + strconv.Itoa(no))
		z, _ := json.Marshal(wantedEvent)
		log.Println(fmt.Sprintf("We assert: %s", z))
		z, _ = json.Marshal(txn.Receipt.EventLogs)
		log.Println(fmt.Sprintf("We have: %s", z))
		t.LogDebug()
		log.Fatalf("💔 TESTS ARE FAILED")
	}
}

func (t *Testing) AddDebug(key, value string) {
	t.debug[key] = value
}

func (t *Testing) LogDebug() {
	var output string
	for key, val := range t.debug {
		output += fmt.Sprintf("% 20s: %s\n", key, val)
	}
	log.Println("Debug info\n" + strings.Trim(output, "\n"))
}

func convertParams(pmap deploy.ParamsMap) []core.ContractValue {
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

func convertEventLog(el interface{}) deploy.EventLog {
	//TODO: correct way to get txEvent EventLog structure
	b, err := json.Marshal(el)
	if err != nil {
		log.Fatal(err)
	}
	var txEvent deploy.EventLog // it's strange, but undefined: core.EventLog
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
