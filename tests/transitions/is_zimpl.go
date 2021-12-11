package transitions

func (t *Testing) IsZimpl() {

	t.LogStart("IsAimpl")

	_, _, _, Buffer, Holder := t.DeployAndUpgrade()

	// Use random user for Buffer
	Buffer.UpdateWallet(key2)

	tx, err := Buffer.AddFunds(zil(10))
	t.AssertError(tx, err, -407)
	tx, err = Buffer.WithdrawStakeRewardsSuccessCallBack(addr2, zil(10))
	t.AssertError(tx, err, -407)
	tx, err = Buffer.DelegateStakeSuccessCallBack(addr2, zil(10))
	t.AssertError(tx, err, -407)

	// Use random user for Buffer
	Holder.UpdateWallet(key2)
	tx, err = Holder.AddFunds(zil(10))
	t.AssertError(tx, err, -307)
	tx, err = Holder.WithdrawStakeAmtSuccessCallBack(addr2, zil(10))
	t.AssertError(tx, err, -307)
	tx, err = Holder.WithdrawStakeRewardsSuccessCallBack(addr2, zil(10))
	t.AssertError(tx, err, -307)
	tx, err = Holder.CompleteWithdrawalSuccessCallBack(zil(10))
	t.AssertError(tx, err, -307)
	tx, err = Holder.CompleteWithdrawalNoUnbondedStakeCallBack(zil(10))
	t.AssertError(tx, err, -307)
}
