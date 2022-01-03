package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) WithdrawStakeAmount() {

	Start("WithdrawStakeAmount")

	// deploy smart contract
	p := tr.DeployAndUpgrade()

	/*******************************************************************************
	 * 0. delegator (sdk.Cfg.Addr2) delegate 15 zil
	 *******************************************************************************/
	p.Aproxy.UpdateWallet(sdk.Cfg.Key2)
	AssertSuccess(p.Aproxy.DelegateStake(ToZil(15)))

	/*******************************************************************************
	 * 1. non delegator(sdk.Cfg.Addr4) try to withdraw stake, should fail
	 *******************************************************************************/
	Start("WithdwarStakeAmount, step 1")
	p.Aproxy.UpdateWallet(sdk.Cfg.Key3)
	txn, err := p.Aproxy.WithdrawStakeAmt(ToAzil(10))

	AssertError(txn, err, -7)

	/*******************************************************************************
	 * 2. Check withdrawal under delegator
	 *******************************************************************************/

	p.Aproxy.UpdateWallet(sdk.Cfg.Key2)

	/*******************************************************************************
	 * 2A. delegator trying to withdraw in the current cycle where he has a buffered deposit
	 *******************************************************************************/

	Start("WithdwarStakeAmount, step 2A")
	txn, err = p.Aproxy.WithdrawStakeAmt(ToAzil(1))

	AssertError(txn, err, -111)
	AssertEqual(p.Aimpl.Field("totaltokenamount"), ToAzil(1015))

	// Trigger switch to the next cycle
	p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare)

	/*******************************************************************************
	 * 2B. delegator trying to withdraw more than staked, should fail
	 *******************************************************************************/

	Start("WithdwarStakeAmount, step 2A")
	txn, err = p.Aproxy.WithdrawStakeAmt(ToAzil(100))

	AssertError(txn, err, -13)
	AssertEqual(p.Aimpl.Field("totaltokenamount"), ToAzil(1015))

	/*******************************************************************************
	 * 2C. delegator send withdraw request, but it should fail because mindelegatestake
	 * TODO: how to be sure about size of mindelegatestake here?
	 *******************************************************************************/
	Start("WithdwarStakeAmount, step 2C")
	txn, err = p.Aproxy.WithdrawStakeAmt(ToAzil(10))

	AssertError(txn, err, -15)
	AssertEqual(p.Aimpl.Field("totaltokenamount"), ToAzil(1015))

	/*******************************************************************************
	 * 3A. delegator withdrawing part of his deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * balances field should be correct
	 *******************************************************************************/
	Start("WithdwarStakeAmount, step 3A")

	sdk.IncreaseBlocknum(10)
	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))
	p.Aimpl.UpdateWallet(sdk.Cfg.AdminKey)
	AssertSuccess(p.Aimpl.DrainBuffer(p.GetBuffer().Addr))

	p.Aproxy.UpdateWallet(sdk.Cfg.Key2)
	txn, err = p.Aproxy.WithdrawStakeAmt(ToAzil(5))
	AssertTransition(txn, Transition{
		p.Aimpl.Addr,
		"WithdrawStakeAmt",
		p.Holder.Addr,
		"0",
		ParamsMap{"amount": ToZil(5)},
	})
	bnum1 := txn.Receipt.EpochNum

	newDelegBalanceZil, err := p.Aproxy.ZilBalanceOf(sdk.Cfg.Addr2)
	//TODO: we can check this only in local testing environment,
	//and even in this case we need to monitor all incoming balances, including Holder initial delegate
	//t.AssertEqual(p.Zproxy.Field("totalstakeamount"), newDelegBalanceZil)
	AssertEqual(p.Aimpl.Field("totalstakeamount"), StrAdd(ToZil(1000), newDelegBalanceZil))
	AssertEqual(p.Aimpl.Field("totaltokenamount"), ToAzil(1010))
	AssertEqual(p.Aimpl.Field("balances", "0x"+sdk.Cfg.Addr2), ToAzil(10))
	AssertEqual(p.Aimpl.Field("withdrawal_pending", bnum1, "0x"+sdk.Cfg.Addr2, "0"), ToAzil(5))
	AssertEqual(p.Aimpl.Field("withdrawal_pending", bnum1, "0x"+sdk.Cfg.Addr2, "1"), ToZil(5))

	/*******************************************************************************
	 * 3B. delegator withdrawing all remaining deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * Balances should be empty
	 *******************************************************************************/
	Start("WithdrawStakeAmount, step 3B")
	txn, _ = p.Aproxy.WithdrawStakeAmt(ToAzil(10))
	bnum2 := txn.Receipt.EpochNum
	AssertEvent(txn, Event{p.Aimpl.Addr, "WithdrawStakeAmt",
		ParamsMap{"withdraw_amount": ToAzil(10), "withdraw_stake_amount": ToZil(10)}})
	AssertEqual(p.Aimpl.Field("totalstakeamount"), ToZil(1000))  //0
	AssertEqual(p.Aimpl.Field("totaltokenamount"), ToAzil(1000)) //0
	//t.AssertEqual(p.Aimpl.Field("balances"), "empty")
	AssertEqual(p.Aimpl.Field("balances", "0x"+sdk.Cfg.Admin), ToAzil(1000))
	//there is holder's initial stake
	//t.AssertEqual(p.Zproxy.Field("totalstakeamount"), "0")
	if bnum1 == bnum2 {
		AssertEqual(p.Aimpl.Field("withdrawal_pending", bnum1, "0x"+sdk.Cfg.Addr2, "0"), ToAzil(15))
		AssertEqual(p.Aimpl.Field("withdrawal_pending", bnum1, "0x"+sdk.Cfg.Addr2, "1"), ToZil(15))
	} else {
		//second withdrawal happened in next block
		AssertEqual(p.Aimpl.Field("withdrawal_pending", bnum2, "0x"+sdk.Cfg.Addr2, "0"), ToAzil(10))
		AssertEqual(p.Aimpl.Field("withdrawal_pending", bnum2, "0x"+sdk.Cfg.Addr2, "1"), ToZil(10))
	}
}
