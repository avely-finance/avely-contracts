package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Pause() {
	Start("Pause/Unpause")

	tr.PauseEmptyBuffers()

	p := tr.DeployAndUpgrade()
	p.StZIL.SetSigner(celestials.Owner)

	tr.PausedIn()
	tr.PausedOut()
	tr.PausedZrc2()
	tr.PauseUnpauseNonAdmin()
	tr.PauseUnpauseAdmin()

	tr.EvmOn()
	//tr.PausedIn()
	//tr.PausedOut()
	tr.PausedZrc2()
	tr.EvmOff()
}

func (tr *Transitions) PauseEmptyBuffers() {
	owner := celestials.Owner
	admin := celestials.Admin

	p := contracts.Deploy(sdk, utils.GetAddressByWallet(owner), admin, GetLog())
	p.StZIL.SetSigner(celestials.Owner)

	tx, _ := p.StZIL.UnpauseIn()
	AssertError(tx, p.StZIL.ErrorCode("BuffersEmpty"))

	p.SyncBufferAndHolder(celestials.Owner)

	tx, _ = p.StZIL.UnpauseIn()
	AssertError(tx, p.StZIL.ErrorCode("SsnAddressesEmpty"))
}

func (tr *Transitions) PauseUnpauseNonAdmin() {
	p := tr.p

	//call pause/unpause admin transitions with non-admin user; expecting errors
	p.StZIL.SetSigner(alice)

	tx, _ := p.StZIL.PauseIn()
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	tx, _ = p.StZIL.UnpauseIn()
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	tx, _ = p.StZIL.PauseOut()
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	tx, _ = p.StZIL.UnpauseOut()
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	tx, _ = p.StZIL.PauseZrc2()
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	tx, _ = p.StZIL.UnpauseZrc2()
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))
}
func (tr *Transitions) PauseUnpauseAdmin() {
	p := tr.p

	// make sure we work under owner account
	p.StZIL.SetSigner(celestials.Owner)
	//pause-in contract, expecting success
	tx, _ := AssertSuccess(p.StZIL.PauseIn())
	AssertEvent(tx, Event{p.StZIL.Addr, "PauseIn", ParamsMap{"is_paused_in": "1"}})
	AssertEqual(Field(p.StZIL, "is_paused_in"), "True")

	//unpause-in contract, expecting success
	tx, _ = AssertSuccess(p.StZIL.UnpauseIn())
	AssertEvent(tx, Event{p.StZIL.Addr, "UnPauseIn", ParamsMap{"is_paused_in": "0"}})
	AssertEqual(Field(p.StZIL, "is_paused_in"), "False")

	//pause-out contract, expecting success
	tx, _ = AssertSuccess(p.StZIL.PauseOut())
	AssertEvent(tx, Event{p.StZIL.Addr, "PauseOut", ParamsMap{"is_paused_out": "1"}})
	AssertEqual(Field(p.StZIL, "is_paused_out"), "True")

	//unpause-out contract, expecting success
	tx, _ = AssertSuccess(p.StZIL.UnpauseOut())
	AssertEvent(tx, Event{p.StZIL.Addr, "UnPauseOut", ParamsMap{"is_paused_out": "0"}})
	AssertEqual(Field(p.StZIL, "is_paused_out"), "False")

	//pause-zrc2 contract, expecting success
	tx, _ = AssertSuccess(p.StZIL.PauseZrc2())
	AssertEvent(tx, Event{p.StZIL.Addr, "PauseZrc2", ParamsMap{"is_paused_zrc2": "1"}})
	AssertEqual(Field(p.StZIL, "is_paused_zrc2"), "True")

	//unpause-zrc2 contract, expecting success
	tx, _ = AssertSuccess(p.StZIL.UnpauseZrc2())
	AssertEvent(tx, Event{p.StZIL.Addr, "UnPauseZrc2", ParamsMap{"is_paused_zrc2": "0"}})
	AssertEqual(Field(p.StZIL, "is_paused_zrc2"), "False")
}

func (tr *Transitions) PausedIn() {
	p := tr.p

	//call transitions, when contract is paused-in; expecting errors
	AssertSuccess(p.StZIL.PauseIn())

	tx, _ := p.StZIL.PauseIn()
	AssertError(tx, p.StZIL.ErrorCode("PausedIn"))

	tx, _ = p.StZIL.DelegateStake(ToZil(10))
	AssertError(tx, p.StZIL.ErrorCode("PausedIn"))

	aliceAddr := utils.GetAddressByWallet(alice)

	p.StZIL.SetSigner(alice)
	tx, _ = p.StZIL.ChownStakeConfirmSwap(aliceAddr)
	AssertError(tx, p.StZIL.ErrorCode("PausedIn"))

	p.StZIL.SetSigner(celestials.Owner)
	AssertSuccess(p.StZIL.UnpauseIn())
}

func (tr *Transitions) PausedOut() {
	p := tr.p

	p.StZIL.SetSigner(celestials.Owner)
	//call transitions, when contract is paused-out; expecting errors
	AssertSuccess(p.StZIL.PauseOut())

	tx, _ := p.StZIL.PauseOut()
	AssertError(tx, p.StZIL.ErrorCode("PausedOut"))

	tx, _ = p.StZIL.WithdrawTokensAmt(ToStZil(10))
	AssertError(tx, p.StZIL.ErrorCode("PausedOut"))

	tx, _ = p.StZIL.CompleteWithdrawal()
	AssertError(tx, p.StZIL.ErrorCode("PausedOut"))

	AssertSuccess(p.StZIL.UnpauseOut())
}

func (tr *Transitions) PausedZrc2() {
	p := tr.p

	//call transitions, when contract is paused-zrc2; expecting errors
	p.StZIL.SetSigner(celestials.Owner)
	AssertSuccess(p.StZIL.PauseZrc2())

	tx, _ := p.StZIL.PauseZrc2()
	AssertError(tx, p.StZIL.ErrorCode("PausedZrc2"))

	from := tr.GetAddressByWallet(alice)
	to := tr.GetAddressByWallet(bob)

	tx1, _ := tr.GetStZIL().TransferFrom(from, to, ToQA(10000))
	AssertError(tx1, p.StZIL.ErrorCode("PausedZrc2"))

	tx1, _ = tr.GetStZIL().Transfer(to, ToQA(10000))
	AssertError(tx1, p.StZIL.ErrorCode("PausedZrc2"))

	tx1, _ = tr.GetStZIL().IncreaseAllowance(from, ToQA(10000))
	AssertError(tx1, p.StZIL.ErrorCode("PausedZrc2"))

	tx1, _ = tr.GetStZIL().DecreaseAllowance(from, ToQA(10000))
	AssertError(tx1, p.StZIL.ErrorCode("PausedZrc2"))

	AssertSuccess(p.StZIL.UnpauseZrc2())
}
