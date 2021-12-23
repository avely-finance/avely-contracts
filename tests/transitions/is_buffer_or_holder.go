package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsBufferOrHolder() {

	t.Start("IsBufferOrHolder")

	p := tr.DeployAndUpgrade()

	tx, err := p.Aimpl.ClaimRewardsSuccessCallBack()
	t.AssertError(tx, err, -112)

	tx, err = p.Aimpl.DelegateStakeSuccessCallBack(Zil(1))
	t.AssertError(tx, err, -112)

	tx, err = p.Aimpl.CompleteWithdrawalSuccessCallBack()
	t.AssertError(tx, err, -112)
}
