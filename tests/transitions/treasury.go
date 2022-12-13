package transitions

import (
	"math/big"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
	"github.com/tyler-smith/go-bip39"
)

func (tr *Transitions) Treasury() {
	treasuryFunds(tr)
	treasuryChangeOwner(tr)
	treasuryRequireOwner(tr)
}

func treasuryChangeOwner(tr *Transitions) {
	txIdLocal1 := 0
	txIdLocal2 := 0

	//deploy multisig
	owner := alice
	owners := []string{utils.GetAddressByWallet(alice)}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	//deploy treasury, set owner to multisig contract
	init_owner := multisig.Addr
	treasury := tr.DeployTreasury(init_owner)

	//deploy other multisig contract
	newSignCount := 1
	newOwner := bob
	newOwners := []string{utils.GetAddressByWallet(newOwner)}
	newMultisig := tr.DeployMultisigWallet(newOwners, newSignCount)

	//try to claim owner, expect error
	newMultisig.SetSigner(bob)
	AssertMultisigSuccess(newMultisig.SubmitClaimOwnerTransaction(treasury.Addr))
	tx, _ := newMultisig.ExecuteTransaction(txIdLocal2)
	AssertError(tx, treasury.ErrorCode("StagingOwnerNotExists"))

	//initiate owner change
	multisig.SetSigner(owner)
	AssertMultisigSuccess(multisig.SubmitChangeOwnerTransaction(treasury.Addr, newMultisig.Addr))
	AssertMultisigSuccess(multisig.ExecuteTransaction(txIdLocal1))
	AssertEqual(Field(treasury, "staging_owner"), newMultisig.Addr)

	//try to claim owner with wrong user, expect error
	treasury.SetSigner(bob)
	tx, _ = treasury.ClaimOwner()
	AssertError(tx, treasury.ErrorCode("StagingOwnerValidationFailed"))

	//claim owner
	txIdLocal2++
	newMultisig.SetSigner(bob)
	AssertMultisigSuccess(newMultisig.SubmitClaimOwnerTransaction(treasury.Addr))
	AssertMultisigSuccess(newMultisig.ExecuteTransaction(txIdLocal2))
	AssertEqual(Field(treasury, "owner"), newMultisig.Addr)

}

func treasuryFunds(tr *Transitions) {
	//deploy multisig
	owner := alice
	owners := []string{utils.GetAddressByWallet(alice)}
	signCount := 1
	multisig := tr.DeployMultisigWallet(owners, signCount)

	//deploy treasury, set owner to multisig contract
	init_owner := multisig.Addr
	treasury := tr.DeployTreasury(init_owner)

	txIdLocal := 0

	//add funds
	treasury.SetSigner(celestials.Admin)
	treasury.AddFunds(ToQA(100))

	//try to withdraw amount exceeding _balance, expect error
	admin := utils.GetAddressByWallet(celestials.Admin)
	multisig.SetSigner(owner)
	AssertMultisigSuccess(multisig.SubmitWithdrawTransaction(treasury.Addr, admin, ToQA(12345)))
	tx, _ := multisig.ExecuteTransaction(txIdLocal)
	AssertError(tx, treasury.ErrorCode("InsufficientFunds"))

	//withdraw valid amount, expect success
	txIdLocal++

	//Generate a mnemonic for memorization or user-friendly seeds
	entropy, _ := bip39.NewEntropy(128) //256
	mnemonic, _ := bip39.NewMnemonic(entropy)

	//mnemonic := "bug feature framework lava jelly keep device journey bean mango rocket festival"
	account1, _ := account.NewDefaultHDAccount(mnemonic, uint32(1))
	RcptAddr1 := "0x" + account1.Address
	//RcptKey1 := util.EncodeHex(account1.PrivateKey)

	//add some funds to newly created account
	sdk.AddFunds(celestials.Admin, RcptAddr1, ToQA(1000))

	recipient := RcptAddr1
	balanceBefore, _ := new(big.Int).SetString(sdk.GetBalance(recipient), 10)
	GetLog().Info(balanceBefore.String())

	multisig.SetSigner(owner)
	AssertMultisigSuccess(multisig.SubmitWithdrawTransaction(treasury.Addr, recipient, ToQA(25)))
	tx, _ = AssertMultisigSuccess(multisig.ExecuteTransaction(txIdLocal))
	AssertTransition(tx, Transition{
		treasury.Addr, //sender
		"AddFunds",
		recipient,
		ToQA(25),
		ParamsMap{},
	})

	//data, _ := json.MarshalIndent(tx, "", "     ")
	//GetLog().Info(string(data))
	AssertEqual(Field(treasury, "_balance"), ToQA(75))
	withdraw, _ := new(big.Int).SetString(ToQA(25), 10)
	balanceAfter := big.NewInt(0).Add(balanceBefore, withdraw)
	AssertEqual(sdk.GetBalance(recipient), balanceAfter.String())
}

func treasuryRequireOwner(tr *Transitions) {

	Start("treasuryRequireOwner")

	p := tr.DeployAndUpgrade()

	// Use non-owner user, expecting errors
	p.Treasury.SetSigner(bob)
	randomAddr := utils.GetAddressByWallet(eve)

	tx, _ := p.Treasury.ChangeOwner(randomAddr)
	AssertError(tx, p.Treasury.ErrorCode("OwnerValidationFailed"))

	tx, _ = p.Treasury.Withdraw(core.ZeroAddr, "123")
	AssertError(tx, p.Treasury.ErrorCode("OwnerValidationFailed"))
}
