package transitions

import (
	. "Azil/test/helpers"
)

func (tr *Transitions) ZilBalanceOf() {

	t.Start("ZilBalanceOf")

	// deploy smart contract
	_, _, Aimpl, _, _ := tr.DeployAndUpgrade()

	/*******************************************************************************
	 * 1. Non-delegator address (tr.cfg.Addr2) should have empty balance
	 *******************************************************************************/
	t.Start("ZilBalanceOf, step 1")
	Aimpl.UpdateWallet(tr.cfg.Key2)
	balance2, _ := Aimpl.ZilBalanceOf(tr.cfg.Addr2)
	t.AssertEqual(balance2, zil(0))

	/*******************************************************************************
	 * 2. After delegate (tr.cfg.Addr2) should have its balance updated
	 *******************************************************************************/
	t.Start("ZilBalanceOf, step 2")
	t.AssertSuccess(Aimpl.DelegateStake(zil(15)))
	balance2, _ = Aimpl.ZilBalanceOf(tr.cfg.Addr2)
	t.AssertEqual(balance2, zil(15))

	/*******************************************************************************
	 * 3. User balance in zil should be updated after restaking rewards
	 * because contract got more zils w/o azil minting
	 * so azil/zil exchange rate changed, azil now costs more zils than before
	 * so balance of addr2 in zil should be more
	 *******************************************************************************/
	t.Start("ZilBalanceOf, step 3")
	Aimpl.UpdateWallet(tr.cfg.AdminKey)
	t.AssertSuccess(Aimpl.IncreaseAutoRestakeAmount(zil(10)))
	t.AssertSuccess(Aimpl.PerformAutoRestake())
	balance2, _ = Aimpl.ZilBalanceOf(tr.cfg.Addr2)
	t.AssertEqual(balance2, StrMulDiv(
		Aimpl.Field("balances", "0x"+tr.cfg.Addr2),
		Aimpl.Field("totalstakeamount"),
		Aimpl.Field("totaltokenamount")))

	/*******************************************************************************
	 * 4. New user (tr.cfg.Addr3) delegating stake
	 * He'll get azils by new azil/zil rate, so zilBalanceOf should be equal to the delegated zils amount.
	 * First user's (tr.cfg.Addr2) zil balance should not be changed.
	 *******************************************************************************/
	t.Start("ZilBalanceOf, step 4")
	Aimpl.UpdateWallet(tr.cfg.Key3)
	t.AssertSuccess(Aimpl.DelegateStake(zil(10)))
	balance2, _ = Aimpl.ZilBalanceOf(tr.cfg.Addr2)
	t.AssertEqual(balance2, StrMulDiv(
		Aimpl.Field("balances", "0x"+tr.cfg.Addr2),
		Aimpl.Field("totalstakeamount"),
		Aimpl.Field("totaltokenamount")))
	balance3, _ := Aimpl.ZilBalanceOf(tr.cfg.Addr3)
	t.AssertEqual(balance3, StrMulDiv(
		Aimpl.Field("balances", "0x"+tr.cfg.Addr3),
		Aimpl.Field("totalstakeamount"),
		Aimpl.Field("totaltokenamount")))
}
