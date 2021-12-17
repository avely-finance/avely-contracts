package helpers

import (
	"encoding/json"
	transaction2 "github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/fatih/color"
	"log"
	"reflect"
	"sort"
	"strings"
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

func (t *Testing) AddShortcut(key, value string) {
	t.shortcuts[key] = value
}

func (t *Testing) HighlightShortcuts(str string) string {

	colors := [...]color.Attribute{
		color.FgRed,
		color.FgGreen,
		color.FgYellow,
		color.FgBlue,
		color.FgMagenta,
		color.FgCyan,
		color.FgHiRed,
		color.FgHiGreen,
		color.FgHiYellow,
		color.FgHiBlue,
		color.FgHiMagenta,
		color.FgHiCyan,
	}

	//sort shortcuts
	keys := make([]string, 0, len(t.shortcuts))
	for k := range t.shortcuts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	l := len(colors)
	i := 0
	for _, k := range keys {
		colorFunc := color.New(colors[i%l]).SprintFunc()
		replacement := colorFunc(strings.ToUpper(k) + " " + t.shortcuts[k])
		str = strings.ReplaceAll(str, t.shortcuts[k], replacement)
		i++
	}
	return str
}
