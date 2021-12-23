package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsZimpl() {
	t.Start("IsAimpl")

	p := tr.DeployAndUpgrade()

	// Use random user for Buffer
	p.Buffer.UpdateWallet(sdk.Cfg.Key2)

	tx, err := p.Buffer.AddFunds(Zil(10))
	t.AssertError(tx, err, -407)
	tx, err = p.Buffer.WithdrawStakeRewardsSuccessCallBack(sdk.Cfg.Addr2, Zil(10))
	t.AssertError(tx, err, -407)
	tx, err = p.Buffer.DelegateStakeSuccessCallBack(sdk.Cfg.Addr2, Zil(10))
	t.AssertError(tx, err, -407)

	// Use random user for Buffer
	p.Holder.UpdateWallet(sdk.Cfg.Key2)
	tx, err = p.Holder.AddFunds(Zil(10))
	t.AssertError(tx, err, -307)
	tx, err = p.Holder.DelegateStakeSuccessCallBack(sdk.Cfg.AzilSsnAddress, Zil(10))
	t.AssertError(tx, err, -307)
	tx, err = p.Holder.WithdrawStakeAmtSuccessCallBack(sdk.Cfg.Addr2, Zil(10))
	t.AssertError(tx, err, -307)
	tx, err = p.Holder.WithdrawStakeRewardsSuccessCallBack(sdk.Cfg.Addr2, Zil(10))
	t.AssertError(tx, err, -307)
	tx, err = p.Holder.CompleteWithdrawalSuccessCallBack(Zil(10))
	t.AssertError(tx, err, -307)
	tx, err = p.Holder.CompleteWithdrawalNoUnbondedStakeCallBack(Zil(10))
	t.AssertError(tx, err, -307)
}
