package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsOwner() {

	Start("IsOwner")

	p := tr.DeployAndUpgrade()

	// Use non-owner user for StZIL, expecting errors
	p.StZIL.UpdateWallet(sdk.Cfg.Key2)

	tx, _ := p.StZIL.ChangeAdmin(sdk.Cfg.Addr3)
	AssertError(tx, "OwnerValidationFailed")

	tx, _ = p.StZIL.ChangeOwner(sdk.Cfg.Addr3)
	AssertError(tx, "OwnerValidationFailed")

	new_buffers := []string{p.GetBuffer().Addr, p.GetBuffer().Addr}
	tx, _ = p.StZIL.ChangeBuffers(new_buffers)
	AssertError(tx, "OwnerValidationFailed")

	tx, _ = p.StZIL.AddSSN(sdk.Cfg.Addr3)
	AssertError(tx, "OwnerValidationFailed")

	tx, _ = p.StZIL.RemoveSSN(sdk.Cfg.Addr3)
	AssertError(tx, "OwnerValidationFailed")

	tx, _ = p.StZIL.SetHolderAddress(sdk.Cfg.Addr3)
	AssertError(tx, "OwnerValidationFailed")

	tx, _ = p.StZIL.ChangeRewardsFee("100")
	AssertError(tx, "OwnerValidationFailed")

	tx, _ = p.StZIL.ChangeTreasuryAddress(sdk.Cfg.Addr3)
	AssertError(tx, "OwnerValidationFailed")

	tx, _ = p.StZIL.ChangeZimplAddress(sdk.Cfg.Addr3)
	AssertError(tx, "OwnerValidationFailed")

	tx, _ = p.StZIL.UpdateStakingParameters(ToZil(100))
	AssertError(tx, "OwnerValidationFailed")

}
