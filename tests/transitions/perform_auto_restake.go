package transitions

import (
	. "Azil/test/helpers"
)

func (tr *Transitions) PerformAuoRestake() {
	Zproxy, _, Aimpl, Buffer, _ := tr.DeployAndUpgrade()

	Aimpl.UpdateWallet(tr.cfg.AdminKey)

	t.AssertEqual(Aimpl.Field("autorestakeamount"), zil(0))

	t.AssertSuccess(Aimpl.IncreaseAutoRestakeAmount(zil(1)))
	txn, err := Aimpl.PerformAutoRestake()
	t.AssertError(txn, err, -15)

	// increases to 100
	t.AssertSuccess(Aimpl.IncreaseAutoRestakeAmount(zil(99)))
	restakeAmount := zil(100)
	t.AssertEqual(Aimpl.Field("autorestakeamount"), restakeAmount)

	txn, _ = Aimpl.PerformAutoRestake()

	// should return to 0
	t.AssertEqual(Aimpl.Field("autorestakeamount"), zil(0))

	t.AssertTransition(txn, Transition{
		Buffer.Addr, //sender
		"DelegateStake",
		Zproxy.Addr,
		restakeAmount,
		ParamsMap{},
	})
}
