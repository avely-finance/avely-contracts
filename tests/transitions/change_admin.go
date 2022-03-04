package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) ChangeAdmin() {
	Start("Azil contract admin transitions")

	p := tr.DeployAndUpgrade()

	newAdminAddr := sdk.Cfg.Addr3

	//change admin, expecting success
	p.Azil.UpdateWallet(sdk.Cfg.OwnerKey)
	tx, _ := AssertSuccess(p.Azil.ChangeAdmin(newAdminAddr))
	AssertEvent(tx, Event{
		Sender:    p.Azil.Addr,
		EventName: "ChangeAdmin",
		Params:    ParamsMap{"old_admin": sdk.Cfg.Admin, "new_admin": newAdminAddr},
	})
	AssertEqual(Field(p.Azil, "admin_address"), newAdminAddr)
}
