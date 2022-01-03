package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsAdmin() {

	Start("IsAdmin")

	p := tr.DeployAndUpgrade()

	// Use non-admin user for p.GetBuffer()
	p.GetBuffer().UpdateWallet(sdk.Cfg.Key3)

	tx, _ := p.GetBuffer().ChangeAzilSSNAddress(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.GetBuffer().ChangeAimplAddress(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.GetBuffer().ChangeZproxyAddress(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.GetBuffer().ChangeZimplAddress(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")

	// Use non-admin user for p.Holder
	p.Holder.UpdateWallet(sdk.Cfg.Key2)

	tx, _ = p.Holder.DelegateStake(ToZil(1))
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Holder.ChangeAzilSSNAddress(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Holder.ChangeAimplAddress(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Holder.ChangeZproxyAddress(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Holder.ChangeZimplAddress(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")

	// Use non-admin user for Aimpl
	p.Aimpl.UpdateWallet(sdk.Cfg.Key2)

	tx, _ = p.Aimpl.ChangeZimplAddress(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Aimpl.ChangeHolderAddress(sdk.Cfg.Addr3)
	AssertError(tx, "AdminValidationFailed")

	new_buffers := []string{"0x" + p.GetBuffer().Addr, "0x" + p.GetBuffer().Addr}
	tx, _ = p.Aimpl.ChangeBuffers(new_buffers)
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Aimpl.PerformAutoRestake()
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Aimpl.UpdateStakingParameters(ToZil(100))
	AssertError(tx, "AdminValidationFailed")
	tx, _ = p.Aimpl.DrainBuffer(p.GetBuffer().Addr)
	AssertError(tx, "AdminValidationFailed")
	readyBlocks := []string{}
	tx, _ = p.Aimpl.ClaimWithdrawal(readyBlocks)
	AssertError(tx, "AdminValidationFailed")
}
