package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) TransferStakeSuccess() {
	Start("Transfer Stake Success")

	p := tr.DeployAndUpgrade()

	transferStakeSetupSSN(p)

	key1, addr1, ssn1, ssn2, _, userStake := transferStakeDefineParams(p)
	totaltokenamount := p.Aimpl.Field("totaltokenamount")
	totalstakeamount := p.Aimpl.Field("totalstakeamount")
	userStakeThroughAimpl := userStake

	//key1 delegates through Aimpl
	AssertSuccess(p.Aimpl.WithUser(key1).DelegateStake(userStakeThroughAimpl))
	transferStakeNextRewardCycle(p)
	transferStakeNextRewardCycle(p)
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))

	depositAmtDeleg := p.Zimpl.Field("deposit_amt_deleg", p.Holder.Addr, sdk.Cfg.AzilSsnAddress)
	ssnDelegAmt := p.Zimpl.Field("ssn_deleg_amt", sdk.Cfg.AzilSsnAddress, p.Holder.Addr)

	//key1 delegates to main contract
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn1, userStake))
	userStake2 := StrAdd(userStake, userStake)
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn2, userStake2))

	//key1 waits 2 reward cycles
	transferStakeNextRewardCycle(p)
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))
	transferStakeNextRewardCycle(p)
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))

	//key1 claims rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn1))
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn2))

	//key1 withdraws some amount, then requests swap with holder
	tx, _ := AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(p.Holder.Addr))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": p.Holder.Addr}})

	//call CompleteTransfer, expecting success
	tx, _ = AssertSuccess(p.Aimpl.WithUser(key1).CompleteTransfer(addr1))
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
	AssertEqual(p.Zimpl.Field("deposit_amt_deleg", addr1), "")
	AssertEqual(p.Zimpl.Field("deposit_amt_deleg", p.Holder.Addr, sdk.Cfg.AzilSsnAddress), StrAdd(depositAmtDeleg, userStake, userStake2))
	AssertEqual(p.Zimpl.Field("ssn_deleg_amt", sdk.Cfg.AzilSsnAddress, p.Holder.Addr), StrAdd(ssnDelegAmt, userStake, userStake2))
	AssertEqual(p.Aimpl.Field("totalstakeamount"), StrAdd(totalstakeamount, userStakeThroughAimpl, userStake, userStake2))
	AssertEqual(p.Aimpl.Field("totaltokenamount"), StrAdd(totaltokenamount, p.Aimpl.Field("balances", addr1)))
}

func (tr *Transitions) TransferStakeAimplErrors() {
	Start("Transfer Stake Aimpl Errors")

	p := tr.DeployAndUpgrade()

	transferStakeSetupSSN(p)

	key1, addr1, ssn1, ssn2, _, userStake := transferStakeDefineParams(p)

	//key1 delegates to main contract
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn1, userStake))

	//key1 waits 2 reward cycles
	transferStakeNextRewardCycle(p)
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))
	transferStakeNextRewardCycle(p)
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))

	//key1 claims rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn1))

	//call CompleteTransfer for addr1, expecting error
	tx, _ := p.Aimpl.WithUser(key1).CompleteTransfer(addr1)
	AssertError(tx, "CompleteTransferSwapRequestNotFound")

	//key1 requests swap with NOT Holder address
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(ssn2))

	//call CompleteTransfer for addr1, expecting error
	tx, _ = p.Aimpl.WithUser(key1).CompleteTransfer(addr1)
	AssertError(tx, "CompleteTransferSwapRequestNotHolder")

	//key1 withdraws some amount, then requests swap with holder
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeAmt(ssn1, userStake))
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(p.Holder.Addr))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": p.Holder.Addr}})

	//call CompleteTransfer for addr1, expecting error
	tx, _ = p.Aimpl.WithUser(key1).CompleteTransfer(addr1)
	AssertError(tx, "CompleteTransferPendingWithdrawal")
}

func (tr *Transitions) TransferStakeZimplErrors() {
	Start("Transfer Stake Zimpl Errors")

	p := tr.DeployAndUpgrade()

	transferStakeSetupSSN(p)
	key1, addr1, ssn1, _, _, userStake := transferStakeDefineParams(p)

	//key1 delegates to main contract, expecting success
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn1, userStake))

	//key1 requests delegator swap, but he has buffered deposit, expecting DelegHasBufferedDeposit
	tx, _ := p.Zproxy.RequestDelegatorSwap(p.Holder.Addr)
	AssertZimplError(tx, -8)

	transferStakeNextRewardCycle(p)
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))

	//key1 requests delegator swap, but he has buffered deposit in previous cycle, expecting DelegHasBufferedDeposit
	tx, _ = p.Zproxy.WithUser(key1).RequestDelegatorSwap(p.Holder.Addr)
	AssertZimplError(tx, -8)

	transferStakeNextRewardCycle(p)

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

func transferStakeDefineParams(p *contracts.Protocol) (string, string, string, string, string, string) {
	key1 := sdk.Cfg.Key1
	addr1 := sdk.Cfg.Addr1
	ssn1 := "0x0000000000000000000000000000000000000001"
	ssn2 := "0x0000000000000000000000000000000000000002"
	minStake := p.Zimpl.Field("minstake")
	userStake := ToZil(10)
	return key1, addr1, ssn1, ssn2, minStake, userStake
}

func transferStakeSetupSSN(p *contracts.Protocol) {
	_, _, ssn1, ssn2, minStake, _ := transferStakeDefineParams(p)

	prevWallet := p.Zproxy.Contract.Wallet

	//add test SSNs to main staking contract
	p.Zproxy.UpdateWallet(sdk.Cfg.AdminKey)
	AssertSuccess(p.Zproxy.AddSSN(ssn1, "SSN 1"))
	AssertSuccess(p.Zproxy.AddSSN(ssn2, "SSN 2"))
	AssertSuccess(p.Zproxy.DelegateStake(ssn1, minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn2, minStake))

	p.Zproxy.Contract.Wallet = prevWallet

	//ssns will become active on the next cycle
	transferStakeNextRewardCycle(p)
	AssertSuccess(p.Aimpl.WithUser(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))
}

func transferStakeNextRewardCycle(p *contracts.Protocol) {
	_, _, ssn1, ssn2, _, _ := transferStakeDefineParams(p)
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
