package transitions

import (
	"reflect"

	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
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

	p.ChangeSSNs()
	p.ChangeTreasuryAddress()
	p.SyncBufferAndHolder()

	p.Unpause()
	p.SetupZProxy()
	p.SetupShortcuts(log)

	return p
}

func (tr *Transitions) DeployZilSwap() *ZilSwap {
	log := GetLog()
	zilSwap, err := NewZilSwap(sdk)
	if err != nil {
		log.Fatal("deploy zilSwap error = " + err.Error())
	}

	_, err = zilSwap.Initialize()
	if err != nil {
		log.Fatal("deploy zilSwap error = " + err.Error())
	}

	log.Info("deploy zilSwap succeed, address = " + zilSwap.Addr)

	return zilSwap
}

func (tr *Transitions) DeployMultisigWallet(owners []string, signCount int) *MultisigWallet {
	log := GetLog()
	multisig, err := NewMultisigContract(sdk, owners, signCount)
	if err != nil {
		log.Fatal("deploy MultisigContract error = " + err.Error())
	}

	return multisig
}

func (tr *Transitions) FocusOn(focus string) {
	st := reflect.TypeOf(tr)
	_, exists := st.MethodByName(focus)
	if exists {
		reflect.ValueOf(tr).MethodByName(focus).Call([]reflect.Value{})
	} else {
		GetLog().Fatal("A focus test suite does not exist")
	}
}

func (tr *Transitions) RunAll() {
	tr.Owner()
	tr.DelegateStakeSuccess()
	tr.DelegateStakeBuffersRotation()
	tr.IsAdmin()
	tr.IsOwner()
	tr.IsAzil()
	tr.IsZimpl()
	tr.IsBufferOrHolder()
	tr.Pause()
	tr.PerformAutoRestake()
	tr.ChownStakeSuccess()
	tr.ChownStakeManySsnSuccess()
	tr.ChownStakeZimplErrors()
	tr.ChownStakeAzilErrors()
	tr.ChownStakeRequireDrainBuffer()
	tr.AddToSwap()
	tr.TransferFrom()
	tr.MultisigWalletTests()

	if !IsCI() {
		tr.DrainBuffer()
		tr.CompleteWithdrawalSuccess()
		tr.WithdrawStakeAmount()
	}
}
