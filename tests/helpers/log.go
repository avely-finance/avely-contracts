package helpers

import (
	"encoding/json"
	transaction2 "github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/fatih/color"
	golog "log"
	"reflect"
	"sort"
	"strings"
)

type Log struct {
	shortcuts map[string]string
}

var log *Log

func init() {
	log = &Log{
		shortcuts: make(map[string]string),
	}
}

func GetLog() *Log {
	return log
}

func (mylog *Log) Debug(v ...interface{}) {
	golog.Println(v)
}

func (mylog *Log) Info(v ...interface{}) {
	v = mylog.nice(v)
	v = append([]interface{}{"🔵"}, v...)
	golog.Println(v)
}

func (mylog *Log) Error(v ...interface{}) {
	v = mylog.nice(v)
	v = append([]interface{}{"🔴"}, v...)
	golog.Println(v)
}

func (mylog *Log) Success(v ...interface{}) {
	v = mylog.nice(v)
	v = append([]interface{}{"🟢"}, v...)
	golog.Println(v)
}

func (mylog *Log) Fatal(v ...interface{}) {
	v = mylog.nice(v)
	v = append([]interface{}{"💔"}, v...)
	golog.Fatal(v)
}

func (mylog *Log) Start(tag string) {
	golog.Printf("⚙️  === Start to test %s === \n", tag)
}

func (mylog *Log) End() {
	golog.Println("🏁 TESTS PASSED SUCCESSFULLY")
}

func (mylog *Log) nice(params []interface{}) []interface{} {
	for i, value := range params {
		typ := reflect.ValueOf(value).Type().String()
		switch typ {
		case "string":
			params[i] = mylog.highlightShortcuts(value.(string))
			break
		case "*transaction.Transaction":
			receipt, _ := json.MarshalIndent(value.(*transaction2.Transaction).Receipt, "", "     ")
			params[i] = mylog.highlightShortcuts(string(receipt))
			break
		default:
			break
		}
	}
	return params
}

func (mylog *Log) AddShortcut(key, value string) {
	mylog.shortcuts[key] = value
}

func (mylog *Log) AddShortcuts(shortcuts map[string]string) {
	for key, val := range shortcuts {
		mylog.AddShortcut(key, val)
	}
}

func (mylog *Log) highlightShortcuts(str string) string {

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
	keys := make([]string, 0, len(mylog.shortcuts))
	for k := range mylog.shortcuts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	length := len(colors)
	i := 0
	for _, k := range keys {
		colorFunc := color.New(colors[i%length]).SprintFunc()
		replacement := colorFunc(strings.ToUpper(k) + " " + mylog.shortcuts[k])
		str = strings.ReplaceAll(str, mylog.shortcuts[k], replacement)
		i++
	}
	return str
}
