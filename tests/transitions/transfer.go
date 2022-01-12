package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Transfer() {
	Start("Transfer")

	p := tr.DeployAndUpgrade()

	transferSetupSSN(p)
	transferNextRewardCycle(p)
	transferErrors(p)

	//key1 tries to transfer, expecting success
	//e = { _eventname: "RequestDelegatorSwap"; initial_deleg: initiator; new_deleg: new_deleg_addr };

	//key1 have withdraw requests and tries to transfer, expecting RejectDelegatorSwap

	//key1 tries to transfer, expecting success
	//p.Aimpl.UpdateWallet(key1)
	//AssertSuccess(p.Aimpl.CompleteTransfer(ssn1, userStake))

	//case when key1 called CompleteTransfer without initial call of RequestDelegatorSwap, or requested swap with otner address (not holder)

	//compare previous user's balance at Zimpl with user's current balance at Aimpl

	//tests for user, which was previously registered at Azil

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

func transferErrors(p *contracts.Protocol) {
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
