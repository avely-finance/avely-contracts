package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) PerformAuoRestake() {
	p := tr.DeployAndUpgrade()

	p.Aimpl.UpdateWallet(sdk.Cfg.AdminKey)

	t.AssertEqual(p.Aimpl.Field("autorestakeamount"), Zil(0))

	t.AssertSuccess(p.Aimpl.IncreaseAutoRestakeAmount(Zil(1)))
	txn, err := p.Aimpl.PerformAutoRestake()
	t.AssertError(txn, err, -15)

	// increases to 100
	t.AssertSuccess(p.Aimpl.IncreaseAutoRestakeAmount(Zil(99)))
	restakeAmount := Zil(100)
	t.AssertEqual(p.Aimpl.Field("autorestakeamount"), restakeAmount)

	txn, _ = p.Aimpl.PerformAutoRestake()

	// should return to 0
	t.AssertEqual(p.Aimpl.Field("autorestakeamount"), Zil(0))

	t.AssertTransition(txn, Transition{
		p.Buffer.Addr, //sender
		"DelegateStake",
		p.Zproxy.Addr,
		restakeAmount,
		ParamsMap{},
	})
}
