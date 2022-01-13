package contracts

import (
	"encoding/json"
	"errors"

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

type Contract struct {
	Sdk      *AvelySDK
	Provider provider2.Provider
	Addr     string
	Bech32   string
	Wallet   *account.Wallet
}

func (c *Contract) UpdateWallet(newKey string) {
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(newKey)
	c.Wallet = wallet
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
		return tx, errors.New("transaction didn't get confirmed")
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
	rsp, _ := c.Provider.GetSmartContractSubState(c.Addr, params...)
	// result, _ := json.MarshalIndent(rsp, "", "     ")
	result, _ := json.Marshal(rsp)
	state := string(result)
	return state
}
