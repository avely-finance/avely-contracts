package transitions

import (
	//. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsOwner() {

	Start("IsOwner")

	p := tr.DeployAndUpgrade()

	// Use non-owner user for Azil, expecting errors
	p.Azil.UpdateWallet(sdk.Cfg.Key2)

	tx, _ := p.Azil.ChangeAzilSSNAddress(sdk.Cfg.Addr3)
	AssertError(tx, "OwnerValidationFailed")

	tx, _ = p.Azil.ChangeTreasuryAddress(sdk.Cfg.Addr3)
	AssertError(tx, "OwnerValidationFailed")

}
