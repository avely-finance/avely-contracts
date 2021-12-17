package transitions

import (
	"Azil/test/helpers"
)

func (tr *Transitions) WithdrawStakeAmount() {

	t.LogStart("WithdrawStakeAmount")

	// deploy smart contract
	Zproxy, _, Aimpl, Buffer, Holder := tr.DeployAndUpgrade()

	/*******************************************************************************
	 * 0. delegator (addr2) delegate 15 zil
	 *******************************************************************************/
	Aimpl.UpdateWallet(key2)
	t.AssertSuccess(Aimpl.DelegateStake(zil(15)))

	/*******************************************************************************
	 * 1. non delegator(addr4) try to withdraw stake, should fail
	 *******************************************************************************/
	t.LogStart("WithdwarStakeAmount, step 1")
	Aimpl.UpdateWallet(key4)
	txn, err := Aimpl.WithdrawStakeAmt(azil(10))

	t.AssertError(txn, err, -7)

	/*******************************************************************************
	 * 2. Check withdrawal under delegator
	 *******************************************************************************/

	Aimpl.UpdateWallet(key2)

	/*******************************************************************************
	 * 2A. delegator trying to withdraw in the current cycle where he has a buffered deposit
	 *******************************************************************************/

	t.LogStart("WithdwarStakeAmount, step 2A")
	txn, err = Aimpl.WithdrawStakeAmt(azil(1))

	t.AssertError(txn, err, -110)
	t.AssertEqual(Aimpl.Field("totaltokenamount"), azil(1015))

	// Trigger switch to the next cycle
	Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT)

	/*******************************************************************************
	 * 2B. delegator trying to withdraw more than staked, should fail
	 *******************************************************************************/

	t.LogStart("WithdwarStakeAmount, step 2A")
	txn, err = Aimpl.WithdrawStakeAmt(azil(100))

	t.AssertError(txn, err, -13)
	t.AssertEqual(Aimpl.Field("totaltokenamount"), azil(1015))

	/*******************************************************************************
	 * 2C. delegator send withdraw request, but it should fail because mindelegatestake
	 * TODO: how to be sure about size of mindelegatestake here?
	 *******************************************************************************/
	t.LogStart("WithdwarStakeAmount, step 2C")
	txn, err = Aimpl.WithdrawStakeAmt(azil(10))

	t.AssertError(txn, err, -15)
	t.AssertEqual(Aimpl.Field("totaltokenamount"), azil(1015))

	/*******************************************************************************
	 * 3A. delegator withdrawing part of his deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * balances field should be correct
	 *******************************************************************************/
	t.LogStart("WithdwarStakeAmount, step 3A")

	helpers.IncreaseBlocknum(10)
	t.AssertSuccess(Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT))
	Aimpl.UpdateWallet(adminKey)
	t.AssertSuccess(Aimpl.DrainBuffer(Buffer.Addr))

	Aimpl.UpdateWallet(key2)
	txn, err = Aimpl.WithdrawStakeAmt(azil(5))
	t.AssertTransition(txn, helpers.Transition{
		Aimpl.Addr,
		"WithdrawStakeAmt",
		Holder.Addr,
		"0",
		helpers.ParamsMap{"amount": zil(5)},
	})
	bnum1 := txn.Receipt.EpochNum

	newDelegBalanceZil, err := Aimpl.ZilBalanceOf(addr2)
	//TODO: we can check this only in local testing environment,
	//and even in this case we need to monitor all incoming balances, including Holder initial delegate
	//t.AssertEqual(Zproxy.Field("totalstakeamount"), newDelegBalanceZil)
	t.AssertEqual(Aimpl.Field("totalstakeamount"), helpers.StrAdd(zil(1000), newDelegBalanceZil))
	t.AssertEqual(Aimpl.Field("totaltokenamount"), azil(1010))
	t.AssertEqual(Aimpl.Field("balances", "0x"+addr2), azil(10))
	t.AssertEqual(Aimpl.Field("withdrawal_pending", bnum1, "0x"+addr2, "0"), azil(5))
	t.AssertEqual(Aimpl.Field("withdrawal_pending", bnum1, "0x"+addr2, "1"), zil(5))

	/*******************************************************************************
	 * 3B. delegator withdrawing all remaining deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * Balances should be empty
	 *******************************************************************************/
	t.LogStart("WithdrawStakeAmount, step 3B")
	txn, _ = Aimpl.WithdrawStakeAmt(azil(10))
	bnum2 := txn.Receipt.EpochNum
	t.AssertEvent(txn, helpers.Event{Aimpl.Addr, "WithdrawStakeAmt",
		helpers.ParamsMap{"withdraw_amount": azil(10), "withdraw_stake_amount": zil(10)}})
	t.AssertEqual(Aimpl.Field("totalstakeamount"), zil(1000))  //0
	t.AssertEqual(Aimpl.Field("totaltokenamount"), azil(1000)) //0
	//t.AssertEqual(Aimpl.Field("balances"), "empty")
	t.AssertEqual(Aimpl.Field("balances", "0x"+admin), azil(1000))
	//there is holder's initial stake
	//t.AssertEqual(Zproxy.Field("totalstakeamount"), "0")
	if bnum1 == bnum2 {
		t.AssertEqual(Aimpl.Field("withdrawal_pending", bnum1, "0x"+addr2, "0"), azil(15))
		t.AssertEqual(Aimpl.Field("withdrawal_pending", bnum1, "0x"+addr2, "1"), zil(15))
	} else {
		//second withdrawal happened in next block
		t.AssertEqual(Aimpl.Field("withdrawal_pending", bnum2, "0x"+addr2, "0"), azil(10))
		t.AssertEqual(Aimpl.Field("withdrawal_pending", bnum2, "0x"+addr2, "1"), zil(10))
	}
}
