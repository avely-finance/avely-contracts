package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
	"strconv"
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
	tx, err := p.Aimpl.ClaimWithdrawal(readyBlocks)
	AssertError(tx, err, -105)

	delta, _ := strconv.ParseInt(StrAdd(p.Zimpl.Field("bnum_req"), "1"), 10, 32)
	sdk.IncreaseBlocknum(int32(delta))
	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	p.Aimpl.UpdateWallet(sdk.Cfg.AdminKey)
	tx, _ = p.Aimpl.ClaimWithdrawal(readyBlocks)
	AssertTransition(tx, Transition{
		p.Aimpl.Addr,         //sender
		"CompleteWithdrawal", //tag
		p.Holder.Addr,        //recipient
		"0",                  //amount
		ParamsMap{},
	})
	AssertEvent(tx, Event{p.Holder.Addr, "AddFunds", ParamsMap{"funder": "0x" + p.Zimpl.Addr, "amount": ToZil(10)}})

	AssertTransition(tx, Transition{
		p.Holder.Addr,                       //sender
		"CompleteWithdrawalSuccessCallBack", //tag
		p.Aimpl.Addr,                        //recipient
		ToZil(10),                           //amount
		ParamsMap{},
	})

	p.Aimpl.UpdateWallet(sdk.Cfg.Key1)
	tx, _ = p.Aimpl.CompleteWithdrawal()
	AssertEvent(tx, Event{p.Aimpl.Addr, "CompleteWithdrawal", ParamsMap{"amount": ToZil(10), "delegator": "0x" + sdk.Cfg.Addr1}})
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

	AssertEqual(ToZil(1000), p.Aimpl.Field("totalstakeamount"))
	AssertEqual(ToAzil(1000), p.Aimpl.Field("totaltokenamount"))
	AssertEqual("0", p.Aimpl.Field("tmp_complete_withdrawal_available"))
	AssertEqual(p.Aimpl.Field("balances", "0x"+sdk.Cfg.Admin), ToAzil(1000))
	AssertEqual("empty", p.Aimpl.Field("withdrawal_unbonded"))
	AssertEqual("empty", p.Aimpl.Field("withdrawal_pending"))
}
