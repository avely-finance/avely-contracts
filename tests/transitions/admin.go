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
    checkChangeZproxyAddress(p)
    checkChangeAimplAddress(p)
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

func checkChangeAzilSSNAddress(p *contracts.Protocol) {
    tx, _ := AssertSuccess(p.Aimpl.ChangeAzilSSNAddress(core.ZeroAddr))
    AssertEvent(tx, Event{p.Aimpl.Addr, "ChangeAzilSSNAddress", ParamsMap{"address": "0x" + core.ZeroAddr}})
    AssertEqual(p.Aimpl.Field("azil_ssn_address"), "0x"+core.ZeroAddr)
    AssertSuccess(p.Aimpl.ChangeAzilSSNAddress(sdk.Cfg.AzilSsnAddress))

    tx, _ = AssertSuccess(p.GetBuffer().ChangeAzilSSNAddress(core.ZeroAddr))
    AssertEvent(tx, Event{p.GetBuffer().Addr, "ChangeAzilSSNAddress", ParamsMap{"address": "0x" + core.ZeroAddr}})
    AssertEqual(p.GetBuffer().Field("azil_ssn_address"), "0x"+core.ZeroAddr)
    AssertSuccess(p.GetBuffer().ChangeAzilSSNAddress(sdk.Cfg.AzilSsnAddress))

    tx, _ = AssertSuccess(p.Holder.ChangeAzilSSNAddress(core.ZeroAddr))
    AssertEvent(tx, Event{p.Holder.Addr, "ChangeAzilSSNAddress", ParamsMap{"address": "0x" + core.ZeroAddr}})
    AssertEqual(p.Holder.Field("azil_ssn_address"), "0x"+core.ZeroAddr)
    AssertSuccess(p.Holder.ChangeAzilSSNAddress(sdk.Cfg.AzilSsnAddress))
}

func checkChangeHolderAddress(p *contracts.Protocol) {
    holderAddr := p.Holder.Addr

    tx, _ := AssertSuccess(p.Aimpl.ChangeHolderAddress(core.ZeroAddr))
    AssertEvent(tx, Event{p.Aimpl.Addr, "ChangeHolderAddress", ParamsMap{"address": "0x" + core.ZeroAddr}})
    AssertEqual(p.Aimpl.Field("holder_address"), "0x"+core.ZeroAddr)
    AssertSuccess(p.Aimpl.ChangeHolderAddress(holderAddr))
}

func checkChangeZimplAddress(p *contracts.Protocol) {
    zimplAddr := p.Zimpl.Addr

    tx, _ := AssertSuccess(p.Aimpl.ChangeZimplAddress(core.ZeroAddr))
    AssertEvent(tx, Event{p.Aimpl.Addr, "ChangeZimplAddress", ParamsMap{"address": "0x" + core.ZeroAddr}})
    AssertEqual(p.Aimpl.Field("zimpl_address"), "0x"+core.ZeroAddr)
    AssertSuccess(p.Aimpl.ChangeZimplAddress(zimplAddr))

    tx, _ = AssertSuccess(p.GetBuffer().ChangeZimplAddress(core.ZeroAddr))
    AssertEvent(tx, Event{p.GetBuffer().Addr, "ChangeZimplAddress", ParamsMap{"address": "0x" + core.ZeroAddr}})
    AssertEqual(p.GetBuffer().Field("zimpl_address"), "0x"+core.ZeroAddr)
    AssertSuccess(p.GetBuffer().ChangeZimplAddress(zimplAddr))

    tx, _ = AssertSuccess(p.Holder.ChangeZimplAddress(core.ZeroAddr))
    AssertEvent(tx, Event{p.Holder.Addr, "ChangeZimplAddress", ParamsMap{"address": "0x" + core.ZeroAddr}})
    AssertEqual(p.Holder.Field("zimpl_address"), "0x"+core.ZeroAddr)
    AssertSuccess(p.Holder.ChangeZimplAddress(zimplAddr))
}

func checkChangeAimplAddress(p *contracts.Protocol) {
    aimplAddr := p.Aimpl.Addr

    tx, _ := AssertSuccess(p.GetBuffer().ChangeAimplAddress(core.ZeroAddr))
    AssertEvent(tx, Event{p.GetBuffer().Addr, "ChangeAimplAddress", ParamsMap{"address": "0x" + core.ZeroAddr}})
    AssertEqual(p.GetBuffer().Field("aimpl_address"), "0x"+core.ZeroAddr)
    AssertSuccess(p.GetBuffer().ChangeAimplAddress(aimplAddr))

    tx, _ = AssertSuccess(p.Holder.ChangeAimplAddress(core.ZeroAddr))
    AssertEvent(tx, Event{p.Holder.Addr, "ChangeAimplAddress", ParamsMap{"address": "0x" + core.ZeroAddr}})
    AssertEqual(p.Holder.Field("aimpl_address"), "0x"+core.ZeroAddr)
    AssertSuccess(p.Holder.ChangeAimplAddress(aimplAddr))
}

func checkChangeZproxyAddress(p *contracts.Protocol) {
    zproxyAddr := p.Zproxy.Addr

    tx, _ := AssertSuccess(p.GetBuffer().ChangeZproxyAddress(core.ZeroAddr))
    AssertEvent(tx, Event{p.GetBuffer().Addr, "ChangeZproxyAddress", ParamsMap{"address": "0x" + core.ZeroAddr}})
    AssertEqual(p.GetBuffer().Field("zproxy_address"), "0x"+core.ZeroAddr)
    AssertSuccess(p.GetBuffer().ChangeZproxyAddress(zproxyAddr))

    tx, _ = AssertSuccess(p.Holder.ChangeZproxyAddress(core.ZeroAddr))
    AssertEvent(tx, Event{p.Holder.Addr, "ChangeZproxyAddress", ParamsMap{"address": "0x" + core.ZeroAddr}})
    AssertEqual(p.Holder.Field("zproxy_address"), "0x"+core.ZeroAddr)
    AssertSuccess(p.Holder.ChangeZproxyAddress(zproxyAddr))
}

func checkUpdateStakingParameters(p *contracts.Protocol) {
    prevValue := p.Aimpl.Field("mindelegstake")
    testValue := utils.ToZil(54321)
    tx, _ := AssertSuccess(p.Aimpl.UpdateStakingParameters(testValue))
    AssertEvent(tx, Event{p.Aimpl.Addr, "UpdateStakingParameters", ParamsMap{"min_deleg_stake": testValue}})
    AssertEqual(p.Aimpl.Field("mindelegstake"), testValue)
    AssertSuccess(p.Aimpl.UpdateStakingParameters(prevValue))
}
