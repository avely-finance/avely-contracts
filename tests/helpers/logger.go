package helpers

import (
	sdk "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/sirupsen/logrus"
)

var log *sdk.Log

func init() {
	log = sdk.NewLog()
	addColorIcons(log)
	disableTimestamps(log)
}

func GetLog() *sdk.Log {
	return log
}

type ColorIconHook struct {
}

func (hook *ColorIconHook) Fire(entry *logrus.Entry) error {
	switch entry.Level {
	case logrus.DebugLevel:
		entry.Message = "⚪️ " + entry.Message
	case logrus.InfoLevel:
		entry.Message = "🟢 " + entry.Message
	case logrus.FatalLevel:
		entry.Message = "💔 " + entry.Message
	case logrus.ErrorLevel:
		entry.Message = "🔴 " + entry.Message
	}
	return nil
}

func (hook *ColorIconHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.FatalLevel, logrus.ErrorLevel, logrus.InfoLevel, logrus.DebugLevel}
}

func addColorIcons(log *sdk.Log) {
	colorIconHook := new(ColorIconHook)
	log.AddHook(colorIconHook)
}

func disableTimestamps(log *sdk.Log) {
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
}
