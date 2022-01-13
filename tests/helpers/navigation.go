package helpers

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
)

func Dig(c contracts.ProtocolContract, path ...string) *contracts.StateItem{
	return contracts.NewState(c.State()).Dig(path...)
}
