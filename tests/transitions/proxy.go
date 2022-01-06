package transitions

import (
    "github.com/avely-finance/avely-contracts/sdk/contracts"
    "github.com/avely-finance/avely-contracts/sdk/core"
    . "github.com/avely-finance/avely-contracts/sdk/utils"
    . "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Proxy() {

    Start("Proxy")

    p := tr.DeployAndUpgrade()

    callAimplDirectly(p)
    callNonAdmin(p)

    newAdminAddr := sdk.Cfg.Addr3
    newAdminKey := sdk.Cfg.Key3
    newImplAddr := core.ZeroAddr

    //claim not existent staging admin, expecting error
    p.Aproxy.UpdateWallet(newAdminKey)
    tx, _ := p.Aimpl.ClaimAdmin()
    AssertError(tx, "StagingAdminNotExists")

    //change admin, expecting success
    p.Aproxy.UpdateWallet(sdk.Cfg.AdminKey)
    tx, _ = AssertSuccess(p.Aproxy.ChangeAdmin(newAdminAddr))
    AssertEvent(tx, Event{p.Aproxy.Addr, "ChangeAdmin", ParamsMap{"current_admin": sdk.Cfg.Admin, "new_admin": newAdminAddr}})
    AssertEqual(p.Aproxy.Field("staging_admin_address"), newAdminAddr)

    //claim admin with wrong user, expecting error
    wrongActor := sdk.Cfg.Key1
    p.Aproxy.UpdateWallet(wrongActor)
    tx, _ = p.Aproxy.ClaimAdmin()
    AssertError(tx, "StagingAdminValidationFailed")

    //claim admin with correct user, expecting success
    p.Aproxy.UpdateWallet(newAdminKey)
    tx, _ = AssertSuccess(p.Aproxy.ClaimAdmin())
    AssertEvent(tx, Event{p.Aproxy.Addr, "ClaimAdmin", ParamsMap{"new_admin": newAdminAddr}})
    AssertEqual(p.Aproxy.Field("admin_address"), newAdminAddr)

    //call UpgradeTo with new admin, expecting success
    tx, _ = AssertSuccess(p.Aproxy.UpgradeTo(newImplAddr))
    AssertEvent(tx, Event{p.Aproxy.Addr, "UpgradeTo", ParamsMap{"aimpl_address": newImplAddr}})
    AssertEqual(p.Aproxy.Field("aimpl_address"), newImplAddr)
    AssertEqual(p.Aproxy.Field("staging_admin_address"), "")
}

func callAimplDirectly(p *contracts.Protocol) {
    //call aimpl transitions, which are supposed to call through proxy, directly; expecting errors
    initiator := sdk.Cfg.Addr3
    tx, _ := p.Aimpl.DelegateStake(ToZil(10), initiator)
    AssertError(tx, "ProxyValidationFailed")

    tx, _ = p.Aimpl.ZilBalanceOf(initiator, initiator)
    AssertError(tx, "ProxyValidationFailed")

    tx, _ = p.Aimpl.WithdrawStakeAmt(ToZil(10), initiator)
    AssertError(tx, "ProxyValidationFailed")

    tx, _ = p.Aimpl.CompleteWithdrawal(initiator)
    AssertError(tx, "ProxyValidationFailed")
}

func callNonAdmin(p *contracts.Protocol) {
    //call proxy admin transitions with non-admin user; expecting errors
    p.Aproxy.UpdateWallet(sdk.Cfg.Key1)
    tx, _ := p.Aproxy.UpgradeTo(core.ZeroAddr)
    AssertError(tx, "AdminValidationFailed")

    tx, _ = p.Aproxy.ChangeAdmin(core.ZeroAddr)
    AssertError(tx, "AdminValidationFailed")
}
