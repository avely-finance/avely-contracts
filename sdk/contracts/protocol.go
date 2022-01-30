package contracts

import (
	"log"
	"strconv"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	avelycore "github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
)

type Protocol struct {
	Zproxy  *Zproxy
	Zimpl   *Zimpl
	Aimpl   *AZil
	Buffers []*BufferContract
	Holder  *HolderContract
}

func NewProtocol(zproxy *Zproxy, zimpl *Zimpl, aimpl *AZil, buffers []*BufferContract, holder *HolderContract) *Protocol {
	if len(buffers) == 0 {
		log.Fatal("Protocol should have at least one buffer")
	}

	return &Protocol{
		Zproxy:  zproxy,
		Zimpl:   zimpl,
		Aimpl:   aimpl,
		Buffers: buffers,
		Holder:  holder,
	}
}

func (p *Protocol) GetAzilSsnAddress() string {
	return strings.ToLower(p.Aimpl.Sdk.Cfg.AzilSsnAddress)
}

func (p *Protocol) DeployBuffer() (*BufferContract, error) {
	return NewBufferContract(p.Aimpl.Sdk, p.Aimpl.Addr, p.Zproxy.Addr, p.Zimpl.Addr)
}

func (p *Protocol) GetBuffer() *BufferContract {
	return p.Buffers[0]
}

func (p *Protocol) GetActiveBuffer() *BufferContract {
	return p.GetBufferByOffset(0)
}

func (p *Protocol) GetBufferToDrain() *BufferContract {
	return p.GetBufferByOffset(-2)
}

func (p *Protocol) GetBufferToSwapWith() *BufferContract {
	return p.GetBufferByOffset(1)
}

func (p *Protocol) GetBufferByOffset(offset int) *BufferContract {
	lrc := p.GetLastRewardCycle()
	lrc = lrc + offset
	buffers := p.Buffers
	i := int(lrc) % len(buffers)
	return buffers[i]
}

func (p *Protocol) GetLastRewardCycle() int {
	partialState := p.Zimpl.Contract.SubState("lastrewardcycle", []string{})

	state := NewState(partialState)

	lrc := state.Dig("result.lastrewardcycle").Int()

	return int(lrc)
}

func (p *Protocol) GetBlockHeight() int {
	result := ""
	if p.Aimpl.Sdk.Cfg.Chain == "local" {
		tx, _ := p.Zproxy.UpdateVerifier(p.Aimpl.Sdk.Cfg.Verifier)
		result = tx.Receipt.EpochNum
	} else {
		result, _ = p.Aimpl.Contract.Provider.GetNumTxBlocks()
	}
	blockHeight, _ := strconv.Atoi(result)
	return blockHeight
}

func (p *Protocol) GetUnbondedWithdrawalsBlocks() []int {
	curBlockNum := p.GetBlockHeight()
	bnumReq := p.Zimpl.GetBnumReq()

	//get all blocks with pending withdrawals
	partialState := p.Aimpl.Contract.SubState("withdrawal_pending", []string{})
	state := NewState(partialState)
	allWithdrawBlocks := state.Dig("result.withdrawal_pending|@keys")
	blocks := allWithdrawBlocks.ArrayInt()

	//see leave_unbonded function in azil
	unbonded := []int{}
	for _, bnum := range blocks {
		if bnum+bnumReq < curBlockNum {
			unbonded = append(unbonded, bnum)
		}
	}

	return unbonded
}

func (p *Protocol) GetSwapRequestsForBuffer(bufferAddr string) []string {
	partialState := p.Zimpl.Contract.SubState("deleg_swap_request", []string{})
	state := NewState(partialState)
	allSwapRequests := state.Dig("result.deleg_swap_request").Map()
	bufferSwapRequests := []string{}
	for initiator, newDeleg := range allSwapRequests {
		newDelegStr := newDeleg.String()
		if newDelegStr == bufferAddr {
			bufferSwapRequests = append(bufferSwapRequests, initiator)
		}
	}
	return bufferSwapRequests
}

func (p *Protocol) InitHolder() (*transaction.Transaction, error) {
	return p.Holder.DelegateStake(ToZil(p.Aimpl.Sdk.Cfg.HolderInitialDelegateZil))
}

func (p *Protocol) SyncBufferAndHolder() {
	new_buffers := []string{}

	for _, b := range p.Buffers {
		new_buffers = append(new_buffers, b.Addr)
	}

	check(p.Aimpl.ChangeBuffers(new_buffers))
	check(p.Aimpl.ChangeHolderAddress(p.Holder.Addr))
}

func (p *Protocol) Unpause() {
	check(p.Aimpl.Unpause())
}

func (p *Protocol) SetupZProxy() {
	sdk := p.Aimpl.Sdk
	args := []core.ContractValue{
		{
			"newImplementation",
			"ByStr20",
			p.Zimpl.Addr,
		},
	}
	check(p.Zproxy.Call("UpgradeTo", args, "0"))
	check(p.Zproxy.AddSSN(sdk.Cfg.AzilSsnAddress, "aZil SSN"))
	check(p.Zproxy.UpdateVerifierRewardAddr(sdk.Cfg.Verifier))
	check(p.Zproxy.UpdateVerifier(sdk.Cfg.Verifier))
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
	log.AddShortcut("Zproxy", p.Zproxy.Addr)
	log.AddShortcut("Zimpl", p.Zimpl.Addr)
	log.AddShortcut("Aimpl", p.Aimpl.Addr)
	log.AddShortcut("Holder", p.Holder.Addr)

	for i, b := range p.Buffers {
		title := "Buffer" + strconv.Itoa(i)
		log.AddShortcut(title, b.Addr)
	}
}
