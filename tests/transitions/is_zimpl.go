package transitions

func (t *Testing) IsZimpl() {

	t.LogStart("IsAimpl")

	_, _, bufferContract, holderContract := t.DeployAndUpgrade()

	// Use random user for Buffer
	bufferContract.UpdateWallet(key2)

	tx, err := bufferContract.AddFunds(zil(10))
	t.AssertError(tx, err, -407)
	tx, err = bufferContract.WithdrawStakeRewardsSuccessCallBack(addr2, zil(10))
	t.AssertError(tx, err, -407)
	tx, err = bufferContract.DelegateStakeSuccessCallBack(addr2, zil(10))
	t.AssertError(tx, err, -407)

	// Use random user for Buffer
	holderContract.UpdateWallet(key2)
	tx, err = holderContract.AddFunds(zil(10))
	t.AssertError(tx, err, -307)
	tx, err = holderContract.WithdrawStakeAmtSuccessCallBack(addr2, zil(10))
	t.AssertError(tx, err, -307)
	tx, err = holderContract.WithdrawStakeRewardsSuccessCallBack(addr2, zil(10))
	t.AssertError(tx, err, -307)
	tx, err = holderContract.CompleteWithdrawalSuccessCallBack(zil(10))
	t.AssertError(tx, err, -307)
	tx, err = holderContract.CompleteWithdrawalNoUnbondedStakeCallBack(zil(10))
	t.AssertError(tx, err, -307)
}
