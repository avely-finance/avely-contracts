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
	// "github.com/Zilliqa/gozilliqa-sdk/keytools"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
)

type AZil struct {
	Code   string
	Init   []core.ContractValue
	Addr   string
	Bech32 string
	Wallet *account.Wallet
}

func (a *AZil) LogContractStateJson() string {
	provider := provider2.NewProvider("http://zilliqa_server:5555")
	rsp, _ := provider.GetSmartContractState(a.Addr)
	j, _ := json.Marshal(rsp)
	a.LogPrettyStateJson(rsp)
	return string(j)
}

func (a *AZil) LogPrettyStateJson(data interface{}) {
	j, _ := json.MarshalIndent(data, "", "   ")
	log.Println(string(j))
}

func (a *AZil) GetBalance() string {
	provider := provider2.NewProvider("http://zilliqa_server:5555")
	balAndNonce, _ := provider.GetBalance(a.Addr)
	return balAndNonce.Balance
}

func (a *AZil) ChangeBufferAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return a.Call("ChangeBufferAddress", args, "0")
}

func (a *AZil) DelegateStake(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{}

	return a.Call("DelegateStake", args, amount)
}

func (a *AZil) Call(transition string, params []core.ContractValue, amount string) (*transaction.Transaction, error) {
	contract := contract2.Contract{
		Address: a.Bech32,
		Signer:  a.Wallet,
	}

	tx, err := CallFor(&contract, transition, params, false, amount)
	if err != nil {
		return tx, err
	}
	tx.Confirm(tx.ID, 1, 1, contract.Provider)
	if tx.Status != core.Confirmed {
		return tx, errors.New("transaction didn't get confirmed")
	}
	if !tx.Receipt.Success {
		return tx, errors.New("transaction failed")
	}
	return tx, nil
}

func NewAZilContract(key string, aZilSSNAddress string, stubStakingAddr string) (*AZil, error) {
	code, _ := ioutil.ReadFile("../contracts/aZil.scilla")
	// adminAddr := keytools.GetAddressFromPrivateKey(util.DecodeHex(key))

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

		return &AZil{
			Code:   string(code),
			Init:   init,
			Addr:   tx.ContractAddress,
			Bech32: b32,
			Wallet: wallet,
		}, nil
	} else {
		return nil, errors.New("deploy failed")
	}
}
