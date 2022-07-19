package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

const ASwapOutput = "9881383789778"

func (tr *Transitions) ASwap() {
	Start("Swap via ASwap")

	p := tr.DeployAndUpgrade()

	init_owner_addr := sdk.Cfg.Admin
	init_owner_key := sdk.Cfg.AdminKey
	operators := []string{sdk.Cfg.Admin}
	aswap := tr.DeployASwap(init_owner_addr, operators)
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

	//toggle pause
	AssertSuccess(aswap.WithUser(init_owner_key).TogglePause())
	AssertEqual(Field(aswap, "pause"), "1")

	//set treasury fee
	new_fee := "12345"
	AssertEqual(Field(aswap, "treasury_fee"), "500")
	AssertSuccess(aswap.WithUser(init_owner_key).SetTreasuryFee(new_fee))
	AssertEqual(Field(aswap, "treasury_fee"), new_fee)

	//set liquidity fee
	new_fee = "23456"
	AssertEqual(Field(aswap, "liquidity_fee"), "10000")
	AssertSuccess(aswap.WithUser(init_owner_key).SetLiquidityFee(new_fee))
	AssertEqual(Field(aswap, "liquidity_fee"), new_fee)
}
