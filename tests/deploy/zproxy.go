package deploy

import (
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

type Zproxy struct {
	Contract
}

type SSNRewardShare struct {
	SSNAddress       string
	RewardPercentage string
}

func (p *Zproxy) AssignStakeReward(ssn, percent string) (*transaction.Transaction, error) {
	args := []core.ContractValue{{
		VName: "ssnreward_list",
		Type:  "List SsnRewardShare",
		Value: []core.ParamConstructor{
			{
				"SsnRewardShare",
				make([]interface{}, 0),
				[]string{ssn, percent},
			},
		},
	}}

	return p.Call("AssignStakeReward", args, percent)
}

func (p *Zproxy) AddSSN(addr string, name string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			addr,
		},
		{
			"name",
			"String",
			name,
		},
		{
			"urlraw",
			"String",
			"fakeurl",
		},
		{
			"urlapi",
			"String",
			"fakeapi",
		},
		{
			"comm",
			"Uint128",
			"0",
		},
	}

	return p.Call("AddSSN", args, "0")
}

func (p *Zproxy) Unpause() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return p.Call("UnPause", args, "0")
}

func (p *Zproxy) UpdateStakingParameters(min, delegmin string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"min_stake",
			"Uint128",
			min,
		},
		{
			"min_deleg_stake",
			"Uint128",
			delegmin,
		},
		{
			"max_comm_change_rate",
			"Uint128",
			"20",
		},
	}
	return p.Call("UpdateStakingParameters", args, "0")
}

func NewZproxy(key string) (*Zproxy, error) {
	code, _ := ioutil.ReadFile("../contracts/zilliqa_staking/proxy.scilla")

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_admin",
			Type:  "ByStr20",
			Value: "0x" + getAddressFromPrivateKey(key),
		}, {
			VName: "init_implementation",
			Type:  "ByStr20",
			Value: "0x" + getAddressFromPrivateKey(key),
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

		contract := Contract{
			Code:            string(code),
			Init:            init,
			Addr:            tx.ContractAddress,
			Bech32:          b32,
			Wallet:          wallet,
			StateFieldTypes: stateFieldTypes,
		}
		TxIdLast = tx.ID

		return &Zproxy{Contract: contract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		log.Println(string(data))
		return nil, errors.New("deploy failed")
	}
}
