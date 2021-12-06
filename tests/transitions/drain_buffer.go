package transitions

import (
	// "log"
	"Azil/test/deploy"
	//"math/big"
)

func (t *Testing) DrainBuffer() {
	t.LogStart("CompleteWithdrawal - success")

	stubStakingContract, aZilContract, bufferContract, holderContract := t.DeployAndUpgrade()

	aZilContract.DelegateStake(zil10)

	txn, err := aZilContract.DrainBuffer(aZilContract.Addr)
	t.AssertError(txn, err, -107)

	txn, _ = aZilContract.DrainBuffer(bufferContract.Addr)

	t.AssertTransition(txn, deploy.Transition{
		aZilContract.Addr,   //sender
		"ClaimRewards",      //tag
		bufferContract.Addr, //recipient
		"0",                 //amount
		deploy.ParamsMap{},
	})

	// Send funds and call a callback via Buffer
	t.AssertTransition(txn, deploy.Transition{
		stubStakingContract.Addr, //sender
		"AddFunds",
		bufferContract.Addr,
		"1000000000000", // 1 ZIL
		deploy.ParamsMap{},
	})

	t.AssertTransition(txn, deploy.Transition{
		stubStakingContract.Addr, //sender
		"WithdrawStakeRewardsSuccessCallBack",
		bufferContract.Addr,
		"0",
		deploy.ParamsMap{"rewards": zil(1)},
	})

	// Send funds and call a callback via Holder
	t.AssertTransition(txn, deploy.Transition{
		stubStakingContract.Addr, //sender
		"AddFunds",
		holderContract.Addr,
		"1000000000000", // 1 ZIL
		deploy.ParamsMap{},
	})

	t.AssertTransition(txn, deploy.Transition{
		stubStakingContract.Addr, //sender
		"WithdrawStakeRewardsSuccessCallBack",
		holderContract.Addr,
		"0",
		deploy.ParamsMap{"rewards": zil(1)},
	})

	// Check aZIL balance
	aZilContractState := aZilContract.LogContractStateJson()

	// 1 ZIL from Buffer + 1 ZIL from Holder
	t.AssertContain(aZilContractState, "_balance\":\""+zil(2))
}
