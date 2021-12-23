package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
	"log"
	"reflect"
)

var t *Testing
var sdk *AvelySDK

func InitTransitions(sdkValue *AvelySDK, testingValue *Testing) *Transitions {
	t = testingValue
	sdk = sdkValue

	return NewTransitions()
}

type Transitions struct {
}

func NewTransitions() *Transitions {
	return &Transitions{}
}

func (tr *Transitions) DeployAndUpgrade() *Protocol {
	p := Deploy(sdk, t.Log)

	p.SyncBufferAndHolder()
	p.SetupZProxy()
	p.SetupShortcuts(t.Log)

	return p
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
	tr.IsBufferOrHolder()
	tr.DrainBuffer()
	tr.PerformAuoRestake()
}
