package deploy

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	transaction2 "github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	// "github.com/Zilliqa/gozilliqa-sdk/keytools"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"strconv"
)

type StubStakingContract struct {
	Code string
	Init []core.ContractValue
	Addr string
}

func (s *StubStakingContract) LogContractStateJson() string {
	provider := provider2.NewProvider("https://zilliqa-isolated-server.zilliqa.com/")
	rsp, _ := provider.GetSmartContractState(s.Addr)
	j, _ := json.Marshal(rsp)
	s.LogPrettyStateJson(rsp)
	return string(j)
}

func (s *StubStakingContract) LogPrettyStateJson(data interface{}) {
	j, _ := json.MarshalIndent(data, "", "   ")
	log.Println(string(j))
}

func (s *StubStakingContract) GetBalance() string {
	provider := provider2.NewProvider("https://zilliqa-isolated-server.zilliqa.com/")
	balAndNonce, _ := provider.GetBalance(s.Addr)
	return balAndNonce.Balance
}

func NewStubStakingContract(key string) (*StubStakingContract, error) {
	code, _ := ioutil.ReadFile("../contracts/stubStakingContract.scilla")
	// adminAddr := keytools.GetAddressFromPrivateKey(util.DecodeHex(key))

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		},
	}
	// 		VName: "init_admin",
	// 		Type:  "ByStr20",
	// 		Value: "0x" + adminAddr,
	// 	}, {
	// 		VName: "init_proxy_address",
	// 		Type:  "ByStr20",
	// 		Value: "0x" + proxy,
	// 	},
	// 	{
	// 		VName: "init_gzil_address",
	// 		Type:  "ByStr20",
	// 		Value: "0x" + adminAddr,
	// 	},
	// }

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
		return &StubStakingContract{
			Code: string(code),
			Init: init,
			Addr: tx.ContractAddress,
		}, nil
	} else {
		return nil, errors.New("deploy failed")
	}
}

func DeployTo(c *contract2.Contract) (*transaction2.Transaction, error) {
	c.Provider = provider2.NewProvider("http://zilliqa_server:5555")
	gasPrice, err := c.Provider.GetMinimumGasPrice()
	if err != nil {
		return nil, err
	}
	parameter := contract2.DeployParams{
		Version:      strconv.FormatInt(int64(util.Pack(222, 1)), 10),
		Nonce:        "",
		GasPrice:     gasPrice,
		GasLimit:     "40000",
		SenderPubKey: "",
	}
	return c.Deploy(parameter)
}
