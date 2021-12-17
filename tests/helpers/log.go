package helpers

import (
	"encoding/json"
	transaction2 "github.com/Zilliqa/gozilliqa-sdk/transaction"
	"log"
	"reflect"
)

func (t *Testing) LogStart(tag string) {
	log.Printf("⚙️  === Start to test %s === \n", tag)
}

func (t *Testing) LogEnd() {
	log.Println("🏁 TESTS PASSED SUCCESSFULLY")
}

func (t *Testing) LogError(tag string, err error) {
	log.Fatalf("🔴 Failed at %s, err = %s\n", tag, err.Error())
}

func (t *Testing) LogNice(param interface{}) {
	typ := reflect.ValueOf(param).Type().String()
	txt := ""
	switch typ {
	case "string":
		txt = param.(string)
		break
	case "*transaction.Transaction":
		receipt, _ := json.MarshalIndent(param.(*transaction2.Transaction).Receipt, "", "     ")
		txt = string(receipt)
		break
	default:
		panic("Unknown type " + typ)
		break
	}
	result := t.HighlightShortcuts(txt)
	log.Println(result)
}
