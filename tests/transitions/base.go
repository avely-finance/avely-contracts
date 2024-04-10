package transitions

import (
	"context"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/Zilliqa/gozilliqa-sdk/v3/account"
	"github.com/Zilliqa/gozilliqa-sdk/v3/transaction"
	"github.com/avely-finance/avely-contracts/sdk/actions"
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/contracts/evm"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	"github.com/avely-finance/avely-contracts/tests/helpers"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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
		evmOn:    false,
	}
}

type Transitions struct {
	Alice        *account.Wallet
	Bob          *account.Wallet
	Eve          *account.Wallet
	Verifier     *account.Wallet
	p            *Protocol
	evmOn        bool
	adapterStzil StZILContract
}

func (tr *Transitions) EvmOn() {
	tr.evmOn = true
	if tr.p != nil && tr.p.EvmStZIL == nil {
		// protocol is deployed, but evm-bridge is not
		tr.p.EvmStZIL = tr.DeployEvm()
		// re-init adapter
		tr.adapterStzil = NewStZILAdapter(tr.p.StZIL, tr.p.EvmStZIL, tr.evmOn)
	}
}

func (tr *Transitions) EvmOff() {
	tr.evmOn = false
}

func (tr *Transitions) GetStZIL() StZILContract {
	tr.adapterStzil.SetEvm(tr.evmOn)
	return tr.adapterStzil
}

func (tr *Transitions) GetAddressByWallet(signer interface{}) string {

	if tr.evmOn {
		// return evm-encoded adresses
		if acc, ok := signer.(*accounts.Account); ok {
			return acc.Address.Hex()
		} else if wallet, ok := signer.(*account.Wallet); ok {
			pkstr := util.EncodeHex(wallet.DefaultAccount.PrivateKey)
			privateKey, err := crypto.HexToECDSA(pkstr)
			if err != nil {
				helpers.GetLog().Fatal(err)
			}
			// zilliqa tx receipts have addresses in lowercase
			return strings.ToLower(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
		}
	} else {
		if wallet, ok := signer.(*account.Wallet); ok {
			return "0x" + wallet.DefaultAccount.Address
		}
	}
	helpers.GetLog().Fatal("can't get address by wallet")
	return ""
}

func (tr *Transitions) GetBlockNumber(tx interface{}) string {
	if txn, ok := tx.(*transaction.Transaction); ok {
		return txn.Receipt.EpochNum
	} else if txn, ok := tx.(*types.Transaction); ok {
		receipt, err := sdk.Evm.Client.TransactionReceipt(context.Background(), txn.Hash())
		if err != nil {
			GetLog().Fatal(err)
		}
		return receipt.BlockNumber.String()
	}
	GetLog().Fatal("Unknown transaction type")
	return "0"
}

func (tr *Transitions) DeployAndUpgrade() *Protocol {
	log := GetLog()
	owner := celestials.Owner
	admin := celestials.Admin

	p := Deploy(sdk, utils.GetAddressByWallet(owner), admin, log)
	sdk.Cfg.ZproxyAddr = p.Zproxy.Addr
	sdk.Cfg.ZimplAddr = p.Zimpl.Addr
	SetupZilliqaStaking(sdk, admin, verifier, log)

	//add buffers to protocol, we need 3
	buffer2, _ := p.DeployBuffer(celestials.Admin)
	buffer3, _ := p.DeployBuffer(celestials.Admin)
	p.Buffers = append(p.Buffers, buffer2, buffer3)

	p.AddSSNs(celestials.Owner)
	p.ChangeTreasuryAddress(celestials.Owner)
	p.SyncBufferAndHolder(celestials.Owner)
	p.Unpause(celestials.Owner)
	p.InitHolder()

	// deploy evm only if evmOn flag is set
	tr.p = p
	if tr.evmOn {
		p.EvmStZIL = tr.DeployEvm()
	}
	tr.adapterStzil = NewStZILAdapter(p.StZIL, p.EvmStZIL, tr.evmOn)

	tr.NextCycle(p)

	p.SetupShortcuts(log)

	return p
}

func (tr *Transitions) DeployEvm() *evm.StZIL {
	log := GetLog()

	stzilEvm, err := evm.NewStZILContract(sdk, tr.p.StZIL.Addr, celestials.EvmDeployer)
	if err != nil {
		log.Fatal("deploy Evm-StZIL error = " + err.Error())
	}
	log.Info("deploy Evm-StZIL succeed, address = " + stzilEvm.Addr)

	// pre-fill celestials.EvmDeployer with stzil tokens
	// we'll StZIL.Transfer(these tokens) instead of EvmStZIL.Delegate(these tokens)
	// see explanation in StZILAdapter.DelegateStake()
	curSigner := tr.p.StZIL.Wallet
	tr.p.StZIL.SetSigner(celestials.Verifier)
	tr.p.StZIL.DelegateStake(utils.ToQA(10_000*10 ^ 12))
	tr.p.StZIL.SetSigner(curSigner)

	return stzilEvm
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

func (tr *Transitions) DeployFraction() *FractionContract {
	log := GetLog()
	fraction, err := NewFractionContract(sdk, celestials.Admin)
	if err != nil {
		log.Fatal("deploy Fraction error = " + err.Error())
	}

	log.Info("deploy Fraction succeed, address = " + fraction.Addr)

	return fraction
}

func (tr *Transitions) NextCycle(p *contracts.Protocol) {
	sdk.IncreaseBlocknum(2)
	tools := actions.NewAdminActions(GetLog())
	prevWallet := p.Zproxy.Contract.Wallet

	p.Zproxy.SetSigner(verifier)

	tools.NextCycleWithAmount(p, 400)

	p.Zproxy.SetSigner(prevWallet)
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
		tr.SlashSSN()
		tr.WithdrawTokenAmount()
		tr.WithdrawTokenAmountWithFee()
	}
}
