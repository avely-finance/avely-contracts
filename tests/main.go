package main

import (
	"Azil/test/helpers"
	"Azil/test/transitions"
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
		"addr4":    "0x" + config.Addr4,
		"admin":    "0x" + config.Admin,
		"verifier": "0x" + config.Verifier,
	}
	log.AddShortcuts(shortcuts)

	helpers.Blockchain.ApiUrl = config.ApiUrl

	go increaseBlocknum()
	tr := transitions.NewTransitions(config)
	tr.DelegateStakeSuccess()
	tr.DelegateStakeBuffersRotation()
	tr.WithdrawStakeAmount()
	tr.CompleteWithdrawalSuccess()
	tr.ZilBalanceOf()
	tr.IsAdmin()
	tr.IsAimpl()
	tr.IsZimpl()
	tr.DrainBuffer()
	tr.PerformAuoRestake()
	log.Info("🏁 TESTS PASSED SUCCESSFULLY")
}
