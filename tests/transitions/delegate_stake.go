package transitions

import (
	"strconv"

	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) DelegateStakeSuccess() {
	Start("DelegateStake: Stake 10 ZIL")

	p := tr.DeployAndUpgrade()
	//totalSsnInitialDelegateZil := len(sdk.Cfg.SsnAddrs) * sdk.Cfg.SsnInitialDelegateZil
	//for now to activate SSNs we delegate required stakes through Zproxy as admin
	totalSsnInitialDelegateZil := 0

	delegateStakeHolder(p)

	p.StZIL.SetSigner(alice)

	// Because of DelegHasNoSufficientAmt
	tx, _ := p.StZIL.DelegateStake(ToZil(1))
	AssertError(tx, p.StZIL.ErrorCode("DelegStakeNotEnough"))

	// Success delegate
	ssnIn := p.GetSsnAddressForInput()
	tx, _ = AssertSuccess(p.StZIL.DelegateStake(ToZil(20)))

	lastrewardcycle := strconv.Itoa(p.Zimpl.GetLastRewardCycle())
	aliceAddr := utils.GetAddressByWallet(alice)

	AssertEqual(Field(p.Zimpl, "buff_deposit_deleg", p.GetActiveBuffer().Addr, ssnIn, lastrewardcycle), ToZil(20))
	AssertEqual(Field(p.StZIL, "_balance"), "0")

	AssertEqual(Field(p.StZIL, "totalstakeamount"), ToZil(totalSsnInitialDelegateZil+20))
	AssertEqual(Field(p.StZIL, "total_supply"), ToStZil(totalSsnInitialDelegateZil+20))
	AssertEvent(tx, Event{p.StZIL.Addr, "Minted", ParamsMap{
		"minter":    p.StZIL.Addr,
		"recipient": aliceAddr,
		"amount":    ToStZil(20),
	}})

	admin := utils.GetAddressByWallet(celestials.Admin)
	if totalSsnInitialDelegateZil == 0 {
		AssertEqual(Field(p.StZIL, "balances", admin), "")
	} else {
		AssertEqual(Field(p.StZIL, "balances", admin), ToStZil(totalSsnInitialDelegateZil))
	}
	AssertEqual(Field(p.StZIL, "balances", aliceAddr), ToStZil(20))

	// Check delegate to the next cycle
	p.Zproxy.AssignStakeReward(sdk.Cfg.StZilSsnAddress, sdk.Cfg.StZilSsnRewardShare)
	p.StZIL.DelegateStake(ToZil(20))
}

func (tr *Transitions) DelegateStakeBuffersRotation() {
	Start("DelegateStake: Buffers rotation")

	p := tr.DeployAndUpgrade()

	anotherBuffer, err := contracts.NewBufferContract(sdk, p.StZIL.Addr, p.Zproxy.Addr, p.Zimpl.Addr, celestials.Admin)
	if err != nil {
		GetLog().Fatal("Deploy buffer error = " + err.Error())
	}

	new_buffers := []string{p.GetBuffer().Addr, p.GetBuffer().Addr, anotherBuffer.Addr}
	p.StZIL.SetSigner(celestials.Owner)
	AssertSuccess(p.StZIL.ChangeBuffers(new_buffers))

	p.StZIL.SetSigner(celestials.Admin) //back to admin

	ssnForInput := p.GetSsnAddressForInput()
	AssertSuccess(p.StZIL.DelegateStake(ToZil(10)))
	lrc := p.Zimpl.GetLastRewardCycle()
	activeBufferAddr := calcActiveBufferAddr(lrc, new_buffers)

	AssertEqual(Field(p.Zimpl, "buff_deposit_deleg", activeBufferAddr, ssnForInput, Field(p.Zimpl, "lastrewardcycle")), ToZil(10))

	//next reward cycle
	p.Zproxy.SetSigner(verifier)
	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.StZilSsnAddress, sdk.Cfg.StZilSsnRewardShare))

	ssnIn := p.GetSsnAddressForInput()
	AssertSuccess(p.StZIL.DelegateStake(ToZil(10)))
	activeBufferAddr = calcActiveBufferAddr(3, new_buffers)
	AssertEqual(Field(p.Zimpl, "buff_deposit_deleg", activeBufferAddr, ssnIn, Field(p.Zimpl, "lastrewardcycle")), ToZil(10))
}

func delegateStakeHolder(p *contracts.Protocol) {
	tx, _ := p.Holder.DelegateStake(sdk.Cfg.StZilSsnAddress, ToZil(1))
	AssertError(tx, p.Holder.ErrorCode("HolderAlreadyInitialized"))
}

func calcActiveBufferAddr(cycle int, buffers []string) string {
	index := int(cycle) % len(buffers)
	return buffers[index]
}
