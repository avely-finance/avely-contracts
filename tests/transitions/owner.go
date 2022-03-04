package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Owner() {

	Start("Azil contract owner transitions")

	p := tr.DeployAndUpgrade()
	p.Azil.UpdateWallet(sdk.Cfg.OwnerKey)

	checkChangeAdmin(p)
	checkChangeAzilSSNAddress(p)
	checkChangeBuffersEmpty(p)
	checkSetHolderAddress(p)
	checkChangeRewardsFee(p)
	checkChangeTreasuryAddress(p)
	checkChangeZimplAddress(p)
	checkUpdateStakingParameters(p)

	newOwnerAddr := sdk.Cfg.Addr3
	newOwnerKey := sdk.Cfg.Key3

	//claim not existent staging owner, expecting error
	p.Azil.UpdateWallet(newOwnerKey)
	tx, _ := p.Azil.ClaimOwner()
	AssertError(tx, "StagingOwnerNotExists")

	//change owner, expecting success
	p.Azil.UpdateWallet(sdk.Cfg.OwnerKey)
	tx, _ = AssertSuccess(p.Azil.ChangeOwner(newOwnerAddr))
	AssertEvent(tx, Event{p.Azil.Addr, "ChangeOwner", ParamsMap{"current_owner": sdk.Cfg.Owner, "new_owner": newOwnerAddr}})
	AssertEqual(Field(p.Azil, "staging_owner_address"), newOwnerAddr)

	//claim owner with wrong user, expecting error
	wrongActor := sdk.Cfg.Key1
	p.Azil.UpdateWallet(wrongActor)
	tx, _ = p.Azil.ClaimOwner()
	AssertError(tx, "StagingOwnerValidationFailed")

	//claim owner with correct user, expecting success
	p.Azil.UpdateWallet(newOwnerKey)
	tx, _ = AssertSuccess(p.Azil.ClaimOwner())
	AssertEvent(tx, Event{p.Azil.Addr, "ClaimOwner", ParamsMap{"new_owner": newOwnerAddr}})
	AssertEqual(Field(p.Azil, "owner_address"), newOwnerAddr)
	AssertEqual(Field(p.Azil, "staging_owner_address"), "")
}

func checkChangeAdmin(p *contracts.Protocol) {
	newAdminAddr := sdk.Cfg.Addr3

	//change admin, expecting success
	p.Azil.UpdateWallet(sdk.Cfg.OwnerKey)
	tx, _ := AssertSuccess(p.Azil.ChangeAdmin(newAdminAddr))
	AssertEvent(tx, Event{
		Sender:    p.Azil.Addr,
		EventName: "ChangeAdmin",
		Params:    ParamsMap{"old_admin": sdk.Cfg.Admin, "new_admin": newAdminAddr},
	})
	AssertEqual(Field(p.Azil, "admin_address"), newAdminAddr)
}

func checkChangeAzilSSNAddress(p *contracts.Protocol) {
	tx, _ := AssertSuccess(p.Azil.ChangeAzilSSNAddress(core.ZeroAddr))
	AssertEvent(tx, Event{p.Azil.Addr, "ChangeAzilSSNAddress", ParamsMap{"address": core.ZeroAddr}})
	AssertEqual(Field(p.Azil, "azil_ssn_address"), core.ZeroAddr)
	AssertSuccess(p.Azil.ChangeAzilSSNAddress(sdk.Cfg.AzilSsnAddress))
}

func checkChangeBuffersEmpty(p *contracts.Protocol) {
	new_buffers := []string{}
	tx, _ := p.Azil.ChangeBuffers(new_buffers)
	AssertError(tx, "BuffersEmpty")
}

func checkSetHolderAddress(p *contracts.Protocol) {
	AssertEqual(Field(p.Azil, "holder_address"), p.Holder.Addr)
	tx, _ := p.Azil.SetHolderAddress(core.ZeroAddr)
	AssertError(tx, "HolderAlreadySet")
}

func checkChangeRewardsFee(p *contracts.Protocol) {
	prevValue := Field(p.Azil, "rewards_fee")
	//try to change fee, expecting error, because fee_denom=10000
	tx, _ := p.Azil.ChangeRewardsFee("12345")
	AssertError(tx, "InvalidRewardsFee")
	goodValue := "2345"
	AssertSuccess(p.Azil.ChangeRewardsFee(goodValue))
	AssertEqual(Field(p.Azil, "rewards_fee"), goodValue)
	AssertSuccess(p.Azil.ChangeRewardsFee(prevValue))
}

func checkChangeTreasuryAddress(p *contracts.Protocol) {
	AssertSuccess(p.Azil.ChangeTreasuryAddress(core.ZeroAddr))
	AssertEqual(Field(p.Azil, "treasury_address"), core.ZeroAddr)
	AssertSuccess(p.Azil.ChangeTreasuryAddress(sdk.Cfg.AzilSsnAddress))
}

func checkChangeZimplAddress(p *contracts.Protocol) {
	zimplAddr := p.Zimpl.Addr
	tx, _ := AssertSuccess(p.Azil.ChangeZimplAddress(core.ZeroAddr))
	AssertEvent(tx, Event{p.Azil.Addr, "ChangeZimplAddress", ParamsMap{"address": core.ZeroAddr}})
	AssertEqual(Field(p.Azil, "zimpl_address"), core.ZeroAddr)
	AssertSuccess(p.Azil.ChangeZimplAddress(zimplAddr))
}

func checkUpdateStakingParameters(p *contracts.Protocol) {
	prevValue := Field(p.Azil, "mindelegstake")
	testValue := utils.ToZil(54321)
	tx, _ := AssertSuccess(p.Azil.UpdateStakingParameters(testValue))
	AssertEvent(tx, Event{p.Azil.Addr, "UpdateStakingParameters", ParamsMap{"min_deleg_stake": testValue}})
	AssertEqual(Field(p.Azil, "mindelegstake"), testValue)
	AssertSuccess(p.Azil.UpdateStakingParameters(prevValue))
}
