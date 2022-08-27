package transitions

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/contract"
	core2 "github.com/Zilliqa/gozilliqa-sdk/core"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
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
	expectedTreasuryRewards := "133333333333"
	expectedSwapOutput := "9887919312466"
	AssertSuccess(aswap.WithUser(init_owner_key).TogglePause())
	tx, _ := AssertSuccess(aswap.SwapExactZILForTokens(ToQA(100), stzil.Contract.Addr, "90", sdk.Cfg.Addr1, blockNum))
	AssertTransition(tx, Transition{
		aswap.Addr, //sender
		"AddFunds",
		treasury_address,
		expectedTreasuryRewards,
		ParamsMap{},
	})
	AssertEqual(stzil.BalanceOf(sdk.Cfg.Addr1).String(), expectedSwapOutput)

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

	addFunds(Addr1, ToQA(5000))
	addFunds(Addr2, ToQA(5000))
	addFunds(Addr3, ToQA(5000))

	return Proto, Aswap, Stzil
}

func (tr *Transitions) ASwapGolden() {

	//TODO: to make this test suite completely correct, we need to maintain pools/balances/contributions here

	Start("ASwap Golden Flow")

	tr.setupGoldenFlow()
	blockNum := Proto.GetBlockHeight()

	//mint some stzil and transfer it to user 1
	liqPair := LiquidityPair{zil: ToQA(1000), stzil: ToQA(1000)}
	AssertSuccess(Stzil.WithUser(sdk.Cfg.AdminKey).DelegateStake(liqPair.zil))
	AssertSuccess(Stzil.WithUser(sdk.Cfg.AdminKey).Transfer(Addr1, liqPair.stzil))

	recordBalance(-1, nil)

	//1) user 1 add liquidity 1000:1000
	user1Balance, _ := new(big.Int).SetString(sdk.GetBalance(Addr1), 10)
	tx, _ := AssertSuccess(Stzil.WithUser(Key1).IncreaseAllowance(Aswap.Contract.Addr, liqPair.stzil))
	user1TxFee := getTxFee(tx)
	recordBalance(1, tx)
	//allowances[_sender][spender] := new_allowance;
	AssertEqual(Field(Stzil, "allowances", Addr1, Aswap.Addr), liqPair.stzil)
	//func (a *Aswap) AddLiquidity(_amount, tokenAddr, minContributionAmount, tokenAmount string, blockNum int)
	tx, _ = AssertSuccess(Aswap.WithUser(Key1).AddLiquidity(liqPair.zil, Stzil.Contract.Addr, liqPair.zil, liqPair.stzil, blockNum))
	user1TxFee = user1TxFee.Add(user1TxFee, getTxFee(tx))
	recordBalance(1, tx)
	//pools[Aswap.addr] := [zil, stzil];
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "0"), liqPair.zil)
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), liqPair.stzil)
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), ToQA(1000))
	AssertEqual(Stzil.BalanceOf(Addr1).String(), "0")
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), ToQA(1000))

	//2) user 2 swap 500 stzil for zil
	user2Balance, _ := new(big.Int).SetString(sdk.GetBalance(Addr2), 10)
	amountIn := calcSwapInputZilForTokens(ToQA(500))
	expectedSwapOutput := calcSwapOutputZilForTokens(amountIn)
	minTokenAmount := "1"
	tx, _ = AssertSuccess(Aswap.WithUser(Key2).SwapExactZILForTokens(amountIn, Stzil.Contract.Addr, minTokenAmount, Addr2, blockNum+1))
	user2TxFee := getTxFee(tx)
	recordBalance(2, tx)
	AssertTransition(tx, Transition{
		Aswap.Addr, //sender
		"AddFunds",
		Field(Aswap, "treasury_address"),
		calcTreasuryReward(amountIn),
		ParamsMap{},
	})
	AssertEqual(Stzil.BalanceOf(Addr2).String(), ToQA(500))
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), ToQA(500))
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), ToQA(500))
	user2StzilBalanceStr := expectedSwapOutput

	//3a) user 3 swap 250 stzil for zil
	user3Balance, _ := new(big.Int).SetString(sdk.GetBalance(Addr3), 10)
	amountIn = calcSwapInputZilForTokens(ToQA(250))
	expectedSwapOutput = calcSwapOutputZilForTokens(amountIn)
	minTokenAmount = "1"
	tx, _ = AssertSuccess(Aswap.WithUser(sdk.Cfg.Key3).SwapExactZILForTokens(amountIn, Stzil.Contract.Addr, minTokenAmount, Addr3, blockNum+1))
	user3TxFee := getTxFee(tx)
	recordBalance(3, tx)
	AssertTransition(tx, Transition{
		Aswap.Addr, //sender
		"AddFunds",
		Field(Aswap, "treasury_address"),
		calcTreasuryReward(amountIn),
		ParamsMap{},
	})
	AssertEqual(Stzil.BalanceOf(Addr3).String(), ToQA(250))
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), ToQA(250))
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), ToQA(250))

	//3b) user 3 add 250 stzil to pool
	pairIn := calcAddLiquidity(LiquidityPair{stzil: ToQA(250)})
	tx, _ = AssertSuccess(Stzil.WithUser(Key3).IncreaseAllowance(Aswap.Contract.Addr, pairIn.stzil))
	user3TxFee = user3TxFee.Add(user3TxFee, getTxFee(tx))
	recordBalance(3, tx)
	//allowances[_sender][spender] := new_allowance;
	AssertEqual(Field(Stzil, "allowances", Addr3, Aswap.Addr), pairIn.stzil)
	//func (a *Aswap) AddLiquidity(_amount, tokenAddr, minContributionAmount, tokenAmount string, blockNum int)
	tx, _ = AssertSuccess(Aswap.WithUser(Key3).AddLiquidity(pairIn.zil, Stzil.Addr, "1", pairIn.stzil, blockNum+1))
	user3TxFee = user3TxFee.Add(user3TxFee, getTxFee(tx))
	recordBalance(3, tx)
	AssertEqual(Stzil.BalanceOf(Addr3).String(), "0")
	//pools[Aswap.addr] := [zil, stzil];
	AssertEqual(Stzil.BalanceOf(Addr3).String(), "0")
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), ToQA(500))
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), ToQA(500))

	//4) user 2 "buy" 100 stzil at swap
	amountIn = calcSwapInputZilForTokens(ToQA(100))
	expectedSwapOutput = calcSwapOutputZilForTokens(amountIn)
	minTokenAmount = "1"
	tx, _ = AssertSuccess(Aswap.WithUser(Key2).SwapExactZILForTokens(amountIn, Stzil.Contract.Addr, minTokenAmount, Addr2, blockNum+1))
	user2TxFee = user2TxFee.Add(user2TxFee, getTxFee(tx))
	recordBalance(2, tx)
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
	AssertEqual(Stzil.BalanceOf(Addr2).String(), user2StzilBalance.String())
	AssertEqual(Stzil.BalanceOf(Addr2).String(), ToQA(600))
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), ToQA(400))
	AssertEqual(Field(Aswap, "pools", Stzil.Addr, "arguments", "1"), ToQA(400))

	//5) user1 remove liquidity
	pairOut := calcRemoveLiquidityOutput(Addr1, ToQA(1000))
	//RemoveLiquidity(tokenAddress, contributionAmount, minZilAmount, minTokenAmount string, blockNum int) (*transaction.Transaction, error)
	tx, _ = AssertSuccess(Aswap.WithUser(Key1).RemoveLiquidity(Stzil.Addr, ToQA(1000), "1", "1", blockNum+1))
	user1TxFee = user1TxFee.Add(user1TxFee, getTxFee(tx))
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
	AssertEqual(Stzil.BalanceOf(Addr1).String(), pairOut.stzil)
	AssertEqual(Stzil.BalanceOf(Addr1).String(), ToQA(200))

	//6) user3 remove liquidity
	pairOut3 := calcRemoveLiquidityOutput(Addr3, "all")
	user3Contribution := Field(Aswap, "balances", Addr3, Stzil.Addr)
	//RemoveLiquidity(tokenAddress, contributionAmount, minZilAmount, minTokenAmount string, blockNum int) (*transaction.Transaction, error)
	tx, _ = AssertSuccess(Aswap.WithUser(Key3).RemoveLiquidity(Stzil.Addr, user3Contribution, "1", "1", blockNum+1))
	user3TxFee = user3TxFee.Add(user3TxFee, getTxFee(tx))
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
	AssertEqual(Stzil.BalanceOf(Addr1).String(), pairOut3.stzil)
	AssertEqual(Stzil.BalanceOf(Aswap.Addr).String(), "0")

	curBalance1, _ := new(big.Int).SetString(sdk.GetBalance(Addr1), 10)
	profit1 := new(big.Int).Sub(curBalance1, user1Balance)
	profit1 = profit1.Add(profit1, user1TxFee)

	GetLog().Info("curBalance1=" + curBalance1.String())
	GetLog().Info("user1Balance=" + user1Balance.String())
	GetLog().Info("user1TxFee=" + user1TxFee.String())

	//panic(1)

	curBalance2, _ := new(big.Int).SetString(sdk.GetBalance(Addr2), 10)
	profit2 := new(big.Int).Sub(curBalance2, user2Balance)
	profit2 = profit2.Add(profit2, user2TxFee)

	GetLog().Info("curBalance2=" + curBalance2.String())
	GetLog().Info("user2Balance=" + user2Balance.String())
	GetLog().Info("user2TxFee=" + user2TxFee.String())

	curBalance3, _ := new(big.Int).SetString(sdk.GetBalance(Addr3), 10)
	profit3 := new(big.Int).Sub(curBalance3, user3Balance)
	profit3 = profit3.Add(profit3, user3TxFee)

	GetLog().Info("curBalance3=" + curBalance3.String())
	GetLog().Info("user3Balance=" + user3Balance.String())
	GetLog().Info("user3TxFee=" + user3TxFee.String())

	fmt.Println(recapBalance())
}

func addFunds(recipient, amount string) (*transaction.Transaction, error) {
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(sdk.Cfg.AdminKey)
	provider := provider2.NewProvider(sdk.Cfg.Api.HttpUrl)

	gasPrice, _ := provider.GetMinimumGasPrice()

	if recipient[0:2] == "0x" {
		recipient = recipient[2:]
	}

	b32, _ := bech32.ToBech32Address("0x" + recipient)

	tx := &transaction.Transaction{
		Version:      strconv.FormatInt(int64(util.Pack(sdk.Cfg.ChainId, 1)), 10),
		SenderPubKey: "",
		ToAddr:       b32,
		Amount:       amount,
		GasPrice:     gasPrice,
		GasLimit:     "40000",
		Code:         "",
		Data:         "",
		Priority:     false,
		Nonce:        "",
	}
	wallet.Sign(tx, *provider)
	rsp, _ := provider.CreateTransaction(tx.ToTransactionPayload())
	resMap := rsp.Result.(map[string]interface{})
	hash := resMap["TranID"].(string)
	//fmt.Printf("hash is %s\n", hash)
	tx.Confirm(hash, 1000, 0, provider)
	if tx.Status == core2.Confirmed {
		return tx, nil
	}
	return nil, errors.New("Can't confirm transaction")
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
	Balances   []UserBalance
}
type BalanceSheet []BalanceRecord

func recapBalance() string {
	res := fmt.Sprintf("%22s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|%6s|\n",
		"Tag", "AZil", "AStzil",
		"1-ΔZil", "1-Zil", "1Δ-ST", "1-ST", "1-ΔFee", "1-Fee",
		"2-ΔZil", "2-Zil", "2Δ-ST", "2-ST", "2-ΔFee", "2-Fee",
		"3-ΔZil", "3-Zil", "3Δ-ST", "3-ST", "3-ΔFee", "3-Fee",
	)
	for _, record := range Archive {
		res += fmt.Sprintf("%22s|%6s|%6s|", record.TxTag,
			FromQA(record.ASwapZil.String()),
			FromQA(record.ASwapStzil.String()),
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

	balanceRecord := BalanceRecord{
		TxTag:      tag,
		ASwapZil:   aswapZil,
		ASwapStzil: aswapStzil,
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
