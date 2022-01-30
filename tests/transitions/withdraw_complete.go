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

	p.Aimpl.UpdateWallet(sdk.Cfg.Key1)
	AssertSuccess(p.Aimpl.DelegateStake(ToZil(10)))

	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	sdk.IncreaseBlocknum(10)

	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	p.Aimpl.UpdateWallet(sdk.Cfg.AdminKey)
	AssertSuccess(p.Aimpl.DrainBuffer(p.GetBuffer().Addr))

	p.Aimpl.UpdateWallet(sdk.Cfg.Key1)
	tx, _ := AssertSuccess(p.Aimpl.WithdrawStakeAmt(ToAzil(10)))

	block1 := tx.Receipt.EpochNum
	tx, _ = p.Aimpl.CompleteWithdrawal()
	AssertEvent(tx, Event{p.Aimpl.Addr, "NoUnbondedStake", ParamsMap{}})

	p.Aimpl.UpdateWallet(sdk.Cfg.Key2)
	tx, _ = p.Aimpl.CompleteWithdrawal()
	AssertEvent(tx, Event{p.Aimpl.Addr, "NoUnbondedStake", ParamsMap{}})

	p.Aimpl.UpdateWallet(sdk.Cfg.AdminKey)
	readyBlocks = append(readyBlocks, block1)
	tx, _ = p.Aimpl.ClaimWithdrawal(readyBlocks)
	AssertError(tx, "ClaimWithdrawalNoUnbonded")

	delta, _ := strconv.ParseInt(StrAdd(Field(p.Zimpl, "bnum_req"), "1"), 10, 32)
	sdk.IncreaseBlocknum(int32(delta))

	tx, err := p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare)
	AssertSuccess(tx, err)

	blockNumStr, _ := strconv.Atoi(tx.Receipt.EpochNum)
	unbondedWithdrawalsBlocks := p.GetUnbondedWithdrawalsBlocks(blockNumStr)
	AssertEqual(readyBlocks[0], strconv.Itoa(unbondedWithdrawalsBlocks[0]))

	p.Aimpl.UpdateWallet(sdk.Cfg.AdminKey)
	tx, _ = p.Aimpl.ClaimWithdrawal(readyBlocks)
	AssertTransition(tx, Transition{
		p.Aimpl.Addr,         //sender
		"CompleteWithdrawal", //tag
		p.Holder.Addr,        //recipient
		"0",                  //amount
		ParamsMap{},
	})
	AssertEvent(tx, Event{p.Holder.Addr, "AddFunds", ParamsMap{"funder": p.Zimpl.Addr, "amount": ToZil(10)}})

	AssertTransition(tx, Transition{
		p.Holder.Addr,                       //sender
		"CompleteWithdrawalSuccessCallBack", //tag
		p.Aimpl.Addr,                        //recipient
		ToZil(10),                           //amount
		ParamsMap{},
	})

	p.Aimpl.UpdateWallet(sdk.Cfg.Key1)
	tx, _ = p.Aimpl.CompleteWithdrawal()
	AssertEvent(tx, Event{p.Aimpl.Addr, "CompleteWithdrawal", ParamsMap{"amount": ToZil(10), "delegator": sdk.Cfg.Addr1}})
	AssertTransition(tx, Transition{
		p.Aimpl.Addr,
		"CompleteWithdrawalSuccessCallBack",
		sdk.Cfg.Addr1,
		"0",
		ParamsMap{"amount": ToZil(10)},
	})
	AssertTransition(tx, Transition{
		p.Aimpl.Addr,
		"AddFunds",
		sdk.Cfg.Addr1,
		ToZil(10),
		ParamsMap{},
	})

	AssertEqual(Field(p.Aimpl, "totalstakeamount"), ToZil(1000))
	AssertEqual(Field(p.Aimpl, "totaltokenamount"), ToAzil(1000))
	AssertEqual(Field(p.Aimpl, "tmp_complete_withdrawal_available"), "0")
	AssertEqual(Field(p.Aimpl, "balances", sdk.Cfg.Admin), ToAzil(1000))

	AssertEqual(Field(p.Aimpl, "withdrawal_unbonded"), "{}")
	AssertEqual(Field(p.Aimpl, "withdrawal_pending"), "{}")
}
