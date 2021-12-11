package transitions

import (
	"Azil/test/deploy"
)

func (t *Testing) ZilBalanceOf() {

	t.LogStart("ZilBalanceOf")

	// deploy smart contract
	_, _, aZilContract, _, _ := t.DeployAndUpgrade()

	/*******************************************************************************
	 * 1. Non-delegator address (addr2) should have empty balance
	 *******************************************************************************/
	t.LogStart("ZilBalanceOf, step 1")
	aZilContract.UpdateWallet(key2)
	balance2, _ := aZilContract.ZilBalanceOf(addr2)
	t.AssertEqual(balance2, zil(0))

	/*******************************************************************************
	 * 2. After delegate (addr2) should have its balance updated
	 *******************************************************************************/
	t.LogStart("ZilBalanceOf, step 2")
	aZilContract.DelegateStake(zil(15))
	balance2, _ = aZilContract.ZilBalanceOf(addr2)
	t.AssertEqual(balance2, zil(15))

	/*******************************************************************************
	 * 3. User balance in zil should be updated after restaking rewards
	 * because contract got more zils w/o azil minting
	 * so azil/zil exchange rate changed, azil now costs more zils than before
	 * so balance of addr2 in zil should be more
	 *******************************************************************************/
	t.LogStart("ZilBalanceOf, step 3")
	aZilContract.UpdateWallet(adminKey)
	aZilContract.IncreaseAutoRestakeAmount(zil(10))
	aZilContract.PerformAutoRestake()
	balance2, _ = aZilContract.ZilBalanceOf(addr2)
	t.AssertEqual(balance2, deploy.StrSum(zil(15), zil(10)))

	/*******************************************************************************
	 * 4. New user (addr3) delegating stake
	 * He'll get azils by new azil/zil rate, so zilBalanceOf should be equal to the delegated zils amount.
	 * First user's (addr2) zil balance should not be changed.
	 *******************************************************************************/
	t.LogStart("ZilBalanceOf, step 4")
	aZilContract.UpdateWallet(key3)
	aZilContract.DelegateStake(zil(10))
	balance2, _ = aZilContract.ZilBalanceOf(addr2)
	t.AssertEqual(balance2, deploy.StrSum(zil(15), zil(10)))
	balance3, _ := aZilContract.ZilBalanceOf(addr3)
	t.AssertEqual(balance3, zil(10))
}
