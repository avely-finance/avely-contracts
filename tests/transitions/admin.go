package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Admin() {

	Start("Azil contract admin transitions")

	p := tr.DeployAndUpgrade()

	newAdminAddr := sdk.Cfg.Addr3
	newAdminKey := sdk.Cfg.Key3

	//claim not existent staging admin, expecting error
	p.Azil.UpdateWallet(newAdminKey)
	tx, _ := p.Azil.ClaimAdmin()
	AssertError(tx, "StagingAdminNotExists")

	//change admin, expecting success
	p.Azil.UpdateWallet(sdk.Cfg.OwnerKey)
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
