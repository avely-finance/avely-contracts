package core

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strings"

	logrus "github.com/sirupsen/logrus"

	transaction2 "github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/fatih/color"
)

type Log struct {
	shortcuts map[string]string
}

func NewLog() *Log {
	log := &Log{
		shortcuts: make(map[string]string),
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
	return log
}

func (mylog *Log) SetOutputStdout() {
	logrus.SetOutput(os.Stdout)
}

func (mylog *Log) Info(v ...interface{}) {
	logrus.Info(mylog.nice(v)...)
}

func (mylog *Log) Infof(format string, v ...interface{}) {
	logrus.Infof(format, mylog.nice(v)...)
}

func (mylog *Log) Success(v ...interface{}) {
	out := []interface{}{"🟢"}
	logrus.Info(append(out, mylog.nice(v)...)...)
}

func (mylog *Log) Successf(format string, v ...interface{}) {
	logrus.Infof("🟢 "+format, mylog.nice(v)...)
}

func (mylog *Log) Error(v ...interface{}) {
	out := []interface{}{"🔴"}
	logrus.Error(append(out, mylog.nice(v)...)...)
}

func (mylog *Log) Errorf(format string, v ...interface{}) {
	logrus.Errorf("🔴 "+format, mylog.nice(v)...)
}

func (mylog *Log) Fatal(v ...interface{}) {
	out := []interface{}{"💔"}
	logrus.Fatal(append(out, mylog.nice(v)...)...)
}

func (mylog *Log) Fatalf(format string, v ...interface{}) {
	logrus.Fatalf("💔 "+format, mylog.nice(v)...)
}

func (mylog *Log) nice(params []interface{}) []interface{} {
	for i, value := range params {
		if value == nil {
			continue
		}
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
		if mylog.shortcuts[k] == "" {
			continue
		}
		colorFunc := color.New(colors[i%length]).SprintFunc()
		replacement := colorFunc(strings.ToUpper(k) + " " + mylog.shortcuts[k])
		str = strings.ReplaceAll(str, mylog.shortcuts[k], replacement)
		i++
	}
	return str
}
