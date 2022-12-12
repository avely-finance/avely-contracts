package transitions

import (
	"log"
	"os"
	"reflect"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/avely-finance/avely-contracts/sdk/actions"
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
	"github.com/joho/godotenv"
)

var sdk *AvelySDK
var celestials *Celestials

var alice *account.Wallet
var bob *account.Wallet
var eve *account.Wallet
var verifier *account.Wallet

func LoadUsersFromEnv(chain string) {
	path := ".env." + chain
	err := godotenv.Load(path)
	if err != nil {
		log.Printf("WARNING! There is no '%s' file. Please, make sure you set up the correct ENV manually", path)
	}

	alice = account.NewWallet()
	alice.AddByPrivateKey(os.Getenv("KEY1"))

	bob = account.NewWallet()
	bob.AddByPrivateKey(os.Getenv("KEY2"))

	eve = account.NewWallet()
	eve.AddByPrivateKey(os.Getenv("KEY3"))

	verifier = account.NewWallet()
	verifier.AddByPrivateKey(os.Getenv("VERIFIERKEY"))
}

func InitTransitions(sdkValue *AvelySDK, celestialsValue *Celestials) *Transitions {
	sdk = sdkValue
	celestials = celestialsValue
	LoadUsersFromEnv("local")

	return &Transitions{
		Alice:    alice,
		Bob:      bob,
		Eve:      eve,
		Verifier: verifier,
	}
}

type Transitions struct {
	Alice    *account.Wallet
	Bob      *account.Wallet
	Eve      *account.Wallet
	Verifier *account.Wallet
}

func (tr *Transitions) DeployAndUpgrade() *Protocol {
	log := GetLog()
	p := Deploy(sdk, celestials, log)
	sdk.Cfg.ZproxyAddr = p.Zproxy.Addr
	sdk.Cfg.ZimplAddr = p.Zimpl.Addr
	SetupZilliqaStaking(sdk, celestials, verifier, log)

	//add buffers to protocol, we need 3
	buffer2, _ := p.DeployBuffer(celestials.Admin)
	buffer3, _ := p.DeployBuffer(celestials.Admin)
	p.Buffers = append(p.Buffers, buffer2, buffer3)

	p.AddSSNs(celestials.Owner)
	p.ChangeTreasuryAddress(celestials.Owner)
	p.SyncBufferAndHolder(celestials.Owner)
	p.Unpause(celestials.Owner)
	p.InitHolder()

	tr.NextCycle(p)

	p.SetupShortcuts(log)

	return p
}

func (tr *Transitions) DeployZilSwap() *ZilSwap {
	log := GetLog()
	zilSwap, err := NewZilSwap(sdk, celestials.Admin)
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

func (tr *Transitions) DeployASwap(init_owner string) *ASwap {
	log := GetLog()
	aswap, err := NewASwap(sdk, init_owner, celestials.Admin)
	if err != nil {
		log.Fatal("deploy ASwap error = " + err.Error())
	}

	log.Info("deploy ASwap succeed, address = " + aswap.Addr)

	return aswap
}

func (tr *Transitions) DeployTreasury(init_owner string) *TreasuryContract {
	log := GetLog()
	treasury, err := NewTreasuryContract(sdk, init_owner, celestials.Admin)
	if err != nil {
		log.Fatal("deploy Treasury error = " + err.Error())
	}

	log.Info("deploy Treasury succeed, address = " + treasury.Addr)

	return treasury
}

func (tr *Transitions) DeploySsn(init_owner, init_zproxy string) *SsnContract {
	log := GetLog()
	ssn, err := NewSsnContract(sdk, init_owner, init_zproxy, celestials.Admin)
	if err != nil {
		log.Fatal("deploy SSN contract error = " + err.Error())
	}

	log.Info("deploy SSN contract succeed, address = " + ssn.Addr)

	return ssn
}

func (tr *Transitions) DeployMultisigWallet(owners []string, signCount int) *MultisigWallet {
	log := GetLog()
	multisig, err := NewMultisigContract(sdk, owners, signCount, celestials.Admin)
	if err != nil {
		log.Fatal("deploy MultisigContract error = " + err.Error())
	}

	return multisig
}

func (tr *Transitions) NextCycle(p *contracts.Protocol) {
	tr.NextCycleWithAmount(p, 400)
}

func (tr *Transitions) NextCycleWithAmount(p *contracts.Protocol, amountPerSSN int) {
	sdk.IncreaseBlocknum(2)
	totalAmount := 0
	prevWallet := p.Zproxy.Contract.Wallet

	p.Zproxy.SetSigner(verifier)

	zimplSsnList := p.Zimpl.GetSsnList()
	ssnRewardFactor := make(map[string]string)
	for _, ssn := range zimplSsnList {
		totalAmount += amountPerSSN
		ssnRewardFactor[ssn] = utils.ToQA(amountPerSSN)
	}
	ssnRewardFactor[sdk.Cfg.StZilSsnAddress] = sdk.Cfg.StZilSsnRewardShare
	AssertSuccess(p.Zproxy.AssignStakeRewardList(ssnRewardFactor, utils.ToQA(totalAmount)))

	p.Zproxy.Contract.Wallet = prevWallet
}

func (tr *Transitions) NextCycleOffchain(p *contracts.Protocol, options ...bool) *actions.AdminActions {
	tools := actions.NewAdminActions(GetLog())
	tools.TxLogMode(true)
	tools.TxLogClear()
	prevWallet := p.StZIL.Contract.Wallet
	p.StZIL.SetSigner(celestials.Admin)

	err := tools.DrainBufferAuto(p)
	if err != nil {
		GetLog().Fatal("Can't drain buffer")
	}
	showOnly := false
	tools.ChownStakeReDelegate(p, showOnly)

	//sometimes we need to disable autorestake in order to simplify calculations in tests
	enableAutorestake := true
	if len(options) > 0 {
		enableAutorestake = options[0]
	}
	if enableAutorestake {
		tools.AutoRestake(p)
	}

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
	tr.Ssn()
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
	tr.Fungible()
	tr.MultisigWalletTests()
	tr.ASwap()
	tr.Treasury()

	if !IsCI() {
		tr.DrainBuffer()
		tr.CompleteWithdrawalSuccess()
		tr.CompleteWithdrawalMultiSsn()
		tr.WithdrawStakeAmount()
	}
}
