package main

import (
	// "Azil/test/helpers"
	// "Azil/test/transitions"
	"github.com/avely-finance/avely-contracts/tests/helpers"
	. "github.com/avely-finance/avely-contracts/tests/transitions"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"flag"
	"time"
)

func increaseBlocknum(sdk *AvelySDK) {
	for {
		time.Sleep(10 * time.Second)
		sdk.IncreaseBlocknum(1)
	}
}

const CHAIN = "local"

// var t *Testing
// var log *Log
// var sdk *AvelySDK

// func init() {
// 	t = NewTesting()
// 	log = GetLog()
// }

func main() {
	config := NewConfig(CHAIN)
	sdk := NewAvelySDK(*config)
	log := helpers.GetLog()

	testing := helpers.NewTesting(log)

	// config := helpers.LoadConfig(CHAIN)

	// log := GetLog()
	shortcuts := map[string]string{
		"azilssn":  config.AzilSsnAddress,
		"addr1":    "0x" + config.Addr1,
		"addr2":    "0x" + config.Addr2,
		"addr3":    "0x" + config.Addr3,
		"admin":    "0x" + config.Admin,
		"verifier": "0x" + config.Verifier,
	}
	log.AddShortcuts(shortcuts)

	// helpers.Blockchain.ApiUrl = config.ApiUrl

	go increaseBlocknum(sdk)
	tr := InitTransitions(sdk, testing)

	// Example: go run main.go -focus=DelegateStakeSuccess
	focusPtr := flag.String("focus", "default", "a focus test suite")

	flag.Parse()

	focus := string(*focusPtr)

	if focus != "default" {
		log.Info("🏁 Focus on " + focus)
		tr.FocusOn(focus)
	} else {
		log.Info("🏁 Run All Suites ")
		tr.RunAll()
	}

	log.Info("🏁 TESTS PASSED SUCCESSFULLY")
}
