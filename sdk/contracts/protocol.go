package contracts

import (
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	avelycore "github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	"log"
	"runtime"
	"strconv"
)

type Protocol struct {
	Zproxy *Zproxy
	Zimpl  *Zimpl
	Aimpl  *AZil
	Buffers []*BufferContract
	Holder *HolderContract
}

func NewProtocol(zproxy *Zproxy, zimpl *Zimpl, azil *AZil, buffers []*BufferContract, holder *HolderContract) *Protocol {
	if len(buffers) == 0 {
		log.Fatal("Protocol should have at least one buffer")
	}

	return &Protocol{
		Zproxy: zproxy,
		Zimpl:  zimpl,
		Aimpl:  azil,
		Buffers: buffers,
		Holder: holder,
	}
}

func (p *Protocol) GetBuffer() *BufferContract {
	return p.Buffers[0]
}

func (p *Protocol) SyncBufferAndHolder() {
	new_buffers := []string{}

	for _, b := range p.Buffers {
		new_buffers = append(new_buffers, "0x" + b.Addr)
	}

	check(p.Aimpl.ChangeBuffers(new_buffers))
	check(p.Aimpl.ChangeHolderAddress(p.Holder.Addr))
}

func (p *Protocol) SetupZProxy() {
	sdk := p.Aimpl.Sdk
	args := []core.ContractValue{
		{
			"newImplementation",
			"ByStr20",
			"0x" + p.Zimpl.Addr,
		},
	}
	check(p.Zproxy.Call("UpgradeTo", args, "0"))
	check(p.Zproxy.AddSSN(sdk.Cfg.AzilSsnAddress, "aZil SSN"))
	check(p.Zproxy.UpdateVerifierRewardAddr("0x" + sdk.Cfg.Verifier))
	check(p.Zproxy.UpdateVerifier("0x" + sdk.Cfg.Verifier))
	check(p.Zproxy.UpdateStakingParameters(ToZil(1000), ToZil(10))) //minstake (ssn not active if less), mindelegstake
	check(p.Zproxy.Unpause())

	//we need our SSN to be active, so delegating some stake
	check(p.Aimpl.DelegateStake(ToZil(1000)))

	//we need to delegate something from Holder, in order to make Zimpl know holder's address
	check(p.Holder.DelegateStake(ToZil(sdk.Cfg.HolderInitialDelegateZil)))

	p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)

	// SSN will become active on next cycle
	//we need to increase blocknum, in order to Gzil won't mint anything. Really minting is over.
	sdk.IncreaseBlocknum(10)
	check(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))
}

func (p *Protocol) SetupShortcuts(log *avelycore.Log) {
	log.AddShortcut("Zproxy", "0x"+p.Zproxy.Addr)
	log.AddShortcut("Zimpl", "0x"+p.Zimpl.Addr)
	log.AddShortcut("Aimpl", "0x"+p.Aimpl.Addr)
	log.AddShortcut("Holder", "0x"+p.Holder.Addr)

	for i, b := range p.Buffers {
		title := "Buffer" + strconv.Itoa(i)
		log.AddShortcut(title, "0x"+b.Addr)
	}
}

func check(tx *transaction.Transaction, err error) (*transaction.Transaction, error) {
	if err != nil {
		_, file, no, _ := runtime.Caller(1)
		log.Fatal("TRANSACTION FAILED, " + file + ":" + strconv.Itoa(no))
	}
	return tx, err
}
