package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsStZil() {
	Start("IsStZil")

	p := tr.DeployAndUpgrade()

	// Use non-admin user for Buffer
	buffer := p.GetBuffer()
	buffer.SetSigner(bob)

	tx, _ := buffer.DelegateStake(sdk.Cfg.StZilSsnAddress, ToZil(1))
	AssertError(tx, buffer.ErrorCode("StZILValidationFailed"))
	tx, _ = buffer.ClaimRewards(sdk.Cfg.StZilSsnAddress)
	AssertError(tx, buffer.ErrorCode("StZILValidationFailed"))
	tx, _ = buffer.RequestDelegatorSwap(p.Holder.Addr)
	AssertError(tx, buffer.ErrorCode("StZILValidationFailed"))
	tx, _ = buffer.ConfirmDelegatorSwap(p.Holder.Addr)
	AssertError(tx, buffer.ErrorCode("StZILValidationFailed"))
	tx, _ = buffer.RejectDelegatorSwap(p.Holder.Addr)
	AssertError(tx, buffer.ErrorCode("StZILValidationFailed"))
	tx, _ = buffer.ReDelegateStake(p.Holder.Addr, sdk.Cfg.StZilSsnAddress, ToZil(1))
	AssertError(tx, buffer.ErrorCode("StZILValidationFailed"))

	// Use non-admin user for p.Holder
	p.Holder.SetSigner(bob)

	tx, _ = p.Holder.WithdrawStakeAmt(sdk.Cfg.StZilSsnAddress, ToZil(1))
	AssertError(tx, p.Holder.ErrorCode("StZILValidationFailed"))
	tx, _ = p.Holder.CompleteWithdrawal()
	AssertError(tx, p.Holder.ErrorCode("StZILValidationFailed"))
	tx, _ = p.Holder.ClaimRewards(sdk.Cfg.StZilSsnAddress)
	AssertError(tx, p.Holder.ErrorCode("StZILValidationFailed"))
	tx, _ = p.Holder.ConfirmDelegatorSwap(buffer.Addr)
	AssertError(tx, p.Holder.ErrorCode("StZILValidationFailed"))
}
