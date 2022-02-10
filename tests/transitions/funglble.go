package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

const swapOutput = "9871580343970"

func (tr *Transitions) AddToSwap() {
	Start("Swap via ZilSwap")

	p := tr.DeployAndUpgrade()
	zilSwap := tr.DeployZilSwap()
	azil := p.Aimpl

	liquidityAmount := ToQA(1000)

	AssertSuccess(azil.DelegateStake(liquidityAmount))

	blockNum := p.GetBlockHeight()

	// Add AddLiquidity
	AssertSuccess(azil.IncreaseAllowance(zilSwap.Contract.Addr, ToQA(10000)))
	AssertSuccess(zilSwap.AddLiquidity(azil.Contract.Addr, liquidityAmount, liquidityAmount, blockNum))

	// Do Swap
	recipient := sdk.Cfg.Addr1
	AssertSuccess(zilSwap.SwapExactZILForTokens(azil.Contract.Addr, ToQA(10), "1", recipient, blockNum))
	AssertEqual(azil.BalanceOf(recipient).String(), swapOutput)
}

func (tr *Transitions) Transfer() {
	Start("Transfer")

	p := tr.DeployAndUpgrade()
	azil := p.Aimpl

	from := sdk.Cfg.Addr1
	to := sdk.Cfg.Addr2
	amount := ToQA(100)

	tx, _ := azil.WithUser(sdk.Cfg.Key1).Transfer(to, amount)
	AssertError(tx, "InsufficientFunds")

	AssertSuccess(azil.WithUser(sdk.Cfg.Key1).DelegateStake(amount))

	AssertEqual(azil.BalanceOf(from).String(), amount)
	AssertEqual(azil.BalanceOf(to).String(), ToQA(0))

	AssertSuccess(azil.WithUser(sdk.Cfg.Key1).Transfer(to, amount))

	AssertEqual(azil.BalanceOf(from).String(), ToQA(0))
	AssertEqual(azil.BalanceOf(to).String(), amount)
}

func (tr *Transitions) TransferFrom() {
	Start("TransferFrom")

	p := tr.DeployAndUpgrade()
	azil := p.Aimpl

	from := sdk.Cfg.Addr1
	to := sdk.Cfg.Addr2

	amount := ToQA(100)

	AssertSuccess(azil.WithUser(sdk.Cfg.Key1).DelegateStake(amount))

	AssertEqual(azil.BalanceOf(from).String(), amount)
	AssertEqual(azil.BalanceOf(to).String(), ToQA(0))

	tx, _ := azil.TransferFrom(from, to, amount)
	AssertError(tx, "InsufficientAllowance")

	// Allow admin user to spend User1 money
	AssertSuccess(azil.WithUser(sdk.Cfg.Key1).IncreaseAllowance(sdk.Cfg.Admin, amount))
	AssertSuccess(azil.WithUser(sdk.Cfg.AdminKey).TransferFrom(from, to, amount))

	AssertEqual(azil.BalanceOf(from).String(), ToQA(0))
	AssertEqual(azil.BalanceOf(to).String(), amount)
}
