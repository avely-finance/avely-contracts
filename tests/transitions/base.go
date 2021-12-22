package transitions

import (
	. "Azil/test/helpers"
	"reflect"
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

func (tr *Transitions) FocusOn(focus string) {
	st := reflect.TypeOf(tr)
	_, exists := st.MethodByName(focus)
	if exists {
		reflect.ValueOf(tr).MethodByName(focus).Call([]reflect.Value{})
	} else {
		log.Fatal(" A focus test suite does not exist")
	}
}

func (tr *Transitions) RunAll() {
	tr.DelegateStakeSuccess()
	tr.DelegateStakeBuffersRotation()
	tr.WithdrawStakeAmount()
	tr.CompleteWithdrawalSuccess()
	tr.ZilBalanceOf()
	tr.IsAdmin()
	tr.IsAimpl()
	tr.IsZimpl()
	tr.DrainBuffer()
	tr.PerformAuoRestake()
}
