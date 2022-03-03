package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) MultisigWalletTests() {
	Start("MultisigWalletTests contract transitions")

	owners := []string{sdk.Cfg.Addr1, sdk.Cfg.Addr2}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	address := sdk.Cfg.Admin // as random address

	AssertSuccess(multisig.SubmitChangeAzilSSNAddressTransaction(address, address))
}
