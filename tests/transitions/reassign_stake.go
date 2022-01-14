package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) ReAssignStakeSuccess() {
	Start("Transfer Stake Success")

	p := tr.DeployAndUpgrade()

	reassignStakeSetupSSN(p)

	key1, addr1, ssn1, ssn2, _, userStake := reassignStakeDefineParams(p)
	totaltokenamount := Field(p.Aimpl, "totaltokenamount")
	totalstakeamount := Field(p.Aimpl, "totalstakeamount")
	userStakeThroughAimpl := userStake

	//key1 delegates through Aimpl
	AssertSuccess(p.Aimpl.WithUser(key1).DelegateStake(userStakeThroughAimpl))
	reassignStakeNextRewardCycle(p)
	reassignStakeNextRewardCycle(p)
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))

	depositAmtDeleg := Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, sdk.Cfg.AzilSsnAddress)
	ssnDelegAmt := Field(p.Zimpl, "ssn_deleg_amt", sdk.Cfg.AzilSsnAddress, p.Holder.Addr)

	//key1 delegates to main contract
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn1, userStake))
	userStake2 := StrAdd(userStake, userStake)
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn2, userStake2))

	//key1 waits 2 reward cycles
	reassignStakeNextRewardCycle(p)
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))
	reassignStakeNextRewardCycle(p)
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))

	//key1 claims rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn1))
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn2))

	//key1 withdraws some amount, then requests swap with holder
	tx, _ := AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(p.Holder.Addr))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": p.Holder.Addr}})

	//call ReAssignStake, expecting success
	tx, _ = AssertSuccess(p.Aimpl.WithUser(key1).ReAssignStake(addr1))
	AssertEvent(tx, Event{p.Zimpl.Addr, "ConfirmDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": p.Holder.Addr}})
	AssertTransition(tx, Transition{
		p.Zimpl.Addr, //sender
		"ReDelegateStakeSuccessCallBack",
		p.Holder.Addr,
		"0",
		ParamsMap{"ssnaddr": ssn1, "tossn": sdk.Cfg.AzilSsnAddress, "amount": userStake},
	})
	AssertTransition(tx, Transition{
		p.Zimpl.Addr, //sender
		"ReDelegateStakeSuccessCallBack",
		p.Holder.Addr,
		"0",
		ParamsMap{"ssnaddr": ssn2, "tossn": sdk.Cfg.AzilSsnAddress, "amount": userStake2},
	})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", addr1), "")
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", p.Holder.Addr, sdk.Cfg.AzilSsnAddress), StrAdd(depositAmtDeleg, userStake, userStake2))
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", sdk.Cfg.AzilSsnAddress, p.Holder.Addr), StrAdd(ssnDelegAmt, userStake, userStake2))
	AssertEqual(Field(p.Aimpl, "totalstakeamount"), StrAdd(totalstakeamount, userStakeThroughAimpl, userStake, userStake2))
	AssertEqual(Field(p.Aimpl, "totaltokenamount"), StrAdd(totaltokenamount, Field(p.Aimpl, "balances", addr1)))
}

func (tr *Transitions) ReAssignStakeAimplErrors() {
	Start("Transfer Stake Aimpl Errors")

	p := tr.DeployAndUpgrade()

	reassignStakeSetupSSN(p)

	key1, addr1, ssn1, ssn2, _, userStake := reassignStakeDefineParams(p)

	//key1 delegates to main contract
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn1, userStake))

	//key1 waits 2 reward cycles
	reassignStakeNextRewardCycle(p)
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))
	reassignStakeNextRewardCycle(p)
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))

	//key1 claims rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn1))

	//call ReAssignStake for addr1, expecting error
	tx, _ := p.Aimpl.WithUser(key1).ReAssignStake(addr1)
	AssertError(tx, "ReAssignStakeSwapRequestNotFound")

	//key1 requests swap with NOT Holder address
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(ssn2))

	//call ReAssignStake for addr1, expecting error
	tx, _ = p.Aimpl.WithUser(key1).ReAssignStake(addr1)
	AssertError(tx, "ReAssignStakeSwapRequestNotHolder")

	//key1 withdraws some amount, then requests swap with holder
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeAmt(ssn1, userStake))
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(p.Holder.Addr))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": p.Holder.Addr}})

	//call ReAssignStake for addr1, expecting error
	tx, _ = p.Aimpl.WithUser(key1).ReAssignStake(addr1)
	AssertError(tx, "ReAssignStakePendingWithdrawal")
}

func (tr *Transitions) ReAssignStakeZimplErrors() {
	Start("Transfer Stake Zimpl Errors")

	p := tr.DeployAndUpgrade()

	reassignStakeSetupSSN(p)
	key1, addr1, ssn1, _, _, userStake := reassignStakeDefineParams(p)

	//key1 delegates to main contract, expecting success
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn1, userStake))

	//key1 requests delegator swap, but he has buffered deposit, expecting DelegHasBufferedDeposit
	tx, _ := p.Zproxy.RequestDelegatorSwap(p.Holder.Addr)
	AssertZimplError(tx, -8)

	reassignStakeNextRewardCycle(p)
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))

	//key1 requests delegator swap, but he has buffered deposit in previous cycle, expecting DelegHasBufferedDeposit
	tx, _ = p.Zproxy.WithUser(key1).RequestDelegatorSwap(p.Holder.Addr)
	AssertZimplError(tx, -8)

	reassignStakeNextRewardCycle(p)

	//key1 requests delegator swap, but he has unclaimed rewards, expecting DelegHasUnwithdrawRewards
	tx, _ = p.Zproxy.WithUser(key1).RequestDelegatorSwap(p.Holder.Addr)
	AssertZimplError(tx, -12)

	//key1 claims rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn1))

	//key1 requests delegator swap, but Holder has unclaimed rewards, expecting DelegHasUnwithdrawRewards
	//workflow of this use case is: Verifier->AssignStakeReward, User->RequestDelegatorSwap, Aimpl->DrainBuffer
	tx, _ = p.Zproxy.WithUser(key1).RequestDelegatorSwap(p.Holder.Addr)
	AssertZimplError(tx, -12)

	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))

	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(p.Holder.Addr))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": p.Holder.Addr}})
}

func reassignStakeDefineParams(p *contracts.Protocol) (string, string, string, string, string, string) {
	key1 := sdk.Cfg.Key1
	addr1 := sdk.Cfg.Addr1
	ssn1 := "0x0000000000000000000000000000000000000001"
	ssn2 := "0x0000000000000000000000000000000000000002"
	minStake := Field(p.Zimpl, "minstake")
	userStake := ToZil(10)
	return key1, addr1, ssn1, ssn2, minStake, userStake
}

func reassignStakeSetupSSN(p *contracts.Protocol) {
	_, _, ssn1, ssn2, minStake, _ := reassignStakeDefineParams(p)

	prevWallet := p.Zproxy.Contract.Wallet

	//add test SSNs to main staking contract
	p.Zproxy.UpdateWallet(sdk.Cfg.AdminKey)
	AssertSuccess(p.Zproxy.AddSSN(ssn1, "SSN 1"))
	AssertSuccess(p.Zproxy.AddSSN(ssn2, "SSN 2"))
	AssertSuccess(p.Zproxy.DelegateStake(ssn1, minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn2, minStake))

	p.Zproxy.Contract.Wallet = prevWallet

	//ssns will become active on the next cycle
	reassignStakeNextRewardCycle(p)
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))
}

func reassignStakeNextRewardCycle(p *contracts.Protocol) {
	_, _, ssn1, ssn2, _, _ := reassignStakeDefineParams(p)
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
