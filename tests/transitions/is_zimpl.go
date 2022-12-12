package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsZimpl() {
	Start("IsZimpl")

	p := tr.DeployAndUpgrade()
	bobAddr := utils.GetAddressByWallet(bob)

	// Use random user for Buffer

	buffer := p.GetBuffer()
	buffer.SetSigner(bob)

	tx, _ := buffer.AddFunds(ToZil(10))
	AssertError(tx, buffer.ErrorCode("ZimplValidationFailed"))
	tx, _ = buffer.WithdrawStakeRewardsSuccessCallBack(bobAddr, ToZil(10))
	AssertError(tx, buffer.ErrorCode("ZimplValidationFailed"))
	tx, _ = buffer.DelegateStakeSuccessCallBack(bobAddr, ToZil(10))
	AssertError(tx, buffer.ErrorCode("ZimplValidationFailed"))
	tx, _ = buffer.ReDelegateStakeSuccessCallBack(sdk.Cfg.StZilSsnAddress, sdk.Cfg.StZilSsnAddress, ToZil(10))
	AssertError(tx, buffer.ErrorCode("ZimplValidationFailed"))

	// Use random user for Buffer
	p.Holder.SetSigner(bob)
	tx, _ = p.Holder.AddFunds(ToZil(10))
	AssertError(tx, p.Holder.ErrorCode("ZimplValidationFailed"))
	tx, _ = p.Holder.DelegateStakeSuccessCallBack(sdk.Cfg.StZilSsnAddress, ToZil(10))
	AssertError(tx, p.Holder.ErrorCode("ZimplValidationFailed"))
	tx, _ = p.Holder.WithdrawStakeAmtSuccessCallBack(bobAddr, ToZil(10))
	AssertError(tx, p.Holder.ErrorCode("ZimplValidationFailed"))
	tx, _ = p.Holder.WithdrawStakeRewardsSuccessCallBack(bobAddr, ToZil(10))
	AssertError(tx, p.Holder.ErrorCode("ZimplValidationFailed"))
	tx, _ = p.Holder.CompleteWithdrawalSuccessCallBack(ToZil(10))
	AssertError(tx, p.Holder.ErrorCode("ZimplValidationFailed"))
	tx, _ = p.Holder.CompleteWithdrawalNoUnbondedStakeCallBack(ToZil(10))
	AssertError(tx, p.Holder.ErrorCode("ZimplValidationFailed"))
}
