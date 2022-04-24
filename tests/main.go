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
		"stzilssn": config.StZilSsnAddress,
		"addr1":    config.Addr1,
		"addr2":    config.Addr2,
		"addr3":    config.Addr3,
		"admin":    config.Admin,
		"verifier": config.Verifier,
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
