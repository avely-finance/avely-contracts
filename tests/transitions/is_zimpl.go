package transitions

func (tr *Transitions) IsZimpl() {

	t.Start("IsAimpl")

	_, _, _, Buffer, Holder := tr.DeployAndUpgrade()

	// Use random user for Buffer
	Buffer.UpdateWallet(tr.cfg.Key2)

	tx, err := Buffer.AddFunds(zil(10))
	t.AssertError(tx, err, -407)
	tx, err = Buffer.WithdrawStakeRewardsSuccessCallBack(tr.cfg.Addr2, zil(10))
	t.AssertError(tx, err, -407)
	tx, err = Buffer.DelegateStakeSuccessCallBack(tr.cfg.Addr2, zil(10))
	t.AssertError(tx, err, -407)

	// Use random user for Buffer
	Holder.UpdateWallet(tr.cfg.Key2)
	tx, err = Holder.AddFunds(zil(10))
	t.AssertError(tx, err, -307)
	tx, err = Holder.DelegateStakeSuccessCallBack(tr.cfg.AzilSsnAddress, zil(10))
	t.AssertError(tx, err, -307)
	tx, err = Holder.WithdrawStakeAmtSuccessCallBack(tr.cfg.Addr2, zil(10))
	t.AssertError(tx, err, -307)
	tx, err = Holder.WithdrawStakeRewardsSuccessCallBack(tr.cfg.Addr2, zil(10))
	t.AssertError(tx, err, -307)
	tx, err = Holder.CompleteWithdrawalSuccessCallBack(zil(10))
	t.AssertError(tx, err, -307)
	tx, err = Holder.CompleteWithdrawalNoUnbondedStakeCallBack(zil(10))
	t.AssertError(tx, err, -307)
}
