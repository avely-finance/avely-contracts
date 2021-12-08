package transitions

import (
	"Azil/test/deploy"
)

func (t *Testing) CompleteWithdrawalSuccess() {

	t.LogStart("CompleteWithdrawal - success")
	readyBlocks := []string{}

	stubStakingContract, aZilContract, _, holderContract := t.DeployAndUpgrade()
	t.AddDebug("addr1", "0x"+addr1)

	aZilContract.UpdateWallet(key1)
	aZilContract.DelegateStake(zil(10))

	stubStakingContract.AssignStakeReward()

	tx, err := aZilContract.WithdrawStakeAmt(azil(10))
	block1 := tx.Receipt.EpochNum
	tx, _ = aZilContract.CompleteWithdrawal()
	t.AssertEvent(tx, deploy.Event{aZilContract.Addr, "NoUnbondedStake", deploy.ParamsMap{}})

	aZilContract.UpdateWallet(key2)
	tx, _ = aZilContract.CompleteWithdrawal()
	t.AssertEvent(tx, deploy.Event{aZilContract.Addr, "NoUnbondedStake", deploy.ParamsMap{}})

	readyBlocks = append(readyBlocks, block1)
	tx, err = aZilContract.ClaimWithdrawal(readyBlocks)
	t.AssertError(tx, err, -105)

	deploy.IncreaseBlocknum(stubStakingContract.GetBnumReq() + 1)
	stubStakingContract.AssignStakeReward()

	aZilContract.UpdateWallet(adminKey)
	tx, err = aZilContract.ClaimWithdrawal(readyBlocks)
	t.AssertTransition(tx, deploy.Transition{
		aZilContract.Addr,    //sender
		"CompleteWithdrawal", //tag
		holderContract.Addr,  //recipient
		"0",                  //amount
		deploy.ParamsMap{},
	})
	t.AssertEvent(tx, deploy.Event{holderContract.Addr, "AddFunds", deploy.ParamsMap{"funder": "0x" + stubStakingContract.Addr, "amount": zil(10)}})

	t.AssertTransition(tx, deploy.Transition{
		holderContract.Addr,                 //sender
		"CompleteWithdrawalSuccessCallBack", //tag
		aZilContract.Addr,                   //recipient
		zil(10),                             //amount
		deploy.ParamsMap{},
	})

	aZilContract.UpdateWallet(key1)
	tx, _ = aZilContract.CompleteWithdrawal()
	t.AssertEvent(tx, deploy.Event{aZilContract.Addr, "CompleteWithdrawal", deploy.ParamsMap{"amount": zil(10), "delegator": "0x" + addr1}})
	t.AssertTransition(tx, deploy.Transition{
		aZilContract.Addr,
		"CompleteWithdrawalSuccessCallBack",
		addr1,
		"0",
		deploy.ParamsMap{"amount": zil(10)},
	})
	t.AssertTransition(tx, deploy.Transition{
		aZilContract.Addr,
		"AddFunds",
		addr1,
		zil(10),
		deploy.ParamsMap{},
	})

	t.AssertEqual("0", aZilContract.StateField("totalstakeamount"))
	t.AssertEqual("0", aZilContract.StateField("totaltokenamount"))
	t.AssertEqual("0", aZilContract.StateField("tmp_complete_withdrawal_available"))
	t.AssertEqual("empty", aZilContract.StateField("balances"))
	t.AssertEqual("empty", aZilContract.StateField("withdrawal_unbonded"))
	t.AssertEqual("empty", aZilContract.StateField("withdrawal_pending"))

}
