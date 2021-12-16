package transitions

import (
	"Azil/test/deploy"
)

func (t *Testing) ZilBalanceOf() {

	t.LogStart("ZilBalanceOf")

	// deploy smart contract
	_, _, Aimpl, _, _ := t.DeployAndUpgrade()

	/*******************************************************************************
	 * 1. Non-delegator address (addr2) should have empty balance
	 *******************************************************************************/
	t.LogStart("ZilBalanceOf, step 1")
	Aimpl.UpdateWallet(key2)
	balance2, _ := Aimpl.ZilBalanceOf(addr2)
	t.AssertEqual(balance2, zil(0))

	/*******************************************************************************
	 * 2. After delegate (addr2) should have its balance updated
	 *******************************************************************************/
	t.LogStart("ZilBalanceOf, step 2")
	t.AssertSuccess(Aimpl.DelegateStake(zil(15)))
	balance2, _ = Aimpl.ZilBalanceOf(addr2)
	t.AssertEqual(balance2, zil(15))

	/*******************************************************************************
	 * 3. User balance in zil should be updated after restaking rewards
	 * because contract got more zils w/o azil minting
	 * so azil/zil exchange rate changed, azil now costs more zils than before
	 * so balance of addr2 in zil should be more
	 *******************************************************************************/
	t.LogStart("ZilBalanceOf, step 3")
	Aimpl.UpdateWallet(adminKey)
	t.AssertSuccess(Aimpl.IncreaseAutoRestakeAmount(zil(10)))
	t.AssertSuccess(Aimpl.PerformAutoRestake())
	balance2, _ = Aimpl.ZilBalanceOf(addr2)
	t.AssertEqual(balance2, deploy.MulDiv(
		Aimpl.Field("balances", "0x"+addr2),
		Aimpl.Field("totalstakeamount"),
		Aimpl.Field("totaltokenamount")))

	/*******************************************************************************
	 * 4. New user (addr3) delegating stake
	 * He'll get azils by new azil/zil rate, so zilBalanceOf should be equal to the delegated zils amount.
	 * First user's (addr2) zil balance should not be changed.
	 *******************************************************************************/
	t.LogStart("ZilBalanceOf, step 4")
	Aimpl.UpdateWallet(key3)
	t.AssertSuccess(Aimpl.DelegateStake(zil(10)))
	balance2, _ = Aimpl.ZilBalanceOf(addr2)
	t.AssertEqual(balance2, deploy.MulDiv(
		Aimpl.Field("balances", "0x"+addr2),
		Aimpl.Field("totalstakeamount"),
		Aimpl.Field("totaltokenamount")))
	balance3, _ := Aimpl.ZilBalanceOf(addr3)
	t.AssertEqual(balance3, deploy.MulDiv(
		Aimpl.Field("balances", "0x"+addr3),
		Aimpl.Field("totalstakeamount"),
		Aimpl.Field("totaltokenamount")))
}
