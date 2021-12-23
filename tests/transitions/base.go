package transitions

import (
	. "github.com/avely-finance/avely-contracts/tests/helpers"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"reflect"
	"log"
)

var t *Testing
// var log Log
var sdk *AvelySDK

func InitTransitions(sdkValue *AvelySDK, testingValue *Testing) *Transitions {
	t = testingValue
	// log = testingValue.log
	sdk = sdkValue

	return NewTransitions()
}

type Transitions struct {
	// cfg Config
}

func NewTransitions() *Transitions {
	return &Transitions{
		// cfg: config,
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
	// tr.DelegateStakeSuccess()
	// tr.DelegateStakeBuffersRotation()
	// tr.WithdrawStakeAmount()
	// tr.CompleteWithdrawalSuccess()
	// tr.ZilBalanceOf()
	// tr.IsAdmin()
	tr.IsAimpl()
	// tr.IsZimpl()
	// tr.IsBufferOrHolder()
	// tr.DrainBuffer()
	// tr.PerformAuoRestake()
}

func DeployAndUpgrade() (*contracts.Protocol) {
	return Deploy(sdk, t.Log)
}