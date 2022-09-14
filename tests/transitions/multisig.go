package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

const txId = 0 // the test transition should be the first

func (tr *Transitions) MultisigWalletTests() {
	multisigGoldenFlowTest(tr)
	multisigUpdateOwner(tr)
	multisigChangeAdminTest(tr)
	multisigChangeBuffersTest(tr)
	multisigAddRemoveSSNTest(tr)
	multisigManagableActions(tr)
}

func multisigGoldenFlowTest(tr *Transitions) {
	Start("MultisigWalletTests contract transitions")

	admin := sdk.Cfg.AdminKey
	owner1 := sdk.Cfg.Key1
	owner2 := sdk.Cfg.Key2

	owners := []string{sdk.Cfg.Addr1, sdk.Cfg.Addr2}
	signCount := 2
	multisig := tr.DeployMultisigWallet(owners, signCount)

	p := tr.DeployAndUpgrade()

	stzil, _ := NewStZILContract(sdk, multisig.Addr, p.Zimpl.Addr)

	newAddr := sdk.Cfg.Admin // could be any random address

	// after submitting transaction it automatically signed by the _sender
	AssertMultisigSuccess(multisig.WithUser(owner1).SubmitChangeTreasuryAddressTransaction(stzil.Addr, newAddr))

	txId := 0 // the test transition should be the first

	tx, _ := multisig.WithUser(admin).SignTransaction(txId)
	AssertMultisigError(tx, "-1") // NonOwnerCannotSign

	tx, _ = multisig.WithUser(owner2).SignTransaction(txId + 1)
	AssertMultisigError(tx, "-2") // UnknownTransactionId

	AssertMultisigSuccess(multisig.WithUser(owner2).SignTransaction(txId))

	// revoke and sign again
	AssertMultisigSuccess(multisig.WithUser(owner1).RevokeSignature(txId))
	tx, _ = multisig.WithUser(owner2).ExecuteTransaction(txId)
	AssertMultisigError(tx, "-9") // NotEnoughSignatures

	// treasury address should be changed after execution
	AssertMultisigSuccess(multisig.WithUser(owner1).SignTransaction(txId))
	//field treasury_address      : ByStr20 = init_admin_address
	AssertEqual(stzil.GetTreasuryAddress(), sdk.Cfg.Admin)
	AssertMultisigSuccess(multisig.WithUser(owner2).ExecuteTransaction(txId))
	AssertEqual(stzil.GetTreasuryAddress(), newAddr)
}

func multisigChangeAdminTest(tr *Transitions) {
	owner := sdk.Cfg.Key1

	owners := []string{sdk.Cfg.Addr1}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	p := tr.DeployAndUpgrade()

	stzil, _ := NewStZILContract(sdk, multisig.Addr, p.Zimpl.Addr)

	newAdmin := sdk.Cfg.Addr1

	// after submitting transaction it automatically signed by the _sender
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitChangeAdminTransaction(stzil.Addr, newAdmin))
	AssertMultisigSuccess(multisig.WithUser(owner).ExecuteTransaction(txId))
	AssertEqual(Field(stzil, "admin_address"), newAdmin)
}

func multisigUpdateOwner(tr *Transitions) {
	signCount := 1

	owner1 := sdk.Cfg.Key1
	owners := []string{sdk.Cfg.Addr1}
	multisig := tr.DeployMultisigWallet(owners, signCount)

	// deploy new multisig with new owners
	owner2 := sdk.Cfg.Key2
	newOwners := []string{sdk.Cfg.Addr2}
	newMultisig := tr.DeployMultisigWallet(newOwners, signCount)

	newOwner := newMultisig.Addr
	p := tr.DeployAndUpgrade()

	stzil, _ := NewStZILContract(sdk, multisig.Addr, p.Zimpl.Addr)

	// after submitting transaction it automatically signed by the _sender
	AssertMultisigSuccess(multisig.WithUser(owner1).SubmitChangeOwnerTransaction(stzil.Addr, newOwner))
	AssertMultisigSuccess(multisig.WithUser(owner1).ExecuteTransaction(txId))
	AssertEqual(Field(stzil, "staging_owner_address"), newOwner)

	// claim owner using; newMultisig has also the first tx in order
	AssertMultisigSuccess(newMultisig.WithUser(owner2).SubmitClaimOwnerTransaction(stzil.Addr))
	AssertMultisigSuccess(newMultisig.WithUser(owner2).ExecuteTransaction(txId))
	AssertEqual(Field(stzil, "owner_address"), newOwner)
}

func multisigManagableActions(tr *Transitions) {
	owner := sdk.Cfg.Key1

	owners := []string{sdk.Cfg.Addr1}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	p := tr.DeployAndUpgrade()

	stzil, _ := NewStZILContract(sdk, multisig.Addr, p.Zimpl.Addr)

	newAddr := sdk.Cfg.Admin // could be any random address

	AssertMultisigSuccess(multisig.WithUser(owner).SubmitChangeTreasuryAddressTransaction(stzil.Addr, newAddr))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitChangeZimplAddressTransaction(stzil.Addr, newAddr))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitChangeRewardsFeeTransaction(stzil.Addr, "1"))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitUpdateStakingParametersTransaction(stzil.Addr, "1"))

	// pause actions
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitPauseInTransaction(stzil.Addr))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitPauseOutTransaction(stzil.Addr))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitPauseZrc2Transaction(stzil.Addr))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitUnPauseInTransaction(stzil.Addr))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitUnPauseOutTransaction(stzil.Addr))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitUnPauseZrc2Transaction(stzil.Addr))

	AssertMultisigSuccess(multisig.WithUser(owner).SubmitSetHolderAddressTransaction(stzil.Addr, newAddr))
}

func multisigChangeBuffersTest(tr *Transitions) {
	owner := sdk.Cfg.Key1

	owners := []string{sdk.Cfg.Addr1}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	p := tr.DeployAndUpgrade()

	stzil, _ := NewStZILContract(sdk, multisig.Addr, p.Zimpl.Addr)

	newBuffers := []string{sdk.Cfg.Addr1} // could be any random addresses

	// after submitting transaction it automatically signed by the _sender
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitChangeBuffersTransaction(stzil.Addr, newBuffers))
	AssertMultisigSuccess(multisig.WithUser(owner).ExecuteTransaction(txId))
	AssertContain(Field(stzil, "buffers_addresses"), sdk.Cfg.Addr1)
}

func multisigAddRemoveSSNTest(tr *Transitions) {
	owner := sdk.Cfg.Key1

	owners := []string{sdk.Cfg.Addr1}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	p := tr.DeployAndUpgrade()

	stzil, _ := NewStZILContract(sdk, multisig.Addr, p.Zimpl.Addr)

	newAddress := sdk.Cfg.Addr1 // could be any random addresses

	// after submitting transaction it automatically signed by the _sender
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitAddSSNTransaction(stzil.Addr, newAddress))
	AssertMultisigSuccess(multisig.WithUser(owner).ExecuteTransaction(txId))
	AssertContain(Field(stzil, "ssn_addresses"), sdk.Cfg.Addr1)

	//try to remove non-existent ssn, expect SsnAddressDoesNotExist error
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitRemoveSSNTransaction(stzil.Addr, core.ZeroAddr))
	tx, _ := multisig.WithUser(owner).ExecuteTransaction(txId + 1)
	AssertError(tx, "SsnAddressDoesNotExist")

	//remove ssn, added before; expect success
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitRemoveSSNTransaction(stzil.Addr, newAddress))
	AssertMultisigSuccess(multisig.WithUser(owner).ExecuteTransaction(txId + 2))
	AssertEqual(Field(stzil, "ssn_addresses"), "[]")
}
