package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsAimpl() {
	Start("IsAimpl")

	p := tr.DeployAndUpgrade()

	// Use non-admin user for Buffer
	p.GetBuffer().UpdateWallet(sdk.Cfg.Key2)

	tx, err := p.GetBuffer().DelegateStake()
	AssertError(tx, err, -401)
	tx, err = p.GetBuffer().ClaimRewards()
	AssertError(tx, err, -401)
	tx, err = p.GetBuffer().RequestDelegatorSwap(p.Holder.Addr)
	AssertError(tx, err, -401)

	// Use non-admin user for p.Holder
	p.Holder.UpdateWallet(sdk.Cfg.Key2)

	tx, err = p.Holder.WithdrawStakeAmt(ToZil(1))
	AssertError(tx, err, -301)
	tx, err = p.Holder.CompleteWithdrawal()
	AssertError(tx, err, -301)
	tx, err = p.Holder.ClaimRewards()
	AssertError(tx, err, -301)
	tx, err = p.Holder.ConfirmDelegatorSwap(p.GetBuffer().Addr)
	AssertError(tx, err, -301)
}
