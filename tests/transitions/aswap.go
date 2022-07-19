package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

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

	//toggle pause
	AssertSuccess(aswap.WithUser(init_owner_key).TogglePause())
	AssertEqual(Field(aswap, "pause"), "1")

	//set treasury fee
	new_fee := "12345"
	AssertEqual(Field(aswap, "treasury_fee"), "500")
	AssertSuccess(aswap.WithUser(init_owner_key).SetTreasuryFee(new_fee))
	AssertEqual(Field(aswap, "treasury_fee"), new_fee)

	//set treasury address
	treasury_address := sdk.Cfg.Addr2
	AssertEqual(Field(aswap, "treasury_address"), core.ZeroAddr)
	AssertSuccess(aswap.WithUser(init_owner_key).SetTreasuryAddress(treasury_address))
	AssertEqual(Field(aswap, "treasury_address"), treasury_address)

	//set liquidity fee
	new_fee = "1000"
	AssertEqual(Field(aswap, "liquidity_fee"), "10000")
	AssertSuccess(aswap.WithUser(init_owner_key).SetLiquidityFee(new_fee))
	AssertEqual(Field(aswap, "liquidity_fee"), new_fee)

	//do swap
	recipient := sdk.Cfg.Addr1
	expectedTreasuryRewards := "8100445524"
	expectedSwapOutput := "9900196014898"
	AssertSuccess(aswap.WithUser(init_owner_key).TogglePause())
	tx, _ := AssertSuccess(aswap.SwapExactZILForTokens(stzil.Contract.Addr, ToQA(100), "90", recipient, blockNum))
	AssertTransition(tx, Transition{
		aswap.Addr, //sender
		"AddFunds",
		treasury_address,
		expectedTreasuryRewards,
		ParamsMap{},
	})
	AssertEqual(stzil.BalanceOf(recipient).String(), expectedSwapOutput)
}
