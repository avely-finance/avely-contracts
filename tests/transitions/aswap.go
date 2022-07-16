package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

const ASwapOutput = "9900990099009"

func (tr *Transitions) ASwap() {
	Start("Swap via ASwap")

	p := tr.DeployAndUpgrade()

	init_owner := sdk.Cfg.Admin
	operators := []string{sdk.Cfg.Admin}
	aswap := tr.DeployASwap(init_owner, operators)
	stzil := p.StZIL

	liquidityAmount := ToQA(1000)

	AssertSuccess(stzil.DelegateStake(liquidityAmount))

	blockNum := p.GetBlockHeight()

	//add liquidity
	AssertSuccess(stzil.IncreaseAllowance(aswap.Contract.Addr, ToQA(10000)))
	AssertSuccess(aswap.AddLiquidity(stzil.Contract.Addr, liquidityAmount, liquidityAmount, blockNum))

	//do swap
	recipient := sdk.Cfg.Addr1
	AssertSuccess(aswap.SwapExactZILForTokens(stzil.Contract.Addr, ToQA(10), "1", recipient, blockNum))
	AssertEqual(stzil.BalanceOf(recipient).String(), ASwapOutput)
}
