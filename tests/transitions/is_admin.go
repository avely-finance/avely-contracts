package transitions

func (tr *Transitions) IsAdmin() {

	t.Start("IsAdmin")

	_, _, Aimpl, Buffer, Holder := tr.DeployAndUpgrade()

	// Use non-admin user for Buffer
	Buffer.UpdateWallet(tr.cfg.Key3)

	tx, err := Buffer.ChangeAzilSSNAddress(tr.cfg.Addr3)
	t.AssertError(tx, err, -402)
	tx, err = Buffer.ChangeAimplAddress(tr.cfg.Addr3)
	t.AssertError(tx, err, -402)
	tx, err = Buffer.ChangeZproxyAddress(tr.cfg.Addr3)
	t.AssertError(tx, err, -402)
	tx, err = Buffer.ChangeZimplAddress(tr.cfg.Addr3)
	t.AssertError(tx, err, -402)

	// Use non-admin user for Holder
	Holder.UpdateWallet(tr.cfg.Key2)

	tx, err = Holder.DelegateStake(zil(1))
	t.AssertError(tx, err, -305)
	tx, err = Holder.ChangeAzilSSNAddress(tr.cfg.Addr3)
	t.AssertError(tx, err, -305)
	tx, err = Holder.ChangeAimplAddress(tr.cfg.Addr3)
	t.AssertError(tx, err, -305)
	tx, err = Holder.ChangeZproxyAddress(tr.cfg.Addr3)
	t.AssertError(tx, err, -305)
	tx, err = Holder.ChangeZimplAddress(tr.cfg.Addr3)
	t.AssertError(tx, err, -305)

	// Use non-admin user for Aimpl
	Aimpl.UpdateWallet(tr.cfg.Key2)

	tx, err = Aimpl.ChangeZimplAddress(tr.cfg.Addr3)
	t.AssertError(tx, err, -106)
	tx, err = Aimpl.ChangeHolderAddress(tr.cfg.Addr3)
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
