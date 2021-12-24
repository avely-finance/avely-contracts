package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
)

func (tr *Transitions) IsZimpl() {
	Start("IsAimpl")

	p := tr.DeployAndUpgrade()

	// Use random user for Buffer
	p.Buffer.UpdateWallet(sdk.Cfg.Key2)

	tx, err := p.Buffer.AddFunds(ToZil(10))
	AssertError(tx, err, -407)
	tx, err = p.Buffer.WithdrawStakeRewardsSuccessCallBack(sdk.Cfg.Addr2, ToZil(10))
	AssertError(tx, err, -407)
	tx, err = p.Buffer.DelegateStakeSuccessCallBack(sdk.Cfg.Addr2, ToZil(10))
	AssertError(tx, err, -407)

	// Use random user for Buffer
	p.Holder.UpdateWallet(sdk.Cfg.Key2)
	tx, err = p.Holder.AddFunds(ToZil(10))
	AssertError(tx, err, -307)
	tx, err = p.Holder.DelegateStakeSuccessCallBack(sdk.Cfg.AzilSsnAddress, ToZil(10))
	AssertError(tx, err, -307)
	tx, err = p.Holder.WithdrawStakeAmtSuccessCallBack(sdk.Cfg.Addr2, ToZil(10))
	AssertError(tx, err, -307)
	tx, err = p.Holder.WithdrawStakeRewardsSuccessCallBack(sdk.Cfg.Addr2, ToZil(10))
	AssertError(tx, err, -307)
	tx, err = p.Holder.CompleteWithdrawalSuccessCallBack(ToZil(10))
	AssertError(tx, err, -307)
	tx, err = p.Holder.CompleteWithdrawalNoUnbondedStakeCallBack(ToZil(10))
	AssertError(tx, err, -307)
}
