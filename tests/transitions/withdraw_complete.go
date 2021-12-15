package transitions

import (
	"Azil/test/deploy"
	"strconv"
)

func (t *Testing) CompleteWithdrawalSuccess() {

	t.LogStart("CompleteWithdrawal - success")
	readyBlocks := []string{}

	Zproxy, Zimpl, Aimpl, Buffer, Holder := t.DeployAndUpgrade()
	t.AddDebug("addr1", "0x"+addr1)

	Aimpl.UpdateWallet(key1)
	Aimpl.DelegateStake(zil(10))

	Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT)

	deploy.IncreaseBlocknum(10)
	Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT)
	Aimpl.UpdateWallet(adminKey)
	tx, err := Aimpl.DrainBuffer(Buffer.Addr)
	if err != nil {
		t.LogError("Aimpl.DrainBuffer(Buffer.Addr) error = ", err)
	}

	Aimpl.UpdateWallet(key1)
	tx, err = Aimpl.WithdrawStakeAmt(azil(10))
	if err != nil {
		t.LogError(" Aimpl.WithdrawStakeAmt() error = ", err)
	}
	block1 := tx.Receipt.EpochNum
	tx, _ = Aimpl.CompleteWithdrawal()
	t.AssertEvent(tx, deploy.Event{Aimpl.Addr, "NoUnbondedStake", deploy.ParamsMap{}})

	Aimpl.UpdateWallet(key2)
	tx, _ = Aimpl.CompleteWithdrawal()
	t.AssertEvent(tx, deploy.Event{Aimpl.Addr, "NoUnbondedStake", deploy.ParamsMap{}})

	readyBlocks = append(readyBlocks, block1)
	tx, err = Aimpl.ClaimWithdrawal(readyBlocks)
	t.AssertError(tx, err, -105)

	delta, _ := strconv.ParseInt(deploy.StrAdd(Zimpl.Field("bnum_req"), "1"), 10, 32)
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
	t.AssertEvent(tx, deploy.Event{Holder.Addr, "AddFunds", deploy.ParamsMap{"funder": "0x" + Zimpl.Addr, "amount": zil(10)}})

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

	t.AssertEqual(zil(1000), Aimpl.Field("totalstakeamount"))
	t.AssertEqual(azil(1000), Aimpl.Field("totaltokenamount"))
	t.AssertEqual("0", Aimpl.Field("tmp_complete_withdrawal_available"))
	//t.AssertEqual("empty", Aimpl.Field("balances"))
	t.AssertEqual(Aimpl.Field("balances", "0x"+admin), azil(1000))
	t.AssertEqual("empty", Aimpl.Field("withdrawal_unbonded"))
	t.AssertEqual("empty", Aimpl.Field("withdrawal_pending"))

}
