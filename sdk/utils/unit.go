package utils

import (
	"fmt"
	"math/big"
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

func FromQA(qa string) string {
	QA, _ := new(big.Int).SetString(qa, 10)
	ZIL := QA.Div(QA, big.NewInt(1000000000000))
	return ZIL.String()
}
