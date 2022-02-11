package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) PerformAutoRestake() {
	p := tr.DeployAndUpgrade()

	p.Azil.UpdateWallet(sdk.Cfg.AdminKey)

	AssertEqual(Field(p.Azil, "autorestakeamount"), ToZil(0))

	AssertSuccess(p.Azil.IncreaseAutoRestakeAmount(ToZil(1)))
	txn, _ := p.Azil.PerformAutoRestake()
	AssertError(txn, "DelegStakeNotEnough")

	// increases to 100
	AssertSuccess(p.Azil.IncreaseAutoRestakeAmount(ToZil(99)))
	restakeAmount := ToZil(100)
	AssertEqual(Field(p.Azil, "autorestakeamount"), restakeAmount)

	txn, _ = p.Azil.PerformAutoRestake()

	// should return to 0
	AssertEqual(Field(p.Azil, "autorestakeamount"), ToZil(0))

	AssertTransition(txn, Transition{
		p.GetBuffer().Addr, //sender
		"DelegateStake",
		p.Zproxy.Addr,
		restakeAmount,
		ParamsMap{},
	})
}
