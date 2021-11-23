package transitions

import (
    "Azil/test/deploy"
)

func (t *Testing) ZilBalanceOf() {

    t.LogStart("ZilBalanceOf")

    // deploy smart contract
    _, aZilContract, _, _ := t.DeployAndUpgrade()

    /*******************************************************************************
     * 1. Non-delegator address (addr2) should have empty balance
     *******************************************************************************/
    t.LogStart("================== ZilBalanceOf, step 1 ===================")
    aZilContract.UpdateWallet(key2)
    balance2, _ := aZilContract.ZilBalanceOf(addr2)
    t.AssertEqual(balance2, zil0)

    /*******************************************************************************
     * 2. After delegate (addr2) should have its balance updated
     *******************************************************************************/
    t.LogStart("================== ZilBalanceOf, step 2 ===================")
    aZilContract.DelegateStake(zil15)
    balance2, _ = aZilContract.ZilBalanceOf(addr2)
    t.AssertEqual(balance2, zil15)

    /*******************************************************************************
     * 3. After IncreaseTotalStakeAmount admin transition user balance in zil should be updated
     * because contract got more zils w/o azil minting
     * so azil/zil exchange rate changed, azil now costs more zils than before
     * so balance of addr2 in zil should be more
     *******************************************************************************/
    t.LogStart("================== ZilBalanceOf, step 3 ===================")
    aZilContract.UpdateWallet(adminKey)
    aZilContract.IncreaseTotalStakeAmount(zil10)
    balance2, _ = aZilContract.ZilBalanceOf(addr2)
    t.AssertEqual(balance2, deploy.StrSum(zil15, zil10))

    /*******************************************************************************
     * 4. New user (addr3) delegating stake
     * He'll get azils by new azil/zil rate, so zilBalanceOf should be equal to the delegated zils amount.
     * First user's (addr2) zil balance should not be changed.
     *******************************************************************************/
    t.LogStart("================== ZilBalanceOf, step 4 ===================")
    aZilContract.UpdateWallet(key3)
    aZilContract.DelegateStake(zil10)
    balance2, _ = aZilContract.ZilBalanceOf(addr2)
    t.AssertEqual(balance2, deploy.StrSum(zil15, zil10))
    balance3, _ := aZilContract.ZilBalanceOf(addr3)
    t.AssertEqual(balance3, zil10)

    t.LogEnd("ZilBalanceOf")
}
