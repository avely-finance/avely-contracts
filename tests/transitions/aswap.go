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
	aswap := tr.DeployASwap(init_owner_addr)
	stzil := p.StZIL

	liquidityAmount := ToQA(1000)

	AssertSuccess(stzil.DelegateStake(liquidityAmount))

	blockNum := p.GetBlockHeight()

	//add liquidity
	AssertSuccess(stzil.IncreaseAllowance(aswap.Contract.Addr, ToQA(10000)))
	AssertSuccess(aswap.AddLiquidity(liquidityAmount, stzil.Contract.Addr, "0", liquidityAmount, blockNum))

	//toggle pause
	AssertSuccess(aswap.WithUser(init_owner_key).TogglePause())
	AssertEqual(Field(aswap, "pause"), "1")

	//set treasury fee
	new_fee := "750"
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
	expectedTreasuryRewards := "133333333333"
	expectedSwapOutput := "9887919312466"
	AssertSuccess(aswap.WithUser(init_owner_key).TogglePause())
	tx, _ := AssertSuccess(aswap.SwapExactZILForTokens(ToQA(100), stzil.Contract.Addr, "90", recipient, blockNum))
	AssertTransition(tx, Transition{
		aswap.Addr, //sender
		"AddFunds",
		treasury_address,
		expectedTreasuryRewards,
		ParamsMap{},
	})
	AssertEqual(stzil.BalanceOf(recipient).String(), expectedSwapOutput)

	//change owner
	new_owner_addr := sdk.Cfg.Addr3
	new_owner_key := sdk.Cfg.Key3

	//CodeStagingOwnerMissing
	tx, _ = aswap.WithUser(new_owner_key).ClaimOwner()
	AssertASwapError(tx, -11)

	AssertEqual(Field(aswap, "owner"), init_owner_addr)
	AssertSuccess(aswap.WithUser(init_owner_key).ChangeOwner(new_owner_addr))
	AssertEqual(Field(aswap, "staging_owner"), new_owner_addr)

	//CodeStagingOwnerInvalid
	tx, _ = aswap.WithUser(init_owner_key).ClaimOwner()
	AssertASwapError(tx, -12)

	//claim owner
	AssertSuccess(aswap.WithUser(new_owner_key).ClaimOwner())
	AssertEqual(Field(aswap, "owner"), new_owner_addr)
}
