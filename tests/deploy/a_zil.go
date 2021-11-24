package deploy

import (
	//"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type AZil struct {
	Contract
}

func (a *AZil) ChangeBufferAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return a.Contract.Call("ChangeBufferAddress", args, "0")
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

func (a *AZil) IncreaseTotalStakeAmount(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			amount,
		},
	}

	return a.Call("IncreaseTotalStakeAmount", args, "0")
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

func NewAZilContract(key string, aZilSSNAddress string, stubStakingAddr string) (*AZil, error) {
	code, _ := ioutil.ReadFile("../contracts/aZil.scilla")

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_azil_ssn_address",
			Type:  "ByStr20",
			Value: aZilSSNAddress,
		}, {
			VName: "init_proxy_staking_contract_address",
			Type:  "ByStr20",
			Value: "0x" + stubStakingAddr,
		}, {
			VName: "init_buffer_address",
			Type:  "ByStr20",
			Value: "0xb2e2c996e6068f4ae11c4cc2c6a189b774819f79",
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
	tx.Confirm(tx.ID, 1, 1, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)
		contract := Contract{
			Code:   string(code),
			Init:   init,
			Addr:   tx.ContractAddress,
			Bech32: b32,
			Wallet: wallet,
		}

		return &AZil{Contract: contract}, nil
	} else {
		return nil, errors.New("deploy failed")
	}
}
