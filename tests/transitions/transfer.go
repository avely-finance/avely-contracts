package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) TransferSuccess() {
	Start("Transfer Success")

	p := tr.DeployAndUpgrade()

	transferSetupSSN(p)

	key1, addr1, ssn1, ssn2, _, userStake := transferDefineParams(p)

	//key1 delegates to main contract
	AssertSuccess(p.Zproxy.Key(key1).DelegateStake(ssn1, userStake))
	AssertSuccess(p.Zproxy.Key(key1).DelegateStake(ssn2, StrAdd(userStake, userStake)))

	//key1 waits 2 reward cycles
	transferNextRewardCycle(p)
	AssertSuccess(p.Aimpl.Key(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))
	transferNextRewardCycle(p)
	AssertSuccess(p.Aimpl.Key(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))

	//key1 claims rewards
	AssertSuccess(p.Zproxy.Key(key1).WithdrawStakeRewards(ssn1))
	AssertSuccess(p.Zproxy.Key(key1).WithdrawStakeRewards(ssn2))

	//key1 withdraws some amount, then requests swap with holder
	tx, _ := AssertSuccess(p.Zproxy.Key(key1).RequestDelegatorSwap(p.Holder.Addr))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": p.Holder.Addr}})

	//call CompleteTransfer, expecting success
	tx, _ = AssertSuccess(p.Aimpl.Key(key1).CompleteTransfer(addr1))
	AssertEvent(tx, Event{p.Zimpl.Addr, "ConfirmDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": p.Holder.Addr}})
	AssertEqual(p.Zimpl.Field("deposit_amt_deleg", p.Holder.Addr, ssn1), userStake)
	AssertEqual(p.Zimpl.Field("deposit_amt_deleg", p.Holder.Addr, ssn2), StrAdd(userStake, userStake))
	AssertEqual(p.Zimpl.Field("ssn_deleg_amt", ssn1, p.Holder.Addr), userStake)
	AssertEqual(p.Zimpl.Field("ssn_deleg_amt", ssn2, p.Holder.Addr), StrAdd(userStake, userStake))

	////////////////////////AssertEvent(tx, Event{p.Aimpl.Addr, "CompleteTransfer", ParamsMap{"initial_deleg": addr1, "new_deleg": p.Holder.Addr}})

	//compare previous user's balance at Zimpl with user's current balance at Aimpl

	//tests for user, which was previously registered at Azil

	//what if user has some redelegate requests?

	//compare azil/holder states for transfered depo with depo, initially delegated throw Azil

}

func (tr *Transitions) TransferAimplErrors() {
	Start("Transfer Aimpl Errors")

	p := tr.DeployAndUpgrade()

	transferSetupSSN(p)

	key1, addr1, ssn1, ssn2, _, userStake := transferDefineParams(p)

	//key1 delegates to main contract
	AssertSuccess(p.Zproxy.Key(key1).DelegateStake(ssn1, userStake))

	//key1 waits 2 reward cycles
	transferNextRewardCycle(p)
	AssertSuccess(p.Aimpl.Key(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))
	transferNextRewardCycle(p)
	AssertSuccess(p.Aimpl.Key(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))

	//key1 claims rewards
	AssertSuccess(p.Zproxy.Key(key1).WithdrawStakeRewards(ssn1))

	//call CompleteTransfer for addr1, expecting error
	tx, _ := p.Aimpl.Key(key1).CompleteTransfer(addr1)
	AssertError(tx, "CompleteTransferSwapRequestNotFound")

	//key1 requests swap with NOT Holder address
	tx, _ = AssertSuccess(p.Zproxy.Key(key1).RequestDelegatorSwap(ssn2))

	//call CompleteTransfer for addr1, expecting error
	tx, _ = p.Aimpl.Key(key1).CompleteTransfer(addr1)
	AssertError(tx, "CompleteTransferSwapRequestNotHolder")

	//key1 withdraws some amount, then requests swap with holder
	AssertSuccess(p.Zproxy.Key(key1).WithdrawStakeAmt(ssn1, userStake))
	tx, _ = AssertSuccess(p.Zproxy.Key(key1).RequestDelegatorSwap(p.Holder.Addr))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": p.Holder.Addr}})

	//call CompleteTransfer for addr1, expecting error
	tx, _ = p.Aimpl.Key(key1).CompleteTransfer(addr1)
	AssertError(tx, "CompleteTransferPendingWithdrawal")
}

func (tr *Transitions) TransferZimplErrors() {
	Start("Transfer Zimpl Errors")

	p := tr.DeployAndUpgrade()

	transferSetupSSN(p)
	transferZimplErrors(p)
}

func transferDefineParams(p *contracts.Protocol) (string, string, string, string, string, string) {
	key1 := sdk.Cfg.Key1
	addr1 := sdk.Cfg.Addr1
	ssn1 := "0x0000000000000000000000000000000000000001"
	ssn2 := "0x0000000000000000000000000000000000000002"
	minStake := p.Zimpl.Field("minstake")
	userStake := ToZil(10)
	return key1, addr1, ssn1, ssn2, minStake, userStake
}

func transferZimplErrors(p *contracts.Protocol) {
	key1, addr1, ssn1, _, _, userStake := transferDefineParams(p)

	//key1 delegates to main contract, expecting success
	AssertSuccess(p.Zproxy.Key(key1).DelegateStake(ssn1, userStake))

	//key1 requests delegator swap, but he has buffered deposit, expecting DelegHasBufferedDeposit
	tx, _ := p.Zproxy.RequestDelegatorSwap(p.Holder.Addr)
	AssertZimplError(tx, -8)

	transferNextRewardCycle(p)
	AssertSuccess(p.Aimpl.Key(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))

	//key1 requests delegator swap, but he has buffered deposit in previous cycle, expecting DelegHasBufferedDeposit
	tx, _ = p.Zproxy.Key(key1).RequestDelegatorSwap(p.Holder.Addr)
	AssertZimplError(tx, -8)

	transferNextRewardCycle(p)

	//key1 requests delegator swap, but he has unclaimed rewards, expecting DelegHasUnwithdrawRewards
	tx, _ = p.Zproxy.Key(key1).RequestDelegatorSwap(p.Holder.Addr)
	AssertZimplError(tx, -12)

	//key1 claims rewards
	AssertSuccess(p.Zproxy.Key(key1).WithdrawStakeRewards(ssn1))

	//key1 requests delegator swap, but Holder has unclaimed rewards, expecting DelegHasUnwithdrawRewards
	//workflow of this use case is: Verifier->AssignStakeReward, User->RequestDelegatorSwap, Aimpl->DrainBuffer
	tx, _ = p.Zproxy.Key(key1).RequestDelegatorSwap(p.Holder.Addr)
	AssertZimplError(tx, -12)

	AssertSuccess(p.Aimpl.Key(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))

	tx, _ = AssertSuccess(p.Zproxy.Key(key1).RequestDelegatorSwap(p.Holder.Addr))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": p.Holder.Addr}})
}

func transferSetupSSN(p *contracts.Protocol) {
	_, _, ssn1, ssn2, minStake, _ := transferDefineParams(p)

	prevWallet := p.Zproxy.Contract.Wallet

	//add test SSNs to main staking contract
	p.Zproxy.UpdateWallet(sdk.Cfg.AdminKey)
	AssertSuccess(p.Zproxy.AddSSN(ssn1, "SSN 1"))
	AssertSuccess(p.Zproxy.AddSSN(ssn2, "SSN 2"))
	AssertSuccess(p.Zproxy.DelegateStake(ssn1, minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn2, minStake))

	p.Zproxy.Contract.Wallet = prevWallet

	//ssns will become active on the next cycle
	transferNextRewardCycle(p)
	AssertSuccess(p.Aimpl.Key(sdk.Cfg.AdminKey).DrainBuffer(p.GetBuffer().Addr))
}

func transferNextRewardCycle(p *contracts.Protocol) {
	_, _, ssn1, ssn2, _, _ := transferDefineParams(p)
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
