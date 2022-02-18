package core

import (
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/johntdyer/slackrus"
	"github.com/sirupsen/logrus"

	"github.com/fatih/color"
)

type Log struct {
	*logrus.Logger
	shortcuts map[string]string
}

func NewLog() *Log {
	shortcuts := make(map[string]string)
	lgr := logrus.New()
	lgr.SetLevel(logrus.DebugLevel)
	log := &Log{
		lgr,
		shortcuts,
	}
	lgr.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	colorAddressHook := &ColorAddressHook{log: log}
	lgr.AddHook(colorAddressHook)

	return log
}

func (mylog *Log) AddSlackHook(hookUrl, level string) {
	_, err := url.ParseRequestURI(hookUrl)
	if err != nil {
		return
	}

	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrusLevel = logrus.ErrorLevel
	}
	mylog.AddHook(&slackrus.SlackrusHook{
		HookURL:        hookUrl,
		AcceptedLevels: slackrus.LevelThreshold(logrusLevel),
		//Channel:        " #watcher-mainnet",
		//IconEmoji:      ":ghost:",
		//Username:       "Watcher",
	})

}

func (mylog *Log) SetOutputStdout() {
	mylog.SetOutput(os.Stdout)
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

type ColorAddressHook struct {
	log *Log
}

func (hook *ColorAddressHook) Fire(entry *logrus.Entry) error {
	entry.Message = hook.log.highlightShortcuts(entry.Message)
	return nil
}

func (hook *ColorAddressHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
