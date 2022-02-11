package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsAdmin() {

	Start("IsAdmin")

	p := tr.DeployAndUpgrade()

	// Use non-admin user for Azil, expecting errors
	p.Azil.UpdateWallet(sdk.Cfg.Key2)

	tx, _ := p.Azil.ChangeAdmin(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Azil.ChangeZimplAddress(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Azil.ChangeAzilSSNAddress(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Azil.ChangeTreasuryAddress(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Azil.ChangeHolderAddress(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")

	new_buffers := []string{p.GetBuffer().Addr, p.GetBuffer().Addr}
	tx, _ = p.Azil.ChangeBuffers(new_buffers)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Azil.IncreaseAutoRestakeAmount(ToZil(1))
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Azil.PerformAutoRestake()
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Azil.UpdateStakingParameters(ToZil(100))
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Azil.ChangeRewardsFee("100")
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Azil.DrainBuffer(p.GetBuffer().Addr)
	AssertError(tx, "AdminValidationFailed")
	readyBlocks := []string{}
	tx, _ = p.Azil.ClaimWithdrawal(readyBlocks)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Azil.ChownStakeConfirmSwap(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Azil.ChownStakeReDelegate(sdk.Cfg.Addr3, ToZil(1))
	AssertError(tx, "AdminValidationFailed")
}
