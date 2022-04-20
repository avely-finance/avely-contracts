package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) WithdrawStakeAmount() {

	Start("WithdrawStakeAmount")

	// deploy smart contract
	p := tr.DeployAndUpgrade()
	//totalSsnInitialDelegateZil := len(sdk.Cfg.SsnAddrs) * sdk.Cfg.SsnInitialDelegateZil
	//for now to activate SSNs we delegate required stakes through Zproxy as admin
	totalSsnInitialDelegateZil := 0

	/*******************************************************************************
	 * 0. delegator (sdk.Cfg.Addr2) delegate 15 zil
	 *******************************************************************************/
	p.Azil.UpdateWallet(sdk.Cfg.Key2)
	AssertSuccess(p.Azil.DelegateStake(ToZil(15)))

	/*******************************************************************************
	 * 1. non delegator(sdk.Cfg.Addr4) try to withdraw stake, should fail
	 *******************************************************************************/
	Start("WithdwarStakeAmount, step 1")
	p.Azil.UpdateWallet(sdk.Cfg.Key3)
	txn, _ := p.Azil.WithdrawStakeAmt(ToAzil(10))

	AssertError(txn, "DelegDoesNotExistAtSSN")

	/*******************************************************************************
	 * 2. Check withdrawal under delegator
	 *******************************************************************************/

	p.Azil.UpdateWallet(sdk.Cfg.Key2)

	/*******************************************************************************
	 * 2A. delegator trying to withdraw more than staked, should fail
	 *******************************************************************************/

	Start("WithdwarStakeAmount, step 2A")
	txn, _ = p.Azil.WithdrawStakeAmt(ToAzil(100))

	AssertError(txn, "DelegHasNoSufficientAmt")
	AssertEqual(Field(p.Azil, "total_supply"), ToAzil(totalSsnInitialDelegateZil+15))

	/*******************************************************************************
	 * 2B. delegator send withdraw request, but it should fail because mindelegatestake
	 * TODO: how to be sure about size of mindelegatestake here?
	 *******************************************************************************/
	Start("WithdwarStakeAmount, step 2C")
	txn, _ = p.Azil.WithdrawStakeAmt(ToAzil(10))

	AssertError(txn, "DelegStakeNotEnough")
	AssertEqual(Field(p.Azil, "total_supply"), ToAzil(totalSsnInitialDelegateZil+15))

	/*******************************************************************************
	 * 3A. delegator withdrawing part of his deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * balances field should be correct.
	 * Delegator able to init withdrawal in same cycle when deposit was done.
	 *******************************************************************************/
	Start("WithdwarStakeAmount, step 3A")

	txn, _ = p.Azil.WithdrawStakeAmt(ToAzil(5))
	AssertTransition(txn, Transition{
		p.Azil.Addr,
		"WithdrawStakeAmt",
		p.Holder.Addr,
		"0",
		ParamsMap{"amount": ToZil(5)},
	})
	bnum1 := txn.Receipt.EpochNum

	newDelegBalanceZil := p.Azil.ZilBalanceOf(sdk.Cfg.Addr2).String()
	AssertEqual(Field(p.Azil, "totalstakeamount"), StrAdd(ToZil(totalSsnInitialDelegateZil), newDelegBalanceZil))
	AssertEqual(Field(p.Azil, "total_supply"), ToAzil(totalSsnInitialDelegateZil+10))
	AssertEqual(Field(p.Azil, "balances", sdk.Cfg.Addr2), ToAzil(10))

	withdrawal := Dig(p.Azil, "withdrawal_pending", bnum1, sdk.Cfg.Addr2).Withdrawal()

	AssertEqual(withdrawal.TokenAmount.String(), ToAzil(5))
	AssertEqual(withdrawal.StakeAmount.String(), ToAzil(5))

	/*******************************************************************************
	 * 3B. delegator withdrawing all remaining deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * Balances should be empty
	 *******************************************************************************/
	Start("WithdrawStakeAmount, step 3B")
	txn, _ = p.Azil.WithdrawStakeAmt(ToAzil(10))
	bnum2 := txn.Receipt.EpochNum
	AssertEvent(txn, Event{p.Azil.Addr, "WithdrawStakeAmt",
		ParamsMap{"withdraw_amount": ToAzil(10), "withdraw_stake_amount": ToZil(10)}})
	AssertEqual(Field(p.Azil, "totalstakeamount"), ToZil(totalSsnInitialDelegateZil))
	AssertEqual(Field(p.Azil, "total_supply"), ToAzil(totalSsnInitialDelegateZil))
	if totalSsnInitialDelegateZil == 0 {
		AssertEqual(Field(p.Azil, "balances"), "{}")
		AssertEqual(Field(p.Azil, "balances", sdk.Cfg.Admin), "")
	} else {
		AssertEqual(Field(p.Azil, "balances", sdk.Cfg.Admin), ToAzil(totalSsnInitialDelegateZil))
	}
	//there is holder's initial stake
	if bnum1 == bnum2 {
		withdrawal := Dig(p.Azil, "withdrawal_pending", bnum1, sdk.Cfg.Addr2).Withdrawal()
		AssertEqual(withdrawal.TokenAmount.String(), ToAzil(15))
		AssertEqual(withdrawal.StakeAmount.String(), ToAzil(15))
	} else {
		//second withdrawal happened in next block
		withdrawal := Dig(p.Azil, "withdrawal_pending", bnum2, sdk.Cfg.Addr2).Withdrawal()
		AssertEqual(withdrawal.TokenAmount.String(), ToAzil(10))
		AssertEqual(withdrawal.StakeAmount.String(), ToAzil(10))
	}
}
