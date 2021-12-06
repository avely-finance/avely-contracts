package main

import (
	"Azil/test/deploy"
	"Azil/test/transitions"
	"time"
)

func increaseBlocknum() {
	for {
		time.Sleep(10 * time.Second)
		deploy.IncreaseBlocknum(1)
	}
}

func main() {
	go increaseBlocknum()
	t := transitions.NewTesting()
	// t.DelegateStakeSuccess()
	// t.DelegateStakeBuffersRotation()
	// t.WithdrawStakeAmount()
	// t.CompleteWithdrawalSuccess()
	// t.ZilBalanceOf()
	// t.IsAdmin()
	// t.IsAimpl()
	t.DrainBuffer()
	t.LogEnd()
}
