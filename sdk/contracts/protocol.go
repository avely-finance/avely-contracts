package contracts

import (
	"log"
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
)

type Protocol struct {
	Zproxy   *Zproxy
	Zimpl    *Zimpl
	StZIL    *StZIL
	Buffers  []*BufferContract
	Holder   *HolderContract
	Treasury *TreasuryContract
}

func NewProtocol(zproxy *Zproxy, zimpl *Zimpl, stzil *StZIL, buffers []*BufferContract, holder *HolderContract, treasury *TreasuryContract) *Protocol {
	if len(buffers) == 0 {
		log.Fatal("Protocol should have at least one buffer")
	}

	return &Protocol{
		Zproxy:   zproxy,
		Zimpl:    zimpl,
		StZIL:    stzil,
		Buffers:  buffers,
		Holder:   holder,
		Treasury: treasury,
	}
}

func (p *Protocol) DeployBuffer(deployer *account.Wallet) (*BufferContract, error) {
	return NewBufferContract(p.StZIL.Sdk, p.StZIL.Addr, p.Zproxy.Addr, deployer)
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
	if p.StZIL.Sdk.Cfg.Chain == "local" {
		//Isolated server has limited set of API methods: https://github.com/Zilliqa/zilliqa-isolated-server#available-apis
		//GetNumTxBlocks is not available.
		//So we'll take BlockNum from receipt of safe transaction.
		tx, _ := p.Zproxy.UpdateVerifier(p.StZIL.Sdk.Cfg.Verifier)
		result := tx.Receipt.EpochNum
		blockHeight, _ := strconv.Atoi(result)
		return blockHeight

	} else {
		result, _ := p.StZIL.Sdk.GetBlockHeight()
		return result
	}
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
	CheckTx(p.StZIL.WithUser(p.StZIL.Sdk.Cfg.OwnerKey).ChangeBuffers(new_buffers))
	CheckTx(p.StZIL.WithUser(p.StZIL.Sdk.Cfg.OwnerKey).SetHolderAddress(p.Holder.Addr))
	p.StZIL.Wallet = prevWallet
}

func (p *Protocol) SyncBuffers() {
	new_buffers := []string{}

	for _, b := range p.Buffers {
		new_buffers = append(new_buffers, b.Addr)
	}

	prevWallet := p.StZIL.Wallet
	CheckTx(p.StZIL.WithUser(p.StZIL.Sdk.Cfg.OwnerKey).ChangeBuffers(new_buffers))
	p.StZIL.Wallet = prevWallet
}

func (p *Protocol) AddSSNs() {
	prevWallet := p.StZIL.Wallet

	//reverse elements to keep order of stzil.ssn_addresses elements same as in config
	for i := len(p.StZIL.Sdk.Cfg.SsnAddrs) - 1; i >= 0; i-- {
		CheckTx(p.StZIL.WithUser(p.StZIL.Sdk.Cfg.OwnerKey).AddSSN(p.StZIL.Sdk.Cfg.SsnAddrs[i]))
	}
	p.StZIL.Wallet = prevWallet
}

func (p *Protocol) ChangeTreasuryAddress() {
	prevWallet := p.StZIL.Wallet
	CheckTx(p.StZIL.WithUser(p.StZIL.Sdk.Cfg.OwnerKey).ChangeTreasuryAddress(p.Treasury.Addr))
	p.StZIL.Wallet = prevWallet
}

func (p *Protocol) Unpause() {
	prevWallet := p.StZIL.Wallet
	p.StZIL.UpdateWallet(p.StZIL.Sdk.Cfg.OwnerKey)
	CheckTx(p.StZIL.UnpauseIn())
	CheckTx(p.StZIL.UnpauseOut())
	CheckTx(p.StZIL.UnpauseZrc2())
	p.StZIL.Wallet = prevWallet
}

func (p *Protocol) SetupShortcuts(log *Log) {
	log.AddShortcut("Zproxy", p.Zproxy.Addr)
	log.AddShortcut("Zimpl", p.Zimpl.Addr)
	log.AddShortcut("StZIL", p.StZIL.Addr)
	log.AddShortcut("Holder", p.Holder.Addr)

	for i, b := range p.Buffers {
		title := "Buffer" + strconv.Itoa(i)
		log.AddShortcut(title, b.Addr)
	}
}
