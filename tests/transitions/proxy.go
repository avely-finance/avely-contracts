package transitions

import (
    . "github.com/avely-finance/avely-contracts/sdk/utils"
    . "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Proxy() {

    Start("Proxy")

    p := tr.DeployAndUpgrade()

    newAddr := sdk.Cfg.Addr3

    //call proxy admin transitions with non-admin user; expecting errors
    p.Aproxy.UpdateWallet(sdk.Cfg.Key3)
    tx, err := p.Aproxy.UpgradeTo(newAddr)
    AssertError(tx, err, -202)

    tx, err = p.Aproxy.ChangeAdmin(newAddr)
    AssertError(tx, err, -202)

    //call proxy admin transitions with admin user; expecting success
    p.Aproxy.UpdateWallet(sdk.Cfg.AdminKey)
    tx, _ = AssertSuccess(p.Aproxy.UpgradeTo(newAddr))
    AssertEvent(tx, Event{p.Aproxy.Addr, "UpgradeTo", ParamsMap{"aimpl_address": "0x" + newAddr}})
    AssertEqual(p.Aproxy.Field("aimpl_address"), "0x"+newAddr)

    tx, _ = AssertSuccess(p.Aproxy.ChangeAdmin(newAddr))
    AssertEvent(tx, Event{p.Aproxy.Addr, "ChangeAdmin", ParamsMap{"currentAdmin": "0x" + sdk.Cfg.Admin, "newAdmin": "0x" + newAddr}})
    AssertEqual(p.Aproxy.Field("staging_admin_address"), "0x"+newAddr)

    //call aimpl transitions, which are supposed to call through proxy, directly; expecting errors
    initiator := sdk.Cfg.Addr3
    tx, err = p.Aimpl.DelegateStake(ToZil(10), initiator)
    AssertError(tx, err, -113)

    tx, err = p.Aimpl.ZilBalanceOf(initiator, initiator)
    AssertError(tx, err, -113)

    tx, err = p.Aimpl.WithdrawStakeAmt(ToZil(10), initiator)
    AssertError(tx, err, -113)

    tx, err = p.Aimpl.CompleteWithdrawal(initiator)
    AssertError(tx, err, -113)
}
