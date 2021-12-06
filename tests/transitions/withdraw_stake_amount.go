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
	aZilContract.DelegateStake(zil15)
	// TODO: if delegator have buffered deposits, withdrawal should fail
	stubStakingContract.AssignStakeReward()

	/*******************************************************************************
	 * 1. non delegator(addr4) try to withdraw stake, should fail
	 *******************************************************************************/
	t.LogStart("WithdwarStakeAmount, step 1")
	aZilContract.UpdateWallet(key4)
	txn, err := aZilContract.WithdrawStakeAmt(azil10)

	t.AssertSuccessCall(err)
	t.AssertContain(t.GetReceiptString(txn), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -7))])")

	/*******************************************************************************
	 * 2A. delegator trying to withdraw more than staked, should fail
	 *******************************************************************************/
	aZilContract.UpdateWallet(key2)
	t.LogStart("WithdwarStakeAmount, step 2A")
	txn, err = aZilContract.WithdrawStakeAmt(azil100)
	t.AssertSuccessCall(err)
	// t.LogPrettyReceipt(txn)
	t.AssertContain(t.GetReceiptString(txn), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -13))])")
	t.AssertState("AimplState", deploy.ParamsMap{"totaltokenamount": azil15})

	/*******************************************************************************
	 * 2B. delegator send withdraw request, but it should fail because mindelegatestake
	 * TODO: how to be sure about size of mindelegatestake here?
	 *******************************************************************************/
	t.LogStart("WithdwarStakeAmount, step 2B")
	txn, err = aZilContract.WithdrawStakeAmt(azil10)
	t.AssertSuccessCall(err)
	t.AssertContain(t.GetReceiptString(txn), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -15))])")
	t.AssertState("AimplState", deploy.ParamsMap{"totaltokenamount": azil15})

	/*******************************************************************************
	 * 3A. delegator withdrawing part of his deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * balances field should be correct
	 *******************************************************************************/
	t.LogStart("WithdwarStakeAmount, step 3A")

	txn, err = aZilContract.WithdrawStakeAmt(azil5)
	t.AssertTransition(txn, deploy.Transition{
		aZilContract.Addr,
		"WithdrawStakeAmt",
		holderContract.Addr,
		"0",
		deploy.ParamsMap{"amount": zil5},
	})
	withdrawBlockNum := txn.Receipt.EpochNum

	newDelegBalanceZil, err := aZilContract.ZilBalanceOf(addr2)

	t.AssertState("ZimplState", deploy.ParamsMap{"totalstakeamount": newDelegBalanceZil})
	t.AssertState("AimplState", deploy.ParamsMap{"totalstakeamount": newDelegBalanceZil, "totaltokenamount": azil10})
	t.AssertState("AimplStateBalance", deploy.ParamsMap{"address": "0x" + addr2, "token": azil10})
	txn, err = FetcherContract.AimplWithdrawalPending(withdrawBlockNum, "0x"+addr2)
	t.AssertEvent(txn, deploy.Event{FetcherContract.Addr, "AimplWithdrawalPending", deploy.ParamsMap{"token": azil5, "stake": zil5}})

	/*******************************************************************************
	 * 3B. delegator withdrawing all remaining deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * Balances should be empty
	 *******************************************************************************/
	t.LogStart("WithdrawStakeAmount, step 3B")
	txn, _ = aZilContract.WithdrawStakeAmt(azil10)
	t.AssertEvent(txn, deploy.Event{aZilContract.Addr, "WithdrawStakeAmt",
		deploy.ParamsMap{"withdraw_amount": azil10, "withdraw_stake_amount": zil10}})
	t.AssertState("AimplState", deploy.ParamsMap{"totalstakeamount": "0", "totaltokenamount": "0", "balances": "empty"})
	t.AssertState("ZimplState", deploy.ParamsMap{"totalstakeamount": "0"})
	/* this assertion is commented, because subsequent withdrawals may go to different block, so it's not trivial to check total withdrawals amount
	   * seems it's enough that we check withdrawal_pending at previous tests and zero-total here
	   //replace epoch number with fake
	   myRegexp = regexp.MustCompile(`\{\"(\d){1,10}\"\:\{\"argtypes\":\[\],`)
	   aZilState = myRegexp.ReplaceAllString(aZilState, "{\"" + FakeEpochNum + "\":{\"argtypes\":[],")
	   t.AssertContain(aZilState,"\"withdrawal_pending\":{\"" + "0x" + addr2 + "\":{\"" + FakeEpochNum + "\":{\"argtypes\":[],\"arguments\":[\"" + azil15 + "\",\"" + azil15 + "\"]")
	*/
}
