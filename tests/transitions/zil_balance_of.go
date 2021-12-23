package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) ZilBalanceOf() {
	t.Start("ZilBalanceOf")

	// deploy smart contract
	p := tr.DeployAndUpgrade()

	/*******************************************************************************
	 * 1. Non-delegator address (sdk.Cfg.Addr2) should have empty balance
	 *******************************************************************************/
	t.Start("ZilBalanceOf, step 1")
	p.Aimpl.UpdateWallet(sdk.Cfg.Key2)
	balance2, _ := p.Aimpl.ZilBalanceOf(sdk.Cfg.Addr2)
	t.AssertEqual(balance2, Zil(0))

	/*******************************************************************************
	 * 2. After delegate (sdk.Cfg.Addr2) should have its balance updated
	 *******************************************************************************/
	t.Start("ZilBalanceOf, step 2")
	t.AssertSuccess(p.Aimpl.DelegateStake(Zil(15)))
	balance2, _ = p.Aimpl.ZilBalanceOf(sdk.Cfg.Addr2)
	t.AssertEqual(balance2, Zil(15))

	/*******************************************************************************
	 * 3. User balance in zil should be updated after restaking rewards
	 * because contract got more zils w/o azil minting
	 * so azil/zil exchange rate changed, azil now costs more zils than before
	 * so balance of addr2 in zil should be more
	 *******************************************************************************/
	t.Start("ZilBalanceOf, step 3")
	p.Aimpl.UpdateWallet(sdk.Cfg.AdminKey)
	t.AssertSuccess(p.Aimpl.IncreaseAutoRestakeAmount(Zil(10)))
	t.AssertSuccess(p.Aimpl.PerformAutoRestake())
	balance2, _ = p.Aimpl.ZilBalanceOf(sdk.Cfg.Addr2)
	t.AssertEqual(balance2, StrMulDiv(
		p.Aimpl.Field("balances", "0x"+sdk.Cfg.Addr2),
		p.Aimpl.Field("totalstakeamount"),
		p.Aimpl.Field("totaltokenamount")))

	/*******************************************************************************
	 * 4. New user (sdk.Cfg.Addr3) delegating stake
	 * He'll get azils by new azil/zil rate, so zilBalanceOf should be equal to the delegated zils amount.
	 * First user's (sdk.Cfg.Addr2) zil balance should not be changed.
	 *******************************************************************************/
	t.Start("ZilBalanceOf, step 4")
	p.Aimpl.UpdateWallet(sdk.Cfg.Key3)
	t.AssertSuccess(p.Aimpl.DelegateStake(Zil(10)))
	balance2, _ = p.Aimpl.ZilBalanceOf(sdk.Cfg.Addr2)
	t.AssertEqual(balance2, StrMulDiv(
		p.Aimpl.Field("balances", "0x"+sdk.Cfg.Addr2),
		p.Aimpl.Field("totalstakeamount"),
		p.Aimpl.Field("totaltokenamount")))
	balance3, _ := p.Aimpl.ZilBalanceOf(sdk.Cfg.Addr3)
	t.AssertEqual(balance3, StrMulDiv(
		p.Aimpl.Field("balances", "0x"+sdk.Cfg.Addr3),
		p.Aimpl.Field("totalstakeamount"),
		p.Aimpl.Field("totaltokenamount")))
}
