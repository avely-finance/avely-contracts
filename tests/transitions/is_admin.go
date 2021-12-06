package transitions

import (
//"Azil/test/deploy"
)

func (t *Testing) IsAdmin() {

	t.LogStart("IsAdmin")

	_, aZilContract, bufferContract, holderContract := t.DeployAndUpgrade()

	bufferContract.UpdateWallet(key3)
	tx, _ := bufferContract.ChangeAzilSSNAddress(addr3)
	t.AssertContain(t.GetReceiptString(tx), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -405))])")
	tx, _ = bufferContract.ChangeAimplAddress(addr3)
	t.AssertContain(t.GetReceiptString(tx), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -405))])")
	tx, _ = bufferContract.ChangeProxyStakingContractAddress(addr3)
	t.AssertContain(t.GetReceiptString(tx), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -405))])")

	holderContract.UpdateWallet(key2)
	tx, _ = holderContract.ChangeAzilSSNAddress(addr3)
	t.AssertContain(t.GetReceiptString(tx), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -305))])")
	tx, _ = holderContract.ChangeAimplAddress(addr3)
	t.AssertContain(t.GetReceiptString(tx), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -305))])")
	tx, _ = holderContract.ChangeProxyStakingContractAddress(addr3)
	t.AssertContain(t.GetReceiptString(tx), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -305))])")

	aZilContract.UpdateWallet(key2)
	tx, _ = aZilContract.ChangeProxyStakingContractAddress(addr3)
	t.AssertContain(t.GetReceiptString(tx), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -106))])")
	tx, _ = aZilContract.ChangeHolderAddress(addr3)
	t.AssertContain(t.GetReceiptString(tx), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -106))])")
	new_buffers := []string{"0x" + bufferContract.Addr, "0x" + bufferContract.Addr}
	aZilContract.ChangeBuffers(new_buffers)
	t.AssertContain(t.GetReceiptString(tx), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -106))])")
	tx, _ = aZilContract.IncreaseTotalStakeAmount(zil100)
	t.AssertContain(t.GetReceiptString(tx), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -106))])")
	tx, _ = aZilContract.UpdateStakingParameters(zil100)
	t.AssertContain(t.GetReceiptString(tx), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -106))])")
}
