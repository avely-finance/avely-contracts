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

	p.Aimpl.UpdateWallet(sdk.Cfg.Key1)

	// Because of DelegHasNoSufficientAmt
	tx, _ := p.Aimpl.DelegateStake(ToZil(1))
	AssertError(tx, "DelegStakeNotEnough")

	// Success delegate
	AssertSuccess(p.Aimpl.DelegateStake(ToZil(20)))

	lastrewardcycle := Field(p.Zimpl, "lastrewardcycle")

	AssertEqual(Field(p.Zimpl, "buff_deposit_deleg", p.GetBuffer().Addr, sdk.Cfg.AzilSsnAddress, lastrewardcycle), ToZil(20))

	AssertEqual(Field(p.Zimpl, "buff_deposit_deleg", p.GetBuffer().Addr, sdk.Cfg.AzilSsnAddress, Field(p.Zimpl, "lastrewardcycle")), ToZil(20))
	AssertEqual(Field(p.Aimpl, "_balance"), "0")

	AssertEqual(Field(p.Aimpl, "totalstakeamount"), ToZil(1020))
	AssertEqual(Field(p.Aimpl, "totaltokenamount"), ToAzil(1020))

	AssertEqual(Field(p.Aimpl, "balances", sdk.Cfg.Admin), ToAzil(1000))
	AssertEqual(Field(p.Aimpl, "balances", sdk.Cfg.Addr1), ToAzil(20))

	AssertEqual(Field(p.Aimpl, "last_buf_deposit_cycle_deleg", sdk.Cfg.Addr1), lastrewardcycle)

	// Check delegate to the next cycle
	p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare)
	p.Aimpl.DelegateStake(ToZil(20))

	nextCycleStr := StrAdd(lastrewardcycle, "1")

	AssertEqual(Field(p.Aimpl, "last_buf_deposit_cycle_deleg", sdk.Cfg.Addr1), nextCycleStr)
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

	AssertSuccess(p.Aimpl.DelegateStake(ToZil(10)))

	activeBufferAddr = calcActiveBufferAddr(p, new_buffers)
	testGetCurrentBuffer(p, activeBufferAddr)
	AssertEqual(Field(p.Zimpl, "buff_deposit_deleg", activeBufferAddr, sdk.Cfg.AzilSsnAddress, Field(p.Zimpl, "lastrewardcycle")), ToZil(10))
}

func calcActiveBufferAddr(p *contracts.Protocol, buffers []string) string {
	lrcInt64, _ := strconv.ParseInt(Field(p.Zimpl, "lastrewardcycle"), 10, 64)
	index := lrcInt64 % int64(len(buffers))
	return buffers[index]
}

func testGetCurrentBuffer(p *contracts.Protocol, activeBufferAddr string) {
	tx, _ := AssertSuccess(p.Aimpl.GetCurrentBuffer())
	AssertEvent(tx, Event{p.Aimpl.Addr, "GetCurrentBufferSuccess", ParamsMap{"buffer": activeBufferAddr}})
}
