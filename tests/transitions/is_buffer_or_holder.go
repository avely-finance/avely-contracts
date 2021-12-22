package transitions

import (
	. "Azil/test/helpers"
)

func (tr *Transitions) IsBufferOrHolder() {

	t.Start("IsBufferOrHolder")

	_, _, Aimpl, _, _ := tr.DeployAndUpgrade()

	tx, err := Aimpl.ClaimRewardsSuccessCallBack()
	t.AssertError(tx, err, -112)

	tx, err = Aimpl.DelegateStakeSuccessCallBack(Zil(1))
	t.AssertError(tx, err, -112)

	tx, err = Aimpl.CompleteWithdrawalSuccessCallBack()
	t.AssertError(tx, err, -112)
}
