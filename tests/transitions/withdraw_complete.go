package transitions

import (
	. "Azil/test/helpers"
	"strconv"
)

func (tr *Transitions) CompleteWithdrawalSuccess() {

	t.Start("CompleteWithdrawal - success")
	readyBlocks := []string{}

	Zproxy, Zimpl, Aimpl, Buffer, Holder := tr.DeployAndUpgrade()

	Aimpl.UpdateWallet(tr.cfg.Key1)
	t.AssertSuccess(Aimpl.DelegateStake(Zil(10)))

	t.AssertSuccess(Zproxy.AssignStakeReward(tr.cfg.AzilSsnAddress, tr.cfg.AzilSsnRewardSharePercent))

	IncreaseBlocknum(10)
	t.AssertSuccess(Zproxy.AssignStakeReward(tr.cfg.AzilSsnAddress, tr.cfg.AzilSsnRewardSharePercent))

	Aimpl.UpdateWallet(tr.cfg.AdminKey)
	t.AssertSuccess(Aimpl.DrainBuffer(Buffer.Addr))

	Aimpl.UpdateWallet(tr.cfg.Key1)
	tx, _ := t.AssertSuccess(Aimpl.WithdrawStakeAmt(Azil(10)))

	block1 := tx.Receipt.EpochNum
	tx, _ = Aimpl.CompleteWithdrawal()
	t.AssertEvent(tx, Event{Aimpl.Addr, "NoUnbondedStake", ParamsMap{}})

	Aimpl.UpdateWallet(tr.cfg.Key2)
	tx, _ = Aimpl.CompleteWithdrawal()
	t.AssertEvent(tx, Event{Aimpl.Addr, "NoUnbondedStake", ParamsMap{}})

	readyBlocks = append(readyBlocks, block1)
	tx, err := Aimpl.ClaimWithdrawal(readyBlocks)
	t.AssertError(tx, err, -105)

	delta, _ := strconv.ParseInt(StrAdd(Zimpl.Field("bnum_req"), "1"), 10, 32)
	IncreaseBlocknum(int32(delta))
	t.AssertSuccess(Zproxy.AssignStakeReward(tr.cfg.AzilSsnAddress, tr.cfg.AzilSsnRewardSharePercent))

	Aimpl.UpdateWallet(tr.cfg.AdminKey)
	tx, _ = Aimpl.ClaimWithdrawal(readyBlocks)
	t.AssertTransition(tx, Transition{
		Aimpl.Addr,           //sender
		"CompleteWithdrawal", //tag
		Holder.Addr,          //recipient
		"0",                  //amount
		ParamsMap{},
	})
	t.AssertEvent(tx, Event{Holder.Addr, "AddFunds", ParamsMap{"funder": "0x" + Zimpl.Addr, "amount": Zil(10)}})

	t.AssertTransition(tx, Transition{
		Holder.Addr,                         //sender
		"CompleteWithdrawalSuccessCallBack", //tag
		Aimpl.Addr,                          //recipient
		Zil(10),                             //amount
		ParamsMap{},
	})

	Aimpl.UpdateWallet(tr.cfg.Key1)
	tx, _ = Aimpl.CompleteWithdrawal()
	t.AssertEvent(tx, Event{Aimpl.Addr, "CompleteWithdrawal", ParamsMap{"amount": Zil(10), "delegator": "0x" + tr.cfg.Addr1}})
	t.AssertTransition(tx, Transition{
		Aimpl.Addr,
		"CompleteWithdrawalSuccessCallBack",
		tr.cfg.Addr1,
		"0",
		ParamsMap{"amount": Zil(10)},
	})
	t.AssertTransition(tx, Transition{
		Aimpl.Addr,
		"AddFunds",
		tr.cfg.Addr1,
		Zil(10),
		ParamsMap{},
	})

	t.AssertEqual(Zil(1000), Aimpl.Field("totalstakeamount"))
	t.AssertEqual(Azil(1000), Aimpl.Field("totaltokenamount"))
	t.AssertEqual("0", Aimpl.Field("tmp_complete_withdrawal_available"))
	t.AssertEqual(Aimpl.Field("balances", "0x"+tr.cfg.Admin), Azil(1000))
	t.AssertEqual("empty", Aimpl.Field("withdrawal_unbonded"))
	t.AssertEqual("empty", Aimpl.Field("withdrawal_pending"))
}
