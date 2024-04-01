package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

const swapOutput = "9881383789778"

func (tr *Transitions) Fungible() {
	tr.FungibleAllowanceErrors()
	tr.AddToSwap()
	tr.TransferFrom()
	tr.Transfer()

	tr.EvmOn()

	tr.FungibleAllowanceErrors()
	//tr.AddToSwap() // aswap bridge not implemented
	tr.TransferFrom()
	tr.Transfer()
}

func (tr *Transitions) FungibleAllowanceErrors() {
	Start("Test ZRC-2 errors")

	p := tr.DeployAndUpgrade()
	tr.GetStZIL().SetSigner(alice)
	recipient := tr.GetAddressByWallet(alice)

	amount := ToQA(1000)
	AssertSuccessAny(tr.GetStZIL().DelegateStake(amount))

	tx, _ := tr.GetStZIL().IncreaseAllowance(recipient, amount)
	AssertError(tx, p.StZIL.ErrorCode("CodeIsSender"))

	tx, _ = tr.GetStZIL().DecreaseAllowance(recipient, amount)
	AssertError(tx, p.StZIL.ErrorCode("CodeIsSender"))
}

func (tr *Transitions) AddToSwap() {
	Start("Swap via aswap")

	p := tr.DeployAndUpgrade()
	init_owner_addr := utils.GetAddressByWallet(celestials.Admin)
	aswap := tr.DeployASwap(init_owner_addr)
	stzil := p.StZIL

	liquidityAmount := ToQA(1000)

	AssertSuccess(stzil.DelegateStake(liquidityAmount))

	blockNum := p.GetBlockHeight()

	// Add AddLiquidity
	AssertSuccess(stzil.IncreaseAllowance(aswap.Contract.Addr, ToQA(10000)))
	AssertSuccess(aswap.AddLiquidity(liquidityAmount, stzil.Contract.Addr, "0", liquidityAmount, blockNum))

	// Do Swap
	recipient := utils.GetAddressByWallet(alice)
	AssertSuccess(aswap.SwapExactZILForTokens(ToQA(10), stzil.Contract.Addr, "1", recipient, blockNum))
	AssertEqual(stzil.BalanceOf(recipient).String(), swapOutput)
}

func (tr *Transitions) Transfer() {
	Start("Transfer")

	p := tr.DeployAndUpgrade()
	stzil := p.StZIL

	from := tr.GetAddressByWallet(alice)
	to := tr.GetAddressByWallet(bob)
	amount := ToQA(100)

	tr.GetStZIL().SetSigner(alice)
	tx, _ := tr.GetStZIL().Transfer(to, amount)
	AssertError(tx, p.StZIL.ErrorCode("CodeInsufficientFunds"))

	AssertSuccessAny(tr.GetStZIL().DelegateStake(amount))

	AssertEqual(stzil.BalanceOf(from).String(), amount)
	AssertEqual(stzil.BalanceOf(to).String(), ToQA(0))

	AssertSuccessAny(tr.GetStZIL().Transfer(to, amount))

	AssertEqual(stzil.BalanceOf(from).String(), ToQA(0))
	AssertEqual(stzil.BalanceOf(to).String(), amount)
}

func (tr *Transitions) TransferFrom() {
	Start("TransferFrom")

	p := tr.DeployAndUpgrade()
	stzil := p.StZIL

	from := tr.GetAddressByWallet(alice)
	to := tr.GetAddressByWallet(bob)

	amount := ToQA(100)

	tr.GetStZIL().SetSigner(alice)
	AssertSuccessAny(tr.GetStZIL().DelegateStake(amount))

	AssertEqual(stzil.BalanceOf(from).String(), amount)
	AssertEqual(stzil.BalanceOf(to).String(), ToQA(0))

	tx, _ := tr.GetStZIL().TransferFrom(from, to, amount)
	AssertError(tx, p.StZIL.ErrorCode("CodeInsufficientAllowance"))

	// Allow admin user to spend User1 money
	admin := celestials.Admin
	AssertSuccessAny(tr.GetStZIL().IncreaseAllowance(tr.GetAddressByWallet(admin), amount))

	tr.GetStZIL().SetSigner(admin)
	AssertSuccessAny(tr.GetStZIL().TransferFrom(from, to, amount))

	AssertEqual(stzil.BalanceOf(from).String(), ToQA(0))
	AssertEqual(stzil.BalanceOf(to).String(), amount)
}
