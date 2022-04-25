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

	p.StZIL.UpdateWallet(sdk.Cfg.Key1)
	AssertSuccess(p.StZIL.DelegateStake(ToZil(10)))

	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	tx, _ := AssertSuccess(p.StZIL.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToStZil(10)))

	block1 := tx.Receipt.EpochNum
	tx, _ = p.StZIL.CompleteWithdrawal()
	AssertEvent(tx, Event{p.StZIL.Addr, "NoUnbondedStake", ParamsMap{}})

	p.StZIL.UpdateWallet(sdk.Cfg.Key2)
	tx, _ = p.StZIL.CompleteWithdrawal()
	AssertEvent(tx, Event{p.StZIL.Addr, "NoUnbondedStake", ParamsMap{}})

	p.StZIL.UpdateWallet(sdk.Cfg.AdminKey)
	readyBlocks = append(readyBlocks, block1)
	tx, _ = p.StZIL.ClaimWithdrawal(readyBlocks)
	AssertError(tx, "ClaimWithdrawalNoUnbonded")

	delta, _ := strconv.ParseInt(StrAdd(Field(p.Zimpl, "bnum_req"), "1"), 10, 32)
	sdk.IncreaseBlocknum(int32(delta))

	tr.NextCycle(p)
	tools := tr.NextCycleOffchain(p)

	unbondedWithdrawalsBlocks := p.GetClaimWithdrawalBlocks()
	AssertEqual(readyBlocks[0], strconv.Itoa(unbondedWithdrawalsBlocks[0]))
	tools.ShowClaimWithdrawal(p)

	withdrawal := Dig(p.StZIL, "withdrawal_pending_of_delegator", sdk.Cfg.Addr1, block1).Withdrawal()
	AssertEqual(withdrawal.TokenAmount.String(), ToStZil(10))
	AssertEqual(withdrawal.StakeAmount.String(), ToStZil(10))

	p.StZIL.UpdateWallet(sdk.Cfg.AdminKey)
	tx, _ = p.StZIL.ClaimWithdrawal(readyBlocks)
	AssertTransition(tx, Transition{
		p.StZIL.Addr,         //sender
		"CompleteWithdrawal", //tag
		p.Holder.Addr,        //recipient
		"0",                  //amount
		ParamsMap{},
	})
	AssertEvent(tx, Event{p.Holder.Addr, "AddFunds", ParamsMap{"funder": p.Zimpl.Addr, "amount": ToZil(10)}})

	AssertTransition(tx, Transition{
		p.Holder.Addr,                       //sender
		"CompleteWithdrawalSuccessCallBack", //tag
		p.StZIL.Addr,                        //recipient
		ToZil(10),                           //amount
		ParamsMap{},
	})

	p.StZIL.UpdateWallet(sdk.Cfg.Key1)
	tx, _ = p.StZIL.CompleteWithdrawal()
	AssertEvent(tx, Event{p.StZIL.Addr, "CompleteWithdrawal", ParamsMap{"amount": ToZil(10), "delegator": sdk.Cfg.Addr1}})
	AssertTransition(tx, Transition{
		p.StZIL.Addr,
		"CompleteWithdrawalSuccessCallBack",
		sdk.Cfg.Addr1,
		"0",
		ParamsMap{"amount": ToZil(10)},
	})
	AssertTransition(tx, Transition{
		p.StZIL.Addr,
		"AddFunds",
		sdk.Cfg.Addr1,
		ToZil(10),
		ParamsMap{},
	})

	withdrawal = Dig(p.StZIL, "withdrawal_pending_of_delegator", sdk.Cfg.Addr1, block1).Withdrawal()
	AssertEqual(withdrawal.TokenAmount.String(), "0")
	AssertEqual(withdrawal.StakeAmount.String(), "0")

	AssertEqual(Field(p.StZIL, "totalstakeamount"), ToZil(totalSsnInitialDelegateZil))
	AssertEqual(Field(p.StZIL, "total_supply"), ToStZil(totalSsnInitialDelegateZil))
	AssertEqual(Field(p.StZIL, "tmp_complete_withdrawal_available"), "0")

	if totalSsnInitialDelegateZil == 0 {
		AssertEqual(Field(p.StZIL, "balances"), "{}")
		AssertEqual(Field(p.StZIL, "balances", sdk.Cfg.Admin), "")
	} else {
		AssertEqual(Field(p.StZIL, "balances", sdk.Cfg.Admin), ToStZil(totalSsnInitialDelegateZil))
	}

	AssertEqual(Field(p.StZIL, "withdrawal_unbonded"), "{}")
	AssertEqual(Field(p.StZIL, "withdrawal_pending"), "{}")
}

func (tr *Transitions) CompleteWithdrawalMultiSsn() {

	Start("CompleteWithdrawal multi ssn")
	readyBlocks := []string{}

	p := tr.DeployAndUpgrade()

	rewardsFee := "1000" //10% of feeDenom=10000
	treasuryAddr := sdk.Cfg.Addr3
	AssertSuccess(p.StZIL.WithUser(sdk.Cfg.OwnerKey).ChangeRewardsFee(rewardsFee))
	AssertSuccess(p.StZIL.WithUser(sdk.Cfg.OwnerKey).ChangeTreasuryAddress(treasuryAddr))
	p.StZIL.UpdateWallet(sdk.Cfg.AdminKey) //back to admin

	//totalSsnInitialDelegateZil := len(sdk.Cfg.SsnAddrs) * sdk.Cfg.SsnInitialDelegateZil
	//for now to activate SSNs we delegate required stakes through Zproxy as admin
	totalSsnInitialDelegateZil := 0

	ssnWhitelistHeavy := p.GetSsnAddressForInput()

	//for current test setup first SSN for input is StZILSSN
	AssertEqual(ssnWhitelistHeavy, sdk.Cfg.StZilSsnAddress)

	AssertSuccess(p.StZIL.WithUser(sdk.Cfg.Key1).DelegateStake(ToZil(5000)))

	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	ssnWhitelistLight := p.GetSsnAddressForInput()
	AssertNotEqual(ssnWhitelistHeavy, ssnWhitelistLight)
	AssertSuccess(p.StZIL.WithUser(sdk.Cfg.Key1).DelegateStake(ToZil(5000)))
	ssnSlashHeavy := p.GetSsnAddressForInput()
	AssertNotEqual(ssnWhitelistLight, ssnSlashHeavy)
	AssertNotEqual(ssnWhitelistHeavy, ssnSlashHeavy)
	AssertSuccess(p.StZIL.WithUser(sdk.Cfg.Key1).DelegateStake(ToZil(4000)))
	ssnSlashLight := p.GetSsnAddressForInput()
	AssertNotEqual(ssnWhitelistLight, ssnSlashLight)
	AssertNotEqual(ssnWhitelistHeavy, ssnSlashLight)
	AssertNotEqual(ssnSlashHeavy, ssnSlashLight)
	AssertSuccess(p.StZIL.WithUser(sdk.Cfg.Key1).DelegateStake(ToZil(3000)))
	AssertEqual(Field(p.StZIL, "totalstakeamount"), ToZil(totalSsnInitialDelegateZil+5000+5000+4000+3000))
	AssertEqual(Field(p.StZIL, "total_supply"), ToStZil(totalSsnInitialDelegateZil+5000+5000+4000+3000))

	AssertEqual(Field(p.StZIL, "balances", sdk.Cfg.Addr1), ToStZil(5000+5000+4000+3000))

	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	//stake is on holder now, splitted between SSNs
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnWhitelistHeavy), ToZil(sdk.Cfg.HolderInitialDelegateZil+5000))
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnWhitelistLight), ToZil(5000))
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnSlashHeavy), ToZil(4000))
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnSlashLight), ToZil(3000))

	//it's impossible to withdraw amount, bigger than amount on heaviest SSN
	tx, _ := p.StZIL.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToStZil(7000))
	AssertError(tx, "WithdrawAmountTooBig")

	//slash SSNs
	AssertSuccess(p.StZIL.WithUser(p.StZIL.Sdk.Cfg.OwnerKey).RemoveSSN(ssnSlashHeavy))
	AssertSuccess(p.StZIL.WithUser(p.StZIL.Sdk.Cfg.OwnerKey).RemoveSSN(ssnSlashLight))

	//withdraw and check from which SSN stake will be withdrawn
	tx, _ = AssertSuccess(p.StZIL.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToStZil(3000)))
	//first is from heaviest slashed SSN
	AssertTransition(tx, Transition{
		p.StZIL.Addr,       //sender
		"WithdrawStakeAmt", //tag
		p.Holder.Addr,      //recipient
		"0",                //amount
		ParamsMap{"amount": ToZil(3000), "ssnaddr": ssnSlashHeavy},
	})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnSlashHeavy), ToZil(1000))

	//there are not enough balance on slashed SSNs now, so withdraw will go from heaviest whitelisted SSN
	tx, _ = AssertSuccess(p.StZIL.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToStZil(5000)))
	AssertTransition(tx, Transition{
		p.StZIL.Addr,       //sender
		"WithdrawStakeAmt", //tag
		p.Holder.Addr,      //recipient
		"0",                //amount
		ParamsMap{"amount": ToZil(5000), "ssnaddr": ssnWhitelistHeavy},
	})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnWhitelistHeavy), ToZil(sdk.Cfg.HolderInitialDelegateZil))

	//next withdraw is going from current heaviest slashed SSN (it was not heaviest before, but now it is)
	tx, _ = AssertSuccess(p.StZIL.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToStZil(3000)))
	AssertTransition(tx, Transition{
		p.StZIL.Addr,       //sender
		"WithdrawStakeAmt", //tag
		p.Holder.Addr,      //recipient
		"0",                //amount
		ParamsMap{"amount": ToZil(3000), "ssnaddr": ssnSlashLight},
	})
	//there is nothing on this SSN now
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnSlashLight), "")

	//withdraw rest from ssnSlashHeavy
	tx, _ = AssertSuccess(p.StZIL.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToStZil(1000)))
	AssertTransition(tx, Transition{
		p.StZIL.Addr,       //sender
		"WithdrawStakeAmt", //tag
		p.Holder.Addr,      //recipient
		"0",                //amount
		ParamsMap{"amount": ToZil(1000), "ssnaddr": ssnSlashHeavy},
	})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnSlashHeavy), "")

	//next withdraw will go from ssnWhitelistLight, because it's heaviest now
	tx, _ = AssertSuccess(p.StZIL.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToStZil(5000)))
	AssertTransition(tx, Transition{
		p.StZIL.Addr,       //sender
		"WithdrawStakeAmt", //tag
		p.Holder.Addr,      //recipient
		"0",                //amount
		ParamsMap{"amount": ToZil(5000), "ssnaddr": ssnWhitelistLight},
	})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnWhitelistLight), "")

	block1 := tx.Receipt.EpochNum
	tx, _ = p.StZIL.CompleteWithdrawal()
	AssertEvent(tx, Event{p.StZIL.Addr, "NoUnbondedStake", ParamsMap{}})

	p.StZIL.UpdateWallet(sdk.Cfg.Key2)
	tx, _ = p.StZIL.CompleteWithdrawal()
	AssertEvent(tx, Event{p.StZIL.Addr, "NoUnbondedStake", ParamsMap{}})

	p.StZIL.UpdateWallet(sdk.Cfg.AdminKey)
	readyBlocks = append(readyBlocks, block1)
	tx, _ = p.StZIL.ClaimWithdrawal(readyBlocks)
	AssertError(tx, "ClaimWithdrawalNoUnbonded")

	delta, _ := strconv.ParseInt(StrAdd(Field(p.Zimpl, "bnum_req"), "1"), 10, 32)
	sdk.IncreaseBlocknum(int32(delta))

	tr.NextCycle(p)
	tools := tr.NextCycleOffchain(p)

	unbondedWithdrawalsBlocks := p.GetClaimWithdrawalBlocks()
	AssertEqual(readyBlocks[0], strconv.Itoa(unbondedWithdrawalsBlocks[0]))
	tools.ShowClaimWithdrawal(p)

	withdrawal := Dig(p.StZIL, "withdrawal_pending_of_delegator", sdk.Cfg.Addr1, block1).Withdrawal()
	AssertEqual(withdrawal.TokenAmount.String(), ToStZil(17000))
	AssertEqual(withdrawal.StakeAmount.String(), ToStZil(17000))

	p.StZIL.UpdateWallet(sdk.Cfg.AdminKey)
	tx, _ = p.StZIL.ClaimWithdrawal(readyBlocks)
	AssertTransition(tx, Transition{
		p.StZIL.Addr,         //sender
		"CompleteWithdrawal", //tag
		p.Holder.Addr,        //recipient
		"0",                  //amount
		ParamsMap{},
	})
	AssertEvent(tx, Event{p.Holder.Addr, "AddFunds", ParamsMap{"funder": p.Zimpl.Addr, "amount": ToZil(17000)}})

	AssertTransition(tx, Transition{
		p.Holder.Addr,                       //sender
		"CompleteWithdrawalSuccessCallBack", //tag
		p.StZIL.Addr,                        //recipient
		ToZil(17000),                        //amount
		ParamsMap{},
	})

	p.StZIL.UpdateWallet(sdk.Cfg.Key1)
	tx, _ = AssertSuccess(p.StZIL.CompleteWithdrawal())
	AssertEvent(tx, Event{p.StZIL.Addr, "CompleteWithdrawal", ParamsMap{"amount": ToZil(17000), "delegator": sdk.Cfg.Addr1}})
	AssertTransition(tx, Transition{
		p.StZIL.Addr,
		"CompleteWithdrawalSuccessCallBack",
		sdk.Cfg.Addr1,
		"0",
		ParamsMap{"amount": ToZil(17000)},
	})
	AssertTransition(tx, Transition{
		p.StZIL.Addr,
		"AddFunds",
		sdk.Cfg.Addr1,
		ToZil(17000),
		ParamsMap{},
	})

	withdrawal = Dig(p.StZIL, "withdrawal_pending_of_delegator", sdk.Cfg.Addr1, block1).Withdrawal()
	AssertEqual(withdrawal.TokenAmount.String(), "0")
	AssertEqual(withdrawal.StakeAmount.String(), "0")

	AssertEqual(Field(p.StZIL, "totalstakeamount"), ToZil(totalSsnInitialDelegateZil))
	AssertEqual(Field(p.StZIL, "total_supply"), ToStZil(totalSsnInitialDelegateZil))
	AssertEqual(Field(p.StZIL, "tmp_complete_withdrawal_available"), "0")

	AssertEqual(Field(p.StZIL, "balances", sdk.Cfg.Addr1), "")
	AssertEqual(Field(p.StZIL, "withdrawal_unbonded"), "{}")
	AssertEqual(Field(p.StZIL, "withdrawal_pending"), "{}")
}
