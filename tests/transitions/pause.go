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
	p.Azil.UpdateWallet(sdk.Cfg.OwnerKey)

	callPausedIn(p)
	callPausedOut(p)
	callPausedZrc2(p)
	callPauseUnpauseNonAdmin(p)

	//pause-in contract, expecting success
	tx, _ := AssertSuccess(p.Azil.WithUser(sdk.Cfg.OwnerKey).PauseIn())
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

	//pause-zrc2 contract, expecting success
	tx, _ = AssertSuccess(p.Azil.PauseZrc2())
	AssertEvent(tx, Event{p.Azil.Addr, "PauseZrc2", ParamsMap{"is_paused_zrc2": "1"}})
	AssertEqual(Field(p.Azil, "is_paused_zrc2"), "True")

	//unpause-zrc2 contract, expecting success
	tx, _ = AssertSuccess(p.Azil.UnpauseZrc2())
	AssertEvent(tx, Event{p.Azil.Addr, "UnPauseZrc2", ParamsMap{"is_paused_zrc2": "0"}})
	AssertEqual(Field(p.Azil, "is_paused_zrc2"), "False")
}

func unPauseEmptyBuffers() {
	p := contracts.Deploy(sdk, GetLog())
	tx, _ := p.Azil.WithUser(sdk.Cfg.OwnerKey).UnpauseIn()
	AssertError(tx, "BuffersEmpty")
}

func callPauseUnpauseNonAdmin(p *contracts.Protocol) {
	//call pause/unpause admin transitions with non-admin user; expecting errors
	p.Azil.UpdateWallet(sdk.Cfg.Key1)

	tx, _ := p.Azil.PauseIn()
	AssertError(tx, "OwnerValidationFailed")

	tx, _ = p.Azil.UnpauseIn()
	AssertError(tx, "OwnerValidationFailed")

	tx, _ = p.Azil.PauseOut()
	AssertError(tx, "OwnerValidationFailed")

	tx, _ = p.Azil.UnpauseOut()
	AssertError(tx, "OwnerValidationFailed")

	tx, _ = p.Azil.PauseZrc2()
	AssertError(tx, "OwnerValidationFailed")

	tx, _ = p.Azil.UnpauseZrc2()
	AssertError(tx, "OwnerValidationFailed")
}

func callPausedIn(p *contracts.Protocol) {
	//call transitions, when contract is paused-in; expecting errors
	AssertSuccess(p.Azil.PauseIn())

	tx, _ := p.Azil.PauseIn()
	AssertError(tx, "PausedIn")

	tx, _ = p.Azil.DelegateStake(ToZil(10))
	AssertError(tx, "PausedIn")

	tx, _ = p.Azil.WithUser(sdk.Cfg.AdminKey).ChownStakeConfirmSwap(sdk.Cfg.Addr1)
	AssertError(tx, "PausedIn")

	AssertSuccess(p.Azil.WithUser(sdk.Cfg.OwnerKey).UnpauseIn())
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

func callPausedZrc2(p *contracts.Protocol) {
	//call transitions, when contract is paused-zrc2; expecting errors
	AssertSuccess(p.Azil.PauseZrc2())

	tx, _ := p.Azil.PauseZrc2()
	AssertError(tx, "PausedZrc2")

	tx, _ = p.Azil.WithUser(sdk.Cfg.OwnerKey).TransferFrom(sdk.Cfg.Addr1, sdk.Cfg.Addr2, ToQA(10000))
	AssertError(tx, "PausedZrc2")

	tx, _ = p.Azil.WithUser(sdk.Cfg.OwnerKey).Transfer(sdk.Cfg.Addr2, ToQA(10000))
	AssertError(tx, "PausedZrc2")

	tx, _ = p.Azil.WithUser(sdk.Cfg.OwnerKey).IncreaseAllowance(sdk.Cfg.Addr1, ToQA(10000))
	AssertError(tx, "PausedZrc2")

	tx, _ = p.Azil.WithUser(sdk.Cfg.OwnerKey).DecreaseAllowance(sdk.Cfg.Addr1, ToQA(10000))
	AssertError(tx, "PausedZrc2")

	AssertSuccess(p.Azil.UnpauseZrc2())
}
