package transitions

func (tr *Transitions) IsAimpl() {

	t.Start("IsAimpl")

	_, _, _, Buffer, Holder := tr.DeployAndUpgrade()

	// Use non-admin user for Buffer
	Buffer.UpdateWallet(key2)

	tx, err := Buffer.DelegateStake()
	t.AssertError(tx, err, -401)
	tx, err = Buffer.ClaimRewards()
	t.AssertError(tx, err, -401)
	tx, err = Buffer.RequestDelegatorSwap(Holder.Addr)
	t.AssertError(tx, err, -401)

	// Use non-admin user for Holder
	Holder.UpdateWallet(key2)

	tx, err = Holder.WithdrawStakeAmt(zil(1))
	t.AssertError(tx, err, -301)
	tx, err = Holder.CompleteWithdrawal()
	t.AssertError(tx, err, -301)
	tx, err = Holder.ClaimRewards()
	t.AssertError(tx, err, -301)
	tx, err = Holder.ConfirmDelegatorSwap(Buffer.Addr)
	t.AssertError(tx, err, -301)
}
