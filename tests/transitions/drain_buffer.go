package transitions

import (
	. "Azil/test/helpers"
)

func (tr *Transitions) DrainBuffer() {
	t.Start("CompleteWithdrawal - success")

	Zproxy, Zimpl, Aimpl, Buffer, Holder := tr.DeployAndUpgrade()

	t.AssertSuccess(Aimpl.DelegateStake(Zil(10)))

	txn, err := Aimpl.DrainBuffer(Aimpl.Addr)
	t.AssertError(txn, err, -107)

	//we need wait 2 reward cycles, in order to pass AssertNoBufferedDepositLessOneCycle, AssertNoBufferedDeposit checks
	Zproxy.UpdateWallet(tr.cfg.VerifierKey)
	IncreaseBlocknum(10)
	t.AssertSuccess(Zproxy.AssignStakeReward(tr.cfg.AzilSsnAddress, tr.cfg.AzilSsnRewardShare))
	IncreaseBlocknum(10)
	t.AssertSuccess(Zproxy.AssignStakeReward(tr.cfg.AzilSsnAddress, tr.cfg.AzilSsnRewardShare))

	txn, _ = Aimpl.DrainBuffer(Buffer.Addr)

	t.AssertTransition(txn, Transition{
		Aimpl.Addr,     //sender
		"ClaimRewards", //tag
		Buffer.Addr,    //recipient
		"0",            //amount
		ParamsMap{},
	})

	// ssnlist#UpdateStakeReward has complex logic based on a fee and comission calculations
	// since we use extra small numbers (not QA 10 ^ 12) all calculations are rounded
	// and all assigned rewards are credited to one SSN node
	bufferRewards := StrAdd(tr.cfg.AzilSsnRewardShare, tr.cfg.AzilSsnRewardShare)
	t.AssertEqual(bufferRewards, "100")

	t.AssertTransition(txn, Transition{
		Zimpl.Addr, //sender
		"AddFunds",
		Buffer.Addr,
		bufferRewards,
		ParamsMap{},
	})

	t.AssertTransition(txn, Transition{
		Zimpl.Addr, //sender
		"WithdrawStakeRewardsSuccessCallBack",
		Buffer.Addr,
		"0",
		ParamsMap{"rewards": bufferRewards},
	})

	// Holder rewards for initial funds
	holderRewards := "49"
	t.AssertTransition(txn, Transition{
		Zimpl.Addr, //sender
		"AddFunds",
		Holder.Addr,
		holderRewards,
		ParamsMap{},
	})

	t.AssertTransition(txn, Transition{
		Zimpl.Addr, //sender
		"WithdrawStakeRewardsSuccessCallBack",
		Holder.Addr,
		"0",
		ParamsMap{"rewards": holderRewards},
	})

	// Check aZIL balance
	totalRewards := "149" // "100" from Buffer + "49" from Holder[]
	t.AssertEqual(Aimpl.Field("_balance"), totalRewards)
	t.AssertEqual(Aimpl.Field("autorestakeamount"), totalRewards)

	// Send Swap transactions
	t.AssertTransition(txn, Transition{
		Buffer.Addr, //sender
		"RequestDelegatorSwap",
		Zproxy.Addr,
		"0",
		ParamsMap{"new_deleg_addr": "0x" + Holder.Addr},
	})

	t.AssertTransition(txn, Transition{
		Holder.Addr, //sender
		"ConfirmDelegatorSwap",
		Zproxy.Addr,
		"0",
		ParamsMap{"requestor": "0x" + Buffer.Addr},
	})

	//try to drain buffer, not existent at main staking contract
	//error should not be thrown
	new_buffers := []string{"0x0000000000000000000000000000000000000000"}
	t.AssertSuccess(Aimpl.ChangeBuffers(new_buffers))
	txn, _ = Aimpl.DrainBuffer("0000000000000000000000000000000000000000")
	t.AssertTransition(txn, Transition{
		Aimpl.Addr, //sender
		"ClaimRewards",
		Holder.Addr,
		"0",
		ParamsMap{},
	})
}
