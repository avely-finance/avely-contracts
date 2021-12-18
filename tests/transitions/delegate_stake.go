package transitions

import (
	"Azil/test/contracts"
	"Azil/test/helpers"
	"strconv"
)

func (tr *Transitions) DelegateStakeSuccess() {
	log.Start("DelegateStake: Stake 10 ZIL")

	Zproxy, Zimpl, Aimpl, Buffer, _ := tr.DeployAndUpgrade()

	Aimpl.UpdateWallet(key1)

	// Because of DelegHasNoSufficientAmt
	tx, err := Aimpl.DelegateStake(zil(1))
	t.AssertError(tx, err, -15)

	// Success delegate
	t.AssertSuccess(Aimpl.DelegateStake(zil(20)))

	lastrewardcycle := Zimpl.Field("lastrewardcycle")

	t.AssertEqual(Zimpl.Field("buff_deposit_deleg", "0x"+Buffer.Addr, AZIL_SSN_ADDRESS, lastrewardcycle), zil(20))

	t.AssertEqual(Zimpl.Field("buff_deposit_deleg", "0x"+Buffer.Addr, AZIL_SSN_ADDRESS, Zimpl.Field("lastrewardcycle")), zil(20))
	t.AssertEqual(Aimpl.Field("_balance"), "0")

	t.AssertEqual(Aimpl.Field("totalstakeamount"), zil(1020))
	t.AssertEqual(Aimpl.Field("totaltokenamount"), azil(1020))

	t.AssertEqual(Aimpl.Field("balances", "0x"+admin), azil(1000))
	t.AssertEqual(Aimpl.Field("balances", "0x"+addr1), azil(20))

	t.AssertEqual(Aimpl.Field("last_buf_deposit_cycle_deleg", "0x"+addr1), lastrewardcycle)

	// Check delegate to the next cycle
	Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT)
	Aimpl.DelegateStake(zil(20))

	nextCycleStr := helpers.StrAdd(lastrewardcycle, "1")

	t.AssertEqual(Aimpl.Field("last_buf_deposit_cycle_deleg", "0x"+addr1), nextCycleStr)
}

func (tr *Transitions) DelegateStakeBuffersRotation() {
	log.Start("DelegateStake: Buffers rotation")

	Zproxy, Zimpl, Aimpl, Buffer, _ := tr.DeployAndUpgrade()

	anotherBuffer, err := contracts.NewBufferContract(adminKey, Aimpl.Addr, AZIL_SSN_ADDRESS, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		log.Fatal("Deploy buffer error = ", err.Error())
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
