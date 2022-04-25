package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsStZil() {
	Start("IsStZil")

	p := tr.DeployAndUpgrade()

	// Use non-admin user for Buffer
	p.GetBuffer().UpdateWallet(sdk.Cfg.Key2)

	tx, _ := p.GetBuffer().DelegateStake(sdk.Cfg.StZilSsnAddress, ToZil(1))
	AssertError(tx, "StZILValidationFailed")
	tx, _ = p.GetBuffer().ClaimRewards(sdk.Cfg.StZilSsnAddress)
	AssertError(tx, "StZILValidationFailed")
	tx, _ = p.GetBuffer().RequestDelegatorSwap(p.Holder.Addr)
	AssertError(tx, "StZILValidationFailed")
	tx, _ = p.GetBuffer().ConfirmDelegatorSwap(p.Holder.Addr)
	AssertError(tx, "StZILValidationFailed")
	tx, _ = p.GetBuffer().RejectDelegatorSwap(p.Holder.Addr)
	AssertError(tx, "StZILValidationFailed")
	tx, _ = p.GetBuffer().ReDelegateStake(p.Holder.Addr, sdk.Cfg.StZilSsnAddress, ToZil(1))
	AssertError(tx, "StZILValidationFailed")

	// Use non-admin user for p.Holder
	p.Holder.UpdateWallet(sdk.Cfg.Key2)

	tx, _ = p.Holder.WithdrawStakeAmt(sdk.Cfg.StZilSsnAddress, ToZil(1))
	AssertError(tx, "StZILValidationFailed")
	tx, _ = p.Holder.CompleteWithdrawal()
	AssertError(tx, "StZILValidationFailed")
	tx, _ = p.Holder.ClaimRewards(sdk.Cfg.StZilSsnAddress)
	AssertError(tx, "StZILValidationFailed")
	tx, _ = p.Holder.ConfirmDelegatorSwap(p.GetBuffer().Addr)
	AssertError(tx, "StZILValidationFailed")
}
