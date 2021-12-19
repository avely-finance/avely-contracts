package transitions

import (
	"Azil/test/helpers"
	"fmt"
)

var t *helpers.Testing
var log *helpers.Log

func init() {
	t = helpers.NewTesting()
	log = helpers.GetLog()
}

type Transitions struct {
	cfg helpers.Config
}

func NewTransitions(config helpers.Config) *Transitions {
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
