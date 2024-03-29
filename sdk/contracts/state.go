package contracts

import (
	"math/big"
	"strings"

	"github.com/tidwall/gjson"
)

type SSNCycleInfoType struct {
	TotalStaking *big.Int
	TotalRewards *big.Int
}

type WithdrawalType struct {
	TokenAmount *big.Int
	StakeAmount *big.Int
}

type StateItem struct {
	gjson.Result
}

type State struct {
	raw string
}

func NewState(raw string) *State {
	return &State{
		raw: raw,
	}
}

func (s *State) Dig(path ...string) *StateItem {
	return &StateItem{gjson.Get(s.raw, strings.Join(path[:], "."))}
}

func (i *StateItem) SSNCycleInfo() *SSNCycleInfoType {
	totalStaking := big.NewInt(0)
	totalRewards := big.NewInt(0)

	if i.Get("arguments").Exists() {
		totalStaking = stringToBigInt(i.Get("arguments.0").String())
		totalRewards = stringToBigInt(i.Get("arguments.1").String())
	}

	return &SSNCycleInfoType{
		TotalStaking: totalStaking,
		TotalRewards: totalRewards,
	}
}

func (i *StateItem) MapAddressAmount() map[string]*big.Int {
	imap := i.Map()
	out := make(map[string]*big.Int)
	for addr, amount := range imap {
		out[addr] = stringToBigInt(amount.String())
	}
	return out
}

func (i *StateItem) ArrayInt() []int {
	iarr := i.Array()
	out := []int{}
	for _, value := range iarr {
		out = append(out, int(value.Int()))
	}
	return out
}

func (i *StateItem) ArrayString() []string {
	iarr := i.Array()
	out := []string{}
	for _, value := range iarr {
		out = append(out, value.String())
	}
	return out
}

func (i *StateItem) Withdrawal() *WithdrawalType {
	tokenAmt := big.NewInt(0)
	stakeAmt := big.NewInt(0)

	if i.Get("arguments").Exists() {
		tokenAmt = stringToBigInt(i.Get("arguments.0").String())
		stakeAmt = stringToBigInt(i.Get("arguments.1").String())
	}

	return &WithdrawalType{
		TokenAmount: tokenAmt,
		StakeAmount: stakeAmt,
	}
}

func (i *StateItem) ToTrue() bool {
	return i.Get("constructor").String() == "True"
}

func (i *StateItem) BigInt() *big.Int {
	return stringToBigInt(i.String())
}

func (i *StateItem) BigFloat() *big.Float {
	return stringToBigFloat(i.String())
}

func stringToBigInt(v string) *big.Int {
	n := new(big.Int)
	n, ok := n.SetString(v, 10)
	if !ok {
		return big.NewInt(0)
	}

	return n
}

func stringToBigFloat(v string) *big.Float {
	n := new(big.Float)
	n, ok := n.SetString(v)
	if !ok {
		return big.NewFloat(0)
	}

	return n
}

func AddBI(a, b *big.Int) *big.Int {
	return big.NewInt(0).Add(a, b)
}

func SubBI(a, b *big.Int) *big.Int {
	return big.NewInt(0).Sub(a, b)
}

func DivBF(a, b *big.Float) *big.Float {
	return big.NewFloat(0).Quo(a, b)
}

func SubOneToZero(a *big.Int) *big.Int {
	if a.Cmp(big.NewInt(1)) == 1 { // a > 1
		return SubBI(a, big.NewInt(1))
	} else {
		return big.NewInt(0)
	}
}
