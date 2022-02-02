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

	checkChangeAzilSSNAddress(p)
	checkChangeHolderAddress(p)
	checkChangeZimplAddress(p)
	checkUpdateStakingParameters(p)
	checkChangeBuffersEmpty(p)

	newAdminAddr := sdk.Cfg.Addr3
	newAdminKey := sdk.Cfg.Key3

	//claim not existent staging admin, expecting error
	p.Aimpl.UpdateWallet(newAdminKey)
	tx, _ := p.Aimpl.ClaimAdmin()
	AssertError(tx, "StagingAdminNotExists")

	//change admin, expecting success
	p.Aimpl.UpdateWallet(sdk.Cfg.AdminKey)
	tx, _ = AssertSuccess(p.Aimpl.ChangeAdmin(newAdminAddr))
	AssertEvent(tx, Event{p.Aimpl.Addr, "ChangeAdmin", ParamsMap{"current_admin": sdk.Cfg.Admin, "new_admin": newAdminAddr}})
	AssertEqual(Field(p.Aimpl, "staging_admin_address"), newAdminAddr)

	//claim admin with wrong user, expecting error
	wrongActor := sdk.Cfg.Key1
	p.Aimpl.UpdateWallet(wrongActor)
	tx, _ = p.Aimpl.ClaimAdmin()
	AssertError(tx, "StagingAdminValidationFailed")

	//claim admin with correct user, expecting success
	p.Aimpl.UpdateWallet(newAdminKey)
	tx, _ = AssertSuccess(p.Aimpl.ClaimAdmin())
	AssertEvent(tx, Event{p.Aimpl.Addr, "ClaimAdmin", ParamsMap{"new_admin": newAdminAddr}})
	AssertEqual(Field(p.Aimpl, "admin_address"), newAdminAddr)
	AssertEqual(Field(p.Aimpl, "staging_admin_address"), "")
}

func checkChangeAzilSSNAddress(p *contracts.Protocol) {
	tx, _ := AssertSuccess(p.Aimpl.ChangeAzilSSNAddress(core.ZeroAddr))
	AssertEvent(tx, Event{p.Aimpl.Addr, "ChangeAzilSSNAddress", ParamsMap{"address": core.ZeroAddr}})
	AssertEqual(Field(p.Aimpl, "azil_ssn_address"), core.ZeroAddr)
	AssertSuccess(p.Aimpl.ChangeAzilSSNAddress(sdk.Cfg.AzilSsnAddress))
}

func checkChangeHolderAddress(p *contracts.Protocol) {
	holderAddr := p.Holder.Addr

	tx, _ := AssertSuccess(p.Aimpl.ChangeHolderAddress(core.ZeroAddr))
	AssertEvent(tx, Event{p.Aimpl.Addr, "ChangeHolderAddress", ParamsMap{"address": core.ZeroAddr}})
	AssertEqual(Field(p.Aimpl, "holder_address"), core.ZeroAddr)
	AssertSuccess(p.Aimpl.ChangeHolderAddress(holderAddr))
}

func checkChangeZimplAddress(p *contracts.Protocol) {
	zimplAddr := p.Zimpl.Addr

	tx, _ := AssertSuccess(p.Aimpl.ChangeZimplAddress(core.ZeroAddr))
	AssertEvent(tx, Event{p.Aimpl.Addr, "ChangeZimplAddress", ParamsMap{"address": core.ZeroAddr}})
	AssertEqual(Field(p.Aimpl, "zimpl_address"), core.ZeroAddr)
	AssertSuccess(p.Aimpl.ChangeZimplAddress(zimplAddr))
}

func checkUpdateStakingParameters(p *contracts.Protocol) {
	prevValue := Field(p.Aimpl, "mindelegstake")
	testValue := utils.ToZil(54321)
	tx, _ := AssertSuccess(p.Aimpl.UpdateStakingParameters(testValue))
	AssertEvent(tx, Event{p.Aimpl.Addr, "UpdateStakingParameters", ParamsMap{"min_deleg_stake": testValue}})
	AssertEqual(Field(p.Aimpl, "mindelegstake"), testValue)
	AssertSuccess(p.Aimpl.UpdateStakingParameters(prevValue))
}

func checkChangeBuffersEmpty(p *contracts.Protocol) {
	new_buffers := []string{}
	tx, _ := p.Aimpl.ChangeBuffers(new_buffers)
	AssertError(tx, "BuffersEmpty")
}
