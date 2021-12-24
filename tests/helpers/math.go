package helpers

import (
	"fmt"
	"math/big"
)

func StrAdd(arg ...string) string {
	if len(arg) < 2 {
		panic("StrAdd needs at least 2 arguments")
	}
	result, _ := new(big.Int).SetString("0", 10)
	for _, v := range arg {
		vInt, ok := new(big.Int).SetString(v, 10)
		if !ok {
			panic(fmt.Sprintf("StrAdd can't get BigInt from argument ", v))
		}
		result = result.Add(result, vInt)
	}
	return result.String()
}

func StrSub(a, b string) string {
	A, _ := new(big.Int).SetString(a, 10)
	B, _ := new(big.Int).SetString(b, 10)
	result := new(big.Int).Sub(A, B)
	return result.String()
}

func StrMulDiv(a, b, c string) string {
	A, _ := new(big.Int).SetString(a, 10)
	B, _ := new(big.Int).SetString(b, 10)
	C, _ := new(big.Int).SetString(c, 10)
	result := new(big.Int).Mul(A, B)
	result = result.Div(result, C)
	return result.String()
}

