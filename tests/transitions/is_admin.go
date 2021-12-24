package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsAdmin() {

	Start("IsAdmin")

	p := tr.DeployAndUpgrade()

	// Use non-admin user for p.Buffer
	p.Buffer.UpdateWallet(sdk.Cfg.Key3)

	tx, err := p.Buffer.ChangeAzilSSNAddress(sdk.Cfg.Addr3)
	AssertError(tx, err, -402)
	tx, err = p.Buffer.ChangeAimplAddress(sdk.Cfg.Addr3)
	AssertError(tx, err, -402)
	tx, err = p.Buffer.ChangeZproxyAddress(sdk.Cfg.Addr3)
	AssertError(tx, err, -402)
	tx, err = p.Buffer.ChangeZimplAddress(sdk.Cfg.Addr3)
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
	AssertError(tx, err, -106)
	tx, err = p.Aimpl.ChangeHolderAddress(sdk.Cfg.Addr3)
	AssertError(tx, err, -106)

	new_buffers := []string{"0x" + p.Buffer.Addr, "0x" + p.Buffer.Addr}
	tx, err = p.Aimpl.ChangeBuffers(new_buffers)
	AssertError(tx, err, -106)
	tx, err = p.Aimpl.PerformAutoRestake()
	AssertError(tx, err, -106)
	tx, err = p.Aimpl.UpdateStakingParameters(ToZil(100))
	AssertError(tx, err, -106)
	tx, err = p.Aimpl.DrainBuffer(p.Buffer.Addr)
	AssertError(tx, err, -106)
	readyBlocks := []string{}
	tx, err = p.Aimpl.ClaimWithdrawal(readyBlocks)
	AssertError(tx, err, -106)
}
