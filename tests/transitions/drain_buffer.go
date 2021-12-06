package transitions

import (
	// "log"
	"Azil/test/deploy"
	//"math/big"
)

func (t *Testing) DrainBuffer() {
	t.LogStart("CompleteWithdrawal - success")

	stubStakingContract, aZilContract, bufferContract, _ := t.DeployAndUpgrade()

	aZilContract.DelegateStake(zil10)

	tx, _ := aZilContract.DrainBuffer(aZilContract.Addr)
	t.AssertError(tx, -106)

	tx, err := aZilContract.DrainBuffer(bufferContract.Addr)
	t.AssertSuccessCall(err)

	t.AssertTransition(tx, deploy.Transition{
		aZilContract.Addr,    //sender
		"ClaimRewards",       //tag
		bufferContract.Addr,  //recipient
		"0",                  //amount
		deploy.ParamsMap{},
	})

	// Send funds and call a callback
	t.AssertTransition(tx, deploy.Transition{
		stubStakingContract.Addr, //sender
		"AddFunds",
		bufferContract.Addr,
		"1000000000000", // 1 ZIL
		deploy.ParamsMap{},
	})

	t.AssertTransition(tx, deploy.Transition{
		stubStakingContract.Addr,  //sender
		"WithdrawStakeRewardsSuccessCallBack",
		bufferContract.Addr,
		"0",
		deploy.ParamsMap{"rewards": zil1},
	})

	// Check aZIL balance
	aZilContractState := aZilContract.LogContractStateJson()
	t.AssertContain(aZilContractState, "_balance\":\""+zil1)
}
