package transitions

import (
	// "log"
	"Azil/test/deploy"
	//"math/big"
)

func (t *Testing) DrainBuffer() {
	t.LogStart("CompleteWithdrawal - success")

	Zproxy, _, Aimpl, Buffer, Holder := t.DeployAndUpgrade()

	t.AssertSuccess(Aimpl.DelegateStake(zil(10)))

	txn, err := Aimpl.DrainBuffer(Aimpl.Addr)
	t.AssertError(txn, err, -107)

	//we need wait 2 reward cycles, in order to pass AssertNoBufferedDepositLessOneCycle, AssertNoBufferedDeposit checks
	Zproxy.UpdateWallet(verifierKey)
	deploy.IncreaseBlocknum(10)
	t.AssertSuccess(Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT))
	deploy.IncreaseBlocknum(10)
	t.AssertSuccess(Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT))

	txn, _ = Aimpl.DrainBuffer(Buffer.Addr)

	t.AssertTransition(txn, deploy.Transition{
		Aimpl.Addr,     //sender
		"ClaimRewards", //tag
		Buffer.Addr,    //recipient
		"0",            //amount
		deploy.ParamsMap{},
	})

	/*
		//In order to check rewards we shoul repeat reward calculation logic from procedure CalcStakeRewards
			// Send funds and call a callback via Buffer
			t.AssertTransition(txn, deploy.Transition{
				Zimpl.Addr, //sender
				"AddFunds",
				Buffer.Addr,
				zil(1),
				deploy.ParamsMap{},
			})

			t.AssertTransition(txn, deploy.Transition{
				Zimpl.Addr, //sender
				"WithdrawStakeRewardsSuccessCallBack",
				Buffer.Addr,
				"0",
				deploy.ParamsMap{"rewards": zil(1)},
			})

			// Send funds and call a callback via Holder
			t.AssertTransition(txn, deploy.Transition{
				Zimpl.Addr, //sender
				"AddFunds",
				Holder.Addr,
				zil(1),
				deploy.ParamsMap{},
			})

			t.AssertTransition(txn, deploy.Transition{
				Zimpl.Addr, //sender
				"WithdrawStakeRewardsSuccessCallBack",
				Holder.Addr,
				"0",
				deploy.ParamsMap{"rewards": zil(1)},
			})

			// Check aZIL balance
			// 1 ZIL from Buffer + 1 ZIL from Holder
			t.AssertEqual(Aimpl.Field("_balance"), zil(2))
			t.AssertEqual(Aimpl.Field("autorestakeamount"), zil(2))
	*/

	// Send Swap transactions
	t.AssertTransition(txn, deploy.Transition{
		Buffer.Addr, //sender
		"RequestDelegatorSwap",
		Zproxy.Addr,
		"0",
		deploy.ParamsMap{"new_deleg_addr": "0x" + Holder.Addr},
	})

	t.AssertTransition(txn, deploy.Transition{
		Holder.Addr, //sender
		"ConfirmDelegatorSwap",
		Zproxy.Addr,
		"0",
		deploy.ParamsMap{"requestor": "0x" + Buffer.Addr},
	})

	//try to drain buffer, not existent at main staking contract
	//error should not be thrown
	new_buffers := []string{"0x0000000000000000000000000000000000000000"}
	t.AssertSuccess(Aimpl.ChangeBuffers(new_buffers))
	Aimpl.DrainBuffer("0000000000000000000000000000000000000000")
	t.AssertTransition(txn, deploy.Transition{
		Aimpl.Addr, //sender
		"ClaimRewards",
		Holder.Addr,
		"0",
		deploy.ParamsMap{},
	})
}
