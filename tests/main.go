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

func main() {

	helpers.Blockchain.ApiUrl = "http://zilliqa_server:5555"

	go increaseBlocknum()
	tr := transitions.NewTransitions()
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
	////////////////////helpers.LogEnd()
}
