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
	tx, _ := p.Aproxy.DelegateStake(ToZil(1))
	AssertError(tx, "DelegStakeNotEnough")

	// Success delegate
	AssertSuccess(p.Aproxy.DelegateStake(ToZil(20)))

	lastrewardcycle := p.Zimpl.Field("lastrewardcycle")

	AssertEqual(p.Zimpl.Field("buff_deposit_deleg", p.GetBuffer().Addr, sdk.Cfg.AzilSsnAddress, lastrewardcycle), ToZil(20))

	AssertEqual(p.Zimpl.Field("buff_deposit_deleg", p.GetBuffer().Addr, sdk.Cfg.AzilSsnAddress, p.Zimpl.Field("lastrewardcycle")), ToZil(20))
	AssertEqual(p.Aimpl.Field("_balance"), "0")

	AssertEqual(p.Aimpl.Field("totalstakeamount"), ToZil(1020))
	AssertEqual(p.Aimpl.Field("totaltokenamount"), ToAzil(1020))

	AssertEqual(p.Aimpl.Field("balances", sdk.Cfg.Admin), ToAzil(1000))
	AssertEqual(p.Aimpl.Field("balances", sdk.Cfg.Addr1), ToAzil(20))

	AssertEqual(p.Aimpl.Field("last_buf_deposit_cycle_deleg", sdk.Cfg.Addr1), lastrewardcycle)

	// Check delegate to the next cycle
	p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare)
	p.Aproxy.DelegateStake(ToZil(20))

	nextCycleStr := StrAdd(lastrewardcycle, "1")

	AssertEqual(p.Aimpl.Field("last_buf_deposit_cycle_deleg", sdk.Cfg.Addr1), nextCycleStr)
}

func (tr *Transitions) DelegateStakeBuffersRotation() {
	Start("DelegateStake: Buffers rotation")

	p := tr.DeployAndUpgrade()

	anotherBuffer, err := contracts.NewBufferContract(sdk, p.Aimpl.Addr, p.Zproxy.Addr, p.Zimpl.Addr)
	if err != nil {
		GetLog().Fatal("Deploy buffer error = " + err.Error())
	}

	new_buffers := []string{p.GetBuffer().Addr, p.GetBuffer().Addr, anotherBuffer.Addr}
	AssertSuccess(p.Aimpl.ChangeBuffers(new_buffers))
	activeBufferAddr := calcActiveBufferAddr(p, new_buffers)
	testGetCurrentBuffer(p, activeBufferAddr)

	//next reward cycle
	p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)
	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	AssertSuccess(p.Aproxy.DelegateStake(ToZil(10)))

	activeBufferAddr = calcActiveBufferAddr(p, new_buffers)
	testGetCurrentBuffer(p, activeBufferAddr)
	AssertEqual(p.Zimpl.Field("buff_deposit_deleg", activeBufferAddr, sdk.Cfg.AzilSsnAddress, p.Zimpl.Field("lastrewardcycle")), ToZil(10))
}

func calcActiveBufferAddr(p *contracts.Protocol, buffers []string) string {
	lrcInt64, _ := strconv.ParseInt(p.Zimpl.Field("lastrewardcycle"), 10, 64)
	index := lrcInt64 % int64(len(buffers))
	return buffers[index]
}

func testGetCurrentBuffer(p *contracts.Protocol, activeBufferAddr string) {
	tx, _ := AssertSuccess(p.Aimpl.GetCurrentBuffer())
	AssertEvent(tx, Event{p.Aimpl.Addr, "GetCurrentBufferSuccess", ParamsMap{"buffer": activeBufferAddr}})
}
