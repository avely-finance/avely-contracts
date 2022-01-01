package transitions

import (
    //. "github.com/avely-finance/avely-contracts/sdk/utils"
    . "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Proxy() {

    Start("Proxy")

    p := tr.DeployAndUpgrade()

    newImpl := sdk.Cfg.Addr3

    //non-admin
    p.Aproxy.UpdateWallet(sdk.Cfg.Key3)
    tx, err := p.Aproxy.UpgradeTo(newImpl)
    AssertError(tx, err, -202)

    //admin
    p.Aproxy.UpdateWallet(sdk.Cfg.AdminKey)
    tx, _ = AssertSuccess(p.Aproxy.UpgradeTo(newImpl))
    AssertEvent(tx, Event{p.Aproxy.Addr, "UpgradeTo", ParamsMap{"aimpl_address": "0x" + newImpl}})
    AssertEqual(p.Aproxy.Field("aimpl_address"), "0x"+newImpl)
}
