package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsZimpl() {
	Start("IsAimpl")

	p := tr.DeployAndUpgrade()

	// Use random user for Buffer
	p.GetBuffer().UpdateWallet(sdk.Cfg.Key2)

	tx, _ := p.GetBuffer().AddFunds(ToZil(10))
	AssertError(tx, "ZimplValidationFailed")
	tx, _ = p.GetBuffer().WithdrawStakeRewardsSuccessCallBack(sdk.Cfg.Addr2, ToZil(10))
	AssertError(tx, "ZimplValidationFailed")
	tx, _ = p.GetBuffer().DelegateStakeSuccessCallBack(sdk.Cfg.Addr2, ToZil(10))
	AssertError(tx, "ZimplValidationFailed")

	// Use random user for Buffer
	p.Holder.UpdateWallet(sdk.Cfg.Key2)
	tx, _ = p.Holder.AddFunds(ToZil(10))
	AssertError(tx, "ZimplValidationFailed")
	tx, _ = p.Holder.DelegateStakeSuccessCallBack(sdk.Cfg.AzilSsnAddress, ToZil(10))
	AssertError(tx, "ZimplValidationFailed")
	tx, _ = p.Holder.ReDelegateStakeSuccessCallBack(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnAddress, ToZil(10))
	AssertError(tx, "ZimplValidationFailed")
	tx, _ = p.Holder.WithdrawStakeAmtSuccessCallBack(sdk.Cfg.Addr2, ToZil(10))
	AssertError(tx, "ZimplValidationFailed")
	tx, _ = p.Holder.WithdrawStakeRewardsSuccessCallBack(sdk.Cfg.Addr2, ToZil(10))
	AssertError(tx, "ZimplValidationFailed")
	tx, _ = p.Holder.CompleteWithdrawalSuccessCallBack(ToZil(10))
	AssertError(tx, "ZimplValidationFailed")
	tx, _ = p.Holder.CompleteWithdrawalNoUnbondedStakeCallBack(ToZil(10))
	AssertError(tx, "ZimplValidationFailed")
}
