package transitions

import (
	// "log"
	"Azil/test/deploy"
	//"math/big"
)

func (t *Testing) DrainBuffer() {
	t.LogStart("CompleteWithdrawal - success")

	Zproxy, _, aZilContract, bufferContract, holderContract := t.DeployAndUpgrade()

	aZilContract.DelegateStake(zil(10))

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
		Zproxy.Addr, //sender
		"AddFunds",
		bufferContract.Addr,
		zil(1),
		deploy.ParamsMap{},
	})

	t.AssertTransition(txn, deploy.Transition{
		Zproxy.Addr, //sender
		"WithdrawStakeRewardsSuccessCallBack",
		bufferContract.Addr,
		"0",
		deploy.ParamsMap{"rewards": zil(1)},
	})

	// Send funds and call a callback via Holder
	t.AssertTransition(txn, deploy.Transition{
		Zproxy.Addr, //sender
		"AddFunds",
		holderContract.Addr,
		zil(1),
		deploy.ParamsMap{},
	})

	t.AssertTransition(txn, deploy.Transition{
		Zproxy.Addr, //sender
		"WithdrawStakeRewardsSuccessCallBack",
		holderContract.Addr,
		"0",
		deploy.ParamsMap{"rewards": zil(1)},
	})

	// Check aZIL balance
	// 1 ZIL from Buffer + 1 ZIL from Holder
	t.AssertEqual(aZilContract.Field("_balance"), zil(2))
	t.AssertEqual(aZilContract.Field("autorestakeamount"), zil(2))

	// Send Swap transactions
	t.AssertTransition(txn, deploy.Transition{
		bufferContract.Addr, //sender
		"RequestDelegatorSwap",
		Zproxy.Addr,
		"0",
		deploy.ParamsMap{"new_deleg_addr": "0x" + holderContract.Addr},
	})

	t.AssertTransition(txn, deploy.Transition{
		holderContract.Addr, //sender
		"ConfirmDelegatorSwap",
		Zproxy.Addr,
		"0",
		deploy.ParamsMap{"requestor": "0x" + bufferContract.Addr},
	})
}
