package contracts

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"

	. "github.com/avely-finance/avely-contracts/sdk/core"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type MultisigWallet struct {
	Contract
}

func (a *MultisigWallet) WithUser(key string) *MultisigWallet {
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)
	a.Contract.Wallet = wallet
	return a
}

func (s *MultisigWallet) SignTransaction(transactionId int) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "transactionId",
			Type:  "Uint32",
			Value: strconv.Itoa(transactionId),
		},
	}

	return s.Call("SignTransaction", args, "0")
}

func (s *MultisigWallet) RevokeSignature(transactionId int) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "transactionId",
			Type:  "Uint32",
			Value: strconv.Itoa(transactionId),
		},
	}

	return s.Call("RevokeSignature", args, "0")
}

func (s *MultisigWallet) ExecuteTransaction(transactionId int) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "transactionId",
			Type:  "Uint32",
			Value: strconv.Itoa(transactionId),
		},
	}

	return s.Call("ExecuteTransaction", args, "0")
}

func (s *MultisigWallet) SubmitChangeAdminTransaction(azilAddr, newAdmin string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
		{
			VName: "new_admin",
			Type:  "ByStr20",
			Value: newAdmin,
		},
	}

	return s.Call("SubmitChangeAdminTransaction", args, "0")
}

func (s *MultisigWallet) SubmitChangeOwnerTransaction(azilAddr, newOwner string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
		{
			VName: "new_owner",
			Type:  "ByStr20",
			Value: newOwner,
		},
	}

	return s.Call("SubmitChangeOwnerTransaction", args, "0")
}

func (s *MultisigWallet) SubmitClaimOwnerTransaction(azilAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
	}

	return s.Call("SubmitClaimOwnerTransaction", args, "0")
}

func (s *MultisigWallet) SubmitChangeTreasuryAddressTransaction(azilAddr string, addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
		{
			VName: "address",
			Type:  "ByStr20",
			Value: addr,
		},
	}

	return s.Call("SubmitChangeTreasuryAddressTransaction", args, "0")
}

func (s *MultisigWallet) SubmitChangeZimplAddressTransaction(azilAddr string, addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
		{
			VName: "address",
			Type:  "ByStr20",
			Value: addr,
		},
	}

	return s.Call("SubmitChangeZimplAddressTransaction", args, "0")
}

func (s *MultisigWallet) SubmitChangeRewardsFeeTransaction(azilAddr string, newFee string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
		{
			VName: "new_fee",
			Type:  "Uint128",
			Value: newFee,
		},
	}

	return s.Call("SubmitChangeRewardsFeeTransaction", args, "0")
}

func (s *MultisigWallet) SubmitUpdateStakingParametersTransaction(azilAddr string, minDelegStake string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
		{
			VName: "min_deleg_stake",
			Type:  "Uint128",
			Value: minDelegStake,
		},
	}

	return s.Call("SubmitUpdateStakingParametersTransaction", args, "0")
}

func (s *MultisigWallet) SubmitSetHolderAddressTransaction(azilAddr string, addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
		{
			VName: "address",
			Type:  "ByStr20",
			Value: addr,
		},
	}

	return s.Call("SubmitSetHolderAddressTransaction", args, "0")
}

func (s *MultisigWallet) SubmitChangeBuffersTransaction(azilAddr string, newBuffers []string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
		{
			VName: "new_buffers",
			Type:  "List ByStr20",
			Value: newBuffers,
		},
	}

	return s.Call("SubmitChangeBuffersTransaction", args, "0")
}

func (s *MultisigWallet) SubmitAddSSNTransaction(azilAddr, ssnaddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
		{
			VName: "ssnaddr",
			Type:  "ByStr20",
			Value: ssnaddr,
		},
	}

	return s.Call("SubmitAddSSNTransaction", args, "0")
}

func (s *MultisigWallet) SubmitPauseInTransaction(azilAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
	}

	return s.Call("SubmitPauseInTransaction", args, "0")
}

func (s *MultisigWallet) SubmitPauseOutTransaction(azilAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
	}

	return s.Call("SubmitPauseOutTransaction", args, "0")
}

func (s *MultisigWallet) SubmitPauseZrc2Transaction(azilAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
	}

	return s.Call("SubmitPauseZrc2Transaction", args, "0")
}

func (s *MultisigWallet) SubmitUnPauseInTransaction(azilAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
	}

	return s.Call("SubmitUnPauseInTransaction", args, "0")
}

func (s *MultisigWallet) SubmitUnPauseOutTransaction(azilAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
	}

	return s.Call("SubmitUnPauseOutTransaction", args, "0")
}

func (s *MultisigWallet) SubmitUnPauseZrc2Transaction(azilAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: azilAddr,
		},
	}

	return s.Call("SubmitUnPauseZrc2Transaction", args, "0")
}

func NewMultisigContract(sdk *AvelySDK, owners []string, requiredSignaturesCount int) (*MultisigWallet, error) {
	// TOOD: add requiredSignaturesCount validation
	contract := buildMultisigContract(sdk, owners, strconv.Itoa(requiredSignaturesCount))

	tx, err := sdk.DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, sdk.Cfg.TxConfrimMaxAttempts, sdk.Cfg.TxConfirmIntervalSec, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)

		sdkContract := Contract{
			Sdk:      sdk,
			Provider: *contract.Provider,
			Addr:     "0x" + tx.ContractAddress,
			Bech32:   b32,
			Wallet:   contract.Signer,
		}
		return &MultisigWallet{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreMultisigContract(sdk *AvelySDK, contractAddress string, owners []string, requiredSignaturesCount int) (*MultisigWallet, error) {
	contract := buildMultisigContract(sdk, owners, strconv.Itoa(requiredSignaturesCount))

	b32, err := bech32.ToBech32Address(contractAddress)

	if err != nil {
		return nil, errors.New("Config has invalid MultisigWallet address")
	}

	sdkContract := Contract{
		Sdk:      sdk,
		Provider: *contract.Provider,
		Addr:     contractAddress,
		Bech32:   b32,
		Wallet:   contract.Signer,
	}
	return &MultisigWallet{Contract: sdkContract}, nil
}

func buildMultisigContract(sdk *AvelySDK, owners []string, requiredSignaturesCount string) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/multisig_wallet.scilla")
	key := sdk.Cfg.AdminKey

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "owners_list",
			Type:  "List ByStr20",
			Value: owners,
		}, {
			VName: "required_signatures",
			Type:  "Uint32",
			Value: requiredSignaturesCount,
		},
	}

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)

	return contract2.Contract{
		Provider: sdk.InitProvider(),
		Code:     string(code),
		Init:     init,
		Signer:   wallet,
	}
}
