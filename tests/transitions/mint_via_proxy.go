package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

const swapOutput = "9871580343970"

func (tr *Transitions) SwapOutput() {
	Start("SwapOutput")

	p := tr.DeployAndUpgrade()
	zilSwap := tr.DeployZilSwap()
	supraToken := tr.DeploySupraToken()

	prepareZilSwap(p, zilSwap, supraToken)

	blockNum := p.GetBlockHeight()

	recipient := sdk.Cfg.Addr1

	AssertSuccess(zilSwap.SwapExactZILForTokens(supraToken.Contract.Addr, ToQA(10), "1", recipient, blockNum))
	AssertEqual(supraToken.BalanceOf(recipient).String(), swapOutput)
}

func (tr *Transitions) MintViaProxy() {
	Start("MintViaProxy")

	p := tr.DeployAndUpgrade()

	zilSwap := tr.DeployZilSwap()
	supraToken := tr.DeploySupraToken()

	prepareZilSwap(p, zilSwap, supraToken)

	recipient := sdk.Cfg.Addr1

	minterProxy := tr.DeployMinterProxy(supraToken.Contract.Addr, zilSwap.Contract.Addr)

	AssertSuccess(minterProxy.WithUser(sdk.Cfg.Key1).Mint(ToQA(10)))
	AssertEqual(supraToken.BalanceOf(recipient).String(), swapOutput)
}

func prepareZilSwap(p *contracts.Protocol, zilSwap *contracts.ZilSwap, supraToken *contracts.SupraToken) {
	zilLiqAmount := ToQA(1000)
	tokenLiqAmount := ToQA(1000)

	blockNum := p.GetBlockHeight()

	AssertSuccess(supraToken.IncreaseAllowance(zilSwap.Contract.Addr, ToQA(10000)))
	AssertSuccess(zilSwap.AddLiquidity(supraToken.Contract.Addr, zilLiqAmount, tokenLiqAmount, blockNum))
}
