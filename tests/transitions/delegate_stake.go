package transitions

import (
	//"log"
	"Azil/test/deploy"
	"strconv"
)

func (t *Testing) DelegateStakeSuccess() {
	t.LogStart("DelegateStake: Stake 10 ZIL")

	_, Zimpl, Aimpl, Buffer, _ := t.DeployAndUpgrade()

	Aimpl.UpdateWallet(key1)
	t.AssertSuccess(Aimpl.DelegateStake(zil(20)))

	t.AssertEqual(Zimpl.Field("buff_deposit_deleg", "0x"+Buffer.Addr, AZIL_SSN_ADDRESS, Zimpl.Field("lastrewardcycle")), zil(20))
	t.AssertEqual(Aimpl.Field("_balance"), "0")
	t.AssertEqual(Aimpl.Field("totalstakeamount"), zil(1020))
	t.AssertEqual(Aimpl.Field("totaltokenamount"), azil(1020))
	t.AssertEqual(Aimpl.Field("balances", "0x"+admin), azil(1000))
	t.AssertEqual(Aimpl.Field("balances", "0x"+addr1), azil(20))
}

func (t *Testing) DelegateStakeBuffersRotation() {
	t.LogStart("DelegateStake: Buffers rotation")

	Zproxy, Zimpl, Aimpl, Buffer, _ := t.DeployAndUpgrade()

	anotherBuffer, err1 := deploy.NewBufferContract(adminKey, Aimpl.Addr, AZIL_SSN_ADDRESS, Zproxy.Addr, Zimpl.Addr)
	if err1 != nil {
		t.LogError("Deploy buffer error = ", err1)
	}

	new_buffers := []string{"0x" + Buffer.Addr, "0x" + Buffer.Addr, "0x" + anotherBuffer.Addr}

	t.AssertSuccess(Aimpl.ChangeBuffers(new_buffers))
	Zproxy.UpdateWallet(verifierKey)
	t.AssertSuccess(Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT))

	t.AssertSuccess(Aimpl.DelegateStake(zil(10)))

	lastRewardCycle, _ := strconv.ParseInt(Zimpl.Field("lastrewardcycle"), 10, 64)
	index := lastRewardCycle % int64(len(new_buffers))
	activeBufferAddr := new_buffers[index]
	t.AssertEqual(Zimpl.Field("buff_deposit_deleg", activeBufferAddr, AZIL_SSN_ADDRESS, strconv.FormatInt(lastRewardCycle, 10)), zil(10))
}
