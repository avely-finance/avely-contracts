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

	callPausedIn(p)
	callPausedOut(p)
	callPauseUnpauseNonAdmin(p)

	p.Azil.UpdateWallet(sdk.Cfg.AdminKey)

	//pause-in contract, expecting success
	tx, _ := AssertSuccess(p.Azil.PauseIn())
	AssertEvent(tx, Event{p.Azil.Addr, "PauseIn", ParamsMap{"is_paused_in": "1"}})
	AssertEqual(Field(p.Azil, "is_paused_in"), "True")

	//unpause-in contract, expecting success
	tx, _ = AssertSuccess(p.Azil.UnpauseIn())
	AssertEvent(tx, Event{p.Azil.Addr, "UnPauseIn", ParamsMap{"is_paused_in": "0"}})
	AssertEqual(Field(p.Azil, "is_paused_in"), "False")

	//pause-out contract, expecting success
	tx, _ = AssertSuccess(p.Azil.PauseOut())
	AssertEvent(tx, Event{p.Azil.Addr, "PauseOut", ParamsMap{"is_paused_out": "1"}})
	AssertEqual(Field(p.Azil, "is_paused_out"), "True")

	//unpause-out contract, expecting success
	tx, _ = AssertSuccess(p.Azil.UnpauseOut())
	AssertEvent(tx, Event{p.Azil.Addr, "UnPauseOut", ParamsMap{"is_paused_out": "0"}})
	AssertEqual(Field(p.Azil, "is_paused_out"), "False")
}

func unPauseEmptyBuffers() {
	p := contracts.Deploy(sdk, GetLog())
	tx, _ := p.Azil.UnpauseIn()
	AssertError(tx, "BuffersEmpty")
}

func callPauseUnpauseNonAdmin(p *contracts.Protocol) {
	//call pause/unpause admin transitions with non-admin user; expecting errors
	p.Azil.UpdateWallet(sdk.Cfg.Key1)

	tx, _ := p.Azil.PauseIn()
	AssertError(tx, "AdminValidationFailed")

	tx, _ = p.Azil.UnpauseIn()
	AssertError(tx, "AdminValidationFailed")

	tx, _ = p.Azil.PauseOut()
	AssertError(tx, "AdminValidationFailed")

	tx, _ = p.Azil.UnpauseOut()
	AssertError(tx, "AdminValidationFailed")
}

func callPausedIn(p *contracts.Protocol) {
	//call transitions, when contract is paused-in; expecting errors
	AssertSuccess(p.Azil.PauseIn())

	tx, _ := p.Azil.PauseIn()
	AssertError(tx, "PausedIn")

	tx, _ = p.Azil.DelegateStake(ToZil(10))
	AssertError(tx, "PausedIn")

	tx, _ = p.Azil.ChownStakeConfirmSwap(sdk.Cfg.Addr1)
	AssertError(tx, "PausedIn")

	AssertSuccess(p.Azil.UnpauseIn())
}

func callPausedOut(p *contracts.Protocol) {
	//call transitions, when contract is paused-out; expecting errors
	AssertSuccess(p.Azil.PauseOut())

	tx, _ := p.Azil.PauseOut()
	AssertError(tx, "PausedOut")

	tx, _ = p.Azil.WithdrawStakeAmt(ToZil(10))
	AssertError(tx, "PausedOut")

	tx, _ = p.Azil.CompleteWithdrawal()
	AssertError(tx, "PausedOut")

	AssertSuccess(p.Azil.UnpauseOut())
}
