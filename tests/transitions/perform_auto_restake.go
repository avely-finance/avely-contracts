package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) PerformAuoRestake() {
	p := tr.DeployAndUpgrade()

	p.Aimpl.UpdateWallet(sdk.Cfg.AdminKey)

	AssertEqual(Field(p.Aimpl, "autorestakeamount"), ToZil(0))

	AssertSuccess(p.Aimpl.IncreaseAutoRestakeAmount(ToZil(1)))
	txn, _ := p.Aimpl.PerformAutoRestake()
	AssertError(txn, "DelegStakeNotEnough")

	// increases to 100
	AssertSuccess(p.Aimpl.IncreaseAutoRestakeAmount(ToZil(99)))
	restakeAmount := ToZil(100)
	AssertEqual(Field(p.Aimpl, "autorestakeamount"), restakeAmount)

	txn, _ = p.Aimpl.PerformAutoRestake()

	// should return to 0
	AssertEqual(Field(p.Aimpl, "autorestakeamount"), ToZil(0))

	AssertTransition(txn, Transition{
		p.GetBuffer().Addr, //sender
		"DelegateStake",
		p.Zproxy.Addr,
		restakeAmount,
		ParamsMap{},
	})
}
