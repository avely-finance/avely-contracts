package transitions

import (
//"Azil/test/deploy"
)

func (t *Testing) IsAimpl() {

	t.LogStart("IsAimpl")

	_, _, bufferContract, holderContract := t.DeployAndUpgrade()

	bufferContract.UpdateWallet(key2)
	tx, _ := bufferContract.DelegateStake()
	t.AssertContain(t.GetReceiptString(tx), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -401))])")

	holderContract.UpdateWallet(key2)
	tx, _ = holderContract.CompleteWithdrawal()
	t.AssertContain(t.GetReceiptString(tx), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -301))])")
}
