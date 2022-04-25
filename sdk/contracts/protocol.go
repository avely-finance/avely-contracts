package contracts

import (
	"log"
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	avelycore "github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
)

type Protocol struct {
	Zproxy  *Zproxy
	Zimpl   *Zimpl
	StZIL   *StZIL
	Buffers []*BufferContract
	Holder  *HolderContract
}

func NewProtocol(zproxy *Zproxy, zimpl *Zimpl, stzil *StZIL, buffers []*BufferContract, holder *HolderContract) *Protocol {
	if len(buffers) == 0 {
		log.Fatal("Protocol should have at least one buffer")
	}

	return &Protocol{
		Zproxy:  zproxy,
		Zimpl:   zimpl,
		StZIL:   stzil,
		Buffers: buffers,
		Holder:  holder,
	}
}

func (p *Protocol) DeployBuffer() (*BufferContract, error) {
	return NewBufferContract(p.StZIL.Sdk, p.StZIL.Addr, p.Zproxy.Addr)
}

func (p *Protocol) GetSsnAddressForInput() string {
	ssnAddrs := p.StZIL.GetSsnWhitelist()
	ssnIndex := p.StZIL.GetSsnIndex()
	i := ssnIndex.Uint64() % uint64(len(ssnAddrs))
	return ssnAddrs[i]
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
	lrc := p.Zimpl.GetLastRewardCycle()
	lrc = lrc + offset
	buffers := p.Buffers
	i := int(lrc) % len(buffers)
	return buffers[i]
}

func (p *Protocol) GetBlockHeight() int {
	result := ""
	if p.StZIL.Sdk.Cfg.Chain == "local" {
		//Isolated server has limited set of API methods: https://github.com/Zilliqa/zilliqa-isolated-server#available-apis
		//GetNumTxBlocks is not available.
		//So we'll take BlockNum from receipt of safe transaction.
		tx, _ := p.Zproxy.UpdateVerifier(p.StZIL.Sdk.Cfg.Verifier)
		result = tx.Receipt.EpochNum
	} else {
		result, _ = p.StZIL.Contract.Provider.GetNumTxBlocks()
	}
	blockHeight, _ := strconv.Atoi(result)
	return blockHeight
}

func (p *Protocol) GetClaimWithdrawalBlocks() []int {
	curBlockNum := p.GetBlockHeight()
	bnumReq := p.Zimpl.GetBnumReq()

	//get all blocks with pending withdrawals
	partialState := p.StZIL.Contract.SubState("withdrawal_pending", []string{})
	state := NewState(partialState)
	allWithdrawBlocks := state.Dig("result.withdrawal_pending|@keys")
	blocks := allWithdrawBlocks.ArrayInt()

	//see leave_unbonded function in stzil
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
	return p.Holder.DelegateStake(p.StZIL.Sdk.Cfg.StZilSsnAddress, ToZil(p.StZIL.Sdk.Cfg.HolderInitialDelegateZil))
}

func (p *Protocol) SyncBufferAndHolder() {
	new_buffers := []string{}

	for _, b := range p.Buffers {
		new_buffers = append(new_buffers, b.Addr)
	}

	prevWallet := p.StZIL.Wallet
	check(p.StZIL.WithUser(p.StZIL.Sdk.Cfg.OwnerKey).ChangeBuffers(new_buffers))
	check(p.StZIL.WithUser(p.StZIL.Sdk.Cfg.OwnerKey).SetHolderAddress(p.Holder.Addr))
	p.StZIL.Wallet = prevWallet
}

func (p *Protocol) SyncBuffers() {
	new_buffers := []string{}

	for _, b := range p.Buffers {
		new_buffers = append(new_buffers, b.Addr)
	}

	prevWallet := p.StZIL.Wallet
	check(p.StZIL.WithUser(p.StZIL.Sdk.Cfg.OwnerKey).ChangeBuffers(new_buffers))
	p.StZIL.Wallet = prevWallet
}

func (p *Protocol) AddSSNs() {
	prevWallet := p.StZIL.Wallet

	//reverse elements to keep order of stzil.ssn_addresses elements same as in config
	for i := len(p.StZIL.Sdk.Cfg.SsnAddrs) - 1; i >= 0; i-- {
		check(p.StZIL.WithUser(p.StZIL.Sdk.Cfg.OwnerKey).AddSSN(p.StZIL.Sdk.Cfg.SsnAddrs[i]))
	}
	p.StZIL.Wallet = prevWallet
}

func (p *Protocol) ChangeTreasuryAddress() {
	prevWallet := p.StZIL.Wallet
	check(p.StZIL.WithUser(p.StZIL.Sdk.Cfg.OwnerKey).ChangeTreasuryAddress(p.StZIL.Sdk.Cfg.TreasuryAddr))
	p.StZIL.Wallet = prevWallet
}

func (p *Protocol) Unpause() {
	prevWallet := p.StZIL.Wallet
	p.StZIL.UpdateWallet(p.StZIL.Sdk.Cfg.OwnerKey)
	check(p.StZIL.UnpauseIn())
	check(p.StZIL.UnpauseOut())
	check(p.StZIL.UnpauseZrc2())
	p.StZIL.Wallet = prevWallet
}

func (p *Protocol) SetupZProxy() {
	sdk := p.StZIL.Sdk
	args := []core.ContractValue{
		{
			"newImplementation",
			"ByStr20",
			p.Zimpl.Addr,
		},
	}
	check(p.Zproxy.Call("UpgradeTo", args, "0"))
	for _, ssnaddr := range sdk.Cfg.SsnAddrs {
		check(p.Zproxy.AddSSN(ssnaddr, ssnaddr))
	}
	check(p.Zproxy.UpdateVerifierRewardAddr(sdk.Cfg.Verifier))
	check(p.Zproxy.UpdateVerifier(sdk.Cfg.Verifier))
	check(p.Zproxy.UpdateStakingParameters(ToZil(sdk.Cfg.SsnInitialDelegateZil), ToZil(10))) //minstake (ssn not active if less), mindelegstake
	check(p.Zproxy.Unpause())

	//we need our SSN to be active, so delegating some stake to each
	for _, ssnaddr := range sdk.Cfg.SsnAddrs {
		check(p.Zproxy.DelegateStake(ssnaddr, ToZil(sdk.Cfg.SsnInitialDelegateZil)))
	}

	//we need to delegate something from Holder, in order to make Zimpl know holder's address
	check(p.Holder.DelegateStake(sdk.Cfg.StZilSsnAddress, ToZil(sdk.Cfg.HolderInitialDelegateZil)))

	p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)

	// SSN will become active on next cycle
	//we need to increase blocknum, in order to Gzil won't mint anything. Really minting is over.
	sdk.IncreaseBlocknum(10)
	check(p.Zproxy.AssignStakeReward(sdk.Cfg.StZilSsnAddress, sdk.Cfg.StZilSsnRewardShare))
}

func (p *Protocol) SetupShortcuts(log *avelycore.Log) {
	log.AddShortcut("Zproxy", p.Zproxy.Addr)
	log.AddShortcut("Zimpl", p.Zimpl.Addr)
	log.AddShortcut("StZIL", p.StZIL.Addr)
	log.AddShortcut("Holder", p.Holder.Addr)

	for i, b := range p.Buffers {
		title := "Buffer" + strconv.Itoa(i)
		log.AddShortcut(title, b.Addr)
	}
}
