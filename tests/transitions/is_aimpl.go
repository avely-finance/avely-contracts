package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsAimpl() {
	Start("IsAimpl")

	p := tr.DeployAndUpgrade()

	// Use non-admin user for Buffer
	p.GetBuffer().UpdateWallet(sdk.Cfg.Key2)

	tx, _ := p.GetBuffer().DelegateStake()
	AssertError(tx, "AimplValidationFailed")
	tx, _ = p.GetBuffer().ClaimRewards()
	AssertError(tx, "AimplValidationFailed")
	tx, _ = p.GetBuffer().RequestDelegatorSwap(p.Holder.Addr)
	AssertError(tx, "AimplValidationFailed")
	tx, _ = p.GetBuffer().ReDelegateStake(p.Holder.Addr, ToZil(1))
	AssertError(tx, "AimplValidationFailed")

	// Use non-admin user for p.Holder
	p.Holder.UpdateWallet(sdk.Cfg.Key2)

	tx, _ = p.Holder.WithdrawStakeAmt(ToZil(1))
	AssertError(tx, "AimplValidationFailed")
	tx, _ = p.Holder.CompleteWithdrawal()
	AssertError(tx, "AimplValidationFailed")
	tx, _ = p.Holder.ClaimRewards()
	AssertError(tx, "AimplValidationFailed")
	tx, _ = p.Holder.ConfirmDelegatorSwap(p.GetBuffer().Addr)
	AssertError(tx, "AimplValidationFailed")
}
