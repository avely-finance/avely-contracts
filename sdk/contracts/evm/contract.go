package evm

import (
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/ethereum/go-ethereum/accounts"
)

// Basic type for all protocol contracts
type ProtocolContract interface {
	//State() string
}

type Contract struct {
	Sdk     *core.AvelySDK
	Account *accounts.Account
	Addr    string
}

func (c *Contract) SetSigner(account *accounts.Account) *Contract {
	c.Account = account

	return c
}
