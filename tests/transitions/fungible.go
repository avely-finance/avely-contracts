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
	p.StZIL.SetSigner(alice)
	recipient := utils.GetAddressByWallet(alice)

	amount := ToQA(1000)
	AssertSuccess(p.StZIL.DelegateStake(amount))

	tx, _ := p.StZIL.IncreaseAllowance(recipient, amount)
	AssertError(tx, p.StZIL.ErrorCode("CodeIsSender"))

	tx, _ = p.StZIL.DecreaseAllowance(recipient, amount)
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
	recipient := utils.GetAddressByWallet(alice)
	AssertSuccess(zilSwap.SwapExactZILForTokens(stzil.Contract.Addr, ToQA(10), "1", recipient, blockNum))
	AssertEqual(stzil.BalanceOf(recipient).String(), swapOutput)
}

func (tr *Transitions) Transfer() {
	Start("Transfer")

	p := tr.DeployAndUpgrade()
	stzil := p.StZIL

	from := utils.GetAddressByWallet(alice)
	to := utils.GetAddressByWallet(bob)
	amount := ToQA(100)

	stzil.SetSigner(alice)
	tx, _ := stzil.Transfer(to, amount)
	AssertError(tx, p.StZIL.ErrorCode("CodeInsufficientFunds"))

	AssertSuccess(stzil.DelegateStake(amount))

	AssertEqual(stzil.BalanceOf(from).String(), amount)
	AssertEqual(stzil.BalanceOf(to).String(), ToQA(0))

	AssertSuccess(stzil.Transfer(to, amount))

	AssertEqual(stzil.BalanceOf(from).String(), ToQA(0))
	AssertEqual(stzil.BalanceOf(to).String(), amount)
}

func (tr *Transitions) TransferFrom() {
	Start("TransferFrom")

	p := tr.DeployAndUpgrade()
	stzil := p.StZIL

	from := utils.GetAddressByWallet(alice)
	to := utils.GetAddressByWallet(bob)

	amount := ToQA(100)

	stzil.SetSigner(alice)
	AssertSuccess(stzil.DelegateStake(amount))

	AssertEqual(stzil.BalanceOf(from).String(), amount)
	AssertEqual(stzil.BalanceOf(to).String(), ToQA(0))

	tx, _ := stzil.TransferFrom(from, to, amount)
	AssertError(tx, p.StZIL.ErrorCode("CodeInsufficientAllowance"))

	// Allow admin user to spend User1 money
	admin := celestials.Admin
	AssertSuccess(stzil.IncreaseAllowance(utils.GetAddressByWallet(admin), amount))

	stzil.SetSigner(admin)
	AssertSuccess(stzil.TransferFrom(from, to, amount))

	AssertEqual(stzil.BalanceOf(from).String(), ToQA(0))
	AssertEqual(stzil.BalanceOf(to).String(), amount)
}
