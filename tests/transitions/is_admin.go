package transitions

import (
//"Azil/test/deploy"
)

func (t *Testing) IsAdmin() {

	t.LogStart("IsAdmin")

	_, aZilContract, bufferContract, holderContract := t.DeployAndUpgrade()

	// Use non-admin user for Buffer
	bufferContract.UpdateWallet(key3)

	tx, err := bufferContract.ChangeAzilSSNAddress(addr3)
	t.AssertError(tx, err, -402)
	tx, err = bufferContract.ChangeAimplAddress(addr3)
	t.AssertError(tx, err, -402)
	tx, err = bufferContract.ChangeZproxyAddress(addr3)
	t.AssertError(tx, err, -402)
	tx, err = bufferContract.ChangeZimplAddress(addr3)
	t.AssertError(tx, err, -402)

	// Use non-admin user for Holder
	holderContract.UpdateWallet(key2)

	tx, err = holderContract.ChangeAzilSSNAddress(addr3)
	t.AssertError(tx, err, -305)
	tx, err = holderContract.ChangeAimplAddress(addr3)
	t.AssertError(tx, err, -305)
	tx, err = holderContract.ChangeProxyStakingContractAddress(addr3)
	t.AssertError(tx, err, -305)

	// Use non-admin user for aZilContract
	aZilContract.UpdateWallet(key2)

	tx, err = aZilContract.ChangeProxyStakingContractAddress(addr3)
	t.AssertError(tx, err, -106)
	tx, err = aZilContract.ChangeHolderAddress(addr3)
	t.AssertError(tx, err, -106)
	new_buffers := []string{"0x" + bufferContract.Addr, "0x" + bufferContract.Addr}
	tx, err = aZilContract.ChangeBuffers(new_buffers)
	t.AssertError(tx, err, -106)
	tx, err = aZilContract.IncreaseTotalStakeAmount(zil(100))
	t.AssertError(tx, err, -106)
	tx, err = aZilContract.UpdateStakingParameters(zil(100))
	t.AssertError(tx, err, -106)
	tx, err = aZilContract.DrainBuffer(bufferContract.Addr)
	t.AssertError(tx, err, -106)
}
