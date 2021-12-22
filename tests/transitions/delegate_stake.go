package transitions

import (
	"Azil/test/contracts"
	. "Azil/test/helpers"
	"strconv"
)

func (tr *Transitions) DelegateStakeSuccess() {
	t.Start("DelegateStake: Stake 10 ZIL")

	Zproxy, Zimpl, Aimpl, Buffer, _ := tr.DeployAndUpgrade()

	Aimpl.UpdateWallet(tr.cfg.Key1)

	// Because of DelegHasNoSufficientAmt
	tx, err := Aimpl.DelegateStake(Zil(1))
	t.AssertError(tx, err, -15)

	// Success delegate
	t.AssertSuccess(Aimpl.DelegateStake(Zil(20)))

	lastrewardcycle := Zimpl.Field("lastrewardcycle")

	t.AssertEqual(Zimpl.Field("buff_deposit_deleg", "0x"+Buffer.Addr, tr.cfg.AzilSsnAddress, lastrewardcycle), Zil(20))

	t.AssertEqual(Zimpl.Field("buff_deposit_deleg", "0x"+Buffer.Addr, tr.cfg.AzilSsnAddress, Zimpl.Field("lastrewardcycle")), Zil(20))
	t.AssertEqual(Aimpl.Field("_balance"), "0")

	t.AssertEqual(Aimpl.Field("totalstakeamount"), Zil(1020))
	t.AssertEqual(Aimpl.Field("totaltokenamount"), Azil(1020))

	t.AssertEqual(Aimpl.Field("balances", "0x"+tr.cfg.Admin), Azil(1000))
	t.AssertEqual(Aimpl.Field("balances", "0x"+tr.cfg.Addr1), Azil(20))

	t.AssertEqual(Aimpl.Field("last_buf_deposit_cycle_deleg", "0x"+tr.cfg.Addr1), lastrewardcycle)

	// Check delegate to the next cycle
	Zproxy.AssignStakeReward(tr.cfg.AzilSsnAddress, tr.cfg.AzilSsnRewardShare)
	Aimpl.DelegateStake(Zil(20))

	nextCycleStr := StrAdd(lastrewardcycle, "1")

	t.AssertEqual(Aimpl.Field("last_buf_deposit_cycle_deleg", "0x"+tr.cfg.Addr1), nextCycleStr)
}

func (tr *Transitions) DelegateStakeBuffersRotation() {
	t.Start("DelegateStake: Buffers rotation")

	Zproxy, Zimpl, Aimpl, Buffer, _ := tr.DeployAndUpgrade()

	anotherBuffer, err := contracts.NewBufferContract(tr.cfg.AdminKey, Aimpl.Addr, tr.cfg.AzilSsnAddress, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		log.Fatal("Deploy buffer error = " + err.Error())
	}

	new_buffers := []string{"0x" + Buffer.Addr, "0x" + Buffer.Addr, "0x" + anotherBuffer.Addr}

	t.AssertSuccess(Aimpl.ChangeBuffers(new_buffers))
	Zproxy.UpdateWallet(tr.cfg.VerifierKey)
	t.AssertSuccess(Zproxy.AssignStakeReward(tr.cfg.AzilSsnAddress, tr.cfg.AzilSsnRewardShare))

	t.AssertSuccess(Aimpl.DelegateStake(Zil(10)))

	lastRewardCycle, _ := strconv.ParseInt(Zimpl.Field("lastrewardcycle"), 10, 64)
	index := lastRewardCycle % int64(len(new_buffers))
	activeBufferAddr := new_buffers[index]
	t.AssertEqual(Zimpl.Field("buff_deposit_deleg", activeBufferAddr, tr.cfg.AzilSsnAddress, strconv.FormatInt(lastRewardCycle, 10)), Zil(10))
}
