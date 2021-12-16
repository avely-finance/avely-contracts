package deploy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type AZil struct {
	Contract
}

func (b *AZil) ChangeZproxyAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return b.Call("ChangeZproxyAddress", args, "0")
}

func (b *AZil) ChangeZimplAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return b.Call("ChangeZimplAddress", args, "0")
}

func (a *AZil) ChangeBuffers(new_buffers []string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_buffers",
			"List ByStr20",
			new_buffers,
		},
	}
	return a.Contract.Call("ChangeBuffers", args, "0")
}

func (a *AZil) ClaimWithdrawal(ready_blocks []string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"blocks_to_withdraw",
			"List BNum",
			ready_blocks,
		},
	}
	return a.Contract.Call("ClaimWithdrawal", args, "0")
}

func (a *AZil) ChangeHolderAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return a.Contract.Call("ChangeHolderAddress", args, "0")
}

func (a *AZil) DelegateStake(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{}

	return a.Call("DelegateStake", args, amount)
}

func (a *AZil) IncreaseAutoRestakeAmount(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("IncreaseAutoRestakeAmount", args, amount)
}

func (a *AZil) PerformAutoRestake() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("PerformAutoRestake", args, "0")
}

func (a *AZil) UpdateStakingParameters(min_deleg_stake string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"min_deleg_stake",
			"Uint128",
			min_deleg_stake,
		},
	}
	return a.Call("UpdateStakingParameters", args, "0")
}

func (a *AZil) WithdrawStakeAmt(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return a.Call("WithdrawStakeAmt", args, "0")
}

func (a *AZil) DrainBuffer(buffer_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"buffer_addr",
			"ByStr20",
			"0x" + buffer_addr,
		},
	}
	return a.Call("DrainBuffer", args, "0")
}

func (a *AZil) CompleteWithdrawal() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("CompleteWithdrawal", args, "0")
}

func (a *AZil) ZilBalanceOf(addr string) (string, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + addr,
		},
	}
	tx, err := a.Contract.Call("ZilBalanceOf", args, "0")
	if err != nil {
		return "", err
	}

	/* for debug
	println("--------------------------------------------")
	data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
	println(string(data))
	println("--------------------------------------------")*/

	for _, transition := range tx.Receipt.Transitions {
		if "ZilBalanceOfCallBack" != transition.Msg.Tag {
			continue
		}
		for _, param := range transition.Msg.Params {
			if param.VName == "address" && param.Value != "0x"+addr {
				//it's balance of some other address, it should not be so
				return "", errors.New("Balance not found for addr=" + addr)
			}
			if param.VName == "balance" {
				//Value interface{} `json:"value"`
				//https://github.com/Zilliqa/gozilliqa-sdk/blob/7a254f739153c0551a327526009b4aaeeb4c9d87/core/types.go#L150
				return fmt.Sprintf("%v", param.Value), nil
			}
		}
		break
	}
	return "", errors.New("Balance not found")
}

func NewAZilContract(key, aZilSSNAddress, zproxyAddr, zimplAddr string) (*AZil, error) {
	code, _ := ioutil.ReadFile("../contracts/aZil.scilla")
	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_admin_address",
			Type:  "ByStr20",
			Value: "0x" + getAddressFromPrivateKey(key),
		}, {
			VName: "init_azil_ssn_address",
			Type:  "ByStr20",
			Value: aZilSSNAddress,
		}, {
			VName: "init_zproxy_address",
			Type:  "ByStr20",
			Value: "0x" + zproxyAddr,
		}, {
			VName: "init_zimpl_address",
			Type:  "ByStr20",
			Value: "0x" + zimplAddr,
		}, {
			VName: "init_holder_address",
			Type:  "ByStr20",
			Value: "0xb2e2c996e6068f4ae11c4cc2c6a189b774819f79",
		},
	}

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)

	contract := contract2.Contract{
		Code:   string(code),
		Init:   init,
		Signer: wallet,
	}

	tx, err := DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, TX_CONFIRM_MAX_ATTEMPTS, TX_CONFIRM_INTERVAL_SEC, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)

		stateFieldTypes := make(StateFieldTypes)
		stateFieldTypes["balances"] = "StateFieldMap"
		stateFieldTypes["buffers_addresses"] = "StateFieldArray"
		stateFieldTypes["withdrawal_pending"] = "StateFieldMapMapWithdrawal"
		stateFieldTypes["withdrawal_unbonded"] = "StateFieldMapWithdrawal"

		contract := Contract{
			Code:            string(code),
			Init:            init,
			Addr:            tx.ContractAddress,
			Bech32:          b32,
			Wallet:          wallet,
			StateFieldTypes: stateFieldTypes,
		}
		TxIdLast = tx.ID
		return &AZil{Contract: contract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		log.Println(string(data))
		return nil, errors.New("deploy failed")
	}
}
