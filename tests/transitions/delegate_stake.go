package transitions

import (
	"strconv"

	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) DelegateStakeSuccess() {
	Start("DelegateStake: Stake 10 ZIL")

	p := tr.DeployAndUpgrade()
	totalSsnInitialDelegateZil := len(sdk.Cfg.SsnAddrs) * sdk.Cfg.SsnInitialDelegateZil
	delegateStakeHolder(p)

	p.Azil.UpdateWallet(sdk.Cfg.Key1)

	// Because of DelegHasNoSufficientAmt
	tx, _ := p.Azil.DelegateStake(ToZil(1))
	AssertError(tx, "DelegStakeNotEnough")

	// Success delegate
	ssnIn := p.GetSsnAddressForInput()
	AssertSuccess(p.Azil.DelegateStake(ToZil(20)))

	lastrewardcycle := strconv.Itoa(p.Zimpl.GetLastRewardCycle())

	AssertEqual(Field(p.Zimpl, "buff_deposit_deleg", p.GetBuffer().Addr, ssnIn, lastrewardcycle), ToZil(20))
	AssertEqual(Field(p.Azil, "_balance"), "0")

	AssertEqual(Field(p.Azil, "totalstakeamount"), ToZil(totalSsnInitialDelegateZil+20))
	AssertEqual(Field(p.Azil, "total_supply"), ToAzil(totalSsnInitialDelegateZil+20))

	AssertEqual(Field(p.Azil, "balances", sdk.Cfg.Admin), ToAzil(totalSsnInitialDelegateZil))
	AssertEqual(Field(p.Azil, "balances", sdk.Cfg.Addr1), ToAzil(20))

	// Check delegate to the next cycle
	p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare)
	p.Azil.DelegateStake(ToZil(20))
}

func (tr *Transitions) DelegateStakeBuffersRotation() {
	Start("DelegateStake: Buffers rotation")

	p := tr.DeployAndUpgrade()

	anotherBuffer, err := contracts.NewBufferContract(sdk, p.Azil.Addr, p.Zproxy.Addr)
	if err != nil {
		GetLog().Fatal("Deploy buffer error = " + err.Error())
	}

	new_buffers := []string{p.GetBuffer().Addr, p.GetBuffer().Addr, anotherBuffer.Addr}
	AssertSuccess(p.Azil.WithUser(sdk.Cfg.OwnerKey).ChangeBuffers(new_buffers))
	p.Azil.UpdateWallet(sdk.Cfg.AdminKey) //back to admin

	AssertSuccess(p.Azil.DelegateStake(ToZil(10)))
	activeBufferAddr := calcActiveBufferAddr(2, new_buffers) // start from second cycle
	AssertEqual(Field(p.Zimpl, "buff_deposit_deleg", activeBufferAddr, sdk.Cfg.AzilSsnAddress, Field(p.Zimpl, "lastrewardcycle")), ToZil(10))

	//next reward cycle
	p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)
	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	ssnIn := p.GetSsnAddressForInput()
	AssertSuccess(p.Azil.DelegateStake(ToZil(10)))
	activeBufferAddr = calcActiveBufferAddr(3, new_buffers)
	AssertEqual(Field(p.Zimpl, "buff_deposit_deleg", activeBufferAddr, ssnIn, Field(p.Zimpl, "lastrewardcycle")), ToZil(10))
}

func delegateStakeHolder(p *contracts.Protocol) {
	tx, _ := p.Holder.DelegateStake(sdk.Cfg.AzilSsnAddress, ToZil(1))
	AssertError(tx, "HolderAlreadyInitialized")
}

func calcActiveBufferAddr(cycle int, buffers []string) string {
	index := int(cycle) % len(buffers)
	return buffers[index]
}
