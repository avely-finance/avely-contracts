package helpers

import (
	"os"

	"github.com/avely-finance/avely-contracts/sdk/contracts"
)

func Field(c contracts.ProtocolContract, path ...string) string {
	item := contracts.NewState(c.State()).Dig(path...)

	if item.Get("constructor").Exists() {
		if item.Get("constructor").String() == "True" {
			return "True"
		}
		if item.Get("constructor").String() == "False" {
			return "False"
		}
	}

	if item.Get("arguments").Exists() {
		return item.Get("arguments.0").String()
	}

	return item.String()
}

func Dig(c contracts.ProtocolContract, path ...string) *contracts.StateItem {
	return contracts.NewState(c.State()).Dig(path...)
}

func IsCI() bool {
	return os.Getenv("CI") == "1"
}
