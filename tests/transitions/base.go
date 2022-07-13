package transitions

import (
	"reflect"

	"github.com/avely-finance/avely-contracts/sdk/actions"
	"github.com/avely-finance/avely-contracts/sdk/contracts"
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
	SetupZilliqaStaking(p)

	//add buffers to protocol, we need 3
	buffer2, _ := p.DeployBuffer()
	buffer3, _ := p.DeployBuffer()
	p.Buffers = append(p.Buffers, buffer2, buffer3)

	p.AddSSNs()
	p.ChangeTreasuryAddress()
	p.SyncBufferAndHolder()
	p.Unpause()
	p.InitHolder()

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

func (tr *Transitions) NextCycle(p *contracts.Protocol) {
	sdk.IncreaseBlocknum(2)
	prevWallet := p.Zproxy.Contract.Wallet

	p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)

	zimplSsnList := p.Zimpl.GetSsnList()
	ssnRewardFactor := make(map[string]string)
	for _, ssn := range zimplSsnList {
		ssnRewardFactor[ssn] = "100"
	}
	ssnRewardFactor[sdk.Cfg.StZilSsnAddress] = sdk.Cfg.StZilSsnRewardShare
	AssertSuccess(p.Zproxy.AssignStakeRewardList(ssnRewardFactor, "10000"))

	p.Zproxy.Contract.Wallet = prevWallet
}

func (tr *Transitions) NextCycleOffchain(p *contracts.Protocol) *actions.AdminActions {
	tools := actions.NewAdminActions(GetLog())
	tools.TxLogMode(true)
	tools.TxLogClear()
	prevWallet := p.StZIL.Contract.Wallet
	p.StZIL.UpdateWallet(sdk.Cfg.AdminKey)
	tools.DrainBufferAuto(p)
	showOnly := false
	tools.ChownStakeReDelegate(p, showOnly)
	//tools.AutoRestake(p)
	p.StZIL.Contract.Wallet = prevWallet
	return tools
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
	tr.IsStZil()
	tr.IsZimpl()
	tr.IsBufferOrHolder()
	tr.Pause()
	tr.PerformAutoRestake()
	tr.ChownStakeAll()
	tr.AddToSwap()
	tr.TransferFrom()
	tr.MultisigWalletTests()

	if !IsCI() {
		tr.DrainBuffer()
		tr.CompleteWithdrawalSuccess()
		tr.CompleteWithdrawalMultiSsn()
		tr.WithdrawStakeAmount()
	}
}
