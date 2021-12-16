package transitions

import (
//"Azil/test/deploy"
)

func (t *Testing) IsAdmin() {

	t.LogStart("IsAdmin")

	_, _, Aimpl, Buffer, Holder := t.DeployAndUpgrade()

	// Use non-admin user for Buffer
	Buffer.UpdateWallet(key3)

	tx, err := Buffer.ChangeAzilSSNAddress(addr3)
	t.AssertError(tx, err, -402)
	tx, err = Buffer.ChangeAimplAddress(addr3)
	t.AssertError(tx, err, -402)
	tx, err = Buffer.ChangeZproxyAddress(addr3)
	t.AssertError(tx, err, -402)
	tx, err = Buffer.ChangeZimplAddress(addr3)
	t.AssertError(tx, err, -402)

	// Use non-admin user for Holder
	Holder.UpdateWallet(key2)

	tx, err = Holder.DelegateStake(zil(1))
	t.AssertError(tx, err, -305)
	tx, err = Holder.ChangeAzilSSNAddress(addr3)
	t.AssertError(tx, err, -305)
	tx, err = Holder.ChangeAimplAddress(addr3)
	t.AssertError(tx, err, -305)
	tx, err = Holder.ChangeZproxyAddress(addr3)
	t.AssertError(tx, err, -305)
	tx, err = Holder.ChangeZimplAddress(addr3)
	t.AssertError(tx, err, -305)

	// Use non-admin user for Aimpl
	Aimpl.UpdateWallet(key2)

	tx, err = Aimpl.ChangeZproxyAddress(addr3)
	t.AssertError(tx, err, -106)
	tx, err = Aimpl.ChangeZimplAddress(addr3)
	t.AssertError(tx, err, -106)
	tx, err = Aimpl.ChangeHolderAddress(addr3)
	t.AssertError(tx, err, -106)
	new_buffers := []string{"0x" + Buffer.Addr, "0x" + Buffer.Addr}
	tx, err = Aimpl.ChangeBuffers(new_buffers)
	t.AssertError(tx, err, -106)
	tx, err = Aimpl.PerformAutoRestake()
	t.AssertError(tx, err, -106)
	tx, err = Aimpl.UpdateStakingParameters(zil(100))
	t.AssertError(tx, err, -106)
	tx, err = Aimpl.DrainBuffer(Buffer.Addr)
	t.AssertError(tx, err, -106)
}
