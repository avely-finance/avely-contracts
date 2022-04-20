package transitions

import (
	"strconv"

	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) CompleteWithdrawalSuccess() {

	Start("CompleteWithdrawal - success")
	readyBlocks := []string{}

	p := tr.DeployAndUpgrade()
	//totalSsnInitialDelegateZil := len(sdk.Cfg.SsnAddrs) * sdk.Cfg.SsnInitialDelegateZil
	//for now to activate SSNs we delegate required stakes through Zproxy as admin
	totalSsnInitialDelegateZil := 0

	p.Azil.UpdateWallet(sdk.Cfg.Key1)
	AssertSuccess(p.Azil.DelegateStake(ToZil(10)))

	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	tx, _ := AssertSuccess(p.Azil.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToAzil(10)))

	block1 := tx.Receipt.EpochNum
	tx, _ = p.Azil.CompleteWithdrawal()
	AssertEvent(tx, Event{p.Azil.Addr, "NoUnbondedStake", ParamsMap{}})

	p.Azil.UpdateWallet(sdk.Cfg.Key2)
	tx, _ = p.Azil.CompleteWithdrawal()
	AssertEvent(tx, Event{p.Azil.Addr, "NoUnbondedStake", ParamsMap{}})

	p.Azil.UpdateWallet(sdk.Cfg.AdminKey)
	readyBlocks = append(readyBlocks, block1)
	tx, _ = p.Azil.ClaimWithdrawal(readyBlocks)
	AssertError(tx, "ClaimWithdrawalNoUnbonded")

	delta, _ := strconv.ParseInt(StrAdd(Field(p.Zimpl, "bnum_req"), "1"), 10, 32)
	sdk.IncreaseBlocknum(int32(delta))

	tr.NextCycle(p)
	tools := tr.NextCycleOffchain(p)

	unbondedWithdrawalsBlocks := p.GetClaimWithdrawalBlocks()
	AssertEqual(readyBlocks[0], strconv.Itoa(unbondedWithdrawalsBlocks[0]))
	tools.ShowClaimWithdrawal(p)

	withdrawal := Dig(p.Azil, "withdrawal_pending_of_delegator", sdk.Cfg.Addr1, block1).Withdrawal()
	AssertEqual(withdrawal.TokenAmount.String(), ToAzil(10))
	AssertEqual(withdrawal.StakeAmount.String(), ToAzil(10))

	p.Azil.UpdateWallet(sdk.Cfg.AdminKey)
	tx, _ = p.Azil.ClaimWithdrawal(readyBlocks)
	AssertTransition(tx, Transition{
		p.Azil.Addr,          //sender
		"CompleteWithdrawal", //tag
		p.Holder.Addr,        //recipient
		"0",                  //amount
		ParamsMap{},
	})
	AssertEvent(tx, Event{p.Holder.Addr, "AddFunds", ParamsMap{"funder": p.Zimpl.Addr, "amount": ToZil(10)}})

	AssertTransition(tx, Transition{
		p.Holder.Addr,                       //sender
		"CompleteWithdrawalSuccessCallBack", //tag
		p.Azil.Addr,                         //recipient
		ToZil(10),                           //amount
		ParamsMap{},
	})

	p.Azil.UpdateWallet(sdk.Cfg.Key1)
	tx, _ = p.Azil.CompleteWithdrawal()
	AssertEvent(tx, Event{p.Azil.Addr, "CompleteWithdrawal", ParamsMap{"amount": ToZil(10), "delegator": sdk.Cfg.Addr1}})
	AssertTransition(tx, Transition{
		p.Azil.Addr,
		"CompleteWithdrawalSuccessCallBack",
		sdk.Cfg.Addr1,
		"0",
		ParamsMap{"amount": ToZil(10)},
	})
	AssertTransition(tx, Transition{
		p.Azil.Addr,
		"AddFunds",
		sdk.Cfg.Addr1,
		ToZil(10),
		ParamsMap{},
	})

	withdrawal = Dig(p.Azil, "withdrawal_pending_of_delegator", sdk.Cfg.Addr1, block1).Withdrawal()
	AssertEqual(withdrawal.TokenAmount.String(), "0")
	AssertEqual(withdrawal.StakeAmount.String(), "0")

	AssertEqual(Field(p.Azil, "totalstakeamount"), ToZil(totalSsnInitialDelegateZil))
	AssertEqual(Field(p.Azil, "total_supply"), ToAzil(totalSsnInitialDelegateZil))
	AssertEqual(Field(p.Azil, "tmp_complete_withdrawal_available"), "0")

	if totalSsnInitialDelegateZil == 0 {
		AssertEqual(Field(p.Azil, "balances"), "{}")
		AssertEqual(Field(p.Azil, "balances", sdk.Cfg.Admin), "")
	} else {
		AssertEqual(Field(p.Azil, "balances", sdk.Cfg.Admin), ToAzil(totalSsnInitialDelegateZil))
	}

	AssertEqual(Field(p.Azil, "withdrawal_unbonded"), "{}")
	AssertEqual(Field(p.Azil, "withdrawal_pending"), "{}")
}

func (tr *Transitions) CompleteWithdrawalMultiSsn() {

	Start("CompleteWithdrawal - success")
	readyBlocks := []string{}

	p := tr.DeployAndUpgrade()

	rewardsFee := "1000" //10% of feeDenom=10000
	treasuryAddr := sdk.Cfg.Addr3
	AssertSuccess(p.Azil.WithUser(sdk.Cfg.OwnerKey).ChangeRewardsFee(rewardsFee))
	AssertSuccess(p.Azil.WithUser(sdk.Cfg.OwnerKey).ChangeTreasuryAddress(treasuryAddr))
	p.Azil.UpdateWallet(sdk.Cfg.AdminKey) //back to admin

	//totalSsnInitialDelegateZil := len(sdk.Cfg.SsnAddrs) * sdk.Cfg.SsnInitialDelegateZil
	//for now to activate SSNs we delegate required stakes through Zproxy as admin
	totalSsnInitialDelegateZil := 0

	ssnForInput1 := p.GetSsnAddressForInput()

	//for current test setup first SSN for input is AzilSSN
	AssertEqual(ssnForInput1, sdk.Cfg.AzilSsnAddress)

	AssertSuccess(p.Azil.WithUser(sdk.Cfg.Key1).DelegateStake(ToZil(5000)))

	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	ssnForInput2 := p.GetSsnAddressForInput()
	AssertNotEqual(ssnForInput1, ssnForInput2)
	AssertSuccess(p.Azil.WithUser(sdk.Cfg.Key1).DelegateStake(ToZil(5000)))
	AssertEqual(Field(p.Azil, "totalstakeamount"), ToZil(totalSsnInitialDelegateZil+5000+5000))
	AssertEqual(Field(p.Azil, "total_supply"), ToAzil(totalSsnInitialDelegateZil+5000+5000))

	//balance of test user is 10k
	AssertEqual(Field(p.Azil, "balances", sdk.Cfg.Addr1), ToAzil(5000+5000))

	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	//stake is on holder now, splitted between SSNs
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnForInput1), ToZil(sdk.Cfg.HolderInitialDelegateZil+5000))
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnForInput2), ToZil(5000))

	//it's impossible to withdraw amount, bigger than amount on heaviest SSN
	tx, _ := p.Azil.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToAzil(7000))
	AssertError(tx, "WithdrawAmountTooBig")

	//withdraw correct amount
	AssertSuccess(p.Azil.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToAzil(5500)))

	block1 := tx.Receipt.EpochNum
	tx, _ = p.Azil.CompleteWithdrawal()
	AssertEvent(tx, Event{p.Azil.Addr, "NoUnbondedStake", ParamsMap{}})

	p.Azil.UpdateWallet(sdk.Cfg.Key2)
	tx, _ = p.Azil.CompleteWithdrawal()
	AssertEvent(tx, Event{p.Azil.Addr, "NoUnbondedStake", ParamsMap{}})

	p.Azil.UpdateWallet(sdk.Cfg.AdminKey)
	readyBlocks = append(readyBlocks, block1)
	tx, _ = p.Azil.ClaimWithdrawal(readyBlocks)
	AssertError(tx, "ClaimWithdrawalNoUnbonded")

	delta, _ := strconv.ParseInt(StrAdd(Field(p.Zimpl, "bnum_req"), "1"), 10, 32)
	sdk.IncreaseBlocknum(int32(delta))

	tr.NextCycle(p)
	tools := tr.NextCycleOffchain(p)

	unbondedWithdrawalsBlocks := p.GetClaimWithdrawalBlocks()
	AssertEqual(readyBlocks[0], strconv.Itoa(unbondedWithdrawalsBlocks[0]))
	tools.ShowClaimWithdrawal(p)

	withdrawal := Dig(p.Azil, "withdrawal_pending_of_delegator", sdk.Cfg.Addr1, block1).Withdrawal()
	AssertEqual(withdrawal.TokenAmount.String(), ToAzil(5500))
	AssertEqual(withdrawal.StakeAmount.String(), ToAzil(5500))

	p.Azil.UpdateWallet(sdk.Cfg.AdminKey)
	tx, _ = p.Azil.ClaimWithdrawal(readyBlocks)
	AssertTransition(tx, Transition{
		p.Azil.Addr,          //sender
		"CompleteWithdrawal", //tag
		p.Holder.Addr,        //recipient
		"0",                  //amount
		ParamsMap{},
	})
	AssertEvent(tx, Event{p.Holder.Addr, "AddFunds", ParamsMap{"funder": p.Zimpl.Addr, "amount": ToZil(5500)}})

	AssertTransition(tx, Transition{
		p.Holder.Addr,                       //sender
		"CompleteWithdrawalSuccessCallBack", //tag
		p.Azil.Addr,                         //recipient
		ToZil(5500),                         //amount
		ParamsMap{},
	})

	p.Azil.UpdateWallet(sdk.Cfg.Key1)
	tx, _ = AssertSuccess(p.Azil.CompleteWithdrawal())
	AssertEvent(tx, Event{p.Azil.Addr, "CompleteWithdrawal", ParamsMap{"amount": ToZil(5500), "delegator": sdk.Cfg.Addr1}})
	AssertTransition(tx, Transition{
		p.Azil.Addr,
		"CompleteWithdrawalSuccessCallBack",
		sdk.Cfg.Addr1,
		"0",
		ParamsMap{"amount": ToZil(5500)},
	})
	AssertTransition(tx, Transition{
		p.Azil.Addr,
		"AddFunds",
		sdk.Cfg.Addr1,
		ToZil(5500),
		ParamsMap{},
	})

	withdrawal = Dig(p.Azil, "withdrawal_pending_of_delegator", sdk.Cfg.Addr1, block1).Withdrawal()
	AssertEqual(withdrawal.TokenAmount.String(), "0")
	AssertEqual(withdrawal.StakeAmount.String(), "0")

	AssertEqual(Field(p.Azil, "totalstakeamount"), ToZil(totalSsnInitialDelegateZil+4500))
	AssertEqual(Field(p.Azil, "total_supply"), ToAzil(totalSsnInitialDelegateZil+4500))
	AssertEqual(Field(p.Azil, "tmp_complete_withdrawal_available"), "0")

	AssertEqual(Field(p.Azil, "balances", sdk.Cfg.Addr1), ToAzil(4500))
	AssertEqual(Field(p.Azil, "withdrawal_unbonded"), "{}")
	AssertEqual(Field(p.Azil, "withdrawal_pending"), "{}")
}
