package contracts

import (
	"encoding/json"
	"errors"
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

func (s *MultisigWallet) SubmitChangeAdminTransaction(stZilAddr, newAdmin string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: stZilAddr,
		},
		{
			VName: "new_admin",
			Type:  "ByStr20",
			Value: newAdmin,
		},
	}

	return s.Call("SubmitChangeAdminTransaction", args, "0")
}

func (s *MultisigWallet) SubmitChangeTreasuryAddressTransaction(stZilAddr string, addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: stZilAddr,
		},
		{
			VName: "address",
			Type:  "ByStr20",
			Value: addr,
		},
	}

	return s.Call("SubmitChangeTreasuryAddressTransaction", args, "0")
}

func (s *MultisigWallet) SubmitChangeZimplAddressTransaction(stZilAddr string, addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: stZilAddr,
		},
		{
			VName: "address",
			Type:  "ByStr20",
			Value: addr,
		},
	}

	return s.Call("SubmitChangeZimplAddressTransaction", args, "0")
}

func (s *MultisigWallet) SubmitUpdateStakingParametersTransaction(stZilAddr, newMinDelegStake, newRewardsFee, newWithdrawalFee string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: stZilAddr,
		},
		{
			VName: "new_mindelegstake",
			Type:  "Uint128",
			Value: newMinDelegStake,
		},
		{
			VName: "new_rewards_fee",
			Type:  "Uint128",
			Value: newRewardsFee,
		},
		{
			VName: "new_withdrawal_fee",
			Type:  "Uint128",
			Value: newWithdrawalFee,
		},
	}

	return s.Call("SubmitUpdateStakingParametersTransaction", args, "0")
}

func (s *MultisigWallet) SubmitSetHolderAddressTransaction(stZilAddr string, addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: stZilAddr,
		},
		{
			VName: "address",
			Type:  "ByStr20",
			Value: addr,
		},
	}

	return s.Call("SubmitSetHolderAddressTransaction", args, "0")
}

func (s *MultisigWallet) SubmitChangeBuffersTransaction(stZilAddr string, newBuffers []string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: stZilAddr,
		},
		{
			VName: "new_buffers",
			Type:  "List ByStr20",
			Value: newBuffers,
		},
	}

	return s.Call("SubmitChangeBuffersTransaction", args, "0")
}

func (s *MultisigWallet) SubmitAddSSNTransaction(stZilAddr, ssnaddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: stZilAddr,
		},
		{
			VName: "ssnaddr",
			Type:  "ByStr20",
			Value: ssnaddr,
		},
	}

	return s.Call("SubmitAddSSNTransaction", args, "0")
}

func (s *MultisigWallet) SubmitRemoveSSNTransaction(stZilAddr, ssnaddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: stZilAddr,
		},
		{
			VName: "ssnaddr",
			Type:  "ByStr20",
			Value: ssnaddr,
		},
	}

	return s.Call("SubmitRemoveSSNTransaction", args, "0")
}

func (s *MultisigWallet) SubmitPauseInTransaction(stZilAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: stZilAddr,
		},
	}

	return s.Call("SubmitPauseInTransaction", args, "0")
}

func (s *MultisigWallet) SubmitPauseOutTransaction(stZilAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: stZilAddr,
		},
	}

	return s.Call("SubmitPauseOutTransaction", args, "0")
}

func (s *MultisigWallet) SubmitPauseZrc2Transaction(stZilAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: stZilAddr,
		},
	}

	return s.Call("SubmitPauseZrc2Transaction", args, "0")
}

func (s *MultisigWallet) SubmitUnPauseInTransaction(stZilAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: stZilAddr,
		},
	}

	return s.Call("SubmitUnPauseInTransaction", args, "0")
}

func (s *MultisigWallet) SubmitUnPauseOutTransaction(stZilAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: stZilAddr,
		},
	}

	return s.Call("SubmitUnPauseOutTransaction", args, "0")
}

func (s *MultisigWallet) SubmitUnPauseZrc2Transaction(stZilAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: stZilAddr,
		},
	}

	return s.Call("SubmitUnPauseZrc2Transaction", args, "0")
}

func (s *MultisigWallet) SubmitSetTreasuryFeeTransaction(aswapAddr string, newFee string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: aswapAddr,
		},
		{
			VName: "new_fee",
			Type:  "Uint128",
			Value: newFee,
		},
	}

	return s.Call("SubmitSetTreasuryFeeTransaction", args, "0")
}

func (s *MultisigWallet) SubmitSetTreasuryAddressTransaction(aswapAddr string, newAddress string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: aswapAddr,
		},
		{
			VName: "new_address",
			Type:  "ByStr20",
			Value: newAddress,
		},
	}

	return s.Call("SubmitSetTreasuryAddressTransaction", args, "0")
}

func (s *MultisigWallet) SubmitSetLiquidityFeeTransaction(aswapAddr string, newFee string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: aswapAddr,
		},
		{
			VName: "new_fee",
			Type:  "Uint256",
			Value: newFee,
		},
	}

	return s.Call("SubmitSetLiquidityFeeTransaction", args, "0")
}

func (s *MultisigWallet) SubmitTogglePauseTransaction(aswapAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: aswapAddr,
		},
	}

	return s.Call("SubmitTogglePauseTransaction", args, "0")
}

func (s *MultisigWallet) SubmitChangeOwnerTransaction(calleeContract, newOwner string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: calleeContract,
		},
		{
			VName: "new_owner",
			Type:  "ByStr20",
			Value: newOwner,
		},
	}

	return s.Call("SubmitChangeOwnerTransaction", args, "0")
}

func (s *MultisigWallet) SubmitClaimOwnerTransaction(calleeContract string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: calleeContract,
		},
	}

	return s.Call("SubmitClaimOwnerTransaction", args, "0")
}

func (s *MultisigWallet) SubmitWithdrawTransaction(treasuryAddr, recipient, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: treasuryAddr,
		},
		{
			VName: "recipient",
			Type:  "ByStr20",
			Value: recipient,
		},
		{
			VName: "amount",
			Type:  "Uint128",
			Value: amount,
		},
	}

	return s.Call("SubmitWithdrawTransaction", args, "0")
}

func (s *MultisigWallet) SubmitChangeZproxyTransaction(ssnContractAddr, new_address string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: ssnContractAddr,
		},
		{
			VName: "new_address",
			Type:  "ByStr20",
			Value: new_address,
		},
	}

	return s.Call("SubmitChangeZproxyTransaction", args, "0")
}

func (s *MultisigWallet) SubmitUpdateReceivingAddrTransaction(ssnContractAddr, new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: ssnContractAddr,
		},
		{
			VName: "new_addr",
			Type:  "ByStr20",
			Value: new_addr,
		},
	}

	return s.Call("SubmitUpdateReceivingAddrTransaction", args, "0")
}

func (s *MultisigWallet) SubmitUpdateCommTransaction(ssnContractAddr string, new_rate int) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: ssnContractAddr,
		},
		{
			VName: "new_rate",
			Type:  "Uint128",
			Value: strconv.Itoa(new_rate),
		},
	}

	return s.Call("SubmitUpdateCommTransaction", args, "0")
}

func (s *MultisigWallet) SubmitWithdrawCommTransaction(ssnContractAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "calleeContract",
			Type:  "ByStr20",
			Value: ssnContractAddr,
		},
	}

	return s.Call("SubmitWithdrawCommTransaction", args, "0")
}

func NewMultisigContract(sdk *AvelySDK, owners []string, requiredSignaturesCount int, deployer *account.Wallet) (*MultisigWallet, error) {
	// TOOD: add requiredSignaturesCount validation

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
			Value: strconv.Itoa(requiredSignaturesCount),
		},
	}

	contract := buildMultisigContract(sdk, init)
	contract.Signer = deployer

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
		sdkContract.ErrorCodes = sdkContract.ParseErrorCodes(contract.Code)
		return &MultisigWallet{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreMultisigContract(sdk *AvelySDK, contractAddress string) (*MultisigWallet, error) {
	contract := buildMultisigContract(sdk, []core.ContractValue{})

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
	sdkContract.ErrorCodes = sdkContract.ParseErrorCodes(contract.Code)
	return &MultisigWallet{Contract: sdkContract}, nil
}

func buildMultisigContract(sdk *AvelySDK, init []core.ContractValue) contract2.Contract {
	return Restore("multisig_wallet", sdk.InitProvider(), init)
}
