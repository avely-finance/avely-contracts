package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsAzil() {
	Start("IsAzil")

	p := tr.DeployAndUpgrade()

	// Use non-admin user for Buffer
	p.GetBuffer().UpdateWallet(sdk.Cfg.Key2)

	tx, _ := p.GetBuffer().DelegateStake(sdk.Cfg.AzilSsnAddress, ToZil(1))
	AssertError(tx, "AzilValidationFailed")
	tx, _ = p.GetBuffer().ClaimRewards(sdk.Cfg.AzilSsnAddress)
	AssertError(tx, "AzilValidationFailed")
	tx, _ = p.GetBuffer().RequestDelegatorSwap(p.Holder.Addr)
	AssertError(tx, "AzilValidationFailed")
	tx, _ = p.GetBuffer().ConfirmDelegatorSwap(p.Holder.Addr)
	AssertError(tx, "AzilValidationFailed")
	tx, _ = p.GetBuffer().RejectDelegatorSwap(p.Holder.Addr)
	AssertError(tx, "AzilValidationFailed")
	tx, _ = p.GetBuffer().ReDelegateStake(p.Holder.Addr, ToZil(1))
	AssertError(tx, "AzilValidationFailed")

	// Use non-admin user for p.Holder
	p.Holder.UpdateWallet(sdk.Cfg.Key2)

	tx, _ = p.Holder.WithdrawStakeAmt(ToZil(1))
	AssertError(tx, "AzilValidationFailed")
	tx, _ = p.Holder.CompleteWithdrawal()
	AssertError(tx, "AzilValidationFailed")
	tx, _ = p.Holder.ClaimRewards(sdk.Cfg.AzilSsnAddress)
	AssertError(tx, "AzilValidationFailed")
	tx, _ = p.Holder.ConfirmDelegatorSwap(p.GetBuffer().Addr)
	AssertError(tx, "AzilValidationFailed")
}
