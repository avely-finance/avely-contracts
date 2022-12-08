package utils

import (
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/account"
)

func ArrayAtoi(input []string) []int {
	out := []int{}
	for _, val := range input {
		res, _ := strconv.Atoi(val)
		out = append(out, res)
	}
	return out
}

func ArrayItoa(input []int) []string {
	out := []string{}
	for _, val := range input {
		out = append(out, strconv.Itoa(val))
	}
	return out
}

func GetAddressByWallet(wallet *account.Wallet) string {
	return "0x" + wallet.DefaultAccount.Address
}
