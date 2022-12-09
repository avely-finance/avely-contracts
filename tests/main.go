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

	shortcuts := map[string]string{
		"stzilssn": config.StZilSsnAddress,
		"addr1":    config.Addr1,
		"addr2":    config.Addr2,
		"addr3":    config.Addr3,
		"admin":    utils.GetAddressByWallet(celestials.Admin),
		"verifier": config.Verifier,
	}

	log.AddShortcuts(shortcuts)

	log.Print(celestials)
	tr := InitTransitions(sdk, celestials)

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
