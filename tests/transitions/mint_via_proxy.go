package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) MintViaProxy() {
	Start("MintViaProxy")

	p := tr.DeployAndUpgrade()

	zilSwap := tr.DeployZilSwap()
	supraToken := tr.DeploySupraToken()

	zilLiqAmount := ToQA(1000)
	tokenLiqAmount := ToQA(1000)

	blockNum := p.GetBlockHeight()

	AssertSuccess(supraToken.IncreaseAllowance(zilSwap.Contract.Addr, ToQA(10000)))
	AssertSuccess(zilSwap.AddLiquidity(supraToken.Contract.Addr, zilLiqAmount, tokenLiqAmount, blockNum))

	recipient := sdk.Cfg.Addr1

	AssertSuccess(zilSwap.SwapExactZILForTokens(supraToken.Contract.Addr, ToQA(10), "1", recipient, blockNum))
	AssertEqual(supraToken.BalanceOf(recipient).String(), "9871580343970")
}
