package helpers

import (
	sdk "github.com/avely-finance/avely-contracts/sdk/core"
)

var log *sdk.Log

func init() {
	log = sdk.NewLog()
}

func GetLog() *sdk.Log {
	return log
}
