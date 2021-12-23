package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
	// "strconv"
)

func (tr *Transitions) DelegateStakeSuccess() {
	t.Start("DelegateStake: Stake 10 ZIL")

	p := DeployAndUpgrade()

	p.Aimpl.UpdateWallet(sdk.Cfg.Key1)

	// Because of DelegHasNoSufficientAmt
	tx, err := p.Aimpl.DelegateStake(Zil(1))
	t.AssertError(tx, err, -15)

	// Success delegate
	t.AssertSuccess(p.Aimpl.DelegateStake(Zil(20)))

	lastrewardcycle := p.Zimpl.Field("lastrewardcycle")

	t.AssertEqual(p.Zimpl.Field("buff_deposit_deleg", "0x"+p.Buffer.Addr, sdk.Cfg.AzilSsnAddress, lastrewardcycle), Zil(20))

	t.AssertEqual(p.Zimpl.Field("buff_deposit_deleg", "0x"+p.Buffer.Addr, sdk.Cfg.AzilSsnAddress, p.Zimpl.Field("lastrewardcycle")), Zil(20))
	t.AssertEqual(p.Aimpl.Field("_balance"), "0")

	t.AssertEqual(p.Aimpl.Field("totalstakeamount"), Zil(1020))
	t.AssertEqual(p.Aimpl.Field("totaltokenamount"), Azil(1020))

	t.AssertEqual(p.Aimpl.Field("balances", "0x"+sdk.Cfg.Admin), Azil(1000))
	t.AssertEqual(p.Aimpl.Field("balances", "0x"+sdk.Cfg.Addr1), Azil(20))

	t.AssertEqual(p.Aimpl.Field("last_buf_deposit_cycle_deleg", "0x"+sdk.Cfg.Addr1), lastrewardcycle)

	// Check delegate to the next cycle
	p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare)
	p.Aimpl.DelegateStake(Zil(20))

	nextCycleStr := StrAdd(lastrewardcycle, "1")

	t.AssertEqual(p.Aimpl.Field("last_buf_deposit_cycle_deleg", "0x"+sdk.Cfg.Addr1), nextCycleStr)
}

// func (tr *Transitions) DelegateStakeBuffersRotation() {
// 	t.Start("DelegateStake: Buffers rotation")

// 	p := DeployAndUpgrade()

// 	anotherBuffer, err := contracts.NewBufferContract(sdk.Cfg.AdminKey, p.Aimpl.Addr, sdk.Cfg.AzilSsnAddress, Zproxy.Addr, Zimpl.Addr)
// 	if err != nil {
// 		log.Fatal("Deploy buffer error = " + err.Error())
// 	}

// 	new_buffers := []string{"0x" + Buffer.Addr, "0x" + Buffer.Addr, "0x" + anotherBuffer.Addr}

// 	t.AssertSuccess(p.Aimpl.ChangeBuffers(new_buffers))
// 	Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)
// 	t.AssertSuccess(Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

// 	t.AssertSuccess(p.Aimpl.DelegateStake(Zil(10)))

// 	lastRewardCycle, _ := strconv.ParseInt(Zimpl.Field("lastrewardcycle"), 10, 64)
// 	index := lastRewardCycle % int64(len(new_buffers))
// 	activeBufferAddr := new_buffers[index]
// 	t.AssertEqual(Zimpl.Field("buff_deposit_deleg", activeBufferAddr, sdk.Cfg.AzilSsnAddress, strconv.FormatInt(lastRewardCycle, 10)), Zil(10))
// }
