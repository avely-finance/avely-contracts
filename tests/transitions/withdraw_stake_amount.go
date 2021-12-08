package transitions

import (
	// "log"
	"Azil/test/deploy"
)

func (t *Testing) WithdrawStakeAmount() {

	const FakeEpochNum = "1234567890"

	t.LogStart("WithdrawStakeAmount")

	// deploy smart contract
	stubStakingContract, aZilContract, _, holderContract := t.DeployAndUpgrade()

	/*******************************************************************************
	 * 0. delegator (addr2) delegate 15 zil, and it should enter in buffered deposit,
	 * we need to move buffered deposits to main stake
	 *******************************************************************************/
	aZilContract.UpdateWallet(key2)
	aZilContract.DelegateStake(zil(15))
	// TODO: if delegator have buffered deposits, withdrawal should fail
	stubStakingContract.AssignStakeReward()

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
	t.AssertEqual(aZilContract.StateField("totaltokenamount"), azil(15))

	/*******************************************************************************
	 * 2B. delegator send withdraw request, but it should fail because mindelegatestake
	 * TODO: how to be sure about size of mindelegatestake here?
	 *******************************************************************************/
	t.LogStart("WithdwarStakeAmount, step 2B")
	txn, err = aZilContract.WithdrawStakeAmt(azil(10))

	t.AssertError(txn, err, -15)
	t.AssertEqual(aZilContract.StateField("totaltokenamount"), azil(15))

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
	withdrawBlockNum := txn.Receipt.EpochNum

	newDelegBalanceZil, err := aZilContract.ZilBalanceOf(addr2)
	t.AssertEqual(stubStakingContract.StateField("totalstakeamount"), newDelegBalanceZil)
	t.AssertEqual(aZilContract.StateField("totalstakeamount"), newDelegBalanceZil)
	t.AssertEqual(aZilContract.StateField("totaltokenamount"), azil(10))
	t.AssertEqual(aZilContract.StateField("balances", "0x"+addr2), azil(10))
	t.AssertEqual(aZilContract.StateField("withdrawal_pending", withdrawBlockNum, "0x"+addr2, "0"), azil(5))
	t.AssertEqual(aZilContract.StateField("withdrawal_pending", withdrawBlockNum, "0x"+addr2, "1"), zil(5))

	/*******************************************************************************
	 * 3B. delegator withdrawing all remaining deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * Balances should be empty
	 *******************************************************************************/
	t.LogStart("WithdrawStakeAmount, step 3B")
	txn, _ = aZilContract.WithdrawStakeAmt(azil(10))
	t.AssertEvent(txn, deploy.Event{aZilContract.Addr, "WithdrawStakeAmt",
		deploy.ParamsMap{"withdraw_amount": azil(10), "withdraw_stake_amount": zil(10)}})
	t.AssertEqual(aZilContract.StateField("totalstakeamount"), "0")
	t.AssertEqual(aZilContract.StateField("totaltokenamount"), "0")
	t.AssertEqual(aZilContract.StateField("balances"), "empty")
	t.AssertEqual(stubStakingContract.StateField("totalstakeamount"), "0")
	/* this assertion is commented, because subsequent withdrawals may go to different block, so it's not trivial to check total withdrawals amount
	   * seems it's enough that we check withdrawal_pending at previous tests and zero-total here
	   //replace epoch number with fake
	   myRegexp = regexp.MustCompile(`\{\"(\d){1,10}\"\:\{\"argtypes\":\[\],`)
	   aZilState = myRegexp.ReplaceAllString(aZilState, "{\"" + FakeEpochNum + "\":{\"argtypes\":[],")
	   t.AssertContain(aZilState,"\"withdrawal_pending\":{\"" + "0x" + addr2 + "\":{\"" + FakeEpochNum + "\":{\"argtypes\":[],\"arguments\":[\"" + azil(15) + "\",\"" + azil(15) + "\"]")
	*/
}
