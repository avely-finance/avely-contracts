package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
	"strconv"
)

func (tr *Transitions) DelegateStakeSuccess() {
	Start("DelegateStake: Stake 10 ZIL")

	p := tr.DeployAndUpgrade()

	p.Aproxy.UpdateWallet(sdk.Cfg.Key1)

	// Because of DelegHasNoSufficientAmt
	tx, err := p.Aproxy.DelegateStake(ToZil(1))
	AssertError(tx, err, "DelegStakeNotEnough")

	// Success delegate
	AssertSuccess(p.Aproxy.DelegateStake(ToZil(20)))

	lastrewardcycle := p.Zimpl.Field("lastrewardcycle")

	AssertEqual(p.Zimpl.Field("buff_deposit_deleg", "0x"+p.GetBuffer().Addr, sdk.Cfg.AzilSsnAddress, lastrewardcycle), ToZil(20))

	AssertEqual(p.Zimpl.Field("buff_deposit_deleg", "0x"+p.GetBuffer().Addr, sdk.Cfg.AzilSsnAddress, p.Zimpl.Field("lastrewardcycle")), ToZil(20))
	AssertEqual(p.Aimpl.Field("_balance"), "0")

	AssertEqual(p.Aimpl.Field("totalstakeamount"), ToZil(1020))
	AssertEqual(p.Aimpl.Field("totaltokenamount"), ToAzil(1020))

	AssertEqual(p.Aimpl.Field("balances", "0x"+sdk.Cfg.Admin), ToAzil(1000))
	AssertEqual(p.Aimpl.Field("balances", "0x"+sdk.Cfg.Addr1), ToAzil(20))

	AssertEqual(p.Aimpl.Field("last_buf_deposit_cycle_deleg", "0x"+sdk.Cfg.Addr1), lastrewardcycle)

	// Check delegate to the next cycle
	p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare)
	p.Aproxy.DelegateStake(ToZil(20))

	nextCycleStr := StrAdd(lastrewardcycle, "1")

	AssertEqual(p.Aimpl.Field("last_buf_deposit_cycle_deleg", "0x"+sdk.Cfg.Addr1), nextCycleStr)
}

func (tr *Transitions) DelegateStakeBuffersRotation() {
	Start("DelegateStake: Buffers rotation")

	p := tr.DeployAndUpgrade()

	anotherBuffer, err := contracts.NewBufferContract(sdk, p.Aimpl.Addr, p.Zproxy.Addr, p.Zimpl.Addr)
	if err != nil {
		GetLog().Fatal("Deploy buffer error = " + err.Error())
	}

	new_buffers := []string{"0x" + p.GetBuffer().Addr, "0x" + p.GetBuffer().Addr, "0x" + anotherBuffer.Addr}

	AssertSuccess(p.Aimpl.ChangeBuffers(new_buffers))
	p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)
	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	AssertSuccess(p.Aproxy.DelegateStake(ToZil(10)))

	lastRewardCycle, _ := strconv.ParseInt(p.Zimpl.Field("lastrewardcycle"), 10, 64)
	index := lastRewardCycle % int64(len(new_buffers))
	activeBufferAddr := new_buffers[index]
	AssertEqual(p.Zimpl.Field("buff_deposit_deleg", activeBufferAddr, sdk.Cfg.AzilSsnAddress, strconv.FormatInt(lastRewardCycle, 10)), ToZil(10))
}
