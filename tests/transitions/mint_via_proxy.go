package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) MintViaProxy() {
	Start("IsAimpl")

	zilSwap := tr.DeployZilSwap()
	supraToken := tr.DeploySupraToken()

	zilAmount := ToQA(100)
	tokenAmount := ToQA(100)

	zilSwap.AddLiquidity(supraToken, zilAmount, tokenAmount)
}
