package transitions

import (
	"strconv"

	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"

	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Ssn() {
	ssnOperations(tr)
	ssnChangeOwner(tr)
	ssnRequireOwner(tr)
	ssnChangeZproxy(tr)
}

func ssnChangeOwner(tr *Transitions) {
	txIdLocal1 := 0
	txIdLocal2 := 0

	//deploy multisig
	owners := []string{utils.GetAddressByWallet(alice)}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	//deploy SSN, set owner to multisig contract
	init_owner := multisig.Addr
	init_zproxy := core.ZeroAddr
	ssn := tr.DeploySsn(init_owner, init_zproxy)

	//deploy other multisig contract
	newSignCount := 1
	newOwner := sdk.Cfg.Key2
	newOwners := []string{sdk.Cfg.Addr2}
	newMultisig := tr.DeployMultisigWallet(newOwners, newSignCount)

	//try to claim owner, expect error
	AssertMultisigSuccess(newMultisig.WithUser(newOwner).SubmitClaimOwnerTransaction(ssn.Addr))
	tx, _ := newMultisig.WithUser(newOwner).ExecuteTransaction(txIdLocal2)
	AssertError(tx, ssn.ErrorCode("StagingOwnerNotExists"))

	//initiate owner change
	multisig.SetSigner(alice)
	AssertMultisigSuccess(multisig.SubmitChangeOwnerTransaction(ssn.Addr, newMultisig.Addr))
	AssertMultisigSuccess(multisig.ExecuteTransaction(txIdLocal1))
	AssertEqual(Field(ssn, "staging_owner"), newMultisig.Addr)

	//try to claim owner with wrong user, expect error
	tx, _ = ssn.WithUser(sdk.Cfg.Key2).ClaimOwner()
	AssertError(tx, ssn.ErrorCode("StagingOwnerValidationFailed"))

	//claim owner
	txIdLocal2++
	AssertMultisigSuccess(newMultisig.WithUser(newOwner).SubmitClaimOwnerTransaction(ssn.Addr))
	AssertMultisigSuccess(newMultisig.WithUser(newOwner).ExecuteTransaction(txIdLocal2))
	AssertEqual(Field(ssn, "owner"), newMultisig.Addr)

}

func ssnRequireOwner(tr *Transitions) {

	Start("ssnRequireOwner")

	//deploy SSN
	init_owner := utils.GetAddressByWallet(celestials.Admin)
	init_zproxy := core.ZeroAddr
	ssn := tr.DeploySsn(init_owner, init_zproxy)

	// Use non-owner user, expecting errors
	ssn.UpdateWallet(sdk.Cfg.Key2)

	tx, _ := ssn.ChangeOwner(sdk.Cfg.Addr3)
	AssertError(tx, ssn.ErrorCode("OwnerValidationFailed"))

	tx, _ = ssn.ChangeZproxy(sdk.Cfg.Addr3)
	AssertError(tx, ssn.ErrorCode("OwnerValidationFailed"))

	tx, _ = ssn.UpdateReceivingAddr(sdk.Cfg.Addr3)
	AssertError(tx, ssn.ErrorCode("OwnerValidationFailed"))

	tx, _ = ssn.UpdateComm(12345)
	AssertError(tx, ssn.ErrorCode("OwnerValidationFailed"))

	tx, _ = ssn.WithdrawComm()
	AssertError(tx, ssn.ErrorCode("OwnerValidationFailed"))
}

func ssnChangeZproxy(tr *Transitions) {

	//deploy multisig
	owners := []string{utils.GetAddressByWallet(celestials.Owner)}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	//deploy ssn, set owner to multisig contract
	init_owner := multisig.Addr
	init_zproxy := core.ZeroAddr
	ssn := tr.DeploySsn(init_owner, init_zproxy)

	//deploy protocol
	p := tr.DeployAndUpgrade()

	txIdLocal := 0

	//change zproxy address, expect success
	AssertEqual(Field(ssn, "zproxy"), core.ZeroAddr)
	new_zproxy := p.Zproxy.Addr
	multisig.SetSigner(celestials.Owner)
	AssertMultisigSuccess(multisig.SubmitChangeZproxyTransaction(ssn.Addr, new_zproxy))
	tx, _ := AssertMultisigSuccess(multisig.ExecuteTransaction(txIdLocal))
	AssertEvent(tx, Event{
		ssn.Addr, //sender
		"ChangeZproxy",
		ParamsMap{"new_address": new_zproxy},
	})
	AssertEqual(Field(ssn, "zproxy"), new_zproxy)
}

func ssnOperations(tr *Transitions) {

	//deploy multisig
	owners := []string{utils.GetAddressByWallet(celestials.Owner)}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	//deploy protocol
	p := tr.DeployAndUpgrade()

	//deploy ssn, set owner to multisig contract
	init_owner := multisig.Addr
	init_zproxy := p.Zproxy.Addr
	ssn := tr.DeploySsn(init_owner, init_zproxy)

	addr := ssn.Addr
	name := ssn.Addr
	//add SSN to ssnlist-contract
	p.Zproxy.AddSSN(addr, name)
	//some amount should be delegated to SSN to make it active
	p.Zproxy.DelegateStake(addr, ToZil(sdk.Cfg.SsnInitialDelegateZil))
	//whitelist SSN
	p.StZIL.SetSigner(celestials.Owner)
	AssertSuccess(p.StZIL.AddSSN(addr))

	//expect SSN added successfully
	//Ssn: active_status stake_amt rewards name urlraw urlapi buffdeposit comm comm_rewards rec_addr
	//name
	AssertEqual(Field(p.Zimpl, "ssnlist", addr, "arguments", "3"), name)
	//rec_addr
	AssertEqual(Field(p.Zimpl, "ssnlist", addr, "arguments", "9"), addr)

	//change receiving address, expect success
	txIdLocal := 0
	treasury_addr := p.Treasury.Addr
	multisig.SetSigner(celestials.Owner)
	AssertMultisigSuccess(multisig.SubmitUpdateReceivingAddrTransaction(ssn.Addr, treasury_addr))
	tx, _ := AssertMultisigSuccess(multisig.ExecuteTransaction(txIdLocal))
	//data, _ := json.MarshalIndent(tx, "", "     ")
	//GetLog().Fatal(string(data))
	AssertEvent(tx, Event{
		p.Zimpl.Addr, //sender
		"UpdateReceivingAddr",
		ParamsMap{"ssn_addr": ssn.Addr, "new_addr": treasury_addr},
	})
	AssertEqual(Field(p.Zimpl, "ssnlist", addr, "arguments", "9"), treasury_addr)

	//change comission, expect success
	//change cycle: comission change allowed once per cycle
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	txIdLocal++
	//Ssn: active_status stake_amt rewards name urlraw urlapi buffdeposit comm comm_rewards rec_addr
	//comm, default value is 0
	AssertEqual(Field(p.Zimpl, "ssnlist", addr, "arguments", "7"), "0")
	new_rate := 123456789
	multisig.SetSigner(celestials.Owner)
	AssertMultisigSuccess(multisig.SubmitUpdateCommTransaction(ssn.Addr, new_rate))
	tx, _ = AssertMultisigSuccess(multisig.ExecuteTransaction(txIdLocal))
	AssertEvent(tx, Event{
		p.Zimpl.Addr, //sender
		"UpdateComm",
		ParamsMap{"ssn_addr": ssn.Addr, "new_rate": strconv.Itoa(new_rate)},
	})
	AssertEqual(Field(p.Zimpl, "ssnlist", addr, "arguments", "7"), strconv.Itoa(new_rate))

	//SSN comission will be available on next cycle
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	txIdLocal++
	//expect comission withdrawn successfully
	expectedCommission := "8230452599999"
	multisig.SetSigner(celestials.Owner)
	AssertMultisigSuccess(multisig.SubmitWithdrawCommTransaction(ssn.Addr))
	tx, _ = AssertMultisigSuccess(multisig.ExecuteTransaction(txIdLocal))
	AssertEvent(tx, Event{
		p.Zimpl.Addr, //sender
		"SSN withdraw reward",
		ParamsMap{"ssn_addr": ssn.Addr},
	})
	AssertTransition(tx, Transition{
		p.Zimpl.Addr, //sender
		"AddFunds",
		treasury_addr,
		expectedCommission,
		ParamsMap{},
	})

}
