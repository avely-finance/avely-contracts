package main

import (
	"flag"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/tests/helpers"
	. "github.com/avely-finance/avely-contracts/tests/transitions"
)

const CHAIN = "local"

func main() {
	config := NewConfig(CHAIN)
	sdk := NewAvelySDK(*config)
	log := helpers.GetLog()

	shortcuts := map[string]string{
		"azilssn":  config.AzilSsnAddress,
		"addr1":    "0x" + config.Addr1,
		"addr2":    "0x" + config.Addr2,
		"addr3":    "0x" + config.Addr3,
		"admin":    "0x" + config.Admin,
		"verifier": "0x" + config.Verifier,
	}

	log.AddShortcuts(shortcuts)

	tr := InitTransitions(sdk)

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
