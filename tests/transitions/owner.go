package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Owner() {
	Start("StZIL & Holder contract owner transitions")

	p := tr.DeployAndUpgrade()
	p.StZIL.SetSigner(celestials.Owner)
	p.Holder.SetSigner(celestials.Owner)
	p.Buffers[0].SetSigner(celestials.Owner)

	checkChangeAdmin(p)
	checkChangeBuffersEmpty(p)
	checkAddSSNExists(p)
	checkSetHolderAddress(p)
	checkChangeRewardsFee(p)
	checkChangeTreasuryAddress(p)
	checkChangeZimplAddress(p)
	checkChangeZproxyAddress(p)
	checkUpdateStakingParameters(p)
	//this test is last because all SSNs will be un-whitelisted
	checkRemoveSSN(p)
	changeOwner(p)
}

func checkChangeAdmin(p *contracts.Protocol) {
	newAdminAddr := utils.GetAddressByWallet(eve)

	//change admin, expecting success
	p.StZIL.SetSigner(celestials.Owner)
	tx, _ := AssertSuccess(p.StZIL.ChangeAdmin(newAdminAddr))
	AssertEvent(tx, Event{
		Sender:    p.StZIL.Addr,
		EventName: "ChangeAdmin",
		Params:    ParamsMap{"old_admin": utils.GetAddressByWallet(celestials.Admin), "new_admin": newAdminAddr},
	})
	AssertEqual(Field(p.StZIL, "admin_address"), newAdminAddr)
}

func checkChangeBuffersEmpty(p *contracts.Protocol) {
	new_buffers := []string{}
	tx, _ := p.StZIL.ChangeBuffers(new_buffers)
	AssertError(tx, p.StZIL.ErrorCode("BuffersEmpty"))
}

func checkAddSSNExists(p *contracts.Protocol) {
	ssnlist := p.StZIL.GetSsnWhitelist()
	tx, _ := p.StZIL.AddSSN(ssnlist[0])
	AssertError(tx, p.StZIL.ErrorCode("SsnAddressExists"))
}

func checkRemoveSSN(p *contracts.Protocol) {
	ssnlist := p.StZIL.GetSsnWhitelist()
	//unwhitelist all ssn addresses except first one (zero-indexed)
	for i := 1; i < len(ssnlist); i++ {
		AssertSuccess(p.StZIL.RemoveSSN(ssnlist[i]))
	}
	//try to remove last whitelisted SSN, expect error
	tx, _ := p.StZIL.RemoveSSN(ssnlist[0])
	AssertError(tx, p.StZIL.ErrorCode("SsnAddressesEmpty"))

	//remove last whitelisted SSN on paused contract, expect success
	AssertSuccess(p.StZIL.PauseIn())
	AssertSuccess(p.StZIL.RemoveSSN(ssnlist[0]))
}

func checkSetHolderAddress(p *contracts.Protocol) {
	AssertEqual(Field(p.StZIL, "holder_address"), p.Holder.Addr)
	tx, _ := p.StZIL.SetHolderAddress(core.ZeroAddr)
	AssertError(tx, p.StZIL.ErrorCode("HolderAlreadySet"))
}

func checkChangeRewardsFee(p *contracts.Protocol) {
	prevValue := Field(p.StZIL, "rewards_fee")
	//try to change fee, expecting error, because fee_denom=10000
	tx, _ := p.StZIL.ChangeRewardsFee("12345")
	AssertError(tx, p.StZIL.ErrorCode("InvalidRewardsFee"))
	goodValue := "2345"
	AssertSuccess(p.StZIL.ChangeRewardsFee(goodValue))
	AssertEqual(Field(p.StZIL, "rewards_fee"), goodValue)
	AssertSuccess(p.StZIL.ChangeRewardsFee(prevValue))
}

func checkChangeTreasuryAddress(p *contracts.Protocol) {
	AssertSuccess(p.StZIL.ChangeTreasuryAddress(core.ZeroAddr))
	AssertEqual(Field(p.StZIL, "treasury_address"), core.ZeroAddr)
	AssertSuccess(p.StZIL.ChangeTreasuryAddress(p.Treasury.Addr))
}

func checkChangeZimplAddress(p *contracts.Protocol) {
	zimplAddr := p.Zimpl.Addr

	// stZIL
	tx, _ := AssertSuccess(p.StZIL.ChangeZimplAddress(core.ZeroAddr))
	AssertEvent(tx, Event{p.StZIL.Addr, "ChangeZimplAddress", ParamsMap{"address": core.ZeroAddr}})
	AssertEqual(Field(p.StZIL, "zimpl_address"), core.ZeroAddr)
	AssertSuccess(p.StZIL.ChangeZimplAddress(zimplAddr))

	// Holder
	tx, _ = AssertSuccess(p.Holder.ChangeZimplAddress(core.ZeroAddr))
	AssertEvent(tx, Event{p.Holder.Addr, "ChangeZimplAddress", ParamsMap{"address": core.ZeroAddr}})
	AssertEqual(Field(p.Holder, "zimpl_address"), core.ZeroAddr)
	AssertSuccess(p.Holder.ChangeZimplAddress(zimplAddr))
}

func checkChangeZproxyAddress(p *contracts.Protocol) {
	zproxyAddr := p.Zproxy.Addr

	tx, _ := AssertSuccess(p.Holder.ChangeZproxyAddress(core.ZeroAddr))
	AssertEvent(tx, Event{p.Holder.Addr, "ChangeZproxyAddress", ParamsMap{"address": core.ZeroAddr}})
	AssertEqual(Field(p.Holder, "zproxy_address"), core.ZeroAddr)
	AssertSuccess(p.Holder.ChangeZproxyAddress(zproxyAddr))
}

func checkUpdateStakingParameters(p *contracts.Protocol) {
	prevValue := Field(p.StZIL, "mindelegstake")
	testValue := utils.ToZil(54321)
	tx, _ := AssertSuccess(p.StZIL.UpdateStakingParameters(testValue))
	AssertEvent(tx, Event{p.StZIL.Addr, "UpdateStakingParameters", ParamsMap{"min_deleg_stake": testValue}})
	AssertEqual(Field(p.StZIL, "mindelegstake"), testValue)
	AssertSuccess(p.StZIL.UpdateStakingParameters(prevValue))
}

func changeOwner(p *contracts.Protocol) {
	// stZIL
	newOwner := eve
	newOwnerAddr := utils.GetAddressByWallet(eve)

	//claim not existent staging owner, expecting error
	p.StZIL.SetSigner(newOwner)
	tx, _ := p.StZIL.ClaimOwner()
	AssertError(tx, p.StZIL.ErrorCode("StagingOwnerNotExists"))

	//change owner, expecting success
	p.StZIL.SetSigner(celestials.Owner)
	tx, _ = AssertSuccess(p.StZIL.ChangeOwner(newOwnerAddr))
	AssertEvent(tx, Event{p.StZIL.Addr, "ChangeOwner", ParamsMap{"current_owner": utils.GetAddressByWallet(celestials.Owner), "new_owner": newOwnerAddr}})
	AssertEqual(Field(p.StZIL, "staging_owner_address"), newOwnerAddr)

	//claim owner with wrong user, expecting error
	p.StZIL.SetSigner(alice)
	tx, _ = p.StZIL.ClaimOwner()
	AssertError(tx, p.StZIL.ErrorCode("StagingOwnerValidationFailed"))

	//claim owner with correct user, expecting success
	p.StZIL.SetSigner(newOwner)
	tx, _ = AssertSuccess(p.StZIL.ClaimOwner())
	AssertEvent(tx, Event{p.StZIL.Addr, "ClaimOwner", ParamsMap{"new_owner": newOwnerAddr}})
	AssertEqual(Field(p.StZIL, "owner_address"), newOwnerAddr)
	AssertEqual(Field(p.StZIL, "staging_owner_address"), "")

	// Holder

	//claim not existent staging owner, expecting error
	p.Holder.SetSigner(newOwner)
	tx, _ = p.Holder.ClaimOwner()
	AssertError(tx, p.Holder.ErrorCode("StagingOwnerNotExists"))

	//change owner, expecting success
	p.Holder.SetSigner(celestials.Owner)
	tx, _ = AssertSuccess(p.Holder.ChangeOwner(newOwnerAddr))
	AssertEvent(tx, Event{p.Holder.Addr, "ChangeOwner", ParamsMap{"current_owner": utils.GetAddressByWallet(celestials.Owner), "new_owner": newOwnerAddr}})
	AssertEqual(Field(p.Holder, "staging_owner_address"), newOwnerAddr)

	//claim owner with wrong user, expecting error
	p.Holder.SetSigner(alice)
	tx, _ = p.Holder.ClaimOwner()
	AssertError(tx, p.Holder.ErrorCode("StagingOwnerValidationFailed"))

	//claim owner with correct user, expecting success
	p.Holder.SetSigner(newOwner)
	tx, _ = AssertSuccess(p.Holder.ClaimOwner())
	AssertEvent(tx, Event{p.Holder.Addr, "ClaimOwner", ParamsMap{"new_owner": newOwnerAddr}})
	AssertEqual(Field(p.Holder, "owner_address"), newOwnerAddr)
	AssertEqual(Field(p.Holder, "staging_owner_address"), "")
}
