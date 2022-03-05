package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

const txId = 0 // the test transition should be the first

func (tr *Transitions) MultisigWalletTests() {
	// multisigGoldenFlowTest(tr)
	// multisigChangeAdminTest(tr)
	// multisigUpdateOwner(tr)
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

	azil, _ := NewAZilContract(sdk, multisig.Addr, p.Zimpl.Addr)

	newSsnAddr := sdk.Cfg.Admin // could be any random address

	// after submitting transaction it automatically signed by the _sender
	AssertMultisigSuccess(multisig.WithUser(owner1).SubmitChangeAzilSSNAddressTransaction(azil.Addr, newSsnAddr))

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

	// should be changed after execution
	AssertMultisigSuccess(multisig.WithUser(owner1).SignTransaction(txId))
	AssertEqual(azil.GetAzilSsnAddress(), sdk.Cfg.AzilSsnAddress)
	AssertMultisigSuccess(multisig.WithUser(owner2).ExecuteTransaction(txId))
	AssertEqual(azil.GetAzilSsnAddress(), newSsnAddr)
}

func multisigChangeAdminTest(tr *Transitions) {
	owner1 := sdk.Cfg.Key1

	owners := []string{sdk.Cfg.Addr1}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	p := tr.DeployAndUpgrade()

	azil, _ := NewAZilContract(sdk, multisig.Addr, p.Zimpl.Addr)

	newAdmin := sdk.Cfg.Addr1

	// after submitting transaction it automatically signed by the _sender
	AssertMultisigSuccess(multisig.WithUser(owner1).SubmitChangeAdminTransaction(azil.Addr, newAdmin))
	AssertMultisigSuccess(multisig.WithUser(owner1).ExecuteTransaction(txId))
	AssertEqual(Field(azil, "admin_address"), newAdmin)
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

	azil, _ := NewAZilContract(sdk, multisig.Addr, p.Zimpl.Addr)

	// after submitting transaction it automatically signed by the _sender
	AssertMultisigSuccess(multisig.WithUser(owner1).SubmitChangeOwnerTransaction(azil.Addr, newOwner))
	AssertMultisigSuccess(multisig.WithUser(owner1).ExecuteTransaction(txId))
	AssertEqual(Field(azil, "staging_owner_address"), newOwner)

	// claim owner using; newMultisig has also the first tx in order
	AssertMultisigSuccess(newMultisig.WithUser(owner2).SubmitClaimOwnerTransaction(azil.Addr))
	AssertMultisigSuccess(newMultisig.WithUser(owner2).ExecuteTransaction(txId))
	AssertEqual(Field(azil, "owner_address"), newOwner)
}

func multisigManagableActions(tr *Transitions) {
	owner := sdk.Cfg.Key1

	owners := []string{sdk.Cfg.Addr1}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	p := tr.DeployAndUpgrade()

	azil, _ := NewAZilContract(sdk, multisig.Addr, p.Zimpl.Addr)

	newAddr := sdk.Cfg.Admin // could be any random address

	AssertMultisigSuccess(multisig.WithUser(owner).SubmitChangeTreasuryAddressTransaction(azil.Addr, newAddr))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitChangeZimplAddressTransaction(azil.Addr, newAddr))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitChangeRewardsFeeTransaction(azil.Addr, "1"))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitUpdateStakingParametersTransaction(azil.Addr, "1"))

	// pause actions
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitPauseInTransaction(azil.Addr))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitPauseOutTransaction(azil.Addr))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitPauseZrc2Transaction(azil.Addr))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitUnPauseInTransaction(azil.Addr))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitUnPauseOutTransaction(azil.Addr))
	AssertMultisigSuccess(multisig.WithUser(owner).SubmitUnPauseZrc2Transaction(azil.Addr))
}
