package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) ReAssignStakeSuccess() {
	Start("Transfer Stake Success")

	p := tr.DeployAndUpgrade()

	reassignStakeSetup(p)

	key1, addr1, key2, addr2, ssn1, ssn2, _, userStake := reassignStakeDefineParams(p)
	totaltokenamount := Field(p.Aimpl, "totaltokenamount")
	totalstakeamount := Field(p.Aimpl, "totalstakeamount")
	userStakeThroughAimpl := userStake
	userStake2x := StrAdd(userStake, userStake)
	userStake4x := StrAdd(userStake2x, userStake2x)

	//key1 delegates through Aimpl (this isn't a part of transfer process)
	AssertSuccess(p.Aimpl.WithUser(key1).DelegateStake(userStakeThroughAimpl))
	reassignStakeNextCycle(p)
	reassignStakeNextCycleOffchain(p)
	reassignStakeNextCycle(p)
	reassignStakeNextCycleOffchain(p)

	//key1, key2 delegate to main contract
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn1, userStake))
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn2, userStake2x))
	AssertSuccess(p.Zproxy.WithUser(key2).DelegateStake(ssn1, userStake4x))

	//key1, key2 wait 2 reward cycles (they should have no buffered depo in current/prev cycles, else swap request will fail)
	reassignStakeNextCycle(p)
	reassignStakeNextCycleOffchain(p)
	reassignStakeNextCycle(p)
	reassignStakeNextCycleOffchain(p)
	swapAddr := reassignStakeGetSwapAddr(p)

	//key1, key2 claim rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn1))
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn2))
	AssertSuccess(p.Zproxy.WithUser(key2).WithdrawStakeRewards(ssn1))

	//key1 requests swap
	tx, _ := AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(swapAddr))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": swapAddr}})

	//key2 requests swap
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key2).RequestDelegatorSwap(swapAddr))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr2, "new_deleg": swapAddr}})

	//offchain-tool calls ReAssignStake(addr1), expecting success
	tx, _ = AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ReAssignStake(addr1))
	AssertEvent(tx, Event{p.Zimpl.Addr, "ConfirmDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": swapAddr}})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", addr1), "")
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", swapAddr, ssn1), userStake)
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", swapAddr, ssn2), userStake2x)
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", ssn1, swapAddr), userStake)
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", ssn2, swapAddr), userStake2x)
	AssertEqual(Field(p.Aimpl, "totalstakeamount"), StrAdd(totalstakeamount, userStakeThroughAimpl, userStake, userStake2x))
	AssertEqual(Field(p.Aimpl, "totaltokenamount"), StrAdd(totaltokenamount, Field(p.Aimpl, "balances", addr1)))

	//offchain-tool calls ReAssignStake(addr2), expecting success
	tx, _ = AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ReAssignStake(addr2))
	AssertEvent(tx, Event{p.Zimpl.Addr, "ConfirmDelegatorSwap", ParamsMap{"initial_deleg": addr2, "new_deleg": swapAddr}})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", addr2), "")
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", swapAddr, ssn1), StrAdd(userStake, userStake4x))
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", ssn1, swapAddr), StrAdd(userStake, userStake4x))
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", ssn2, swapAddr), userStake2x)
	AssertEqual(Field(p.Aimpl, "totalstakeamount"), StrAdd(totalstakeamount, userStakeThroughAimpl, userStake, userStake2x, userStake4x))
	AssertEqual(Field(p.Aimpl, "totaltokenamount"), StrAdd(totaltokenamount, Field(p.Aimpl, "balances", addr1), Field(p.Aimpl, "balances", addr2)))

	//call of ReAssignStakeReDelegate() will not break transfer process
	tx, _ = AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ReAssignStakeReDelegate())

	reassignStakeNextCycle(p)
	reassignStakeNextCycleOffchain(p)

	//key2 delegates through Aimpl
	//this isn't a part of transfer process, but delegate can happen before offchain-tool calls
	AssertSuccess(p.Aimpl.WithUser(key2).DelegateStake(userStakeThroughAimpl))

	//offchain tool calls ReAssignStakeReDelegate once when new cycle starts
	tx, _ = AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ReAssignStakeReDelegate())
	AssertTransition(tx, Transition{
		p.Zimpl.Addr, //sender
		"ReDelegateStakeSuccessCallBack",
		swapAddr,
		"0",
		ParamsMap{"ssnaddr": ssn1, "tossn": sdk.Cfg.AzilSsnAddress, "amount": StrAdd(userStake, userStake4x)},
	})
	AssertTransition(tx, Transition{
		p.Zimpl.Addr, //sender
		"ReDelegateStakeSuccessCallBack",
		swapAddr,
		"0",
		ParamsMap{"ssnaddr": ssn2, "tossn": sdk.Cfg.AzilSsnAddress, "amount": userStake2x},
	})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", swapAddr, sdk.Cfg.AzilSsnAddress), StrAdd(userStake, userStake2x, userStake4x, userStakeThroughAimpl))
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", sdk.Cfg.AzilSsnAddress, swapAddr), StrAdd(userStake, userStake2x, userStake4x, userStakeThroughAimpl))
}

func (tr *Transitions) ReAssignStakeAimplErrors() {
	Start("Transfer Stake Aimpl Errors")

	p := tr.DeployAndUpgrade()

	reassignStakeSetup(p)

	key1, addr1, _, _, ssn1, ssn2, _, userStake := reassignStakeDefineParams(p)

	//key1 delegates to main contract
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn1, userStake))

	//key1 waits 2 reward cycles
	reassignStakeNextCycle(p)
	reassignStakeNextCycleOffchain(p)
	reassignStakeNextCycle(p)
	reassignStakeNextCycleOffchain(p)
	swapAddr := reassignStakeGetSwapAddr(p)

	//key1 claims rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn1))

	//offchain-tool calls ReAssignStake(addr1), but addr1 didn't called RequestDelegatorSwap before, expecting error
	tx, _ := p.Aimpl.WithUser(sdk.Cfg.AdminKey).ReAssignStake(addr1)
	AssertError(tx, "ReAssignStakeSwapRequestNotFound")

	//key1 requests swap with NOT buffer address
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(ssn2))

	//call ReAssignStake for addr1, expecting error
	tx, _ = p.Aimpl.WithUser(key1).ReAssignStake(addr1)
	AssertError(tx, "BufferAddrUnknown")

	//key1 requests swap with NOT next buffer address
	_, activeBuffer := p.GetActiveBuffer()
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(activeBuffer.Addr))

	//call ReAssignStake for addr1, expecting error
	tx, _ = p.Aimpl.WithUser(key1).ReAssignStake(addr1)
	AssertError(tx, "ReAssignStakeSwapRequestWrongBuffer")

	//key1 withdraws some amount, then requests swap
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeAmt(ssn1, userStake))
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(swapAddr))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": swapAddr}})

	//call ReAssignStake for addr1, expecting error
	tx, _ = p.Aimpl.WithUser(key1).ReAssignStake(addr1)
	AssertError(tx, "ReAssignStakePendingWithdrawal")
}

func (tr *Transitions) ReAssignStakeZimplErrors() {
	Start("Transfer Stake Zimpl Errors")

	p := tr.DeployAndUpgrade()

	reassignStakeSetup(p)
	key1, addr1, _, _, ssn1, _, _, userStake := reassignStakeDefineParams(p)

	//key1 delegates to main contract, expecting success
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn1, userStake))

	//key1 requests delegator swap, but he has buffered deposit, expecting DelegHasBufferedDeposit
	swapAddr := reassignStakeGetSwapAddr(p)
	tx, _ := p.Zproxy.RequestDelegatorSwap(swapAddr)
	AssertZimplError(tx, -8)

	reassignStakeNextCycle(p)
	reassignStakeNextCycleOffchain(p)
	swapAddr = reassignStakeGetSwapAddr(p)

	//key1 requests delegator swap, but he has buffered deposit in previous cycle, expecting DelegHasBufferedDeposit
	tx, _ = p.Zproxy.WithUser(key1).RequestDelegatorSwap(swapAddr)
	AssertZimplError(tx, -8)

	reassignStakeNextCycle(p)
	swapAddr = reassignStakeGetSwapAddr(p)

	//key1 requests delegator swap, but he has unclaimed rewards, expecting DelegHasUnwithdrawRewards
	tx, _ = p.Zproxy.WithUser(key1).RequestDelegatorSwap(swapAddr)
	AssertZimplError(tx, -12)

	//key1 claims rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn1))

	//key1 requests delegator swap, but Holder has unclaimed rewards, expecting DelegHasUnwithdrawRewards
	//workflow of this use case is: Verifier->AssignStakeReward, User->RequestDelegatorSwap, Aimpl->DrainBuffer
	tx, _ = p.Zproxy.WithUser(key1).RequestDelegatorSwap(swapAddr)
	AssertZimplError(tx, -12)

	reassignStakeNextCycleOffchain(p)

	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(swapAddr))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": swapAddr}})
}

func reassignStakeDefineParams(p *contracts.Protocol) (string, string, string, string, string, string, string, string) {
	key1 := sdk.Cfg.Key1
	addr1 := sdk.Cfg.Addr1
	key2 := sdk.Cfg.Key2
	addr2 := sdk.Cfg.Addr2
	ssn1 := "0x0000000000000000000000000000000000000001"
	ssn2 := "0x0000000000000000000000000000000000000002"
	minStake := Field(p.Zimpl, "minstake")
	userStake := ToZil(10)
	return key1, addr1, key2, addr2, ssn1, ssn2, minStake, userStake
}

func reassignStakeSetup(p *contracts.Protocol) {
	_, _, _, _, ssn1, ssn2, minStake, _ := reassignStakeDefineParams(p)

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
	AssertSuccess(p.Zproxy.AddSSN(ssn1, "SSN 1"))
	AssertSuccess(p.Zproxy.AddSSN(ssn2, "SSN 2"))
	AssertSuccess(p.Zproxy.DelegateStake(ssn1, minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn2, minStake))

	p.Zproxy.Contract.Wallet = prevWallet

	//ssns will become active on the next cycle
	reassignStakeNextCycle(p)
	reassignStakeNextCycleOffchain(p)
}

func reassignStakeNextCycle(p *contracts.Protocol) {
	_, _, _, _, ssn1, ssn2, _, _ := reassignStakeDefineParams(p)
	prevWallet := p.Zproxy.Contract.Wallet

	p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)
	ssnRewardFactor := map[string]string{
		ssn1:                   "100",
		ssn2:                   "100",
		sdk.Cfg.AzilSsnAddress: sdk.Cfg.AzilSsnRewardShare,
	}
	AssertSuccess(p.Zproxy.AssignStakeRewardList(ssnRewardFactor, "10000"))

	p.Zproxy.Contract.Wallet = prevWallet
}

func reassignStakeNextCycleOffchain(p *contracts.Protocol) {
	_, buffer := p.GetBufferToDrain()
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(buffer.Addr))
	//AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).ReAssignStakeReDelegate())
}

func reassignStakeGetSwapAddr(p *contracts.Protocol) string {
	_, buffer := p.GetBufferByOffset(1)
	return buffer.Addr
}
