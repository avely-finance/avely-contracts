package main

import (
	"Azil/test/helpers"
	"Azil/test/transitions"
	"flag"
	"time"
)

func increaseBlocknum() {
	for {
		time.Sleep(10 * time.Second)
		helpers.IncreaseBlocknum(1)
	}
}

const CHAIN = "local"

func main() {
	config := helpers.LoadConfig(CHAIN)

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

	helpers.Blockchain.ApiUrl = config.ApiUrl

	go increaseBlocknum()
	tr := transitions.NewTransitions(config)

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
