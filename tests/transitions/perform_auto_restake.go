package transitions

import (
	"Azil/test/deploy"
)

func (t *Testing) PerformAuoRestake() {
	Zproxy, _, Aimpl, Buffer, _ := t.DeployAndUpgrade()

	Aimpl.UpdateWallet(adminKey)

	t.AssertEqual(Aimpl.Field("autorestakeamount"), zil(0))

	restakeAmount := zil(100)
	Aimpl.IncreaseAutoRestakeAmount(restakeAmount)

	// increases to 100
	t.AssertEqual(Aimpl.Field("autorestakeamount"), restakeAmount)

	txn, _ := Aimpl.PerformAutoRestake()

	// should return to 0
	t.AssertEqual(Aimpl.Field("autorestakeamount"), zil(0))

	t.AssertTransition(txn, deploy.Transition{
		Buffer.Addr, //sender
		"DelegateStake",
		Zproxy.Addr,
		restakeAmount,
		deploy.ParamsMap{},
	})
}
