package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsAdmin() {

	Start("IsAdmin")

	p := tr.DeployAndUpgrade()

	// Use non-admin user for Aimpl, expecting errors
	p.Aimpl.UpdateWallet(sdk.Cfg.Key2)

	tx, _ := p.Aimpl.ChangeAdmin(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")
	new_buffers := []string{p.GetBuffer().Addr, p.GetBuffer().Addr}
	tx, _ = p.Aimpl.ChangeBuffers(new_buffers)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Aimpl.IncreaseAutoRestakeAmount(ToZil(1))
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Aimpl.PerformAutoRestake()
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Aimpl.UpdateStakingParameters(ToZil(100))
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Aimpl.GetCurrentBuffer()
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Aimpl.DrainBuffer(p.GetBuffer().Addr)
	AssertError(tx, "AdminValidationFailed")
	readyBlocks := []string{}
	tx, _ = p.Aimpl.ClaimWithdrawal(readyBlocks)
	AssertError(tx, "AdminValidationFailed")
}
