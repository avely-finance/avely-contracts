package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsBufferOrHolder() {

	Start("IsBufferOrHolder")

	p := tr.DeployAndUpgrade()

	tx, _ := p.StZIL.ClaimRewardsSuccessCallBack()
	AssertError(tx, p.StZIL.ErrorCode("BufferOrHolderValidationFailed"))

	tx, _ = p.StZIL.DelegateStakeSuccessCallBack(ToZil(1))
	AssertError(tx, p.StZIL.ErrorCode("BufferOrHolderValidationFailed"))

	tx, _ = p.StZIL.CompleteWithdrawalSuccessCallBack()
	AssertError(tx, p.StZIL.ErrorCode("BufferOrHolderValidationFailed"))
}
