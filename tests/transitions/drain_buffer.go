package transitions

import (
	. "Azil/test/helpers"
)

func (tr *Transitions) DrainBuffer() {
	t.Start("CompleteWithdrawal - success")

	Zproxy, _, Aimpl, Buffer, Holder := tr.DeployAndUpgrade()

	t.AssertSuccess(Aimpl.DelegateStake(zil(10)))

	txn, err := Aimpl.DrainBuffer(Aimpl.Addr)
	t.AssertError(txn, err, -107)

	//we need wait 2 reward cycles, in order to pass AssertNoBufferedDepositLessOneCycle, AssertNoBufferedDeposit checks
	Zproxy.UpdateWallet(tr.cfg.VerifierKey)
	IncreaseBlocknum(10)
	t.AssertSuccess(Zproxy.AssignStakeReward(tr.cfg.AzilSsnAddress, tr.cfg.AzilSsnRewardSharePercent))
	IncreaseBlocknum(10)
	t.AssertSuccess(Zproxy.AssignStakeReward(tr.cfg.AzilSsnAddress, tr.cfg.AzilSsnRewardSharePercent))

	txn, _ = Aimpl.DrainBuffer(Buffer.Addr)

	t.AssertTransition(txn, Transition{
		Aimpl.Addr,     //sender
		"ClaimRewards", //tag
		Buffer.Addr,    //recipient
		"0",            //amount
		ParamsMap{},
	})

	/*
		//In order to check rewards we shoul repeat reward calculation logic from procedure CalcStakeRewards
			// Send funds and call a callback via Buffer
			t.AssertTransition(txn, Transition{
				Zimpl.Addr, //sender
				"AddFunds",
				Buffer.Addr,
				zil(1),
				ParamsMap{},
			})

			t.AssertTransition(txn, Transition{
				Zimpl.Addr, //sender
				"WithdrawStakeRewardsSuccessCallBack",
				Buffer.Addr,
				"0",
				ParamsMap{"rewards": zil(1)},
			})

			// Send funds and call a callback via Holder
			t.AssertTransition(txn, Transition{
				Zimpl.Addr, //sender
				"AddFunds",
				Holder.Addr,
				zil(1),
				ParamsMap{},
			})

			t.AssertTransition(txn, Transition{
				Zimpl.Addr, //sender
				"WithdrawStakeRewardsSuccessCallBack",
				Holder.Addr,
				"0",
				ParamsMap{"rewards": zil(1)},
			})

			// Check aZIL balance
			// 1 ZIL from Buffer + 1 ZIL from Holder
			t.AssertEqual(Aimpl.Field("_balance"), zil(2))
			t.AssertEqual(Aimpl.Field("autorestakeamount"), zil(2))
	*/

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
