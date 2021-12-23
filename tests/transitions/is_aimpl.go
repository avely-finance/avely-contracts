package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsAimpl() {
	t.Start("IsAimpl")

	_, _, _, Buffer, Holder := sdk.DeployAndUpgrade()

	// Use non-admin user for Buffer
	Buffer.UpdateWallet(sdk.Cfg.Key2)

	tx, err := Buffer.DelegateStake()
	t.AssertError(tx, err, -401)
	tx, err = Buffer.ClaimRewards()
	t.AssertError(tx, err, -401)
	tx, err = Buffer.RequestDelegatorSwap(Holder.Addr)
	t.AssertError(tx, err, -401)

	// Use non-admin user for Holder
	Holder.UpdateWallet(sdk.Cfg.Key2)

	tx, err = Holder.WithdrawStakeAmt(Zil(1))
	t.AssertError(tx, err, -301)
	tx, err = Holder.CompleteWithdrawal()
	t.AssertError(tx, err, -301)
	tx, err = Holder.ClaimRewards()
	t.AssertError(tx, err, -301)
	tx, err = Holder.ConfirmDelegatorSwap(Buffer.Addr)
	t.AssertError(tx, err, -301)
}
