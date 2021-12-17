package contracts

import (
	"Azil/test/helpers"

	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type HolderContract struct {
	Contract
}

func (b *HolderContract) AddFunds(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("AddFunds", args, amount)
}

func (b *HolderContract) CompleteWithdrawalNoUnbondedStakeCallBack(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return b.Call("CompleteWithdrawalNoUnbondedStakeCallBack", args, "0")
}

func (b *HolderContract) CompleteWithdrawalSuccessCallBack(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return b.Call("CompleteWithdrawalSuccessCallBack", args, "0")
}

func (b *HolderContract) DelegateStakeSuccessCallBack(ssnaddr, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return b.Call("DelegateStakeSuccessCallBack", args, "0")
}

func (b *HolderContract) WithdrawStakeRewardsSuccessCallBack(ssnaddr, rewards string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			"0x" + ssnaddr,
		},
		{
			"rewards",
			"Uint128",
			rewards,
		},
	}
	return b.Call("WithdrawStakeRewardsSuccessCallBack", args, "0")
}

func (b *HolderContract) WithdrawStakeAmtSuccessCallBack(ssnaddr, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			"0x" + ssnaddr,
		},
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return b.Call("WithdrawStakeAmtSuccessCallBack", args, "0")
}

func (b *HolderContract) WithdrawStakeAmt(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return b.Call("WithdrawStakeAmt", args, "0")
}

func (b *HolderContract) ChangeZproxyAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return b.Call("ChangeZproxyAddress", args, "0")
}

func (b *HolderContract) ChangeZimplAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return b.Call("ChangeZimplAddress", args, "0")
}

func (b *HolderContract) ChangeAzilSSNAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return b.Call("ChangeAzilSSNAddress", args, "0")
}

func (b *HolderContract) ChangeAimplAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return b.Call("ChangeAimplAddress", args, "0")
}

func (b *HolderContract) CompleteWithdrawal() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("CompleteWithdrawal", args, "0")
}

func (b *HolderContract) ClaimRewards() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("ClaimRewards", args, "0")
}

func (b *HolderContract) ConfirmDelegatorSwap(requestor string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"requestor",
			"ByStr20",
			"0x" + requestor,
		},
	}
	return b.Call("ConfirmDelegatorSwap", args, "0")
}

func (b *HolderContract) DelegateStake(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return b.Call("DelegateStake", args, amount)
}

func NewHolderContract(key, aimplAddr, azilSsnAddr, zproxyAddr, zimplAddr string) (*HolderContract, error) {
	code, _ := ioutil.ReadFile("../contracts/holder.scilla")

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_admin_address",
			Type:  "ByStr20",
			Value: "0x" + helpers.GetAddressFromPrivateKey(key),
		}, {
			VName: "init_aimpl_address",
			Type:  "ByStr20",
			Value: "0x" + aimplAddr,
		}, {
			VName: "init_azil_ssn_address",
			Type:  "ByStr20",
			Value: azilSsnAddr,
		}, {
			VName: "init_zproxy_address",
			Type:  "ByStr20",
			Value: "0x" + zproxyAddr,
		}, {
			VName: "init_zimpl_address",
			Type:  "ByStr20",
			Value: "0x" + zimplAddr,
		},
	}

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)

	contract := contract2.Contract{
		Code:   string(code),
		Init:   init,
		Signer: wallet,
	}

	tx, err := helpers.DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, TX_CONFIRM_MAX_ATTEMPTS, TX_CONFIRM_INTERVAL_SEC, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)
		stateFieldTypes := make(StateFieldTypes)
		contract := Contract{
			Provider:        *contract.Provider,
			Code:            string(code),
			Init:            init,
			Addr:            tx.ContractAddress,
			Bech32:          b32,
			Wallet:          wallet,
			StateFieldTypes: stateFieldTypes,
		}
		TxIdLast = tx.ID
		return &HolderContract{Contract: contract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		log.Println(string(data))
		return nil, errors.New("deploy failed")
	}
}
