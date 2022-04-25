package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) PerformAutoRestake() {
	p := tr.DeployAndUpgrade()

	p.StZIL.UpdateWallet(sdk.Cfg.AdminKey)

	AssertEqual(Field(p.StZIL, "autorestakeamount"), ToZil(0))

	AssertSuccess(p.StZIL.IncreaseAutoRestakeAmount(ToZil(1)))
	txn, _ := p.StZIL.PerformAutoRestake()
	AssertError(txn, "DelegStakeNotEnough")

	// increases to 100
	AssertSuccess(p.StZIL.IncreaseAutoRestakeAmount(ToZil(99)))
	restakeAmount := ToZil(100)
	AssertEqual(Field(p.StZIL, "autorestakeamount"), restakeAmount)

	txn, _ = p.StZIL.PerformAutoRestake()

	// should return to 0
	AssertEqual(Field(p.StZIL, "autorestakeamount"), ToZil(0))

	AssertTransition(txn, Transition{
		p.GetActiveBuffer().Addr, //sender
		"DelegateStake",
		p.Zproxy.Addr,
		restakeAmount,
		ParamsMap{},
	})
}
