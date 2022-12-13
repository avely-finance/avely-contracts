package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/utils"
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
	p.StZIL.SetSigner(bob)
	AssertSuccess(p.StZIL.DelegateStake(ToZil(15)))

	/*******************************************************************************
	 * 1. non delegator(sdk.Cfg.Addr4) try to withdraw stake, should fail
	 *******************************************************************************/
	Start("WithdwarStakeAmount, step 1")
	p.StZIL.SetSigner(eve)
	txn, _ := p.StZIL.WithdrawStakeAmt(ToStZil(10))

	AssertError(txn, p.StZIL.ErrorCode("DelegDoesNotExistAtSSN"))

	/*******************************************************************************
	 * 2. Check withdrawal under delegator
	 *******************************************************************************/

	p.StZIL.SetSigner(bob)
	bobAddr := utils.GetAddressByWallet(bob)

	/*******************************************************************************
	 * 2A. delegator trying to withdraw more than staked, should fail
	 *******************************************************************************/

	Start("WithdwarStakeAmount, step 2A")
	txn, _ = p.StZIL.WithdrawStakeAmt(ToStZil(100))

	AssertError(txn, p.StZIL.ErrorCode("DelegHasNoSufficientAmt"))
	AssertEqual(Field(p.StZIL, "total_supply"), ToStZil(totalSsnInitialDelegateZil+15))

	/*******************************************************************************
	 * 2B. delegator send withdraw request, but it should fail because mindelegatestake
	 * TODO: how to be sure about size of mindelegatestake here?
	 *******************************************************************************/
	Start("WithdwarStakeAmount, step 2C")
	txn, _ = p.StZIL.WithdrawStakeAmt(ToStZil(10))

	AssertError(txn, p.StZIL.ErrorCode("DelegStakeNotEnough"))
	AssertEqual(Field(p.StZIL, "total_supply"), ToStZil(totalSsnInitialDelegateZil+15))

	/*******************************************************************************
	 * 3A. delegator withdrawing part of his deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * balances field should be correct.
	 * Delegator able to init withdrawal in same cycle when deposit was done.
	 *******************************************************************************/
	Start("WithdwarStakeAmount, step 3A")

	txn, _ = p.StZIL.WithdrawStakeAmt(ToStZil(5))
	AssertTransition(txn, Transition{
		p.StZIL.Addr,
		"WithdrawStakeAmt",
		p.Holder.Addr,
		"0",
		ParamsMap{"amount": ToZil(5)},
	})
	bnum1 := txn.Receipt.EpochNum

	newDelegBalanceZil := p.StZIL.ZilBalanceOf(bobAddr).String()
	AssertEqual(Field(p.StZIL, "totalstakeamount"), StrAdd(ToZil(totalSsnInitialDelegateZil), newDelegBalanceZil))
	AssertEqual(Field(p.StZIL, "total_supply"), ToStZil(totalSsnInitialDelegateZil+10))
	AssertEqual(Field(p.StZIL, "balances", bobAddr), ToStZil(10))

	AssertEvent(txn, Event{p.StZIL.Addr, "Burnt", ParamsMap{
		"burner":       p.StZIL.Addr,
		"burn_account": bobAddr,
		"amount":       ToStZil(5),
	}})

	withdrawal := Dig(p.StZIL, "withdrawal_pending", bnum1, bobAddr).Withdrawal()

	AssertEqual(withdrawal.TokenAmount.String(), ToStZil(5))
	AssertEqual(withdrawal.StakeAmount.String(), ToStZil(5))

	/*******************************************************************************
	 * 3B. delegator withdrawing all remaining deposit, it should success with "_eventname": "WithdrawStakeAmt"
	 * Also check that withdrawal_pending field contains correct information about requested withdrawal
	 * Balances should be empty
	 *******************************************************************************/
	Start("WithdrawStakeAmount, step 3B")
	txn, _ = p.StZIL.WithdrawStakeAmt(ToStZil(10))
	bnum2 := txn.Receipt.EpochNum
	AssertEvent(txn, Event{p.StZIL.Addr, "WithdrawStakeAmt",
		ParamsMap{"withdraw_amount": ToStZil(10), "withdraw_stake_amount": ToZil(10)}})
	AssertEqual(Field(p.StZIL, "totalstakeamount"), ToZil(totalSsnInitialDelegateZil))
	AssertEqual(Field(p.StZIL, "total_supply"), ToStZil(totalSsnInitialDelegateZil))
	if totalSsnInitialDelegateZil == 0 {
		AssertEqual(Field(p.StZIL, "balances"), "{}")
		AssertEqual(Field(p.StZIL, "balances", utils.GetAddressByWallet(celestials.Admin)), "")
	} else {
		AssertEqual(Field(p.StZIL, "balances", utils.GetAddressByWallet(celestials.Admin)), ToStZil(totalSsnInitialDelegateZil))
	}
	//there is holder's initial stake
	if bnum1 == bnum2 {
		withdrawal := Dig(p.StZIL, "withdrawal_pending", bnum1, bobAddr).Withdrawal()
		AssertEqual(withdrawal.TokenAmount.String(), ToStZil(15))
		AssertEqual(withdrawal.StakeAmount.String(), ToStZil(15))
	} else {
		//second withdrawal happened in next block
		withdrawal := Dig(p.StZIL, "withdrawal_pending", bnum2, bobAddr).Withdrawal()
		AssertEqual(withdrawal.TokenAmount.String(), ToStZil(10))
		AssertEqual(withdrawal.StakeAmount.String(), ToStZil(10))
	}
}
