package transitions

import (
    "github.com/avely-finance/avely-contracts/sdk/contracts"
    "github.com/avely-finance/avely-contracts/sdk/core"
    . "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Owner() {

    Start("Azil contract owner transitions")

    p := tr.DeployAndUpgrade()
    p.Azil.UpdateWallet(sdk.Cfg.OwnerKey)

    checkChangeAzilSSNAddress(p)
}

func checkChangeAzilSSNAddress(p *contracts.Protocol) {
    tx, _ := AssertSuccess(p.Azil.ChangeAzilSSNAddress(core.ZeroAddr))
    AssertEvent(tx, Event{p.Azil.Addr, "ChangeAzilSSNAddress", ParamsMap{"address": core.ZeroAddr}})
    AssertEqual(Field(p.Azil, "azil_ssn_address"), core.ZeroAddr)
    AssertSuccess(p.Azil.ChangeAzilSSNAddress(sdk.Cfg.AzilSsnAddress))
}
