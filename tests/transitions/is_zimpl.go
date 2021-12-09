package transitions

func (t *Testing) IsZimpl() {

    t.LogStart("IsAimpl")

    _, _, bufferContract, _ := t.DeployAndUpgrade()

    // Use random user for Buffer
    bufferContract.UpdateWallet(key2)

    tx, err := bufferContract.AddFunds(zil(10))
    t.AssertError(tx, err, -407)
    tx, err = bufferContract.WithdrawStakeRewardsSuccessCallBack(addr2, zil(10))
    t.AssertError(tx, err, -407)
    tx, err = bufferContract.DelegateStakeSuccessCallBack(addr2, zil(10))
    t.AssertError(tx, err, -407)

}
