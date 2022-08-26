package core

import (
	"log"
	"net/url"
	"runtime"
	"strconv"
	"time"

	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	transaction2 "github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/ybbus/jsonrpc"
)

const ZeroAddr = "0x0000000000000000000000000000000000000000"

type AvelySDK struct {
	Cfg Config
}

func NewAvelySDK(config Config) *AvelySDK {
	return &AvelySDK{
		Cfg: config,
	}
}

func (sdk *AvelySDK) InitProvider() *provider2.Provider {
	return provider2.NewProvider(sdk.Cfg.Api.HttpUrl)
}

func (sdk *AvelySDK) GetWsURL() url.URL {
	return url.URL{
		Scheme: sdk.Cfg.Api.WebsocketSchema,
		Host:   sdk.Cfg.Api.WebsocketUrl,
		Path:   "",
	}
}

// IncreaseBlocknum can be called if isolated server works in "manual" mode:
// https://github.com/Zilliqa/zilliqa-isolated-server#running-the-isolated-server-with-manual-block-increase
func (sdk *AvelySDK) IncreaseBlocknum(delta int) error {
	//local blockchain
	if sdk.Cfg.ChainId == 222 {
		rpcClient := jsonrpc.NewClient(sdk.Cfg.Api.HttpUrl)
		params := []interface{}{delta}
		tx, err := rpcClient.Call("IncreaseBlocknum", params)

		if err != nil {
			return err
		} else if tx.Error != nil {
			return tx.Error
		}

		return nil
	}

	//for testnet/mainnet will wait required delta blocks
	blockHeight := sdk.GetBlockHeight()
	for {
		time.Sleep(5 * time.Second)
		cur := sdk.GetBlockHeight()
		log.Printf("sdk.IncreaseBlocknum(): start block=%d, delta=%d, cur. block=%d", blockHeight, delta, cur)
		if cur >= blockHeight+delta {
			return nil
		}
	}
}

//this function will only works for testnet or mainnet
//it will not work for local server in manual mode
func (sdk *AvelySDK) GetBlockHeight() int {
	provider := sdk.InitProvider()
	result, _ := provider.GetNumTxBlocks()
	blockHeight, _ := strconv.Atoi(result)
	return blockHeight
}

func (sdk *AvelySDK) GetBalance(addr string) string {
	provider := sdk.InitProvider()
	if addr[0:2] == "0x" {
		addr = addr[2:]
	}
	balAndNonce, err := provider.GetBalance(addr)
	if err != nil {
		panic(err)
	}
	return balAndNonce.Balance
}

func (sdk *AvelySDK) GetAddressFromPrivateKey(privateKey string) string {
	publicKey := keytools.GetPublicKeyFromPrivateKey(util.DecodeHex(privateKey), true)
	address := keytools.GetAddressFromPublic(publicKey)
	return "0x" + address
}

func (sdk *AvelySDK) DeployTo(c *contract2.Contract) (*transaction2.Transaction, error) {
	c.Provider = sdk.InitProvider()
	gasPrice, err := c.Provider.GetMinimumGasPrice()
	if err != nil {
		return nil, err
	}
	parameter := contract2.DeployParams{
		Version:      strconv.FormatInt(int64(util.Pack(sdk.Cfg.ChainId, 1)), 10),
		Nonce:        "",
		GasPrice:     gasPrice,
		GasLimit:     "75000",
		SenderPubKey: "",
	}
	tx, err := c.Deploy(parameter)

	return tx, err
}

func (sdk *AvelySDK) CallFor(c *contract2.Contract, transition string, args []core.ContractValue, priority bool, amount string) (*transaction2.Transaction, error) {
	c.Provider = sdk.InitProvider()
	gasPrice, err := c.Provider.GetMinimumGasPrice()
	if err != nil {
		return nil, err
	}
	params := contract2.CallParams{
		Version:      strconv.FormatInt(int64(util.Pack(sdk.Cfg.ChainId, 1)), 10),
		Nonce:        "",
		GasPrice:     gasPrice,
		GasLimit:     "40000",
		Amount:       amount,
		SenderPubKey: "",
	}
	tx, err := c.Call(transition, args, params, priority)

	return tx, err
}

func CheckTx(tx *transaction.Transaction, err error) (*transaction.Transaction, error) {
	if err != nil {
		_, file, no, _ := runtime.Caller(1)
		log.Fatal("TRANSACTION FAILED, " + file + ":" + strconv.Itoa(no))
	}
	return tx, err
}
