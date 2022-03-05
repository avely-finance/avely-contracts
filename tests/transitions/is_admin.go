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

	tx, _ := p.Azil.IncreaseAutoRestakeAmount(ToZil(1))
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Azil.PerformAutoRestake()
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Azil.DrainBuffer(p.GetBuffer().Addr)
	AssertError(tx, "AdminValidationFailed")
	readyBlocks := []string{}
	tx, _ = p.Azil.ClaimWithdrawal(readyBlocks)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Azil.ChownStakeReDelegate(sdk.Cfg.Addr3, ToZil(1))
	AssertError(tx, "AdminValidationFailed")
}
