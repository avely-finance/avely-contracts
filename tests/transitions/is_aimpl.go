package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsAimpl() {
	t.Start("IsAimpl")

	p := DeployAndUpgrade()

	// Use non-admin user for Buffer
	p.Buffer.UpdateWallet(sdk.Cfg.Key2)

	tx, err := p.Buffer.DelegateStake()
	t.AssertError(tx, err, -401)
	tx, err = p.Buffer.ClaimRewards()
	t.AssertError(tx, err, -401)
	tx, err = p.Buffer.RequestDelegatorSwap(p.Holder.Addr)
	t.AssertError(tx, err, -401)

	// Use non-admin user for p.Holder
	p.Holder.UpdateWallet(sdk.Cfg.Key2)

	tx, err = p.Holder.WithdrawStakeAmt(Zil(1))
	t.AssertError(tx, err, -301)
	tx, err = p.Holder.CompleteWithdrawal()
	t.AssertError(tx, err, -301)
	tx, err = p.Holder.ClaimRewards()
	t.AssertError(tx, err, -301)
	tx, err = p.Holder.ConfirmDelegatorSwap(p.Buffer.Addr)
	t.AssertError(tx, err, -301)
}
