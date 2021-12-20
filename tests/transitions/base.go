package transitions

import (
	. "Azil/test/helpers"
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
