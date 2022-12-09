package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

const swapOutput = "9871580343970"

func (tr *Transitions) Fungible() {
	tr.FungibleAllowanceErrors()
	tr.AddToSwap()
	tr.TransferFrom()
	tr.Transfer()
}

func (tr *Transitions) FungibleAllowanceErrors() {
	Start("Test ZRC-2 errors")

	p := tr.DeployAndUpgrade()

	amount := ToQA(1000)
	AssertSuccess(p.StZIL.WithUser(sdk.Cfg.Key1).DelegateStake(amount))

	tx, _ := p.StZIL.IncreaseAllowance(sdk.Cfg.Addr1, amount)
	AssertError(tx, p.StZIL.ErrorCode("CodeIsSender"))

	tx, _ = p.StZIL.DecreaseAllowance(sdk.Cfg.Addr1, amount)
	AssertError(tx, p.StZIL.ErrorCode("CodeIsSender"))
}

func (tr *Transitions) AddToSwap() {
	Start("Swap via ZilSwap")

	p := tr.DeployAndUpgrade()
	zilSwap := tr.DeployZilSwap()
	stzil := p.StZIL

	liquidityAmount := ToQA(1000)

	AssertSuccess(stzil.DelegateStake(liquidityAmount))

	blockNum := p.GetBlockHeight()

	// Add AddLiquidity
	AssertSuccess(stzil.IncreaseAllowance(zilSwap.Contract.Addr, ToQA(10000)))
	AssertSuccess(zilSwap.AddLiquidity(stzil.Contract.Addr, liquidityAmount, liquidityAmount, blockNum))

	// Do Swap
	recipient := sdk.Cfg.Addr1
	AssertSuccess(zilSwap.SwapExactZILForTokens(stzil.Contract.Addr, ToQA(10), "1", recipient, blockNum))
	AssertEqual(stzil.BalanceOf(recipient).String(), swapOutput)
}

func (tr *Transitions) Transfer() {
	Start("Transfer")

	p := tr.DeployAndUpgrade()
	stzil := p.StZIL

	from := sdk.Cfg.Addr1
	to := sdk.Cfg.Addr2
	amount := ToQA(100)

	tx, _ := stzil.WithUser(sdk.Cfg.Key1).Transfer(to, amount)
	AssertError(tx, p.StZIL.ErrorCode("CodeInsufficientFunds"))

	AssertSuccess(stzil.WithUser(sdk.Cfg.Key1).DelegateStake(amount))

	AssertEqual(stzil.BalanceOf(from).String(), amount)
	AssertEqual(stzil.BalanceOf(to).String(), ToQA(0))

	AssertSuccess(stzil.WithUser(sdk.Cfg.Key1).Transfer(to, amount))

	AssertEqual(stzil.BalanceOf(from).String(), ToQA(0))
	AssertEqual(stzil.BalanceOf(to).String(), amount)
}

func (tr *Transitions) TransferFrom() {
	Start("TransferFrom")

	p := tr.DeployAndUpgrade()
	stzil := p.StZIL

	from := sdk.Cfg.Addr1
	to := sdk.Cfg.Addr2

	amount := ToQA(100)

	AssertSuccess(stzil.WithUser(sdk.Cfg.Key1).DelegateStake(amount))

	AssertEqual(stzil.BalanceOf(from).String(), amount)
	AssertEqual(stzil.BalanceOf(to).String(), ToQA(0))

	tx, _ := stzil.TransferFrom(from, to, amount)
	AssertError(tx, p.StZIL.ErrorCode("CodeInsufficientAllowance"))

	// Allow admin user to spend User1 money
	admin := celestials.Admin
	AssertSuccess(stzil.WithUser(sdk.Cfg.Key1).IncreaseAllowance(utils.GetAddressByWallet(admin), amount))

	stzil.SetSigner(admin)
	AssertSuccess(stzil.TransferFrom(from, to, amount))

	AssertEqual(stzil.BalanceOf(from).String(), ToQA(0))
	AssertEqual(stzil.BalanceOf(to).String(), amount)
}
