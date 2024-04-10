package evm

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/avely-finance/avely-contracts/sdk/contracts/evm/bind"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type StZIL struct {
	Contract
	*bind.StZIL
}

func NewStZILContract(sdk *core.AvelySDK, stzil_address string, deployer *accounts.Account) (*StZIL, error) {

	transactOpts, err := sdk.Evm.GetTransactOpts(deployer)
	if err != nil {
		panic(err)
	}

	address /*tx*/, _, stzilEvm, err := bind.DeployStZIL(transactOpts, sdk.Evm.Client, common.HexToAddress(stzil_address))
	if err != nil {
		panic(err)
	}

	//fmt.Printf("Api contract deployed to %s\n", address.Hex())
	//fmt.Printf("Tx: %s\n", tx.Hash().Hex())

	sdkContract := Contract{
		Sdk:     sdk,
		Addr:    address.String(),
		Account: deployer,
	}

	return &StZIL{sdkContract, stzilEvm}, nil
}

func (s *StZIL) ChownStakeConfirmSwap(delegator string) (*types.Transaction, error) {
	txEvm, err := s.StZIL.ChownStakeConfirmSwap(s.Sdk.Evm.GetTransactOptsOrPanic(s.Account),
		common.HexToAddress(delegator))

	return txEvm, err
}

func (s *StZIL) Transfer(to, amount string) (*types.Transaction, error) {
	amtBi, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf("cant't create BigInt from %s string", amount)
	}
	txEvm, err := s.StZIL.Transfer(s.Sdk.Evm.GetTransactOptsOrPanic(s.Account),
		common.HexToAddress(to), amtBi)

	return txEvm, err
}

func (s *StZIL) TransferFrom(from, to, amount string) (*types.Transaction, error) {
	amtBi, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf("cant't create BigInt from %s string", amount)
	}
	txEvm, err := s.StZIL.TransferFrom(s.Sdk.Evm.GetTransactOptsOrPanic(s.Account), common.HexToAddress(from),
		common.HexToAddress(to), amtBi)

	return txEvm, err
}

func (s *StZIL) IncreaseAllowance(to, amount string) (*types.Transaction, error) {
	amtBi, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf("cant't create BigInt from %s string", amount)
	}
	txEvm, err := s.StZIL.IncreaseAllowance(s.Sdk.Evm.GetTransactOptsOrPanic(s.Account),
		common.HexToAddress(to), amtBi)

	return txEvm, err
}

func (s *StZIL) DecreaseAllowance(to, amount string) (*types.Transaction, error) {
	amtBi, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf("cant't create BigInt from %s string", amount)
	}
	txEvm, err := s.StZIL.DecreaseAllowance(s.Sdk.Evm.GetTransactOptsOrPanic(s.Account),
		common.HexToAddress(to), amtBi)

	return txEvm, err
}

func (s *StZIL) DelegateStake(amount string) (*types.Transaction, error) {
	// we can't implement this for now because of bug/feature in zilliqa's evm implementation
	// it's impossible to transfer native funds from an evm contract to a scilla contract
	return nil, errors.New("evm/stzil->DelegateStake() not implemented")
}

func (s *StZIL) WithdrawTokensAmt(amount string) (*types.Transaction, error) {
	amtBi, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf("cant't create BigInt from %s string", amount)
	}
	txEvm, err := s.StZIL.WithdrawTokensAmt(s.Sdk.Evm.GetTransactOptsOrPanic(s.Account), amtBi)

	return txEvm, err
}

func (s *StZIL) CompleteWithdrawal() (*types.Transaction, error) {
	txEvm, err := s.StZIL.CompleteWithdrawal(s.Sdk.Evm.GetTransactOptsOrPanic(s.Account))
	return txEvm, err
}
