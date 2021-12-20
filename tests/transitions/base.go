package transitions

import (
	. "Azil/test/helpers"
	"fmt"
)

var t *Testing
var log *Log

func init() {
	t = NewTesting()
	log = GetLog()
}

type Transitions struct {
	cfg Config
}

func NewTransitions(config Config) *Transitions {
	return &Transitions{
		cfg: config,
	}
}

const qa = "000000000000"

func zil(amount int) string {
	if amount == 0 {
		return "0"
	}
	return fmt.Sprintf("%d%s", amount, qa)
}

func azil(amount int) string {
	if amount == 0 {
		return "0"
	}
	return fmt.Sprintf("%d%s", amount, qa)
}
