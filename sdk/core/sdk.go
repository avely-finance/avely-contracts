package core

import (
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	transaction2 "github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/ybbus/jsonrpc"
	"strconv"
	"log"
)

type AvelySDK struct {
	Cfg Config
	TxIdLast string
}

func NewAvelySDK(config Config) *AvelySDK {
	return &AvelySDK{
		Cfg: config,
		TxIdLast: "",
	}
}

func (sdk *AvelySDK) IncreaseBlocknum(delta int32) {
	//https://raw.githubusercontent.com/Zilliqa/gozilliqa-sdk/7a254f739153c0551a327526009b4aaeeb4c9d87/provider/provider.go

	if (sdk.Cfg.Chain != "local") {
		log.Fatalf("Increasing block number available only for the local blockchain")
	}

	rpcClient := jsonrpc.NewClient(sdk.Cfg.ApiUrl)
	params := []interface{}{delta}
	rpcClient.Call("IncreaseBlocknum", params)
}

func (sdk *AvelySDK) GetBalance(addr string) string {
	provider := provider2.NewProvider(sdk.Cfg.ApiUrl)
	balAndNonce, err := provider.GetBalance(addr)
	if err != nil {
		panic(err)
	}
	return balAndNonce.Balance
}

func (sdk *AvelySDK) GetAddressFromPrivateKey(privateKey string) string {
	publicKey := keytools.GetPublicKeyFromPrivateKey(util.DecodeHex(privateKey), true)
	address := keytools.GetAddressFromPublic(publicKey)
	return address
}

func (sdk *AvelySDK) DeployTo(c *contract2.Contract) (*transaction2.Transaction, error) {
	c.Provider = provider2.NewProvider(sdk.Cfg.ApiUrl)
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
	tx, err := c.Deploy(parameter)

	sdk.TxIdLast = tx.ID

	return tx, err
}

func (sdk *AvelySDK) CallFor(c *contract2.Contract, transition string, args []core.ContractValue, priority bool, amount string) (*transaction2.Transaction, error) {
	c.Provider = provider2.NewProvider(sdk.Cfg.ApiUrl)
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
	tx, err := c.Call(transition, args, params, priority)

	sdk.TxIdLast = tx.ID

	return tx, err
}