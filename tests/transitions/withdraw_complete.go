package transitions

import (
	"Azil/test/helpers"
	"strconv"
)

func (tr *Transitions) CompleteWithdrawalSuccess() {

	log.Start("CompleteWithdrawal - success")
	readyBlocks := []string{}

	Zproxy, Zimpl, Aimpl, Buffer, Holder := tr.DeployAndUpgrade()

	Aimpl.UpdateWallet(key1)
	t.AssertSuccess(Aimpl.DelegateStake(zil(10)))

	t.AssertSuccess(Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT))

	helpers.IncreaseBlocknum(10)
	t.AssertSuccess(Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT))

	Aimpl.UpdateWallet(adminKey)
	t.AssertSuccess(Aimpl.DrainBuffer(Buffer.Addr))

	Aimpl.UpdateWallet(key1)
	tx, _ := t.AssertSuccess(Aimpl.WithdrawStakeAmt(azil(10)))

	block1 := tx.Receipt.EpochNum
	tx, _ = Aimpl.CompleteWithdrawal()
	t.AssertEvent(tx, helpers.Event{Aimpl.Addr, "NoUnbondedStake", helpers.ParamsMap{}})

	Aimpl.UpdateWallet(key2)
	tx, _ = Aimpl.CompleteWithdrawal()
	t.AssertEvent(tx, helpers.Event{Aimpl.Addr, "NoUnbondedStake", helpers.ParamsMap{}})

	readyBlocks = append(readyBlocks, block1)
	tx, err := Aimpl.ClaimWithdrawal(readyBlocks)
	t.AssertError(tx, err, -105)

	delta, _ := strconv.ParseInt(helpers.StrAdd(Zimpl.Field("bnum_req"), "1"), 10, 32)
	helpers.IncreaseBlocknum(int32(delta))
	t.AssertSuccess(Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT))

	Aimpl.UpdateWallet(adminKey)
	tx, _ = Aimpl.ClaimWithdrawal(readyBlocks)
	t.AssertTransition(tx, helpers.Transition{
		Aimpl.Addr,           //sender
		"CompleteWithdrawal", //tag
		Holder.Addr,          //recipient
		"0",                  //amount
		helpers.ParamsMap{},
	})
	t.AssertEvent(tx, helpers.Event{Holder.Addr, "AddFunds", helpers.ParamsMap{"funder": "0x" + Zimpl.Addr, "amount": zil(10)}})

	t.AssertTransition(tx, helpers.Transition{
		Holder.Addr,                         //sender
		"CompleteWithdrawalSuccessCallBack", //tag
		Aimpl.Addr,                          //recipient
		zil(10),                             //amount
		helpers.ParamsMap{},
	})

	Aimpl.UpdateWallet(key1)
	tx, _ = Aimpl.CompleteWithdrawal()
	t.AssertEvent(tx, helpers.Event{Aimpl.Addr, "CompleteWithdrawal", helpers.ParamsMap{"amount": zil(10), "delegator": "0x" + addr1}})
	t.AssertTransition(tx, helpers.Transition{
		Aimpl.Addr,
		"CompleteWithdrawalSuccessCallBack",
		addr1,
		"0",
		helpers.ParamsMap{"amount": zil(10)},
	})
	t.AssertTransition(tx, helpers.Transition{
		Aimpl.Addr,
		"AddFunds",
		addr1,
		zil(10),
		helpers.ParamsMap{},
	})

	t.AssertEqual(zil(1000), Aimpl.Field("totalstakeamount"))
	t.AssertEqual(azil(1000), Aimpl.Field("totaltokenamount"))
	t.AssertEqual("0", Aimpl.Field("tmp_complete_withdrawal_available"))
	t.AssertEqual(Aimpl.Field("balances", "0x"+admin), azil(1000))
	t.AssertEqual("empty", Aimpl.Field("withdrawal_unbonded"))
	t.AssertEqual("empty", Aimpl.Field("withdrawal_pending"))
}
