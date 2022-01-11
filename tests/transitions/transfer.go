package transitions

import (
    //. "github.com/avely-finance/avely-contracts/sdk/utils"
    . "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Transfer() {
    Start("Transfer")

    p := tr.DeployAndUpgrade()

    minStake := p.Zimpl.Field("minstake")
    ssn1 := "0x0000000000000000000000000000000000000001"
    ssn2 := "0x0000000000000000000000000000000000000002"
    ssn3 := "0x0000000000000000000000000000000000000003"

    //add test SSNs to main staking contract
    p.Zproxy.UpdateWallet(sdk.Cfg.AdminKey)
    AssertSuccess(p.Zproxy.AddSSN(ssn1, "SSN 1"))
    AssertSuccess(p.Zproxy.AddSSN(ssn2, "SSN 2"))
    AssertSuccess(p.Zproxy.AddSSN(ssn3, "SSN 3"))
    AssertSuccess(p.Zproxy.DelegateStake(ssn1, minStake))
    AssertSuccess(p.Zproxy.DelegateStake(ssn2, minStake))
    AssertSuccess(p.Zproxy.DelegateStake(ssn3, minStake))

    //assign rewards, ssns will become active
    p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)
    AssertSuccess(p.Zproxy.AssignStakeReward(ssn1, "0"))
    AssertSuccess(p.Zproxy.AssignStakeReward(ssn2, "0"))
    AssertSuccess(p.Zproxy.AssignStakeReward(ssn3, "0"))

    //user delegates to main contract

    //user tries to transfer, but he has unclaimed rewards, expecting errors

    //user claims rewards, expecting success

    //user tries to transfer, but holder has rewards, expecting errors

    //holder claims rewards

    //user tries to transfer, expecting success

    //user have withdraw requests and tries to transfer, expecting RejectDelegatorSwap

    //user tries to transfer, expecting success
    //compare previous user's balance at Zimpl with user's current balance at Aimpl

    //tests for user, which was previously registered at Azil

}
