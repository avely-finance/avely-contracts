package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
	"reflect"
)

var sdk *AvelySDK

func InitTransitions(sdkValue *AvelySDK) *Transitions {
	sdk = sdkValue

	return NewTransitions()
}

type Transitions struct {
}

func NewTransitions() *Transitions {
	return &Transitions{}
}

func (tr *Transitions) DeployAndUpgrade() *Protocol {
	log := GetLog()
	p := Deploy(sdk, log)

	p.SyncBufferAndHolder()
	p.SetupZProxy()
	p.SetupShortcuts(log)
	p.Unpause()

	return p
}

func (tr *Transitions) FocusOn(focus string) {
	st := reflect.TypeOf(tr)
	_, exists := st.MethodByName(focus)
	if exists {
		reflect.ValueOf(tr).MethodByName(focus).Call([]reflect.Value{})
	} else {
		GetLog().Fatal(" A focus test suite does not exist")
	}
}

func (tr *Transitions) RunAll() {
	tr.Admin()
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
	tr.Proxy()
}
