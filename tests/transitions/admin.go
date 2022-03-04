package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Admin() {

	Start("Azil contract admin transitions")

	p := tr.DeployAndUpgrade()

	checkChangeHolderAddress(p)
	checkUpdateStakingParameters(p)
	checkChangeBuffersEmpty(p)

	newAdminAddr := sdk.Cfg.Addr3
	newAdminKey := sdk.Cfg.Key3

	//claim not existent staging admin, expecting error
	p.Azil.UpdateWallet(newAdminKey)
	tx, _ := p.Azil.ClaimAdmin()
	AssertError(tx, "StagingAdminNotExists")

	//change admin, expecting success
	p.Azil.UpdateWallet(sdk.Cfg.AdminKey)
	tx, _ = AssertSuccess(p.Azil.ChangeAdmin(newAdminAddr))
	AssertEvent(tx, Event{p.Azil.Addr, "ChangeAdmin", ParamsMap{"current_admin": sdk.Cfg.Admin, "new_admin": newAdminAddr}})
	AssertEqual(Field(p.Azil, "staging_admin_address"), newAdminAddr)

	//claim admin with wrong user, expecting error
	wrongActor := sdk.Cfg.Key1
	p.Azil.UpdateWallet(wrongActor)
	tx, _ = p.Azil.ClaimAdmin()
	AssertError(tx, "StagingAdminValidationFailed")

	//claim admin with correct user, expecting success
	p.Azil.UpdateWallet(newAdminKey)
	tx, _ = AssertSuccess(p.Azil.ClaimAdmin())
	AssertEvent(tx, Event{p.Azil.Addr, "ClaimAdmin", ParamsMap{"new_admin": newAdminAddr}})
	AssertEqual(Field(p.Azil, "admin_address"), newAdminAddr)
	AssertEqual(Field(p.Azil, "staging_admin_address"), "")
}

func checkChangeHolderAddress(p *contracts.Protocol) {
	holderAddr := p.Holder.Addr

	tx, _ := AssertSuccess(p.Azil.ChangeHolderAddress(core.ZeroAddr))
	AssertEvent(tx, Event{p.Azil.Addr, "ChangeHolderAddress", ParamsMap{"address": core.ZeroAddr}})
	AssertEqual(Field(p.Azil, "holder_address"), core.ZeroAddr)
	AssertSuccess(p.Azil.ChangeHolderAddress(holderAddr))
}

func checkUpdateStakingParameters(p *contracts.Protocol) {
	prevValue := Field(p.Azil, "mindelegstake")
	testValue := utils.ToZil(54321)
	tx, _ := AssertSuccess(p.Azil.UpdateStakingParameters(testValue))
	AssertEvent(tx, Event{p.Azil.Addr, "UpdateStakingParameters", ParamsMap{"min_deleg_stake": testValue}})
	AssertEqual(Field(p.Azil, "mindelegstake"), testValue)
	AssertSuccess(p.Azil.UpdateStakingParameters(prevValue))
}

func checkChangeBuffersEmpty(p *contracts.Protocol) {
	new_buffers := []string{}
	tx, _ := p.Azil.ChangeBuffers(new_buffers)
	AssertError(tx, "BuffersEmpty")
}
