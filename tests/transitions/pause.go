package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Pause() {

	Start("Pause/Unpause")

	p := tr.DeployAndUpgrade()

	callPaused(p)
	callPauseUnpauseNonAdmin(p)

	p.Aimpl.UpdateWallet(sdk.Cfg.AdminKey)

	//pause contract, expecting success
	tx, _ := AssertSuccess(p.Aimpl.Pause())
	AssertEvent(tx, Event{p.Aimpl.Addr, "Pause", ParamsMap{"is_paused": "1"}})
	AssertEqual(p.Aimpl.Field("is_paused"), "True")

	//unpause contract, expecting success
	tx, _ = AssertSuccess(p.Aimpl.Unpause())
	AssertEvent(tx, Event{p.Aimpl.Addr, "UnPause", ParamsMap{"is_paused": "0"}})
	AssertEqual(p.Aimpl.Field("is_paused"), "False")
}

func callPauseUnpauseNonAdmin(p *contracts.Protocol) {
	//call pause/unpause admin transitions with non-admin user; expecting errors
	p.Aimpl.UpdateWallet(sdk.Cfg.Key1)

	tx, _ := p.Aimpl.Pause()
	AssertError(tx, "AdminValidationFailed")

	tx, _ = p.Aimpl.Unpause()
	AssertError(tx, "AdminValidationFailed")
}

func callPaused(p *contracts.Protocol) {
	//call user's transitions, when contract is paused; expecting errors
	AssertSuccess(p.Aimpl.Pause())

	tx, _ := p.Aimpl.Pause()
	AssertError(tx, "Paused")

	tx, _ = p.Aproxy.DelegateStake(ToZil(10))
	AssertError(tx, "Paused")

	p.Aproxy.ZilBalanceOf(sdk.Cfg.Addr1)
	tx = sdk.TxLast
	AssertError(tx, "Paused")

	tx, _ = p.Aproxy.WithdrawStakeAmt(ToZil(10))
	AssertError(tx, "Paused")

	tx, _ = p.Aproxy.CompleteWithdrawal()
	AssertError(tx, "Paused")

	AssertSuccess(p.Aimpl.Unpause())
}
