package deploy

import (
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	transaction2 "github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/ybbus/jsonrpc"
	"log"
	"strconv"
)

func IncreaseBlocknum(delta int32) {
	//https://raw.githubusercontent.com/Zilliqa/gozilliqa-sdk/7a254f739153c0551a327526009b4aaeeb4c9d87/provider/provider.go
	//TODO singleton
	rpcClient := jsonrpc.NewClient("http://zilliqa_server:5555")
	params := []interface{}{delta}
	rpcClient.Call("IncreaseBlocknum", params)
	log.Printf("🔗  === Blocknumber increased by %d === \n", delta)
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

func CallFor(c *contract2.Contract, transition string, args []core.ContractValue, priority bool, amount string) (*transaction2.Transaction, error) {
	c.Provider = provider2.NewProvider("http://zilliqa_server:5555")
	gasPrice, err := c.Provider.GetMinimumGasPrice()
	if err != nil {
		return nil, err
	}
	params := contract2.CallParams{
		Version:      strconv.FormatInt(int64(util.Pack(222, 1)), 10),
		Nonce:        "",
		GasPrice:     gasPrice,
		GasLimit:     "40000",
		Amount:       amount,
		SenderPubKey: "",
	}
	return c.Call(transition, args, params, priority)
}

func StrSum(s1, s2 string) string {
	x, _ := strconv.Atoi(s1)
	y, _ := strconv.Atoi(s2)
	return strconv.Itoa(x + y)
}
