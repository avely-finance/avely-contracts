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
    p.Aimpl.UpdateWallet(newAdminKey)
    tx, _ := p.Aimpl.ClaimAdmin()
    AssertError(tx, "StagingAdminNotExists")

    //change admin, expecting success
    p.Aimpl.UpdateWallet(sdk.Cfg.AdminKey)
    tx, _ = AssertSuccess(p.Aimpl.ChangeAdmin(newAdminAddr))
    AssertEvent(tx, Event{p.Aimpl.Addr, "ChangeAdmin", ParamsMap{"current_admin": "0x" + sdk.Cfg.Admin, "new_admin": "0x" + newAdminAddr}})
    AssertEqual(p.Aimpl.Field("staging_admin_address"), "0x"+newAdminAddr)

    //claim admin with wrong user, expecting error
    wrongActor := sdk.Cfg.Key1
    p.Aimpl.UpdateWallet(wrongActor)
    tx, _ = p.Aimpl.ClaimAdmin()
    AssertError(tx, "StagingAdminValidationFailed")

    //claim admin with correct user, expecting success
    p.Aimpl.UpdateWallet(newAdminKey)
    tx, _ = AssertSuccess(p.Aimpl.ClaimAdmin())
    AssertEvent(tx, Event{p.Aimpl.Addr, "ClaimAdmin", ParamsMap{"new_admin": "0x" + newAdminAddr}})
    AssertEqual(p.Aimpl.Field("admin_address"), "0x"+newAdminAddr)
    AssertEqual(p.Aimpl.Field("staging_admin_address"), "")

}
