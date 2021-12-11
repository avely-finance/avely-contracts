package transitions

import (
//"Azil/test/deploy"
)

func (t *Testing) IsAimpl() {

	t.LogStart("IsAimpl")

	_, _, _, bufferContract, holderContract := t.DeployAndUpgrade()

	// Use non-admin user for Buffer
	bufferContract.UpdateWallet(key2)

	tx, err := bufferContract.DelegateStake()
	t.AssertError(tx, err, -401)
	tx, err = bufferContract.ClaimRewards()
	t.AssertError(tx, err, -401)
	tx, err = bufferContract.RequestDelegatorSwap(holderContract.Addr)
	t.AssertError(tx, err, -401)

	// Use non-admin user for Holder
	holderContract.UpdateWallet(key2)

	tx, err = holderContract.CompleteWithdrawal()
	t.AssertError(tx, err, -301)
	tx, err = holderContract.ClaimRewards()
	t.AssertError(tx, err, -301)
	tx, err = holderContract.ConfirmDelegatorSwap(bufferContract.Addr)
	t.AssertError(tx, err, -301)
}
