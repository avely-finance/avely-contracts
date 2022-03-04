package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) MultisigWalletTests() {
	Start("MultisigWalletTests contract transitions")

	admin := sdk.Cfg.AdminKey
	owner1 := sdk.Cfg.Key1
	owner2 := sdk.Cfg.Key2

	owners := []string{sdk.Cfg.Addr1, sdk.Cfg.Addr2}
	signCount := 2
	multisig := tr.DeployMultisigWallet(owners, signCount)

	p := tr.DeployAndUpgrade()
	newSsnAddr := sdk.Cfg.Admin // any random address

	// after submitting transaction it automatically signed by the _sender
	AssertMultisigSuccess(multisig.WithUser(owner1).SubmitChangeAzilSSNAddressTransaction(p.Azil.Addr, newSsnAddr))

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
	AssertEqual(p.Azil.GetAzilSsnAddress(), sdk.Cfg.AzilSsnAddress)
	AssertMultisigSuccess(multisig.WithUser(owner2).ExecuteTransaction(txId))
	AssertEqual(p.Azil.GetAzilSsnAddress(), newSsnAddr)
}

// RevokeSignature
