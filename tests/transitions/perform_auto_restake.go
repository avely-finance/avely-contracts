package transitions

import (
	"Azil/test/deploy"
)

func (t *Testing) PerformAuoRestake() {
	stubStakingContract, aZilContract, bufferContract, _ := t.DeployAndUpgrade()

	aZilContract.UpdateWallet(adminKey)

	t.AssertEqual(aZilContract.Field("autorestakeamount"), zil(0))

	restakeAmount := zil(100)
	aZilContract.IncreaseAutoRestakeAmount(restakeAmount)

	// increases to 100
	t.AssertEqual(aZilContract.Field("autorestakeamount"), restakeAmount)

	txn, _ := aZilContract.PerformAutoRestake()

	// should return to 0
	t.AssertEqual(aZilContract.Field("autorestakeamount"), zil(0))

	t.AssertTransition(txn, deploy.Transition{
		bufferContract.Addr, //sender
		"DelegateStake",
		stubStakingContract.Addr,
		restakeAmount,
		deploy.ParamsMap{},
	})
}
