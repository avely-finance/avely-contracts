package contracts

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	. "github.com/avely-finance/avely-contracts/sdk/core"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type Zproxy struct {
	Contract
}

func (p *Zproxy) AssignStakeReward(ssn, reward string) (*transaction.Transaction, error) {

	type Constructor struct {
		Constructor string   `json:"constructor"`
		ArgTypes    []string `json:"argtypes"`
		Arguments   []string `json:"arguments"`
	}

	ats := []string{
		"ByStr20",
		"Uint128",
	}

	ars := []string{
		ssn,
		reward,
	}

	args := []core.ContractValue{
		{
			VName: "ssnreward_list",
			Type:  "List (Pair ByStr20 Uint128)",
			Value: []Constructor{
				{
					Constructor: "Pair",
					ArgTypes:    ats,
					Arguments:   ars,
				},
			},
		},
	}

	// we send reward as ZIL amount because AssignStake works with only 1 SSN
	return p.Call("AssignStakeReward", args, reward)
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

func (p *Zproxy) UpdateVerifier(addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{{
		"verif",
		"ByStr20",
		addr,
	}}
	return p.Call("UpdateVerifier", args, "0")

}

func (p *Zproxy) UpdateVerifierRewardAddr(newAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{{
		"addr",
		"ByStr20",
		newAddr,
	}}
	return p.Call("UpdateVerifierRewardAddr", args, "0")
}

func NewZproxy(sdk *AvelySDK) (*Zproxy, error) {
	contract := buildZproxyContract(sdk)

	tx, err := sdk.DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, sdk.Cfg.TxConfrimMaxAttempts, sdk.Cfg.TxConfirmIntervalSec, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)

		stateFieldTypes := make(StateFieldTypes)

		sdkContract := Contract{
			Sdk:             sdk,
			Provider:        *contract.Provider,
			Addr:            "0x" + tx.ContractAddress,
			Bech32:          b32,
			Wallet:          contract.Signer,
			StateFieldTypes: stateFieldTypes,
		}

		return &Zproxy{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreZproxy(sdk *AvelySDK, contractAddress string) (*Zproxy, error) {
	contract := buildZproxyContract(sdk)

	b32, err := bech32.ToBech32Address(contractAddress)

	if err != nil {
		return nil, errors.New("Config has invalid Zproxy address")
	}

	stateFieldTypes := make(StateFieldTypes)

	sdkContract := Contract{
		Sdk:             sdk,
		Provider:        *contract.Provider,
		Addr:            contractAddress,
		Bech32:          b32,
		Wallet:          contract.Signer,
		StateFieldTypes: stateFieldTypes,
	}

	return &Zproxy{Contract: sdkContract}, nil
}

func buildZproxyContract(sdk *AvelySDK) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/zilliqa_staking/proxy.scilla")
	key := sdk.Cfg.AdminKey

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_admin",
			Type:  "ByStr20",
			Value: sdk.GetAddressFromPrivateKey(key),
		}, {
			VName: "init_implementation",
			Type:  "ByStr20",
			Value: sdk.GetAddressFromPrivateKey(key),
		},
	}

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)

	return contract2.Contract{
		Provider: provider2.NewProvider(sdk.Cfg.ApiUrl),
		Code:     string(code),
		Init:     init,
		Signer:   wallet,
	}
}
