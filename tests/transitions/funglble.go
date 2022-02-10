package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

const swapOutput = "9871580343970"

func (tr *Transitions) SwapOutput() {
	Start("Swap via ZilSwap")

	p := tr.DeployAndUpgrade()
	zilSwap := tr.DeployZilSwap()
	token := p.Aimpl

	liquidityAmount := ToQA(1000)

	// supraToken := tr.DeploySupraToken()
	// Success delegate 1000 ZIL for Liquidity
	AssertSuccess(p.Aimpl.DelegateStake(liquidityAmount))

	prepareZilSwap(p, zilSwap, token, liquidityAmount)

	blockNum := p.GetBlockHeight()
	recipient := sdk.Cfg.Addr1

	AssertSuccess(zilSwap.SwapExactZILForTokens(token.Contract.Addr, ToQA(10), "1", recipient, blockNum))
	AssertEqual(token.BalanceOf(recipient).String(), swapOutput)
}

// func (tr *Transitions) MintViaProxy() {
// 	Start("MintViaProxy")

// 	p := tr.DeployAndUpgrade()

// 	zilSwap := tr.DeployZilSwap()
// 	supraToken := tr.DeploySupraToken()

// 	prepareZilSwap(p, zilSwap, supraToken)

// 	recipient := sdk.Cfg.Addr1

// 	minterProxy := tr.DeployMinterProxy(supraToken.Contract.Addr, zilSwap.Contract.Addr)

// 	blockNum := p.GetBlockHeight()

// 	AssertSuccess(minterProxy.WithUser(sdk.Cfg.Key1).Mint(ToQA(10), swapOutput, strconv.Itoa(blockNum+1)))
// 	AssertEqual(supraToken.BalanceOf(recipient).String(), swapOutput)
// }

func prepareZilSwap(p *contracts.Protocol, zilSwap *contracts.ZilSwap, azil *contracts.AZil, liquidityAmount string) {
	blockNum := p.GetBlockHeight()

	AssertSuccess(azil.IncreaseAllowance(zilSwap.Contract.Addr, ToQA(10000)))
	AssertSuccess(zilSwap.AddLiquidity(azil.Contract.Addr, liquidityAmount, liquidityAmount, blockNum))
}
