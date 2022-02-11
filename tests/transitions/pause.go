package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Pause() {
	Start("Pause/Unpause")

	unPauseEmptyBuffers()

	p := tr.DeployAndUpgrade()

	callPaused(p)
	callPauseUnpauseNonAdmin(p)

	p.Azil.UpdateWallet(sdk.Cfg.AdminKey)

	//pause contract, expecting success
	tx, _ := AssertSuccess(p.Azil.Pause())
	AssertEvent(tx, Event{p.Azil.Addr, "Pause", ParamsMap{"is_paused": "1"}})
	AssertEqual(Field(p.Azil, "is_paused"), "True")

	//unpause contract, expecting success
	tx, _ = AssertSuccess(p.Azil.Unpause())
	AssertEvent(tx, Event{p.Azil.Addr, "UnPause", ParamsMap{"is_paused": "0"}})
	AssertEqual(Field(p.Azil, "is_paused"), "False")
}

func unPauseEmptyBuffers() {
	p := contracts.Deploy(sdk, GetLog())
	tx, _ := p.Azil.Unpause()
	AssertError(tx, "BuffersEmpty")
}

func callPauseUnpauseNonAdmin(p *contracts.Protocol) {
	//call pause/unpause admin transitions with non-admin user; expecting errors
	p.Azil.UpdateWallet(sdk.Cfg.Key1)

	tx, _ := p.Azil.Pause()
	AssertError(tx, "AdminValidationFailed")

	tx, _ = p.Azil.Unpause()
	AssertError(tx, "AdminValidationFailed")
}

func callPaused(p *contracts.Protocol) {
	//call user's transitions, when contract is paused; expecting errors
	AssertSuccess(p.Azil.Pause())

	tx, _ := p.Azil.Pause()
	AssertError(tx, "Paused")

	tx, _ = p.Azil.DelegateStake(ToZil(10))
	AssertError(tx, "Paused")

	tx, _ = p.Azil.WithdrawStakeAmt(ToZil(10))
	AssertError(tx, "Paused")

	tx, _ = p.Azil.CompleteWithdrawal()
	AssertError(tx, "Paused")

	AssertSuccess(p.Azil.Unpause())
}
