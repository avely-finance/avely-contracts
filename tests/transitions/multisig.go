package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
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

	owners := []string{utils.GetAddressByWallet(alice), utils.GetAddressByWallet(bob)}
	signCount := 2
	multisig := tr.DeployMultisigWallet(owners, signCount)

	p := tr.DeployAndUpgrade()

	stzil, _ := NewStZILContract(sdk, multisig.Addr, p.Zimpl.Addr, celestials.Admin)

	newAddr := utils.GetAddressByWallet(celestials.Admin) // could be any random address

	// after submitting transaction it automatically signed by the _sender
	multisig.SetSigner(alice)
	AssertMultisigSuccess(multisig.SubmitChangeTreasuryAddressTransaction(stzil.Addr, newAddr))

	txId := 0 // the test transition should be the first

	multisig.SetSigner(celestials.Admin)
	tx, _ := multisig.SignTransaction(txId)
	AssertMultisigError(tx, multisig.ErrorCode("NonOwnerCannotSign"))

	multisig.SetSigner(bob)
	tx, _ = multisig.SignTransaction(txId + 1)
	AssertMultisigError(tx, multisig.ErrorCode("UnknownTransactionId"))

	AssertMultisigSuccess(multisig.SignTransaction(txId))

	// revoke and sign again
	multisig.SetSigner(alice)
	AssertMultisigSuccess(multisig.RevokeSignature(txId))

	multisig.SetSigner(bob)
	tx, _ = multisig.ExecuteTransaction(txId)
	AssertMultisigError(tx, multisig.ErrorCode("NotEnoughSignatures"))

	// treasury address should be changed after execution
	multisig.SetSigner(alice)
	AssertMultisigSuccess(multisig.SignTransaction(txId))
	//field treasury_address      : ByStr20 = init_admin_address

	AssertEqual(stzil.GetTreasuryAddress(), newAddr)

	multisig.SetSigner(bob)
	AssertMultisigSuccess(multisig.ExecuteTransaction(txId))
	AssertEqual(stzil.GetTreasuryAddress(), newAddr)
}

func multisigChangeAdminTest(tr *Transitions) {
	owners := []string{utils.GetAddressByWallet(alice)}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	p := tr.DeployAndUpgrade()

	stzil, _ := NewStZILContract(sdk, multisig.Addr, p.Zimpl.Addr, celestials.Admin)

	newAdmin := utils.GetAddressByWallet(bob)

	// after submitting transaction it automatically signed by the _sender
	multisig.SetSigner(alice)
	AssertMultisigSuccess(multisig.SubmitChangeAdminTransaction(stzil.Addr, newAdmin))
	AssertMultisigSuccess(multisig.ExecuteTransaction(txId))
	AssertEqual(Field(stzil, "admin_address"), newAdmin)
}

func multisigUpdateOwner(tr *Transitions) {
	signCount := 1

	owners := []string{utils.GetAddressByWallet(alice)}
	multisig := tr.DeployMultisigWallet(owners, signCount)

	// deploy new multisig with new owners
	owner2 := sdk.Cfg.Key2
	newOwners := []string{sdk.Cfg.Addr2}
	newMultisig := tr.DeployMultisigWallet(newOwners, signCount)

	newOwner := newMultisig.Addr
	p := tr.DeployAndUpgrade()

	stzil, _ := NewStZILContract(sdk, multisig.Addr, p.Zimpl.Addr, celestials.Admin)

	// after submitting transaction it automatically signed by the _sender
	multisig.SetSigner(alice)
	AssertMultisigSuccess(multisig.SubmitChangeOwnerTransaction(stzil.Addr, newOwner))
	AssertMultisigSuccess(multisig.ExecuteTransaction(txId))
	AssertEqual(Field(stzil, "staging_owner_address"), newOwner)

	// claim owner using; newMultisig has also the first tx in order
	AssertMultisigSuccess(newMultisig.WithUser(owner2).SubmitClaimOwnerTransaction(stzil.Addr))
	AssertMultisigSuccess(newMultisig.WithUser(owner2).ExecuteTransaction(txId))
	AssertEqual(Field(stzil, "owner_address"), newOwner)
}

func multisigManagableActions(tr *Transitions) {
	owners := []string{utils.GetAddressByWallet(alice)}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	p := tr.DeployAndUpgrade()

	stzil, _ := NewStZILContract(sdk, multisig.Addr, p.Zimpl.Addr, celestials.Admin)

	newAddr := utils.GetAddressByWallet(celestials.Admin) // could be any random address

	multisig.SetSigner(alice)
	AssertMultisigSuccess(multisig.SubmitChangeTreasuryAddressTransaction(stzil.Addr, newAddr))
	AssertMultisigSuccess(multisig.SubmitChangeZimplAddressTransaction(stzil.Addr, newAddr))
	AssertMultisigSuccess(multisig.SubmitChangeRewardsFeeTransaction(stzil.Addr, "1"))
	AssertMultisigSuccess(multisig.SubmitUpdateStakingParametersTransaction(stzil.Addr, "1"))

	// pause actions
	AssertMultisigSuccess(multisig.SubmitPauseInTransaction(stzil.Addr))
	AssertMultisigSuccess(multisig.SubmitPauseOutTransaction(stzil.Addr))
	AssertMultisigSuccess(multisig.SubmitPauseZrc2Transaction(stzil.Addr))
	AssertMultisigSuccess(multisig.SubmitUnPauseInTransaction(stzil.Addr))
	AssertMultisigSuccess(multisig.SubmitUnPauseOutTransaction(stzil.Addr))
	AssertMultisigSuccess(multisig.SubmitUnPauseZrc2Transaction(stzil.Addr))

	AssertMultisigSuccess(multisig.SubmitSetHolderAddressTransaction(stzil.Addr, newAddr))
}

func multisigChangeBuffersTest(tr *Transitions) {
	owners := []string{utils.GetAddressByWallet(alice)}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	p := tr.DeployAndUpgrade()

	stzil, _ := NewStZILContract(sdk, multisig.Addr, p.Zimpl.Addr, celestials.Admin)

	newBuffers := []string{"0xf61477D7919478e5AfFe1fbd9A0CDCeee9fdE42d"} // could be any random addresses

	// after submitting transaction it automatically signed by the _sender
	multisig.SetSigner(alice)
	AssertMultisigSuccess(multisig.SubmitChangeBuffersTransaction(stzil.Addr, newBuffers))
	AssertMultisigSuccess(multisig.ExecuteTransaction(txId))
	AssertContain(Field(stzil, "buffers_addresses"), "0xf61477d7919478e5affe1fbd9a0cdceee9fde42d")
}

func multisigAddRemoveSSNTest(tr *Transitions) {
	owners := []string{utils.GetAddressByWallet(alice)}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	p := tr.DeployAndUpgrade()

	stzil, _ := NewStZILContract(sdk, multisig.Addr, p.Zimpl.Addr, celestials.Admin)

	newAddress := "0xf61477D7919478e5AfFe1fbd9A0CDCeee9fdE42d" // could be any random addresses

	multisig.SetSigner(alice)
	// after submitting transaction it automatically signed by the _sender
	AssertMultisigSuccess(multisig.SubmitAddSSNTransaction(stzil.Addr, newAddress))
	AssertMultisigSuccess(multisig.ExecuteTransaction(txId))
	AssertContain(Field(stzil, "ssn_addresses"), "0xf61477d7919478e5affe1fbd9a0cdceee9fde42d")

	//try to remove non-existent ssn, expect SsnAddressDoesNotExist error
	AssertMultisigSuccess(multisig.SubmitRemoveSSNTransaction(stzil.Addr, core.ZeroAddr))
	tx, _ := multisig.ExecuteTransaction(txId + 1)
	AssertError(tx, p.StZIL.ErrorCode("SsnAddressDoesNotExist"))

	//remove ssn, added before; expect success
	AssertMultisigSuccess(multisig.SubmitRemoveSSNTransaction(stzil.Addr, newAddress))
	AssertMultisigSuccess(multisig.ExecuteTransaction(txId + 2))
	AssertEqual(Field(stzil, "ssn_addresses"), "[]")
}
