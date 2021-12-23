package helpers

import (
	"github.com/avely-finance/avely-contracts/sdk"
)

var log *sdk.Log

func init() {
	log = sdk.NewLog()
}

func GetLog() *sdk.Log {
	return log
}
