package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Transfer() {
	Start("Transfer")

	p := tr.DeployAndUpgrade()

	minStake := p.Zimpl.Field("minstake")
	userStake := ToZil(10)
	ssn1 := "0x0000000000000000000000000000000000000001"
	ssn2 := "0x0000000000000000000000000000000000000002"
	ssn3 := "0x0000000000000000000000000000000000000003"
	key1 := sdk.Cfg.Key1
	AdminKey := sdk.Cfg.AdminKey

	nextRewardCycle := func() {
		p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)
		AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))
		AssertSuccess(p.Zproxy.AssignStakeReward(ssn1, "10"))
		AssertSuccess(p.Zproxy.AssignStakeReward(ssn2, "10"))
		AssertSuccess(p.Zproxy.AssignStakeReward(ssn3, "10"))
		AssertSuccess(p.Aimpl.Key(AdminKey).DrainBuffer(p.GetBuffer().Addr))
	}

	//add test SSNs to main staking contract
	p.Zproxy.UpdateWallet(sdk.Cfg.AdminKey)
	AssertSuccess(p.Zproxy.AddSSN(ssn1, "SSN 1"))
	AssertSuccess(p.Zproxy.AddSSN(ssn2, "SSN 2"))
	AssertSuccess(p.Zproxy.AddSSN(ssn3, "SSN 3"))
	AssertSuccess(p.Zproxy.DelegateStake(ssn1, minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn2, minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn3, minStake))

	//assign rewards, ssns will become active
	nextRewardCycle()

	//user delegates to main contract
	AssertSuccess(p.Zproxy.Key(key1).DelegateStake(ssn1, userStake))

	//user tries to request delegator swap, but he has buffered deposit, expecting Zimpl error
	tx, _ := p.Zproxy.RequestDelegatorSwap(p.Holder.Addr)
	AssertZimplError(tx, -8) //DelegHasBufferedDeposit

	nextRewardCycle()
	nextRewardCycle()

	//user tries to request delegator swap, but he has buffered deposit, expecting Zimpl error
	tx, _ = p.Zproxy.Key(key1).RequestDelegatorSwap(p.Holder.Addr)
	AssertZimplError(tx, -12) //DelegHasUnwithdrawRewards

	//key1 claim rewards
	AssertSuccess(p.Zproxy.Key(key1).WithdrawStakeRewards(ssn1))
	AssertSuccess(p.Zproxy.Key(key1).RequestDelegatorSwap(p.Holder.Addr))

	//user tries to transfer, but he has unclaimed rewards, expecting Zimpl errors

	//user claims rewards, expecting success

	//user tries to transfer, but holder has rewards, expecting errors

	//holder claims rewards

	//user tries to transfer, expecting success
	//e = { _eventname: "RequestDelegatorSwap"; initial_deleg: initiator; new_deleg: new_deleg_addr };

	//user have withdraw requests and tries to transfer, expecting RejectDelegatorSwap

	//user tries to transfer, expecting success
	//p.Aimpl.UpdateWallet(key1)
	//AssertSuccess(p.Aimpl.CompleteTransfer(ssn1, userStake))

	//case when user called CompleteTransfer without initial call of RequestDelegatorSwap, or requested swap with otner address (not holder)

	//compare previous user's balance at Zimpl with user's current balance at Aimpl

	//tests for user, which was previously registered at Azil

}
