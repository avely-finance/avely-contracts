package utils

import (
	"fmt"
)

const qa = "000000000000"

func ToZil(amount int) string {
	if amount == 0 {
		return "0"
	}
	return fmt.Sprintf("%d%s", amount, qa)
}

func ToStZil(amount int) string {
	return ToZil(amount)
}

func ToQA(amount int) string {
	return ToZil(amount)
}
