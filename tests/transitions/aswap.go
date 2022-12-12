package transitions

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
	"github.com/tyler-smith/go-bip39"
)

type LiquidityPair struct {
	zil   string
	stzil string
}

var Aswap *contracts.ASwap
var Stzil *contracts.StZIL
var Proto *contracts.Protocol

var Addr1, Addr2, Addr3, Key1, Key2, Key3 string
var Addresses []string
var Archive BalanceSheet

/*
before, _ := new(big.Int).SetString(balance1_before, 10)
after, _ := new(big.Int).SetString(deploy.GetBalance(addr1), 10)
gasPrice, _ := new(big.Int).SetString(tx.GasPrice, 10)
gasUsed, _ := new(big.Int).SetString(tx.Receipt.CumulativeGas, 10)
gasFee := new(big.Int).Mul(gasPrice, gasUsed)
withdrawed, _ := new(big.Int).SetString(zil10, 10)
result := new(big.Int).Sub(before, gasFee)
result = result.Add(result, withdrawed)
result = result.Sub(result, after)
println(fmt.Sprintf("====%s====", result))
t.LogDebug()
//t.AssertEqual(balance1_before, fmt.Sprintf("%s", needed))*/

func (tr *Transitions) ASwap() {
	aswapBasic(tr)
	aswapGolden(tr)
	aswapMultisig(tr)
	aswapOwnerOnly(tr)
}

func aswapBasic(tr *Transitions) {
	Start("Swap via ASwap")

	p := tr.DeployAndUpgrade()

	init_owner_addr := utils.GetAddressByWallet(celestials.Admin)

	aswap := tr.DeployASwap(init_owner_addr)
	stzil := p.StZIL

	liquidityAmount := ToQA(1000)

	AssertSuccess(stzil.DelegateStake(liquidityAmount))

	blockNum := p.GetBlockHeight()

	//add liquidity
	AssertSuccess(stzil.IncreaseAllowance(aswap.Contract.Addr, ToQA(10000)))
	AssertSuccess(aswap.AddLiquidity(liquidityAmount, stzil.Contract.Addr, "0", liquidityAmount, blockNum))

	//toggle pause
	AssertSuccess(aswap.TogglePause())
	AssertEqual(Field(aswap, "pause"), "1")

	//set treasury fee, check zero
	new_fee := "0"
	AssertEqual(Field(aswap, "treasury_fee"), "500")
	AssertSuccess(aswap.SetTreasuryFee(new_fee))
	AssertEqual(Field(aswap, "treasury_fee"), new_fee)

	//set treasury fee
	new_fee = "750"
	AssertSuccess(aswap.SetTreasuryFee(new_fee))
	AssertEqual(Field(aswap, "treasury_fee"), new_fee)

	//set treasury address
	treasury_address := p.Treasury.Addr
	AssertEqual(Field(aswap, "treasury_address"), core.ZeroAddr)
	AssertSuccess(aswap.SetTreasuryAddress(treasury_address))
	AssertEqual(Field(aswap, "treasury_address"), treasury_address)

	//set liquidity fee
	new_fee = "1000"
	AssertEqual(Field(aswap, "liquidity_fee"), "10000")
	AssertSuccess(aswap.SetLiquidityFee(new_fee))
	AssertEqual(Field(aswap, "liquidity_fee"), new_fee)

	//do swap
	expectedTreasuryRewards := "133333333333"
	expectedSwapOutput := "9887919312466"
	AssertSuccess(aswap.TogglePause())
	aliceAddr := utils.GetAddressByWallet(alice)

	tx, _ := AssertSuccess(aswap.SwapExactZILForTokens(ToQA(100), stzil.Contract.Addr, "90", aliceAddr, blockNum))
	AssertTransition(tx, Transition{
		aswap.Addr, //sender
		"AddFunds",
		treasury_address,
		expectedTreasuryRewards,
		ParamsMap{},
	})
	AssertEqual(stzil.BalanceOf(aliceAddr).String(), expectedSwapOutput)

	//change owner
	new_owner_addr := sdk.Cfg.Addr3
	new_owner_key := sdk.Cfg.Key3

	//try to claim owner without staging owner, expect error
	tx, _ = aswap.WithUser(new_owner_key).ClaimOwner()
	AssertASwapError(tx, aswap.ErrorCode("CodeStagingOwnerMissing"))

	AssertEqual(Field(aswap, "owner"), init_owner_addr)

	aswap.SetSigner(celestials.Admin) // use original owner
	AssertSuccess(aswap.ChangeOwner(new_owner_addr))
	AssertEqual(Field(aswap, "staging_owner"), new_owner_addr)

	//try to claim owner with invalid user, expect error
	tx, _ = aswap.ClaimOwner()
	AssertASwapError(tx, aswap.ErrorCode("CodeStagingOwnerInvalid"))

	//claim owner
	AssertSuccess(aswap.WithUser(new_owner_key).ClaimOwner())
	AssertEqual(Field(aswap, "owner"), new_owner_addr)
}

func aswapMultisig(tr *Transitions) {
	txIdLocal := 0

	//deploy multisig
	owners := []string{utils.GetAddressByWallet(alice)}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	//deploy aswap, set owner to multisig contract
	init_owner := multisig.Addr
	aswap := tr.DeployASwap(init_owner)

	//test ASwap.TogglePause
	multisig.SetSigner(alice)
	AssertMultisigSuccess(multisig.SubmitTogglePauseTransaction(aswap.Addr))
	AssertMultisigSuccess(multisig.ExecuteTransaction(txIdLocal))
	AssertEqual(Field(aswap, "pause"), "1")

	//test ASwap.SetTreasuryFee()
	txIdLocal++
	new_fee := "12345"
	AssertEqual(Field(aswap, "treasury_fee"), "500")
	AssertMultisigSuccess(multisig.SubmitSetTreasuryFeeTransaction(aswap.Addr, new_fee))
	AssertMultisigSuccess(multisig.ExecuteTransaction(txIdLocal))
	AssertEqual(Field(aswap, "treasury_fee"), new_fee)

	//test ASwap.SetLiquidityFee()
	txIdLocal++
	new_fee = "23456"
	AssertEqual(Field(aswap, "liquidity_fee"), "10000")
	AssertMultisigSuccess(multisig.SubmitSetLiquidityFeeTransaction(aswap.Addr, new_fee))
	AssertMultisigSuccess(multisig.ExecuteTransaction(txIdLocal))
	AssertEqual(Field(aswap, "liquidity_fee"), new_fee)

	//test ASwap.SetTreasuryAddress()
	txIdLocal++
	new_address := sdk.Cfg.Addr3
	AssertEqual(Field(aswap, "treasury_address"), core.ZeroAddr)
	AssertMultisigSuccess(multisig.SubmitSetTreasuryAddressTransaction(aswap.Addr, new_address))
	AssertMultisigSuccess(multisig.ExecuteTransaction(txIdLocal))
	AssertEqual(Field(aswap, "treasury_address"), new_address)

	//deploy other multisig contract
	newSignCount := 1
	newOwners := []string{utils.GetAddressByWallet(bob)}
	newMultisig := tr.DeployMultisigWallet(newOwners, newSignCount)

	//test ASwap.ChangeOwner()
	txIdLocal++
	AssertMultisigSuccess(multisig.SubmitChangeOwnerTransaction(aswap.Addr, newMultisig.Addr))
	AssertMultisigSuccess(multisig.ExecuteTransaction(txIdLocal))
	AssertEqual(Field(aswap, "staging_owner"), newMultisig.Addr)

	//test ASwap.ClaimOwner()
	//first transaction id is 0 for newly deployed multisig contract
	newMultisig.SetSigner(bob)
	AssertMultisigSuccess(newMultisig.SubmitClaimOwnerTransaction(aswap.Addr))
	AssertMultisigSuccess(newMultisig.ExecuteTransaction(0))
	AssertEqual(Field(aswap, "owner"), newMultisig.Addr)

}

func (tr *Transitions) setupGoldenFlow() (*contracts.Protocol, *contracts.ASwap, *contracts.StZIL) {

	Proto = tr.DeployAndUpgrade()

	init_owner_addr := utils.GetAddressByWallet(celestials.Admin)
	Aswap = tr.DeployASwap(init_owner_addr)
	Stzil = Proto.StZIL

	//toggle pause
	AssertSuccess(Aswap.TogglePause())
	AssertEqual(Field(Aswap, "pause"), "1")

	//set treasury fee 0.3% = 1/333
	new_fee := "333"
	AssertSuccess(Aswap.SetTreasuryFee(new_fee))
	AssertEqual(Field(Aswap, "treasury_fee"), new_fee)

	//set treasury address
	treasury_address := Proto.Treasury.Addr
	AssertEqual(Field(Aswap, "treasury_address"), core.ZeroAddr)
	AssertSuccess(Aswap.SetTreasuryAddress(treasury_address))
	AssertEqual(Field(Aswap, "treasury_address"), treasury_address)

	//set liquidity fee
	new_fee = "9940"
	AssertEqual(Field(Aswap, "liquidity_fee"), "10000")
	AssertSuccess(Aswap.SetLiquidityFee(new_fee))
	AssertEqual(Field(Aswap, "liquidity_fee"), new_fee)

	//toggle pause
	AssertSuccess(Aswap.TogglePause())
	AssertEqual(Field(Aswap, "pause"), "0")

	//Generate a mnemonic for memorization or user-friendly seeds
	entropy, _ := bip39.NewEntropy(128) //256
	mnemonic, _ := bip39.NewMnemonic(entropy)

	//mnemonic := "bug feature framework lava jelly keep device journey bean mango rocket festival"
	account1, _ := account.NewDefaultHDAccount(mnemonic, uint32(1))
	account2, _ := account.NewDefaultHDAccount(mnemonic, uint32(2))
	account3, _ := account.NewDefaultHDAccount(mnemonic, uint32(3))

	Addr1 = "0x" + account1.Address
	Addr2 = "0x" + account2.Address
	Addr3 = "0x" + account3.Address
	Key1 = util.EncodeHex(account1.PrivateKey)
	Key2 = util.EncodeHex(account2.PrivateKey)
	Key3 = util.EncodeHex(account3.PrivateKey)
	Addresses = make([]string, 4)
	Addresses[1] = Addr1
	Addresses[2] = Addr2
	Addresses[3] = Addr3
	Archive = make([]BalanceRecord, 0)

	sdk.AddFunds(celestials.Admin, Addr1, ToQA(5000))
	sdk.AddFunds(celestials.Admin, Addr2, ToQA(5000))
	sdk.AddFunds(celestials.Admin, Addr3, ToQA(5000))

	return Proto, Aswap, Stzil
}

func aswapGolden(tr *Transitions) {

	//TODO: to make this test suite completely correct, we need to maintain pools/balances/contributions here

	Start("ASwap Golden Flow")

	tr.setupGoldenFlow()
	blockNum := Proto.GetBlockHeight()

	//mint some stzil and transfer 1000stzil to user 1, 500stzil to user 2
	Stzil.SetSigner(celestials.Admin)

	AssertSuccess(Stzil.DelegateStake(ToQA(1500)))
	AssertSuccess(Stzil.Transfer(Addr1, ToQA(1000)))
	AssertSuccess(Stzil.Transfer(Addr2, ToQA(500)))

	recordBalance(-1, nil)

	//1) user 1 add liquidity 1000:1000
	liqPair := LiquidityPair{zil: ToQA(1000), stzil: ToQA(1000)}
	tx, _ := AssertSuccess(Stzil.WithUser(Key1).IncreaseAllowance(Aswap.Contract.Addr, liqPair.stzil))
	recordBalance(1, tx)
	//allowances[_sender][spender] := new_allowance;
	AssertEqual(Field(Stzil, "allowances", Addr1, Aswap.Addr), liqPair.stzil)
	//func (a *Aswap) AddLiquidity(_amount, tokenAddr, minContributionAmount, tokenAmount string, blockNum int)
	tx, _ = AssertSuccess(Aswap.WithUser(Key1).AddLiquidity(liqPair.zil, Stzil.Contract.Addr, liqPair.zil, liqPair.stzil, blockNum))
	recordBalance(1, tx)
	//pools[Aswap.addr] := [zil, stzil];
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "0"), liqPair.zil)
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), liqPair.stzil)
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), ToQA(1000))
	AssertEqual(Stzil.BalanceOf(Addr1).String(), "0")
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), ToQA(1000))

	//2) user 2 sell 500 stzil to swap for zil
	exactStzil := ToQA(500)
	tx, _ = AssertSuccess(Stzil.WithUser(Key2).IncreaseAllowance(Aswap.Contract.Addr, exactStzil))
	recordBalance(2, tx)
	expectedSwapOutput, treasuryFee := calcOutputSwapExactTokensForZil(exactStzil)
	minZilAmount := "1"
	//SwapExactTokensForZIL(tokenAddress, tokenAmount, minZilAmount, recipientAddress string, blockNum int)
	tx, _ = AssertSuccess(Aswap.WithUser(Key2).SwapExactTokensForZIL(Stzil.Contract.Addr, exactStzil, minZilAmount, Addr2, blockNum+1))
	recordBalance(2, tx)
	AssertTransition(tx, Transition{
		Aswap.Addr, //sender
		"AddFunds",
		Field(Aswap, "treasury_address"),
		treasuryFee,
		ParamsMap{},
	})
	AssertTransition(tx, Transition{
		Aswap.Addr, //sender
		"AddFunds",
		Addr2,
		expectedSwapOutput,
		ParamsMap{},
	})
	AssertEqual(Stzil.BalanceOf(Addr2).String(), "0")
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), ToQA(1500))
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), ToQA(1500))

	//3a) user 3 buy 250 stzil from swap for zil
	exactZil := calcInputSwapExactZilForTokens(ToQA(250))
	expectedSwapOutput = calcOutputSwapExactZilForTokens(exactZil)
	minTokenAmount := "1"

	tx, _ = AssertSuccess(Aswap.WithUser(Key3).SwapExactZILForTokens(exactZil, Stzil.Contract.Addr, minTokenAmount, Addr3, blockNum+1))
	recordBalance(3, tx)
	AssertTransition(tx, Transition{
		Aswap.Addr, //sender
		"AddFunds",
		Field(Aswap, "treasury_address"),
		calcTreasuryReward(exactZil),
		ParamsMap{},
	})
	AssertEqual(Stzil.BalanceOf(Addr3).String(), ToQA(250))
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), ToQA(1250))
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), ToQA(1250))

	//3b) user 3 add 250 stzil to pool
	pairIn := calcAddLiquidity(LiquidityPair{stzil: ToQA(250)})
	tx, _ = AssertSuccess(Stzil.WithUser(Key3).IncreaseAllowance(Aswap.Contract.Addr, pairIn.stzil))
	recordBalance(3, tx)
	//allowances[_sender][spender] := new_allowance;
	AssertEqual(Field(Stzil, "allowances", Addr3, Aswap.Addr), pairIn.stzil)
	//func (a *Aswap) AddLiquidity(_amount, tokenAddr, minContributionAmount, tokenAmount string, blockNum int)
	tx, _ = AssertSuccess(Aswap.WithUser(Key3).AddLiquidity(pairIn.zil, Stzil.Addr, "1", pairIn.stzil, blockNum+1))
	recordBalance(3, tx)
	AssertEqual(Stzil.BalanceOf(Addr3).String(), "0")
	//pools[Aswap.addr] := [zil, stzil];
	AssertEqual(Stzil.BalanceOf(Addr3).String(), "0")
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), ToQA(1500))
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), ToQA(1500))

	//4) user 2 "buy" 100 stzil at swap
	exactZil = calcInputSwapExactZilForTokens(ToQA(100))
	expectedSwapOutput = calcOutputSwapExactZilForTokens(exactZil)
	minTokenAmount = "1"
	tx, _ = AssertSuccess(Aswap.WithUser(Key2).SwapExactZILForTokens(exactZil, Stzil.Contract.Addr, minTokenAmount, Addr2, blockNum+1))
	recordBalance(2, tx)
	AssertTransition(tx, Transition{
		Aswap.Addr, //sender
		"AddFunds",
		Field(Aswap, "treasury_address"),
		calcTreasuryReward(exactZil),
		ParamsMap{},
	})
	user2StzilBalance := big.NewInt(0)
	newOutput, _ := new(big.Int).SetString(expectedSwapOutput, 10)
	user2StzilBalance = user2StzilBalance.Add(user2StzilBalance, newOutput)
	AssertEqual(Stzil.BalanceOf(Addr2).String(), user2StzilBalance.String())

	//5) user1 remove liquidity
	contribution := ToQA(1000)
	pairOut := calcRemoveLiquidityOutput(Addr1, contribution)
	//RemoveLiquidity(tokenAddress, contributionAmount, minZilAmount, minTokenAmount string, blockNum int) (*transaction.Transaction, error)
	tx, _ = AssertSuccess(Aswap.WithUser(Key1).RemoveLiquidity(Stzil.Addr, contribution, "1", "1", blockNum+1))
	recordBalance(1, tx)
	//data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
	//GetLog().Info(string(data))
	AssertTransition(tx, Transition{
		Aswap.Addr, //sender
		"AddFunds",
		Addr1,
		pairOut.zil,
		ParamsMap{},
	})
	AssertTransition(tx, Transition{
		Aswap.Addr, //sender
		"Transfer",
		Stzil.Addr,
		"0",
		ParamsMap{"to": Addr1, "amount": pairOut.stzil},
	})
	AssertEqual(Stzil.BalanceOf(Addr1).String(), pairOut.stzil)
	AssertEqual(Field(Aswap, "balances", Addr1), "")

	//6) user3 remove liquidity
	pairOut3 := calcRemoveLiquidityOutput(Addr3, "all")
	user3Contribution := Field(Aswap, "balances", Addr3, Stzil.Addr)
	//RemoveLiquidity(tokenAddress, contributionAmount, minZilAmount, minTokenAmount string, blockNum int) (*transaction.Transaction, error)
	tx, _ = AssertSuccess(Aswap.WithUser(Key3).RemoveLiquidity(Stzil.Addr, user3Contribution, "1", "1", blockNum+1))
	recordBalance(3, tx)
	//data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
	//GetLog().Info(string(data))
	AssertTransition(tx, Transition{
		Aswap.Addr, //sender
		"AddFunds",
		Addr3,
		pairOut3.zil,
		ParamsMap{},
	})
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), "0")
	AssertEqual(Field(Aswap, "balances", Addr3), "")
	AssertEqual(Field(Aswap, "balances"), "{}")

	fmt.Println(recapBalance())
}

func getTxFee(tx *transaction.Transaction) *big.Int {
	/*
		before, _ := new(big.Int).SetString(balance1_before, 10)
		after, _ := new(big.Int).SetString(deploy.GetBalance(addr1), 10)
		gasPrice, _ := new(big.Int).SetString(tx.GasPrice, 10)
		gasUsed, _ := new(big.Int).SetString(tx.Receipt.CumulativeGas, 10)
		gasFee := new(big.Int).Mul(gasPrice, gasUsed)
		withdrawed, _ := new(big.Int).SetString(zil10, 10)
		result := new(big.Int).Sub(before, gasFee)
		result = result.Add(result, withdrawed)
		result = result.Sub(result, after)
		println(fmt.Sprintf("====%s====", result))
		t.LogDebug()
		//t.AssertEqual(balance1_before, fmt.Sprintf("%s", needed))
	*/
	gasPrice, _ := new(big.Int).SetString(tx.GasPrice, 10)
	gasUsed, _ := new(big.Int).SetString(tx.Receipt.CumulativeGas, 10)
	gasFee := new(big.Int).Mul(gasPrice, gasUsed)
	return gasFee
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

func calcOutputSwapExactZilForTokens(amntIn string) string {
	//1) substract treasury fee
	amountIn, _ := new(big.Int).SetString(amntIn, 10)
	treasuryFeeCalculated, _ := new(big.Int).SetString(calcTreasuryReward(amntIn), 10)
	afterTreasuryFee := new(big.Int).Sub(amountIn, treasuryFeeCalculated)

	//2) constant product formula
	//see let output_for
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

/*
let output_for =
  fun (input_amount: Uint128) =>
  fun (input_reserve: Uint128) =>
  fun (output_reserve: Uint128) =>
  fun (fee: Uint256) =>
    let exact_amount_u256 = grow_u128 input_amount in
    let input_reserve_u256 = grow_u128 input_reserve in
    let output_reserve_u256 = grow_u128 output_reserve in
    let exact_amount_after_fee = builtin mul exact_amount_u256 fee in
    let numerator = builtin mul exact_amount_after_fee output_reserve_u256 in
    let input_reserve_after_fee = builtin mul input_reserve_u256 fee_demon in
    let denominator = builtin add input_reserve_after_fee exact_amount_after_fee in
      builtin div numerator denominator
*/
func calcOutputSwapExactTokensForZil(amntIn string) (string, string) {
	//1) constant product formula
	//input_amount: exactTokens
	//input_reserve: total tokens
	//output reserve: total zil
	exactAmount, _ := new(big.Int).SetString(amntIn, 10)
	zilReserve, _ := new(big.Int).SetString(Field(Aswap, "pools", Stzil.Addr, "arguments", "0"), 10)
	tokenReserve, _ := new(big.Int).SetString(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), 10)
	lpFee, _ := new(big.Int).SetString(Field(Aswap, "liquidity_fee"), 10)
	lpFeeDenom, _ := new(big.Int).SetString("10000", 10)
	afterLPFee := new(big.Int).Mul(exactAmount, lpFee)
	nominator := new(big.Int).Mul(afterLPFee, zilReserve)

	denom := new(big.Int).Mul(lpFeeDenom, tokenReserve)
	denom = denom.Add(denom, afterLPFee)

	amountOut := new(big.Int).Div(nominator, denom)

	//2) substract treasury fee, SwapOptions, TokenToZil =>...
	//treasury_rewards = get_amount_with_fee trea_fee output_zil_amount_u128;
	//output_zil_amount_u128 = builtin sub output_zil_amount_u128 treasury_rewards;
	treasuryFeeCalculated, _ := new(big.Int).SetString(calcTreasuryReward(amountOut.String()), 10)
	afterTreasuryFee := new(big.Int).Sub(amountOut, treasuryFeeCalculated)
	return afterTreasuryFee.String(), treasuryFeeCalculated.String()
}

func calcInputSwapExactZilForTokens(amntOut string) string {
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
	if treasuryFee.Cmp(big.NewInt(0)) == 0 {
		return afterTreasuryFee.Add(afterTreasuryFee, big.NewInt(1)).String()
	}
	afterTreasuryFee = afterTreasuryFee.Add(afterTreasuryFee, big.NewInt(1))
	beforeTreasuryFee := new(big.Int).Mul(afterTreasuryFee, treasuryFee)
	denom = new(big.Int).Sub(treasuryFee, big.NewInt(1))
	beforeTreasuryFee = beforeTreasuryFee.Div(beforeTreasuryFee, denom)

	return beforeTreasuryFee.String()
}

func calcTreasuryReward(amount string) string {
	treasury_fee := Field(Aswap, "treasury_fee")
	if treasury_fee == "0" {
		return "0"
	}
	biAmt, _ := new(big.Int).SetString(amount, 10)
	biFee, _ := new(big.Int).SetString(treasury_fee, 10)
	feeRes := new(big.Int).Div(biAmt, biFee)
	return feeRes.String()
}

type UserBalance struct {
	Zil        *big.Int
	ZilDelta   *big.Int
	Stzil      *big.Int
	StzilDelta *big.Int
	TxFee      *big.Int
	TxFeeTotal *big.Int
}

type BalanceRecord struct {
	TxTag      string
	ASwapZil   *big.Int
	ASwapStzil *big.Int
	StzilZil   *big.Int
	StzilStzil *big.Int
	Balances   []UserBalance
}
type BalanceSheet []BalanceRecord

func recapBalance() string {
	res := fmt.Sprintf("%22s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|\n",
		"Tag", "AZil", "AStzil", "SZil", "SStzil",
		"1-ΔZil", "1-Zil", "1Δ-ST", "1-ST", "1-ΔFee", "1-Fee",
		"2-ΔZil", "2-Zil", "2Δ-ST", "2-ST", "2-ΔFee", "2-Fee",
		"3-ΔZil", "3-Zil", "3Δ-ST", "3-ST", "3-ΔFee", "3-Fee",
	)
	for _, record := range Archive {
		res += fmt.Sprintf("%22s|%6s|%6s|%6s|%6s|", record.TxTag,
			FromQA(record.ASwapZil.String()),
			FromQA(record.ASwapStzil.String()),
			FromQA(record.StzilZil.String()),
			FromQA(record.StzilStzil.String()),
		)
		for i := 1; i <= 3; i++ {
			res += fmt.Sprintf("%6s|%6s|%6s|%6s|%6s|%6s|",
				FromQA(record.Balances[i].ZilDelta.String()),
				FromQA(record.Balances[i].Zil.String()),
				FromQA(record.Balances[i].StzilDelta.String()),
				FromQA(record.Balances[i].Stzil.String()),
				FromQA(record.Balances[i].TxFee.String()),
				FromQA(record.Balances[i].TxFeeTotal.String()),
			)
		}
		res += "\n"
	}
	return res
}

func recordBalance(userId int, tx *transaction.Transaction) {

	//data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
	//GetLog().Info(string(data))

	var userBalance UserBalance
	tag := ""
	if tx != nil && reflect.TypeOf(tx.Data).String() != "string" {
		tag = tx.Data.(contract.Data).Tag
	}
	var aswapZil, aswapStzil *big.Int
	if Field(Aswap, "pools", Stzil.Addr) == "" {
		aswapZil = big.NewInt(0)
		aswapStzil = big.NewInt(0)
	} else {
		aswapZil, _ = new(big.Int).SetString(Field(Aswap, "pools", Stzil.Addr, "arguments", "0"), 10)
		aswapStzil, _ = new(big.Int).SetString(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), 10)
	}

	stzilZil, _ := new(big.Int).SetString(Field(Stzil, "totalstakeamount"), 10)
	stzilStzil, _ := new(big.Int).SetString(Field(Stzil, "total_supply"), 10)
	balanceRecord := BalanceRecord{
		TxTag:      tag,
		ASwapZil:   aswapZil,
		ASwapStzil: aswapStzil,
		StzilZil:   stzilZil,
		StzilStzil: stzilStzil,
		Balances:   make([]UserBalance, 4),
	}

	alen := len(Archive)
	for i := 1; i <= 3; i++ {
		zilPrev := big.NewInt(0)
		stzilPrev := big.NewInt(0)
		txFeeTotalPrev := big.NewInt(0)
		if alen != 0 {
			zilPrev = Archive[alen-1].Balances[i].Zil
			stzilPrev = Archive[alen-1].Balances[i].Stzil
			txFeeTotalPrev = Archive[alen-1].Balances[i].TxFeeTotal
		}
		//fee
		txFee := big.NewInt(0)
		if tx != nil && i == userId {
			txFee = getTxFee(tx)
		}
		txFeeTotal := big.NewInt(0).Add(txFeeTotalPrev, txFee)

		zil, _ := new(big.Int).SetString(sdk.GetBalance(Addresses[i]), 10)
		stzil := Stzil.BalanceOf(Addresses[i])
		zilDelta := new(big.Int).Sub(zil, zilPrev)
		stzilDelta := new(big.Int).Sub(stzil, stzilPrev)
		userBalance = UserBalance{
			Zil:        zil,
			ZilDelta:   zilDelta,
			Stzil:      stzil,
			StzilDelta: stzilDelta,
			TxFee:      txFee,
			TxFeeTotal: txFeeTotal,
		}
		balanceRecord.Balances[i] = userBalance

	}
	Archive = append(Archive, balanceRecord)
}

func aswapOwnerOnly(tr *Transitions) {

	Start("aswapOwnerOnly")

	init_owner_addr := utils.GetAddressByWallet(celestials.Admin)
	aswap := tr.DeployASwap(init_owner_addr)
	// Use non-owner user for Aswap, expecting errors
	aswap.SetSigner(bob)

	tx, _ := aswap.SetLiquidityFee("12345")
	AssertASwapError(tx, aswap.ErrorCode("CodeNotContractOwner"))

	tx, _ = aswap.SetTreasuryFee("12345")
	AssertASwapError(tx, aswap.ErrorCode("CodeNotContractOwner"))

	tx, _ = aswap.SetTreasuryAddress(core.ZeroAddr)
	AssertASwapError(tx, aswap.ErrorCode("CodeNotContractOwner"))

	tx, _ = aswap.TogglePause()
	AssertASwapError(tx, aswap.ErrorCode("CodeNotContractOwner"))

	tx, _ = aswap.ChangeOwner(core.ZeroAddr)
	AssertASwapError(tx, aswap.ErrorCode("CodeNotContractOwner"))
}
