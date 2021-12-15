package deploy

import (
	"fmt"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	transaction2 "github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/ybbus/jsonrpc"
	"log"
	"math/big"
	"strconv"
)

func IncreaseBlocknum(delta int32) {
	//https://raw.githubusercontent.com/Zilliqa/gozilliqa-sdk/7a254f739153c0551a327526009b4aaeeb4c9d87/provider/provider.go
	//TODO singleton
	rpcClient := jsonrpc.NewClient(API_PROVIDER)
	params := []interface{}{delta}
	rpcClient.Call("IncreaseBlocknum", params)
	log.Printf("🔗  === Blocknumber increased by %d === \n", delta)
}

func GetBalance(addr string) string {
	provider := provider2.NewProvider(API_PROVIDER)
	balAndNonce, err := provider.GetBalance(addr)
	if err != nil {
		panic(err)
	}
	return balAndNonce.Balance
}

func getAddressFromPrivateKey(privateKey string) string {
	publicKey := keytools.GetPublicKeyFromPrivateKey(util.DecodeHex(privateKey), true)
	address := keytools.GetAddressFromPublic(publicKey)
	return address
}

func DeployTo(c *contract2.Contract) (*transaction2.Transaction, error) {
	c.Provider = provider2.NewProvider(API_PROVIDER)
	gasPrice, err := c.Provider.GetMinimumGasPrice()
	if err != nil {
		return nil, err
	}
	parameter := contract2.DeployParams{
		Version:      strconv.FormatInt(int64(util.Pack(222, 1)), 10),
		Nonce:        "",
		GasPrice:     gasPrice,
		GasLimit:     "75000",
		SenderPubKey: "",
	}
	return c.Deploy(parameter)
}

func CallFor(c *contract2.Contract, transition string, args []core.ContractValue, priority bool, amount string) (*transaction2.Transaction, error) {
	c.Provider = provider2.NewProvider(API_PROVIDER)
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

func StrAdd(arg ...string) string {
	if len(arg) < 2 {
		panic("StrAdd needs at least 2 arguments")
	}
	result, _ := new(big.Int).SetString("0", 10)
	for _, v := range arg {
		vInt, ok := new(big.Int).SetString(v, 10)
		if !ok {
			println(v)
			panic(fmt.Sprintf("StrAdd can't get BigInt from argument ", v))
		}
		result = result.Add(result, vInt)
	}
	return result.String()
}

func StrSub(a, b string) string {
	A, _ := new(big.Int).SetString(a, 10)
	B, _ := new(big.Int).SetString(b, 10)
	result := new(big.Int).Sub(A, B)
	return result.String()
}

func StrMulDiv(a, b, c string) string {
	A, _ := new(big.Int).SetString(a, 10)
	B, _ := new(big.Int).SetString(b, 10)
	C, _ := new(big.Int).SetString(c, 10)
	result := new(big.Int).Mul(A, B)
	result = result.Div(result, C)
	return result.String()
}
