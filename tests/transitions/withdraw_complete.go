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

	Start("CompleteWithdrawal multi ssn")
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

	ssnWhitelistHeavy := p.GetSsnAddressForInput()

	//for current test setup first SSN for input is AzilSSN
	AssertEqual(ssnWhitelistHeavy, sdk.Cfg.AzilSsnAddress)

	AssertSuccess(p.Azil.WithUser(sdk.Cfg.Key1).DelegateStake(ToZil(5000)))

	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	ssnWhitelistLight := p.GetSsnAddressForInput()
	AssertNotEqual(ssnWhitelistHeavy, ssnWhitelistLight)
	AssertSuccess(p.Azil.WithUser(sdk.Cfg.Key1).DelegateStake(ToZil(5000)))
	ssnSlashHeavy := p.GetSsnAddressForInput()
	AssertNotEqual(ssnWhitelistLight, ssnSlashHeavy)
	AssertNotEqual(ssnWhitelistHeavy, ssnSlashHeavy)
	AssertSuccess(p.Azil.WithUser(sdk.Cfg.Key1).DelegateStake(ToZil(4000)))
	ssnSlashLight := p.GetSsnAddressForInput()
	AssertNotEqual(ssnWhitelistLight, ssnSlashLight)
	AssertNotEqual(ssnWhitelistHeavy, ssnSlashLight)
	AssertNotEqual(ssnSlashHeavy, ssnSlashLight)
	AssertSuccess(p.Azil.WithUser(sdk.Cfg.Key1).DelegateStake(ToZil(3000)))
	AssertEqual(Field(p.Azil, "totalstakeamount"), ToZil(totalSsnInitialDelegateZil+5000+5000+4000+3000))
	AssertEqual(Field(p.Azil, "total_supply"), ToAzil(totalSsnInitialDelegateZil+5000+5000+4000+3000))

	AssertEqual(Field(p.Azil, "balances", sdk.Cfg.Addr1), ToAzil(5000+5000+4000+3000))

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
	tx, _ := p.Azil.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToAzil(7000))
	AssertError(tx, "WithdrawAmountTooBig")

	//slash SSNs
	AssertSuccess(p.Azil.WithUser(p.Azil.Sdk.Cfg.OwnerKey).RemoveSSN(ssnSlashHeavy))
	AssertSuccess(p.Azil.WithUser(p.Azil.Sdk.Cfg.OwnerKey).RemoveSSN(ssnSlashLight))

	//withdraw and check from which SSN stake will be withdrawn
	tx, _ = AssertSuccess(p.Azil.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToAzil(3000)))
	//first is from heaviest slashed SSN
	AssertTransition(tx, Transition{
		p.Azil.Addr,        //sender
		"WithdrawStakeAmt", //tag
		p.Holder.Addr,      //recipient
		"0",                //amount
		ParamsMap{"amount": ToZil(3000), "ssnaddr": ssnSlashHeavy},
	})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnSlashHeavy), ToZil(1000))

	//next withdraw is going from current heavisest SSN
	tx, _ = AssertSuccess(p.Azil.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToAzil(3000)))
	AssertTransition(tx, Transition{
		p.Azil.Addr,        //sender
		"WithdrawStakeAmt", //tag
		p.Holder.Addr,      //recipient
		"0",                //amount
		ParamsMap{"amount": ToZil(3000), "ssnaddr": ssnSlashLight},
	})
	//there is nothing on this SSN now
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnSlashLight), "")

	//withdraw rest from ssnSlashHeavy
	tx, _ = AssertSuccess(p.Azil.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToAzil(1000)))
	AssertTransition(tx, Transition{
		p.Azil.Addr,        //sender
		"WithdrawStakeAmt", //tag
		p.Holder.Addr,      //recipient
		"0",                //amount
		ParamsMap{"amount": ToZil(1000), "ssnaddr": ssnSlashHeavy},
	})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnSlashHeavy), "")

	//there are no balance on slashed SSNs now, so withdraw will go from heaviest whitelisted SSN
	tx, _ = AssertSuccess(p.Azil.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToAzil(5000)))
	AssertTransition(tx, Transition{
		p.Azil.Addr,        //sender
		"WithdrawStakeAmt", //tag
		p.Holder.Addr,      //recipient
		"0",                //amount
		ParamsMap{"amount": ToZil(5000), "ssnaddr": ssnWhitelistHeavy},
	})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnWhitelistHeavy), ToZil(sdk.Cfg.HolderInitialDelegateZil))

	//next withdraw will go from ssnWhitelistLight, because it's heaviest now
	tx, _ = AssertSuccess(p.Azil.WithUser(sdk.Cfg.Key1).WithdrawStakeAmt(ToAzil(5000)))
	AssertTransition(tx, Transition{
		p.Azil.Addr,        //sender
		"WithdrawStakeAmt", //tag
		p.Holder.Addr,      //recipient
		"0",                //amount
		ParamsMap{"amount": ToZil(5000), "ssnaddr": ssnWhitelistLight},
	})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, ssnWhitelistLight), "")

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
	AssertEqual(withdrawal.TokenAmount.String(), ToAzil(17000))
	AssertEqual(withdrawal.StakeAmount.String(), ToAzil(17000))

	p.Azil.UpdateWallet(sdk.Cfg.AdminKey)
	tx, _ = p.Azil.ClaimWithdrawal(readyBlocks)
	AssertTransition(tx, Transition{
		p.Azil.Addr,          //sender
		"CompleteWithdrawal", //tag
		p.Holder.Addr,        //recipient
		"0",                  //amount
		ParamsMap{},
	})
	AssertEvent(tx, Event{p.Holder.Addr, "AddFunds", ParamsMap{"funder": p.Zimpl.Addr, "amount": ToZil(17000)}})

	AssertTransition(tx, Transition{
		p.Holder.Addr,                       //sender
		"CompleteWithdrawalSuccessCallBack", //tag
		p.Azil.Addr,                         //recipient
		ToZil(17000),                        //amount
		ParamsMap{},
	})

	p.Azil.UpdateWallet(sdk.Cfg.Key1)
	tx, _ = AssertSuccess(p.Azil.CompleteWithdrawal())
	AssertEvent(tx, Event{p.Azil.Addr, "CompleteWithdrawal", ParamsMap{"amount": ToZil(17000), "delegator": sdk.Cfg.Addr1}})
	AssertTransition(tx, Transition{
		p.Azil.Addr,
		"CompleteWithdrawalSuccessCallBack",
		sdk.Cfg.Addr1,
		"0",
		ParamsMap{"amount": ToZil(17000)},
	})
	AssertTransition(tx, Transition{
		p.Azil.Addr,
		"AddFunds",
		sdk.Cfg.Addr1,
		ToZil(17000),
		ParamsMap{},
	})

	withdrawal = Dig(p.Azil, "withdrawal_pending_of_delegator", sdk.Cfg.Addr1, block1).Withdrawal()
	AssertEqual(withdrawal.TokenAmount.String(), "0")
	AssertEqual(withdrawal.StakeAmount.String(), "0")

	AssertEqual(Field(p.Azil, "totalstakeamount"), ToZil(totalSsnInitialDelegateZil))
	AssertEqual(Field(p.Azil, "total_supply"), ToAzil(totalSsnInitialDelegateZil))
	AssertEqual(Field(p.Azil, "tmp_complete_withdrawal_available"), "0")

	AssertEqual(Field(p.Azil, "balances", sdk.Cfg.Addr1), "")
	AssertEqual(Field(p.Azil, "withdrawal_unbonded"), "{}")
	AssertEqual(Field(p.Azil, "withdrawal_pending"), "{}")
}
