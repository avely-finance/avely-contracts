package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
	"strconv"
)

func (tr *Transitions) CompleteWithdrawalSuccess() {

	t.Start("CompleteWithdrawal - success")
	readyBlocks := []string{}

	p := tr.DeployAndUpgrade()

	p.Aimpl.UpdateWallet(sdk.Cfg.Key1)
	t.AssertSuccess(p.Aimpl.DelegateStake(Zil(10)))

	t.AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	sdk.IncreaseBlocknum(10)
	t.AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	p.Aimpl.UpdateWallet(sdk.Cfg.AdminKey)
	t.AssertSuccess(p.Aimpl.DrainBuffer(p.Buffer.Addr))

	p.Aimpl.UpdateWallet(sdk.Cfg.Key1)
	tx, _ := t.AssertSuccess(p.Aimpl.WithdrawStakeAmt(Azil(10)))

	block1 := tx.Receipt.EpochNum
	tx, _ = p.Aimpl.CompleteWithdrawal()
	t.AssertEvent(tx, Event{p.Aimpl.Addr, "NoUnbondedStake", ParamsMap{}})

	p.Aimpl.UpdateWallet(sdk.Cfg.Key2)
	tx, _ = p.Aimpl.CompleteWithdrawal()
	t.AssertEvent(tx, Event{p.Aimpl.Addr, "NoUnbondedStake", ParamsMap{}})

	p.Aimpl.UpdateWallet(sdk.Cfg.AdminKey)
	readyBlocks = append(readyBlocks, block1)
	tx, err := p.Aimpl.ClaimWithdrawal(readyBlocks)
	t.AssertError(tx, err, -105)

	delta, _ := strconv.ParseInt(StrAdd(p.Zimpl.Field("bnum_req"), "1"), 10, 32)
	sdk.IncreaseBlocknum(int32(delta))
	t.AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	p.Aimpl.UpdateWallet(sdk.Cfg.AdminKey)
	tx, _ = p.Aimpl.ClaimWithdrawal(readyBlocks)
	t.AssertTransition(tx, Transition{
		p.Aimpl.Addr,           //sender
		"CompleteWithdrawal", //tag
		p.Holder.Addr,          //recipient
		"0",                  //amount
		ParamsMap{},
	})
	t.AssertEvent(tx, Event{p.Holder.Addr, "AddFunds", ParamsMap{"funder": "0x" + p.Zimpl.Addr, "amount": Zil(10)}})

	t.AssertTransition(tx, Transition{
		p.Holder.Addr,                         //sender
		"CompleteWithdrawalSuccessCallBack", //tag
		p.Aimpl.Addr,                          //recipient
		Zil(10),                             //amount
		ParamsMap{},
	})

	p.Aimpl.UpdateWallet(sdk.Cfg.Key1)
	tx, _ = p.Aimpl.CompleteWithdrawal()
	t.AssertEvent(tx, Event{p.Aimpl.Addr, "CompleteWithdrawal", ParamsMap{"amount": Zil(10), "delegator": "0x" + sdk.Cfg.Addr1}})
	t.AssertTransition(tx, Transition{
		p.Aimpl.Addr,
		"CompleteWithdrawalSuccessCallBack",
		sdk.Cfg.Addr1,
		"0",
		ParamsMap{"amount": Zil(10)},
	})
	t.AssertTransition(tx, Transition{
		p.Aimpl.Addr,
		"AddFunds",
		sdk.Cfg.Addr1,
		Zil(10),
		ParamsMap{},
	})

	t.AssertEqual(Zil(1000), p.Aimpl.Field("totalstakeamount"))
	t.AssertEqual(Azil(1000), p.Aimpl.Field("totaltokenamount"))
	t.AssertEqual("0", p.Aimpl.Field("tmp_complete_withdrawal_available"))
	t.AssertEqual(p.Aimpl.Field("balances", "0x"+sdk.Cfg.Admin), Azil(1000))
	t.AssertEqual("empty", p.Aimpl.Field("withdrawal_unbonded"))
	t.AssertEqual("empty", p.Aimpl.Field("withdrawal_pending"))
}
