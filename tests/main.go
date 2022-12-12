package main

import (
	"flag"

	. "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	"github.com/avely-finance/avely-contracts/tests/helpers"
	. "github.com/avely-finance/avely-contracts/tests/transitions"
)

const CHAIN = "local"

func main() {
	config := NewConfig(CHAIN)
	celestials := LoadCelestialsFromEnv(CHAIN)
	sdk := NewAvelySDK(*config)
	log := helpers.GetLog()

	tr := InitTransitions(sdk, celestials)

	shortcuts := map[string]string{
		"stzilssn": config.StZilSsnAddress,
		"alice":    utils.GetAddressByWallet(tr.Alice),
		"bob":      utils.GetAddressByWallet(tr.Bob),
		"eve":      utils.GetAddressByWallet(tr.Eve),
		"admin":    utils.GetAddressByWallet(celestials.Admin),
		"verifier": utils.GetAddressByWallet(tr.Verifier),
	}

	log.AddShortcuts(shortcuts)

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
