package contracts

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	embedcontracts "github.com/avely-finance/avely-contracts/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

// Basic type for all protocol contracts
type ProtocolContract interface {
	State() string
}

type ContractErrorCodes map[string]string

type Contract struct {
	Sdk        *AvelySDK
	Provider   provider2.Provider
	Addr       string
	Bech32     string
	Wallet     *account.Wallet
	ErrorCodes ContractErrorCodes
}

func Restore(name string, provider *provider2.Provider, init []core.ContractValue) contract2.Contract {
	code, _ := embedcontracts.GetContractFs().ReadFile("source/" + name + ".scilla")

	return contract2.Contract{
		Provider: provider,
		Code:     string(code),
		Init:     init,
		Signer:   nil,
	}
}

func (c *Contract) SetSigner(wallet *account.Wallet) *Contract {
	c.Wallet = wallet

	return c
}

func (c *Contract) Call(transition string, params []core.ContractValue, amount string) (*transaction.Transaction, error) {
	contract := contract2.Contract{
		Address: c.Bech32,
		Signer:  c.Wallet,
	}

	tx, err := c.Sdk.CallFor(&contract, transition, params, false, amount)
	if err != nil {
		return tx, err
	}

	tx.Confirm(tx.ID, c.Sdk.Cfg.TxConfrimMaxAttempts, c.Sdk.Cfg.TxConfirmIntervalSec, contract.Provider)
	if tx.Status != core.Confirmed {
		err := errors.New("transaction didn't get confirmed")
		if len(tx.Receipt.Exceptions) > 0 {
			err = errors.New(tx.Receipt.Exceptions[0].Message)
		}
		return tx, err
	}
	if !tx.Receipt.Success {
		return tx, errors.New("transaction failed")
	}
	return tx, nil
}

func (c *Contract) State() string {
	rsp, _ := c.Provider.GetSmartContractState(c.Addr[2:])
	result, _ := json.MarshalIndent(rsp.Result, "", "     ")
	state := string(result)
	return state
}

func (c *Contract) SubState(params ...interface{}) string {
	rsp, _ := c.Provider.GetSmartContractSubState(c.Bech32, params...)
	state := string(rsp)

	return state
}

func (c *Contract) BuildBatchParams(fields []string) [][]interface{} {
	var params [][]interface{}

	for _, v := range fields {
		params = append(params, []interface{}{v, []string{}})
	}

	return params
}

// Inspired by https://github.com/Zilliqa/gozilliqa-sdk/blob/2ff222c97fc6fa2855ef2c5bffbd56faddd6291f/provider/provider.go#L877
//
// To build params use BuildBatchParams or:
//
//	var params [][]interface{}
//	params = append(params, []interface{}{"total_supply", []string{}})
//	params = append(params, []interface{}{"totalstakeamount", []string{}})
func (c *Contract) BatchSubState(params [][]interface{}) (string, error) {
	//we should hack here for now
	type req struct {
		Id      string      `json:"id"`
		Jsonrpc string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params"`
	}

	reqs := []*req{}

	for i, param := range params {
		p := []interface{}{
			c.Bech32,
		}

		for _, v := range param {
			p = append(p, v)
		}

		r := &req{
			Id:      strconv.Itoa(i + 1),
			Jsonrpc: "2.0",
			Method:  "GetSmartContractSubState",
			Params:  p,
		}

		reqs = append(reqs, r)
	}

	b, _ := json.Marshal(reqs)
	reader := bytes.NewReader(b)
	request, err := http.NewRequest("POST", c.Sdk.Cfg.Api.HttpUrl, reader)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (c *Contract) ParseErrorCodes(contractCode string) ContractErrorCodes {
	codes := make(ContractErrorCodes)
	re := regexp.MustCompilePOSIX(`\| *?([A-Za-z0-9]+) *?=> *?Int32 *?([0-9-]+)\n`)
	results := re.FindAllStringSubmatch(contractCode, -1)
	for _, row := range results {
		codes[row[1]] = row[2]
	}
	return codes
}

func (c *Contract) ErrorCode(errorName string) string {
	return c.ErrorCodes[errorName]
}
