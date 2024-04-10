package transitions

import (
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/Zilliqa/gozilliqa-sdk/v3/account"
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/contracts/evm"
	"github.com/avely-finance/avely-contracts/tests/helpers"
	"github.com/ethereum/go-ethereum/accounts"
)

type StZILContract interface {
	SetEvm(state bool)
	SetSigner(signer interface{}) StZILContract
	IncreaseAllowance(to, amount string) (interface{}, error)
	DecreaseAllowance(to, amount string) (interface{}, error)
	Transfer(to, amount string) (interface{}, error)
	TransferFrom(from, to, amount string) (interface{}, error)
	DelegateStake(amount string) (interface{}, error)
	ChownStakeConfirmSwap(delegator string) (interface{}, error)
	WithdrawTokensAmt(amount string) (interface{}, error)
	CompleteWithdrawal() (interface{}, error)
}

type StZILAdapter struct {
	scillaContract *contracts.StZIL
	evmContract    *evm.StZIL
	evmOn          bool
}

func NewStZILAdapter(scilla *contracts.StZIL, evm *evm.StZIL, evmOn bool) *StZILAdapter {
	return &StZILAdapter{
		scillaContract: scilla,
		evmContract:    evm,
		evmOn:          evmOn,
	}
}

func (a *StZILAdapter) SetEvm(state bool) {
	a.evmOn = state
}

func (a *StZILAdapter) SetSigner(signer interface{}) StZILContract {
	if a.evmOn {
		if acc, ok := signer.(*accounts.Account); ok {
			a.evmContract.SetSigner(acc)
		} else if wallet, ok := signer.(*account.Wallet); ok {
			acc, _ := sdk.Evm.AddAccountByPrivateKey(util.EncodeHex(wallet.DefaultAccount.PrivateKey))
			a.evmContract.SetSigner(acc)
		}
	} else {
		if wallet, ok := signer.(*account.Wallet); ok {
			a.scillaContract.SetSigner(wallet)
		}
	}
	return a
}

func (a *StZILAdapter) ChownStakeConfirmSwap(delegator string) (interface{}, error) {
	if a.evmOn {
		return a.evmContract.ChownStakeConfirmSwap(delegator)
	} else {
		return a.scillaContract.ChownStakeConfirmSwap(delegator)
	}
}

func (a *StZILAdapter) IncreaseAllowance(to, amount string) (interface{}, error) {
	if a.evmOn {
		return a.evmContract.IncreaseAllowance(to, amount)
	} else {
		return a.scillaContract.IncreaseAllowance(to, amount)
	}
}

func (a *StZILAdapter) DecreaseAllowance(to, amount string) (interface{}, error) {
	if a.evmOn {
		return a.evmContract.DecreaseAllowance(to, amount)
	} else {
		return a.scillaContract.DecreaseAllowance(to, amount)
	}
}

func (a *StZILAdapter) Transfer(to, amount string) (interface{}, error) {
	if a.evmOn {
		return a.evmContract.Transfer(to, amount)
	} else {
		return a.scillaContract.Transfer(to, amount)
	}
}

func (a *StZILAdapter) TransferFrom(from, to, amount string) (interface{}, error) {
	if a.evmOn {
		return a.evmContract.TransferFrom(from, to, amount)
	} else {
		return a.scillaContract.TransferFrom(from, to, amount)
	}
}

func (a *StZILAdapter) DelegateStake(amount string) (interface{}, error) {
	if a.evmOn {
		if helpers.Field(a.scillaContract, "is_paused_in") == "True" {
			//In this case we use Stzil.scilla->DelegateStake to get PausedIn error.
			//When ZQ2 will be delivered, this should be replaced with EVM transaction
			return a.scillaContract.DelegateStake(amount)
		}
		// we can't delegate via evm for now because of bug/feature in zilliqa's evm implementation
		// it's impossible to transfer native funds from an evm contract to a scilla contract

		// take current evm signer, get it's address
		// set stzil-prefilled account signer, verifier in our case
		// transfer `amount` of stzil to the address
		// restore previous signer
		// return transfer tx

		curScillaSigner := a.scillaContract.Wallet
		a.scillaContract.SetSigner(celestials.Verifier)
		tx, err := a.scillaContract.Transfer(a.evmContract.Account.Address.Hex(), amount)
		a.scillaContract.SetSigner(curScillaSigner)

		// type of return tx is Scilla transaction
		// if this will brake something, we could transfer via Evm-bridge
		return tx, err
	} else {
		return a.scillaContract.DelegateStake(amount)
	}
}

func (a *StZILAdapter) WithdrawTokensAmt(amount string) (interface{}, error) {
	if a.evmOn {
		return a.evmContract.WithdrawTokensAmt(amount)
	} else {
		return a.scillaContract.WithdrawTokensAmt(amount)
	}
}

func (a *StZILAdapter) CompleteWithdrawal() (interface{}, error) {
	if a.evmOn {
		return a.evmContract.CompleteWithdrawal()
	} else {
		return a.scillaContract.CompleteWithdrawal()
	}
}
