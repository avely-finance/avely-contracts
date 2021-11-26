package main

import "Azil/test/transitions"

func main() {
	t := transitions.NewTesting()
	t.DelegateStakeSuccess()
	t.DelegateStakeBuffersRotation()
	t.WithdrawStakeAmount()
	t.ZilBalanceOf()
	t.LogEnd()
}
