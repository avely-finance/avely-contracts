package transitions

import (
	"Azil/test/deploy"
	"strconv"
)

func (t *Testing) CompleteWithdrawalSuccess() {

	t.LogStart("CompleteWithdrawal - success")
	readyBlocks := []string{}

	Zproxy, _, Aimpl, _, Holder := t.DeployAndUpgrade()
	t.AddDebug("addr1", "0x"+addr1)

	Aimpl.UpdateWallet(key1)
	Aimpl.DelegateStake(zil(10))

	Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT)

	tx, err := Aimpl.WithdrawStakeAmt(azil(10))
	block1 := tx.Receipt.EpochNum
	tx, _ = Aimpl.CompleteWithdrawal()
	t.AssertEvent(tx, deploy.Event{Aimpl.Addr, "NoUnbondedStake", deploy.ParamsMap{}})

	Aimpl.UpdateWallet(key2)
	tx, _ = Aimpl.CompleteWithdrawal()
	t.AssertEvent(tx, deploy.Event{Aimpl.Addr, "NoUnbondedStake", deploy.ParamsMap{}})

	readyBlocks = append(readyBlocks, block1)
	tx, err = Aimpl.ClaimWithdrawal(readyBlocks)
	t.AssertError(tx, err, -105)

	delta, _ := strconv.ParseInt(deploy.StrSum(Zproxy.Field("bnum_req"), "1"), 10, 32)
	deploy.IncreaseBlocknum(int32(delta))
	Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT)

	Aimpl.UpdateWallet(adminKey)
	tx, err = Aimpl.ClaimWithdrawal(readyBlocks)
	t.AssertTransition(tx, deploy.Transition{
		Aimpl.Addr,           //sender
		"CompleteWithdrawal", //tag
		Holder.Addr,          //recipient
		"0",                  //amount
		deploy.ParamsMap{},
	})
	t.AssertEvent(tx, deploy.Event{Holder.Addr, "AddFunds", deploy.ParamsMap{"funder": "0x" + Zproxy.Addr, "amount": zil(10)}})

	t.AssertTransition(tx, deploy.Transition{
		Holder.Addr,                         //sender
		"CompleteWithdrawalSuccessCallBack", //tag
		Aimpl.Addr,                          //recipient
		zil(10),                             //amount
		deploy.ParamsMap{},
	})

	Aimpl.UpdateWallet(key1)
	tx, _ = Aimpl.CompleteWithdrawal()
	t.AssertEvent(tx, deploy.Event{Aimpl.Addr, "CompleteWithdrawal", deploy.ParamsMap{"amount": zil(10), "delegator": "0x" + addr1}})
	t.AssertTransition(tx, deploy.Transition{
		Aimpl.Addr,
		"CompleteWithdrawalSuccessCallBack",
		addr1,
		"0",
		deploy.ParamsMap{"amount": zil(10)},
	})
	t.AssertTransition(tx, deploy.Transition{
		Aimpl.Addr,
		"AddFunds",
		addr1,
		zil(10),
		deploy.ParamsMap{},
	})

	t.AssertEqual("0", Aimpl.Field("totalstakeamount"))
	t.AssertEqual("0", Aimpl.Field("totaltokenamount"))
	t.AssertEqual("0", Aimpl.Field("tmp_complete_withdrawal_available"))
	t.AssertEqual("empty", Aimpl.Field("balances"))
	t.AssertEqual("empty", Aimpl.Field("withdrawal_unbonded"))
	t.AssertEqual("empty", Aimpl.Field("withdrawal_pending"))

}
