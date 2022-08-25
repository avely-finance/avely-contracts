package transitions

import (
	"math/big"

	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

type LiquidityPair struct {
	zil   string
	stzil string
}

var Aswap *contracts.ASwap
var Stzil *contracts.StZIL
var Proto *contracts.Protocol

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

func (tr *Transitions) setupGoldenFlow() (*contracts.Protocol, *contracts.ASwap, *contracts.StZIL) {
	Proto = tr.DeployAndUpgrade()

	init_owner_addr := sdk.Cfg.Admin
	init_owner_key := sdk.Cfg.AdminKey
	Aswap = tr.DeployASwap(init_owner_addr)
	Stzil = Proto.StZIL

	//toggle pause
	AssertSuccess(Aswap.WithUser(init_owner_key).TogglePause())
	AssertEqual(Field(Aswap, "pause"), "1")

	//set treasury fee 0.3% = 1/333
	new_fee := "333"
	AssertSuccess(Aswap.WithUser(init_owner_key).SetTreasuryFee(new_fee))
	AssertEqual(Field(Aswap, "treasury_fee"), new_fee)

	//set treasury address
	treasury_address := sdk.Cfg.Addr2
	AssertEqual(Field(Aswap, "treasury_address"), core.ZeroAddr)
	AssertSuccess(Aswap.WithUser(init_owner_key).SetTreasuryAddress(treasury_address))
	AssertEqual(Field(Aswap, "treasury_address"), treasury_address)

	//set liquidity fee
	new_fee = "9940"
	AssertEqual(Field(Aswap, "liquidity_fee"), "10000")
	AssertSuccess(Aswap.WithUser(init_owner_key).SetLiquidityFee(new_fee))
	AssertEqual(Field(Aswap, "liquidity_fee"), new_fee)

	//toggle pause
	AssertSuccess(Aswap.WithUser(init_owner_key).TogglePause())
	AssertEqual(Field(Aswap, "pause"), "0")

	return Proto, Aswap, Stzil
}

func (tr *Transitions) ASwapGolden() {

	//TODO: to make this test suite completely correct, we need to maintain pools/balances/contributions here

	Start("ASwap Golden Flow")

	tr.setupGoldenFlow()
	blockNum := Proto.GetBlockHeight()

	//user 1 should have 1000 stzil
	liqPair := LiquidityPair{zil: ToQA(1000), stzil: ToQA(1000)}
	AssertSuccess(Stzil.WithUser(sdk.Cfg.Key1).DelegateStake(liqPair.zil))

	//1) user 1 add liquidity 1000:1000
	AssertSuccess(Stzil.WithUser(sdk.Cfg.Key1).IncreaseAllowance(Aswap.Contract.Addr, liqPair.stzil))
	//allowances[_sender][spender] := new_allowance;
	AssertEqual(Field(Stzil, "allowances", sdk.Cfg.Addr1, Aswap.Addr), liqPair.stzil)
	//func (a *Aswap) AddLiquidity(_amount, tokenAddr, minContributionAmount, tokenAmount string, blockNum int)
	AssertSuccess(Aswap.WithUser(sdk.Cfg.Key1).AddLiquidity(liqPair.zil, Stzil.Contract.Addr, liqPair.zil, liqPair.stzil, blockNum))
	//pools[Aswap.addr] := [zil, stzil];
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "0"), liqPair.zil)
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), liqPair.stzil)
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), ToQA(1000))
	AssertEqual(Stzil.BalanceOf(sdk.Cfg.Addr1).String(), "0")
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), ToQA(1000))

	//2) user 2 swap 500 stzil for zil
	amountIn := calcSwapInputZilForTokens(ToQA(500))
	recipient := sdk.Cfg.Addr2
	expectedSwapOutput := calcSwapOutputZilForTokens(amountIn)
	minTokenAmount := "1"
	tx, _ := AssertSuccess(Aswap.WithUser(sdk.Cfg.Key2).SwapExactZILForTokens(amountIn, Stzil.Contract.Addr, minTokenAmount, recipient, blockNum+1))
	AssertTransition(tx, Transition{
		Aswap.Addr, //sender
		"AddFunds",
		Field(Aswap, "treasury_address"),
		calcTreasuryReward(amountIn),
		ParamsMap{},
	})
	AssertEqual(Stzil.BalanceOf(recipient).String(), ToQA(500))
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), ToQA(500))
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), ToQA(500))
	user2StzilBalanceStr := expectedSwapOutput

	//3a) user 3 swap 250 stzil for zil
	amountIn = calcSwapInputZilForTokens(ToQA(250))
	recipient = sdk.Cfg.Addr3
	expectedSwapOutput = calcSwapOutputZilForTokens(amountIn)
	minTokenAmount = "1"
	tx, _ = AssertSuccess(Aswap.WithUser(sdk.Cfg.Key3).SwapExactZILForTokens(amountIn, Stzil.Contract.Addr, minTokenAmount, recipient, blockNum+1))
	AssertTransition(tx, Transition{
		Aswap.Addr, //sender
		"AddFunds",
		Field(Aswap, "treasury_address"),
		calcTreasuryReward(amountIn),
		ParamsMap{},
	})
	AssertEqual(Stzil.BalanceOf(recipient).String(), ToQA(250))
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), ToQA(250))
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), ToQA(250))

	//3b) user 3 add 250 stzil to pool
	pairIn := calcAddLiquidity(LiquidityPair{stzil: ToQA(250)})
	AssertSuccess(Stzil.WithUser(sdk.Cfg.Key3).IncreaseAllowance(Aswap.Contract.Addr, pairIn.stzil))
	//allowances[_sender][spender] := new_allowance;
	AssertEqual(Field(Stzil, "allowances", recipient, Aswap.Addr), pairIn.stzil)
	//func (a *Aswap) AddLiquidity(_amount, tokenAddr, minContributionAmount, tokenAmount string, blockNum int)
	AssertSuccess(Aswap.WithUser(sdk.Cfg.Key3).AddLiquidity(pairIn.zil, Stzil.Addr, "1", pairIn.stzil, blockNum+1))
	AssertEqual(Stzil.BalanceOf(sdk.Cfg.Addr3).String(), "0")
	//pools[Aswap.addr] := [zil, stzil];
	AssertEqual(Stzil.BalanceOf(sdk.Cfg.Addr3).String(), "0")
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), ToQA(500))
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), ToQA(500))

	//4) user 2 "buy" 100 stzil at swap
	amountIn = calcSwapInputZilForTokens(ToQA(100))
	recipient = sdk.Cfg.Addr2
	expectedSwapOutput = calcSwapOutputZilForTokens(amountIn)
	minTokenAmount = "1"
	tx, _ = AssertSuccess(Aswap.WithUser(sdk.Cfg.Key2).SwapExactZILForTokens(amountIn, Stzil.Contract.Addr, minTokenAmount, recipient, blockNum+1))
	AssertTransition(tx, Transition{
		Aswap.Addr, //sender
		"AddFunds",
		Field(Aswap, "treasury_address"),
		calcTreasuryReward(amountIn),
		ParamsMap{},
	})
	user2StzilBalance, _ := new(big.Int).SetString(user2StzilBalanceStr, 10)
	newOutput, _ := new(big.Int).SetString(expectedSwapOutput, 10)
	user2StzilBalance = user2StzilBalance.Add(user2StzilBalance, newOutput)
	AssertEqual(Stzil.BalanceOf(recipient).String(), user2StzilBalance.String())
	AssertEqual(Stzil.BalanceOf(recipient).String(), ToQA(600))
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), ToQA(400))
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), ToQA(400))

	//5) user1 remove liquidity
	pairOut := calcRemoveLiquidityOutput(sdk.Cfg.Addr1, ToQA(1000))
	//RemoveLiquidity(tokenAddress, contributionAmount, minZilAmount, minTokenAmount string, blockNum int) (*transaction.Transaction, error)
	tx, _ = Aswap.WithUser(sdk.Cfg.Key1).RemoveLiquidity(Stzil.Addr, ToQA(1000), "1", "1", blockNum+1)
	//data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
	//GetLog().Info(string(data))
	AssertTransition(tx, Transition{
		Aswap.Addr, //sender
		"AddFunds",
		sdk.Cfg.Addr1,
		pairOut.zil,
		ParamsMap{},
	})
	AssertEqual(Stzil.BalanceOf(sdk.Cfg.Addr1).String(), pairOut.stzil)
	AssertEqual(Stzil.BalanceOf(sdk.Cfg.Addr1).String(), ToQA(200))

	//6) user3 remove liquidity
	pairOut3 := calcRemoveLiquidityOutput(sdk.Cfg.Addr3, "all")
	user3Contribution := Field(Aswap, "balances", sdk.Cfg.Addr3, Stzil.Addr)
	//RemoveLiquidity(tokenAddress, contributionAmount, minZilAmount, minTokenAmount string, blockNum int) (*transaction.Transaction, error)
	tx, _ = Aswap.WithUser(sdk.Cfg.Key3).RemoveLiquidity(Stzil.Addr, user3Contribution, "1", "1", blockNum+1)
	//data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
	//GetLog().Info(string(data))
	AssertTransition(tx, Transition{
		Aswap.Addr, //sender
		"AddFunds",
		sdk.Cfg.Addr3,
		pairOut3.zil,
		ParamsMap{},
	})
	AssertEqual(Stzil.BalanceOf(sdk.Cfg.Addr1).String(), pairOut3.stzil)
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), "0")
}

func calcAddLiquidity(pair LiquidityPair) LiquidityPair {
	if Field(Aswap, "pools", Stzil.Addr) == "" {
		panic("Pool does not exist")
	} else if pair.stzil != "" {
		//dY = dX * Y / X
		//dX is always the QA transferred
		//dX = dY * X / Y, where X - zilReserve, Y - stzilReserve
		dY, _ := new(big.Int).SetString(pair.stzil, 10)
		zilReserve, _ := new(big.Int).SetString(Field(Aswap, "pools", Stzil.Addr, "arguments", "0"), 10)
		stzilReserve, _ := new(big.Int).SetString(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), 10)
		dX := new(big.Int).Mul(dY, zilReserve)
		dX = dX.Div(dX, stzilReserve)
		return LiquidityPair{stzil: pair.stzil, zil: dX.String()}
	}
	panic("Calculation of stzil by zil not implemented")
}

func calcRemoveLiquidityOutput(senderAddr, contribution string) LiquidityPair {
	if Field(Aswap, "pools", Stzil.Addr) == "" {
		panic("Pool does not exist")
	} else if Field(Aswap, "balances", senderAddr, Stzil.Addr) == "" {
		panic("User has no liquidity")
	}
	//zil_amount_u256 = fraction contribution_amount_u256 total_contribution_u256 zil_reserve_u256;
	//fraction(d,x,y)= d * y / x
	zilReserve, _ := new(big.Int).SetString(Field(Aswap, "pools", Stzil.Addr, "arguments", "0"), 10)
	stzilReserve, _ := new(big.Int).SetString(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), 10)
	totalContribution, _ := new(big.Int).SetString(Field(Aswap, "total_contributions", Stzil.Addr), 10)
	var userContribution *big.Int
	if contribution == "all" {
		userContribution, _ = new(big.Int).SetString(Field(Aswap, "balances", senderAddr, Stzil.Addr), 10)
	} else {
		//contribution <= Field(Aswap, "balances", senderAddr, Stzil.Addr) else error
		userContribution, _ = new(big.Int).SetString(contribution, 10)
	}
	zilAmount := new(big.Int).Mul(userContribution, zilReserve)
	zilAmount = zilAmount.Div(zilAmount, totalContribution)

	//token_amount_u256 = fraction contribution_amount_u256 total_contribution_u256 token_reserve_u256;
	stzilAmount := new(big.Int).Mul(userContribution, stzilReserve)
	stzilAmount = stzilAmount.Div(stzilAmount, totalContribution)

	return LiquidityPair{zil: zilAmount.String(), stzil: stzilAmount.String()}
}

func calcSwapOutputZilForTokens(amntIn string) string {
	//1) substract treasury fee
	amountIn, _ := new(big.Int).SetString(amntIn, 10)
	treasuryFeeCalculated, _ := new(big.Int).SetString(calcTreasuryReward(amntIn), 10)
	afterTreasuryFee := new(big.Int).Sub(amountIn, treasuryFeeCalculated)

	//2) constant product formula
	zilReserve, _ := new(big.Int).SetString(Field(Aswap, "pools", Stzil.Addr, "arguments", "0"), 10)
	tokenReserve, _ := new(big.Int).SetString(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), 10)
	lpFee, _ := new(big.Int).SetString(Field(Aswap, "liquidity_fee"), 10)
	lpFeeDenom, _ := new(big.Int).SetString("10000", 10)
	afterLPFee := new(big.Int).Mul(afterTreasuryFee, lpFee)

	nominator := new(big.Int).Mul(afterLPFee, tokenReserve)
	denom := new(big.Int).Mul(lpFeeDenom, zilReserve)
	denom = denom.Add(denom, afterLPFee)

	amountOut := new(big.Int).Div(nominator, denom)
	return amountOut.String()
}

func calcSwapInputZilForTokens(amntOut string) string {
	amountOut, _ := new(big.Int).SetString(amntOut, 10)

	treasuryFee, _ := new(big.Int).SetString(Field(Aswap, "treasury_fee"), 10)
	zilReserve, _ := new(big.Int).SetString(Field(Aswap, "pools", Stzil.Addr, "arguments", "0"), 10)
	tokenReserve, _ := new(big.Int).SetString(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), 10)
	lpFee, _ := new(big.Int).SetString(Field(Aswap, "liquidity_fee"), 10)
	lpFeeDenom, _ := new(big.Int).SetString("10000", 10)

	nominator := new(big.Int).Mul(amountOut, zilReserve)
	nominator = nominator.Mul(nominator, lpFeeDenom)
	denom := new(big.Int).Sub(tokenReserve, amountOut)
	denom = denom.Mul(denom, lpFee)

	afterTreasuryFee := new(big.Int).Div(nominator, denom)
	//fix big num math innacuracy, see https://github.com/Uniswap/v2-periphery/blob/0335e8f7e1bd1e8d8329fd300aea2ef2f36dd19f/contracts/libraries/UniswapV2Library.sol#L58
	afterTreasuryFee = afterTreasuryFee.Add(afterTreasuryFee, big.NewInt(1))
	beforeTreasuryFee := new(big.Int).Mul(afterTreasuryFee, treasuryFee)
	denom = new(big.Int).Sub(treasuryFee, big.NewInt(1))
	beforeTreasuryFee = beforeTreasuryFee.Div(beforeTreasuryFee, denom)

	return beforeTreasuryFee.String()
}

func calcTreasuryReward(amount string) string {
	treasury_fee := Field(Aswap, "treasury_fee")
	biAmt, _ := new(big.Int).SetString(amount, 10)
	biFee, _ := new(big.Int).SetString(treasury_fee, 10)
	feeRes := new(big.Int).Div(biAmt, biFee)
	return feeRes.String()
}
