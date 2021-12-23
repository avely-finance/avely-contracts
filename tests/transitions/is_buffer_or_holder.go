package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsBufferOrHolder() {

	Start("IsBufferOrHolder")

	p := tr.DeployAndUpgrade()

	tx, err := p.Aimpl.ClaimRewardsSuccessCallBack()
	AssertError(tx, err, -112)

	tx, err = p.Aimpl.DelegateStakeSuccessCallBack(Zil(1))
	AssertError(tx, err, -112)

	tx, err = p.Aimpl.CompleteWithdrawalSuccessCallBack()
	AssertError(tx, err, -112)
}
