package transitions

import (
	// "log"
	"Azil/test/deploy"
)

func (t *Testing) WithdrawStakeAmount() {

	const FakeEpochNum = "1234567890"

	t.LogStart("WithdrawStakeAmount")

	// deploy smart contract
	Zproxy, _, aZilContract, _, holderContract := t.DeployAndUpgrade()

	/*******************************************************************************
	 * 0. delegator (addr2) delegate 15 zil, and it should enter in buffered deposit,
	 * we need to move buffered deposits to main stake
	 *******************************************************************************/
	aZilContract.UpdateWallet(key2)
	aZilContract.DelegateStake(zil(15))
	// TODO: if delegator have buffered deposits, withdrawal should fail
	Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT)

	/*******************************************************************************
	 * 1. non delegator(addr4) try to withdraw stake, should fail
	 *******************************************************************************/
	t.LogStart("WithdwarStakeAmount, step 1")
	aZilContract.UpdateWallet(key4)
	txn, err := aZilContract.WithdrawStakeAmt(azil(10))

	t.AssertError(txn, err, -7)

	/*******************************************************************************
	 * 2A. delegator trying to withdraw more than staked, should fail
	 *******************************************************************************/
	aZilContract.UpdateWallet(key2)
	t.LogStart("WithdwarStakeAmount, step 2A")
	txn, err = aZilContract.WithdrawStakeAmt(azil(100))

	t.AssertError(txn, err, -13)
	t.AssertEqual(aZilContract.Field("totaltokenamount"), azil(15))

	/*******************************************************************************
	 * 2B. delegator send withdraw request, but it should fail because mindelegatestake
	 * TODO: how to be sure about size of mindelegatestake here?
	 *******************************************************************************/
	t.LogStart("WithdwarStakeAmount, step 2B")
	txn, err = aZilContract.WithdrawStakeAmt(azil(10))

	t.AssertError(txn, err, -15)
	t.AssertEqual(aZilContract.Field("totaltokenamount"), azil(15))

	/*******************************************************************************
	 * 3A. delegator withdrawing part of his deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * balances field should be correct
	 *******************************************************************************/
	t.LogStart("WithdwarStakeAmount, step 3A")

	txn, err = aZilContract.WithdrawStakeAmt(azil(5))
	t.AssertTransition(txn, deploy.Transition{
		aZilContract.Addr,
		"WithdrawStakeAmt",
		holderContract.Addr,
		"0",
		deploy.ParamsMap{"amount": zil(5)},
	})
	bnum1 := txn.Receipt.EpochNum

	newDelegBalanceZil, err := aZilContract.ZilBalanceOf(addr2)
	t.AssertEqual(Zproxy.Field("totalstakeamount"), newDelegBalanceZil)
	t.AssertEqual(aZilContract.Field("totalstakeamount"), newDelegBalanceZil)
	t.AssertEqual(aZilContract.Field("totaltokenamount"), azil(10))
	t.AssertEqual(aZilContract.Field("balances", "0x"+addr2), azil(10))
	t.AssertEqual(aZilContract.Field("withdrawal_pending", bnum1, "0x"+addr2, "0"), azil(5))
	t.AssertEqual(aZilContract.Field("withdrawal_pending", bnum1, "0x"+addr2, "1"), zil(5))

	/*******************************************************************************
	 * 3B. delegator withdrawing all remaining deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * Balances should be empty
	 *******************************************************************************/
	t.LogStart("WithdrawStakeAmount, step 3B")
	txn, _ = aZilContract.WithdrawStakeAmt(azil(10))
	bnum2 := txn.Receipt.EpochNum
	t.AssertEvent(txn, deploy.Event{aZilContract.Addr, "WithdrawStakeAmt",
		deploy.ParamsMap{"withdraw_amount": azil(10), "withdraw_stake_amount": zil(10)}})
	t.AssertEqual(aZilContract.Field("totalstakeamount"), "0")
	t.AssertEqual(aZilContract.Field("totaltokenamount"), "0")
	t.AssertEqual(aZilContract.Field("balances"), "empty")
	t.AssertEqual(Zproxy.Field("totalstakeamount"), "0")
	if bnum1 == bnum2 {
		t.AssertEqual(aZilContract.Field("withdrawal_pending", bnum1, "0x"+addr2, "0"), azil(15))
		t.AssertEqual(aZilContract.Field("withdrawal_pending", bnum1, "0x"+addr2, "1"), zil(15))
	} else {
		//second withdrawal happened in next block
		t.AssertEqual(aZilContract.Field("withdrawal_pending", bnum2, "0x"+addr2, "0"), azil(10))
		t.AssertEqual(aZilContract.Field("withdrawal_pending", bnum2, "0x"+addr2, "1"), zil(10))
	}
}
