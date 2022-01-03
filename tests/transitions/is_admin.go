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

	tx, err := p.GetBuffer().ChangeAzilSSNAddress(sdk.Cfg.Addr3)
	AssertError(tx, err, -402)
	tx, err = p.GetBuffer().ChangeAimplAddress(sdk.Cfg.Addr3)
	AssertError(tx, err, -402)
	tx, err = p.GetBuffer().ChangeZproxyAddress(sdk.Cfg.Addr3)
	AssertError(tx, err, -402)
	tx, err = p.GetBuffer().ChangeZimplAddress(sdk.Cfg.Addr3)
	AssertError(tx, err, -402)

	// Use non-admin user for p.Holder
	p.Holder.UpdateWallet(sdk.Cfg.Key2)

	tx, err = p.Holder.DelegateStake(ToZil(1))
	AssertError(tx, err, -305)
	tx, err = p.Holder.ChangeAzilSSNAddress(sdk.Cfg.Addr3)
	AssertError(tx, err, -305)
	tx, err = p.Holder.ChangeAimplAddress(sdk.Cfg.Addr3)
	AssertError(tx, err, -305)
	tx, err = p.Holder.ChangeZproxyAddress(sdk.Cfg.Addr3)
	AssertError(tx, err, -305)
	tx, err = p.Holder.ChangeZimplAddress(sdk.Cfg.Addr3)
	AssertError(tx, err, -305)

	// Use non-admin user for Aimpl
	p.Aimpl.UpdateWallet(sdk.Cfg.Key2)

	tx, err = p.Aimpl.ChangeZimplAddress(sdk.Cfg.Addr3)
	AssertError(tx, err, "AdminValidationFailed")
	tx, err = p.Aimpl.ChangeHolderAddress(sdk.Cfg.Addr3)
	AssertError(tx, err, "AdminValidationFailed")

	new_buffers := []string{"0x" + p.GetBuffer().Addr, "0x" + p.GetBuffer().Addr}
	tx, err = p.Aimpl.ChangeBuffers(new_buffers)
	AssertError(tx, err, "AdminValidationFailed")
	tx, err = p.Aimpl.PerformAutoRestake()
	AssertError(tx, err, "AdminValidationFailed")
	tx, err = p.Aimpl.UpdateStakingParameters(ToZil(100))
	AssertError(tx, err, "AdminValidationFailed")
	tx, err = p.Aimpl.DrainBuffer(p.GetBuffer().Addr)
	AssertError(tx, err, "AdminValidationFailed")
	readyBlocks := []string{}
	tx, err = p.Aimpl.ClaimWithdrawal(readyBlocks)
	AssertError(tx, err, "AdminValidationFailed")
}
