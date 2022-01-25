package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) ChownStakeSuccess() {
	Start("Chown Stake Success")

	p := tr.DeployAndUpgrade()

	chownStakeSetup(p)

	key1, addr1, key2, addr2, ssn, _, userStake := chownStakeDefineParams(p)

	totaltokenamount := Field(p.Aimpl, "totaltokenamount")
	totalstakeamount := Field(p.Aimpl, "totalstakeamount")
	stake1_azil := StrMul(userStake, "64")
	stake2_azil := StrMul(userStake, "128")
	stake1_1 := StrMul(userStake, "2")
	stake1_2 := StrMul(userStake, "4")
	stake2_1 := StrMul(userStake, "8")

	//key1 delegates through Aimpl (this isn't a part of transfer process)
	AssertSuccess(p.Aimpl.WithUser(key1).DelegateStake(stake1_azil))
	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)
	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)

	//key1, key2 delegate to main contract
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[1], stake1_1))
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[2], stake1_2))
	AssertSuccess(p.Zproxy.WithUser(key2).DelegateStake(ssn[1], stake2_1))

	//key1, key2 wait 2 reward cycles (they should have no buffered depo in current/prev cycles, else swap request will fail)
	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)
	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)
	nextBuffer := p.GetBufferToSwapWith().Addr

	//key1, key2 claim rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[1]))
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[2]))
	AssertSuccess(p.Zproxy.WithUser(key2).WithdrawStakeRewards(ssn[1]))

	//key1 requests swap
	tx, _ := AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})

	//key2 requests swap
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key2).RequestDelegatorSwap(nextBuffer))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr2, "new_deleg": nextBuffer}})

	//offchain-tool calls ChownStakeConfirmSwap(addr1), expecting success
	tx, _ = AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeConfirmSwap(addr1))
	AssertEvent(tx, Event{p.Zimpl.Addr, "ConfirmDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", addr1), "")
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", nextBuffer, ssn[1]), stake1_1)
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", nextBuffer, ssn[2]), stake1_2)
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", ssn[1], nextBuffer), stake1_1)
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", ssn[2], nextBuffer), stake1_2)
	AssertEqual(Field(p.Aimpl, "totalstakeamount"), StrAdd(totalstakeamount, stake1_azil, stake1_1, stake1_2))
	AssertEqual(Field(p.Aimpl, "totaltokenamount"), StrAdd(totaltokenamount, Field(p.Aimpl, "balances", addr1)))

	//offchain-tool calls ChownStakeConfirmSwap(addr2), expecting success
	tx, _ = AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeConfirmSwap(addr2))
	AssertEvent(tx, Event{p.Zimpl.Addr, "ConfirmDelegatorSwap", ParamsMap{"initial_deleg": addr2, "new_deleg": nextBuffer}})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", addr2), "")
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", nextBuffer, ssn[1]), StrAdd(stake1_1, stake2_1))
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", ssn[1], nextBuffer), StrAdd(stake1_1, stake2_1))
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", ssn[2], nextBuffer), stake1_2)
	AssertEqual(Field(p.Aimpl, "totalstakeamount"), StrAdd(totalstakeamount, stake1_azil, stake1_1, stake1_2, stake2_1))
	AssertEqual(Field(p.Aimpl, "totaltokenamount"), StrAdd(totaltokenamount, Field(p.Aimpl, "balances", addr1), Field(p.Aimpl, "balances", addr2)))

	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)

	//nextBuffer becomes active
	activeBuffer := p.GetActiveBuffer().Addr
	AssertEqual(nextBuffer, activeBuffer)

	//key2 delegates through Aimpl
	//this isn't a part of transfer process, but delegate can happen before offchain-tool calls
	AssertSuccess(p.Aimpl.WithUser(key2).DelegateStake(stake2_azil))

	//offchain tool calls ChownStakeReDelegate for each SSN (excepting AzilSSN) when new cycle starts
	tx, _ = AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeReDelegate(ssn[1], StrAdd(stake1_1, stake2_1)))
	AssertTransition(tx, Transition{
		p.Zimpl.Addr, //sender
		"ReDelegateStakeSuccessCallBack",
		activeBuffer,
		"0",
		ParamsMap{"ssnaddr": ssn[1], "tossn": sdk.Cfg.AzilSsnAddress, "amount": StrAdd(stake1_1, stake2_1)},
	})
	tx, _ = AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeReDelegate(ssn[2], stake1_2))
	AssertTransition(tx, Transition{
		p.Zimpl.Addr, //sender
		"ReDelegateStakeSuccessCallBack",
		activeBuffer,
		"0",
		ParamsMap{"ssnaddr": ssn[2], "tossn": sdk.Cfg.AzilSsnAddress, "amount": stake1_2},
	})

	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", activeBuffer, sdk.Cfg.AzilSsnAddress), StrAdd(stake1_1, stake1_2, stake2_1, stake2_azil))
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", sdk.Cfg.AzilSsnAddress, activeBuffer), StrAdd(stake1_1, stake1_2, stake2_1, stake2_azil))
	AssertEqual(Field(p.Aimpl, "totalstakeamount"), StrAdd(totalstakeamount, stake1_azil, stake1_1, stake1_2, stake2_1, stake2_azil))
	AssertEqual(Field(p.Aimpl, "totaltokenamount"), StrAdd(totaltokenamount, Field(p.Aimpl, "balances", addr1), Field(p.Aimpl, "balances", addr2)))
}

func (tr *Transitions) ChownStakeManySsnSuccess() {
	Start("Chown Stake Success")

	p := tr.DeployAndUpgrade()

	chownStakeSetup(p)

	key1, addr1, _, _, ssn, _, userStake := chownStakeDefineParams(p)
	totaltokenamount := Field(p.Aimpl, "totaltokenamount")
	totalstakeamount := Field(p.Aimpl, "totalstakeamount")

	//key1 delegates to main contract
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(sdk.Cfg.AzilSsnAddress, userStake))
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[1], userStake))
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[2], userStake))
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[3], userStake))
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[4], userStake))
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[5], userStake))

	//key1 waits 2 reward cycles (they should have no buffered depo in current/prev cycles, else swap request will fail)
	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)
	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)
	nextBuffer := p.GetBufferToSwapWith().Addr

	//key1 claims rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(sdk.Cfg.AzilSsnAddress))
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[1]))
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[2]))
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[3]))
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[4]))
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[5]))

	//key1 requests swap
	tx, _ := AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})

	//offchain-tool calls ChownStakeConfirmSwap(addr1), expecting success
	tx, _ = AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeConfirmSwap(addr1))
	AssertEvent(tx, Event{p.Zimpl.Addr, "ConfirmDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})

	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)

	//nextBuffer becomes active
	activeBuffer := p.GetActiveBuffer().Addr
	AssertEqual(nextBuffer, activeBuffer)

	//offchain tool calls ChownStakeReDelegate for each ssn/amount
	tx, _ = p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeReDelegate(sdk.Cfg.AzilSsnAddress, userStake)
	AssertError(tx, "ChownStakeReDelegateAzilSsn")
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeReDelegate(ssn[1], userStake))
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeReDelegate(ssn[2], userStake))
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeReDelegate(ssn[3], userStake))
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeReDelegate(ssn[4], userStake))
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeReDelegate(ssn[5], userStake))

	//check balances
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", activeBuffer, sdk.Cfg.AzilSsnAddress), StrMul(userStake, "6"))
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", sdk.Cfg.AzilSsnAddress, activeBuffer), StrMul(userStake, "6"))
	AssertEqual(Field(p.Aimpl, "totalstakeamount"), StrAdd(totalstakeamount, StrMul(userStake, "6")))
	AssertEqual(Field(p.Aimpl, "totaltokenamount"), StrAdd(totaltokenamount, Field(p.Aimpl, "balances", addr1)))

}

func (tr *Transitions) ChownStakeAimplErrors() {
	Start("Chown Stake Aimpl Errors")

	p := tr.DeployAndUpgrade()

	chownStakeSetup(p)

	key1, addr1, _, _, ssn, _, userStake := chownStakeDefineParams(p)

	//key1 delegates to main contract
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[1], userStake))

	//key1 waits 2 reward cycles
	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)
	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)
	nextBuffer := p.GetBufferToSwapWith().Addr

	//key1 claims rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[1]))

	//offchain-tool calls ChownStakeConfirmSwap(addr1), but addr1 didn't called RequestDelegatorSwap before, expecting error
	tx, _ := p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeConfirmSwap(addr1)
	AssertError(tx, "ChownStakeSwapRequestNotFound")

	//key1 requests swap with NOT buffer address
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(ssn[2]))

	//call ChownStake for addr1, expecting error
	tx, _ = p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeConfirmSwap(addr1)
	AssertError(tx, "BufferAddrUnknown")

	//key1 requests swap with NOT next buffer address
	activeBuffer := p.GetActiveBuffer()
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(activeBuffer.Addr))

	//call ChownStake for addr1, expecting error
	tx, _ = p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeConfirmSwap(addr1)
	AssertError(tx, "ChownStakeSwapRequestWrongBuffer")

	//key1 withdraws some amount, then requests swap
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeAmt(ssn[1], userStake))
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})

	//call ChownStake for addr1, expecting error
	tx, _ = p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeConfirmSwap(addr1)
	AssertError(tx, "ChownStakePendingWithdrawal")
}

func (tr *Transitions) ChownStakeZimplErrors() {
	Start("Chown Stake Zimpl Errors")

	p := tr.DeployAndUpgrade()

	chownStakeSetup(p)
	key1, addr1, _, _, ssn, _, userStake := chownStakeDefineParams(p)

	//key1 delegates to main contract, expecting success
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[1], userStake))

	//key1 requests delegator swap, but he has buffered deposit, expecting DelegHasBufferedDeposit
	nextBuffer := p.GetBufferToSwapWith().Addr
	tx, _ := p.Zproxy.RequestDelegatorSwap(nextBuffer)
	AssertZimplError(tx, -8)

	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)
	nextBuffer = p.GetBufferToSwapWith().Addr

	//key1 requests delegator swap, but he has buffered deposit in previous cycle, expecting DelegHasBufferedDeposit
	tx, _ = p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer)
	AssertZimplError(tx, -8)

	chownStakeNextCycle(p)
	nextBuffer = p.GetBufferToSwapWith().Addr

	//key1 requests delegator swap, but he has unclaimed rewards, expecting DelegHasUnwithdrawRewards
	tx, _ = p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer)
	AssertZimplError(tx, -12)

	//key1 claims rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[1]))

	//key1 requests delegator swap, but Holder has unclaimed rewards, expecting DelegHasUnwithdrawRewards
	//workflow of this use case is: Verifier->AssignStakeReward, User->RequestDelegatorSwap, Aimpl->DrainBuffer
	tx, _ = p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer)
	AssertZimplError(tx, -12)

	chownStakeNextCycleOffchain(p)

	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})
}

func chownStakeDefineParams(p *contracts.Protocol) (string, string, string, string, []string, string, string) {
	key1 := sdk.Cfg.Key1
	addr1 := sdk.Cfg.Addr1
	key2 := sdk.Cfg.Key2
	addr2 := sdk.Cfg.Addr2
	ssn := []string{"0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000001",
		"0x0000000000000000000000000000000000000002", "0x0000000000000000000000000000000000000003",
		"0x0000000000000000000000000000000000000004", "0x0000000000000000000000000000000000000005"}
	minStake := Field(p.Zimpl, "minstake")
	userStake := ToZil(10)
	return key1, addr1, key2, addr2, ssn, minStake, userStake
}

func chownStakeSetup(p *contracts.Protocol) {
	_, _, _, _, ssn, minStake, _ := chownStakeDefineParams(p)

	prevWallet := p.Zproxy.Contract.Wallet

	//TODO: move this to protocol.go
	//add buffers to protocol, we need 3
	buffer2, _ := p.DeployBuffer()
	buffer3, _ := p.DeployBuffer()
	p.Buffers = append(p.Buffers, buffer2, buffer3)
	p.SyncBufferAndHolder()
	p.SetupShortcuts(GetLog())

	//add test SSNs to main staking contract
	p.Zproxy.UpdateWallet(sdk.Cfg.AdminKey)
	AssertSuccess(p.Zproxy.AddSSN(ssn[0], "SSN 0"))
	AssertSuccess(p.Zproxy.AddSSN(ssn[1], "SSN 1"))
	AssertSuccess(p.Zproxy.AddSSN(ssn[2], "SSN 2"))
	AssertSuccess(p.Zproxy.AddSSN(ssn[3], "SSN 3"))
	AssertSuccess(p.Zproxy.AddSSN(ssn[4], "SSN 4"))
	AssertSuccess(p.Zproxy.AddSSN(ssn[5], "SSN 5"))
	AssertSuccess(p.Zproxy.DelegateStake(ssn[0], minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn[1], minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn[2], minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn[3], minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn[4], minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn[5], minStake))

	p.Zproxy.Contract.Wallet = prevWallet

	//ssns will become active on the next cycle
	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)
}

func chownStakeNextCycle(p *contracts.Protocol) {
	_, _, _, _, ssn, _, _ := chownStakeDefineParams(p)
	prevWallet := p.Zproxy.Contract.Wallet

	p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)
	ssnRewardFactor := map[string]string{
		ssn[0]:                 "100",
		ssn[1]:                 "100",
		ssn[2]:                 "100",
		ssn[3]:                 "100",
		ssn[4]:                 "100",
		ssn[5]:                 "100",
		sdk.Cfg.AzilSsnAddress: sdk.Cfg.AzilSsnRewardShare,
	}
	AssertSuccess(p.Zproxy.AssignStakeRewardList(ssnRewardFactor, "10000"))

	p.Zproxy.Contract.Wallet = prevWallet
}

func chownStakeNextCycleOffchain(p *contracts.Protocol) {
	buffer := p.GetBufferToDrain()
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(buffer.Addr))
	//AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeReDelegate())
}

func (tr *Transitions) ChownStakeRequireDrainBuffer() {
	Start("Chown Stake Drain Buffer")

	p := tr.DeployAndUpgrade()

	chownStakeSetup(p)

	key1, addr1, _, _, ssn, _, userStake := chownStakeDefineParams(p)

	//key1 delegates to main contract
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[1], userStake))

	//after 3 cycles all buffers are empty
	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)
	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)
	chownStakeNextCycle(p)
	chownStakeNextCycleOffchain(p)

	//next cycle
	chownStakeNextCycle(p)
	nextBuffer := p.GetBufferToSwapWith().Addr

	//quick swap sequence!

	//key1 claims rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[1]))

	//key1 requests swap
	tx, _ := AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})

	//offchain-tool calls ChownStakeConfirmSwap(addr1) before DrainBuffer(), expecting error
	tx, _ = p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeConfirmSwap(addr1)
	AssertError(tx, "BufferNotDrained")

	//drain buffer
	chownStakeNextCycleOffchain(p)

	//offchain-tool re-calls ChownStakeConfirmSwap(addr1) after DrainBuffer(), expecting success
	tx, _ = AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ChownStakeConfirmSwap(addr1))
	AssertEvent(tx, Event{p.Zimpl.Addr, "ConfirmDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})
}
