package transitions

import (
	. "Azil/test/helpers"
)

func (tr *Transitions) WithdrawStakeAmount() {

	t.Start("WithdrawStakeAmount")

	// deploy smart contract
	Zproxy, _, Aimpl, Buffer, Holder := tr.DeployAndUpgrade()

	/*******************************************************************************
	 * 0. delegator (tr.cfg.Addr2) delegate 15 zil
	 *******************************************************************************/
	Aimpl.UpdateWallet(tr.cfg.Key2)
	t.AssertSuccess(Aimpl.DelegateStake(Zil(15)))

	/*******************************************************************************
	 * 1. non delegator(tr.cfg.Addr4) try to withdraw stake, should fail
	 *******************************************************************************/
	t.Start("WithdwarStakeAmount, step 1")
	Aimpl.UpdateWallet(tr.cfg.Key3)
	txn, err := Aimpl.WithdrawStakeAmt(Azil(10))

	t.AssertError(txn, err, -7)

	/*******************************************************************************
	 * 2. Check withdrawal under delegator
	 *******************************************************************************/

	Aimpl.UpdateWallet(tr.cfg.Key2)

	/*******************************************************************************
	 * 2A. delegator trying to withdraw in the current cycle where he has a buffered deposit
	 *******************************************************************************/

	t.Start("WithdwarStakeAmount, step 2A")
	txn, err = Aimpl.WithdrawStakeAmt(Azil(1))

	t.AssertError(txn, err, -111)
	t.AssertEqual(Aimpl.Field("totaltokenamount"), Azil(1015))

	// Trigger switch to the next cycle
	Zproxy.AssignStakeReward(tr.cfg.AzilSsnAddress, tr.cfg.AzilSsnRewardShare)

	/*******************************************************************************
	 * 2B. delegator trying to withdraw more than staked, should fail
	 *******************************************************************************/

	t.Start("WithdwarStakeAmount, step 2A")
	txn, err = Aimpl.WithdrawStakeAmt(Azil(100))

	t.AssertError(txn, err, -13)
	t.AssertEqual(Aimpl.Field("totaltokenamount"), Azil(1015))

	/*******************************************************************************
	 * 2C. delegator send withdraw request, but it should fail because mindelegatestake
	 * TODO: how to be sure about size of mindelegatestake here?
	 *******************************************************************************/
	t.Start("WithdwarStakeAmount, step 2C")
	txn, err = Aimpl.WithdrawStakeAmt(Azil(10))

	t.AssertError(txn, err, -15)
	t.AssertEqual(Aimpl.Field("totaltokenamount"), Azil(1015))

	/*******************************************************************************
	 * 3A. delegator withdrawing part of his deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * balances field should be correct
	 *******************************************************************************/
	t.Start("WithdwarStakeAmount, step 3A")

	IncreaseBlocknum(10)
	t.AssertSuccess(Zproxy.AssignStakeReward(tr.cfg.AzilSsnAddress, tr.cfg.AzilSsnRewardShare))
	Aimpl.UpdateWallet(tr.cfg.AdminKey)
	t.AssertSuccess(Aimpl.DrainBuffer(Buffer.Addr))

	Aimpl.UpdateWallet(tr.cfg.Key2)
	txn, err = Aimpl.WithdrawStakeAmt(Azil(5))
	t.AssertTransition(txn, Transition{
		Aimpl.Addr,
		"WithdrawStakeAmt",
		Holder.Addr,
		"0",
		ParamsMap{"amount": Zil(5)},
	})
	bnum1 := txn.Receipt.EpochNum

	newDelegBalanceZil, err := Aimpl.ZilBalanceOf(tr.cfg.Addr2)
	//TODO: we can check this only in local testing environment,
	//and even in this case we need to monitor all incoming balances, including Holder initial delegate
	//t.AssertEqual(Zproxy.Field("totalstakeamount"), newDelegBalanceZil)
	t.AssertEqual(Aimpl.Field("totalstakeamount"), StrAdd(Zil(1000), newDelegBalanceZil))
	t.AssertEqual(Aimpl.Field("totaltokenamount"), Azil(1010))
	t.AssertEqual(Aimpl.Field("balances", "0x"+tr.cfg.Addr2), Azil(10))
	t.AssertEqual(Aimpl.Field("withdrawal_pending", bnum1, "0x"+tr.cfg.Addr2, "0"), Azil(5))
	t.AssertEqual(Aimpl.Field("withdrawal_pending", bnum1, "0x"+tr.cfg.Addr2, "1"), Zil(5))

	/*******************************************************************************
	 * 3B. delegator withdrawing all remaining deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * Balances should be empty
	 *******************************************************************************/
	t.Start("WithdrawStakeAmount, step 3B")
	txn, _ = Aimpl.WithdrawStakeAmt(Azil(10))
	bnum2 := txn.Receipt.EpochNum
	t.AssertEvent(txn, Event{Aimpl.Addr, "WithdrawStakeAmt",
		ParamsMap{"withdraw_amount": Azil(10), "withdraw_stake_amount": Zil(10)}})
	t.AssertEqual(Aimpl.Field("totalstakeamount"), Zil(1000))  //0
	t.AssertEqual(Aimpl.Field("totaltokenamount"), Azil(1000)) //0
	//t.AssertEqual(Aimpl.Field("balances"), "empty")
	t.AssertEqual(Aimpl.Field("balances", "0x"+tr.cfg.Admin), Azil(1000))
	//there is holder's initial stake
	//t.AssertEqual(Zproxy.Field("totalstakeamount"), "0")
	if bnum1 == bnum2 {
		t.AssertEqual(Aimpl.Field("withdrawal_pending", bnum1, "0x"+tr.cfg.Addr2, "0"), Azil(15))
		t.AssertEqual(Aimpl.Field("withdrawal_pending", bnum1, "0x"+tr.cfg.Addr2, "1"), Zil(15))
	} else {
		//second withdrawal happened in next block
		t.AssertEqual(Aimpl.Field("withdrawal_pending", bnum2, "0x"+tr.cfg.Addr2, "0"), Azil(10))
		t.AssertEqual(Aimpl.Field("withdrawal_pending", bnum2, "0x"+tr.cfg.Addr2, "1"), Zil(10))
	}
}
