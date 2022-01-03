package transitions

import (
    "github.com/avely-finance/avely-contracts/sdk/contracts"
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
    newImplAddr := "0000000000000000000000000000000000000000"

    //change admin, expecting success
    p.Aproxy.UpdateWallet(sdk.Cfg.AdminKey)
    tx, _ := AssertSuccess(p.Aproxy.ChangeAdmin(newAdminAddr))
    AssertEvent(tx, Event{p.Aproxy.Addr, "ChangeAdmin", ParamsMap{"currentAdmin": "0x" + sdk.Cfg.Admin, "newAdmin": "0x" + newAdminAddr}})
    AssertEqual(p.Aproxy.Field("staging_admin_address"), "0x"+newAdminAddr)

    //claim admin with wrong user, expecting error
    wrongActor := sdk.Cfg.Key1
    p.Aproxy.UpdateWallet(wrongActor)
    tx, err := p.Aproxy.ClaimAdmin()
    AssertError(tx, err, "StagingAdminValidationFailed")

    //claim admin with correct user, expecting success
    p.Aproxy.UpdateWallet(newAdminKey)
    tx, _ = AssertSuccess(p.Aproxy.ClaimAdmin())
    AssertEvent(tx, Event{p.Aproxy.Addr, "ClaimAdmin", ParamsMap{"newAdmin": "0x" + newAdminAddr}})
    AssertEqual(p.Aproxy.Field("admin_address"), "0x"+newAdminAddr)

    //call UpgradeTo with new admin, expecting success
    tx, _ = AssertSuccess(p.Aproxy.UpgradeTo(newImplAddr))
    AssertEvent(tx, Event{p.Aproxy.Addr, "UpgradeTo", ParamsMap{"aimpl_address": "0x" + newImplAddr}})
    AssertEqual(p.Aproxy.Field("aimpl_address"), "0x"+newImplAddr)
}

func callAimplDirectly(p *contracts.Protocol) {
    //call aimpl transitions, which are supposed to call through proxy, directly; expecting errors
    initiator := sdk.Cfg.Addr3
    tx, err := p.Aimpl.DelegateStake(ToZil(10), initiator)
    AssertError(tx, err, "ProxyValidationFailed")

    tx, err = p.Aimpl.ZilBalanceOf(initiator, initiator)
    AssertError(tx, err, "ProxyValidationFailed")

    tx, err = p.Aimpl.WithdrawStakeAmt(ToZil(10), initiator)
    AssertError(tx, err, "ProxyValidationFailed")

    tx, err = p.Aimpl.CompleteWithdrawal(initiator)
    AssertError(tx, err, "ProxyValidationFailed")
}

func callNonAdmin(p *contracts.Protocol) {
    //call proxy admin transitions with non-admin user; expecting errors
    p.Aproxy.UpdateWallet(sdk.Cfg.Key1)
    tx, err := p.Aproxy.UpgradeTo("0000000000000000000000000000000000000000")
    AssertError(tx, err, "AdminValidationFailed")

    tx, err = p.Aproxy.ChangeAdmin("0000000000000000000000000000000000000000")
    AssertError(tx, err, "AdminValidationFailed")
}
