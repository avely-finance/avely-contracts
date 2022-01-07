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

    checkSetHolderAddress(p)
    checkUpdateStakingParameters(p)

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
    AssertEqual(p.Aimpl.Field("staging_admin_address"), newAdminAddr)

    //claim admin with wrong user, expecting error
    wrongActor := sdk.Cfg.Key1
    p.Aimpl.UpdateWallet(wrongActor)
    tx, _ = p.Aimpl.ClaimAdmin()
    AssertError(tx, "StagingAdminValidationFailed")

    //claim admin with correct user, expecting success
    p.Aimpl.UpdateWallet(newAdminKey)
    tx, _ = AssertSuccess(p.Aimpl.ClaimAdmin())
    AssertEvent(tx, Event{p.Aimpl.Addr, "ClaimAdmin", ParamsMap{"new_admin": newAdminAddr}})
    AssertEqual(p.Aimpl.Field("admin_address"), newAdminAddr)
    AssertEqual(p.Aimpl.Field("staging_admin_address"), "")
}

func checkSetHolderAddress(p *contracts.Protocol) {
    AssertEqual(p.Aimpl.Field("holder_address"), p.Holder.Addr)
    tx, _ := p.Aimpl.SetHolderAddress(core.ZeroAddr)
    AssertError(tx, "HolderAlreadySet")
}

func checkUpdateStakingParameters(p *contracts.Protocol) {
    prevValue := p.Aimpl.Field("mindelegstake")
    testValue := utils.ToZil(54321)
    tx, _ := AssertSuccess(p.Aimpl.UpdateStakingParameters(testValue))
    AssertEvent(tx, Event{p.Aimpl.Addr, "UpdateStakingParameters", ParamsMap{"min_deleg_stake": testValue}})
    AssertEqual(p.Aimpl.Field("mindelegstake"), testValue)
    AssertSuccess(p.Aimpl.UpdateStakingParameters(prevValue))
}
