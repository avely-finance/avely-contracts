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

	address := sdk.Cfg.Admin // as random address

	// after submitting transaction it automatically signed by the _sender
	AssertMultisigSuccess(multisig.WithUser(owner1).SubmitChangeAzilSSNAddressTransaction(address, address))

	txId := 0 // the test transition should be the first

	tx, _ := multisig.WithUser(admin).SignTransaction(txId)
	AssertMultisigError(tx, "-1") // NonOwnerCannotSign

	tx, _ = multisig.WithUser(owner2).SignTransaction(txId + 1)
	AssertMultisigError(tx, "-2") // UnknownTransactionId

	AssertMultisigSuccess(multisig.WithUser(owner2).SignTransaction(txId))
	AssertMultisigSuccess(multisig.WithUser(owner2).SubmitChangeAzilSSNAddressTransaction(address, address))

	// GetLog().Info(tx)
}
