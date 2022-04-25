package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

//ChownStakeConfirmSwap transition called with VerifierKey in order to demonstrate that it could be called by any user

func (tr *Transitions) ChownStakeAll() {
	tr.ChownStakeSuccess()
	tr.ChownStakeManySsnSuccess()
	tr.ChownStakeStZilErrors()
	tr.ChownStakeZimplErrors()
	tr.ChownStakeRequireDrainBuffer()
}

func (tr *Transitions) ChownStakeSuccess() {
	Start("Chown Stake Success")

	p := tr.DeployAndUpgrade()

	chownStakeSetup(tr, p)

	key1, addr1, key2, addr2, ssn, _, userStake := chownStakeDefineParams(p)

	total_supply := Field(p.StZIL, "total_supply")
	totalstakeamount := Field(p.StZIL, "totalstakeamount")
	stake1_stzil := StrMul(userStake, "64")
	stake2_stzil := StrMul(userStake, "128")
	stake1_1 := StrMul(userStake, "2")
	stake1_2 := StrMul(userStake, "4")
	stake2_1 := StrMul(userStake, "8")

	//key1 delegates through StZIL (this isn't a part of transfer process)
	AssertSuccess(p.StZIL.WithUser(key1).DelegateStake(stake1_stzil))
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	//key1, key2 delegate to main contract
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[1], stake1_1))
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[2], stake1_2))
	AssertSuccess(p.Zproxy.WithUser(key2).DelegateStake(ssn[1], stake2_1))

	//key1, key2 wait 2 reward cycles (they should have no buffered depo in current/prev cycles, else swap request will fail)
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)
	nextBuffer := p.GetBufferToSwapWith().Addr

	//key1, key2 claim rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[1]))
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[2]))
	AssertSuccess(p.Zproxy.WithUser(key2).WithdrawStakeRewards(ssn[1]))

	//key1 requests swap
	tx, _ := AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})

	//key2 requests swap
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key2).RequestDelegatorSwap(nextBuffer))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr2, "new_deleg": nextBuffer}})

	//offchain-tool calls ChownStakeConfirmSwap(addr1), expecting success
	tx, _ = AssertSuccess(p.StZIL.WithUser(sdk.Cfg.VerifierKey).ChownStakeConfirmSwap(addr1))
	AssertEvent(tx, Event{p.Zimpl.Addr, "ConfirmDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", addr1), "")
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", nextBuffer, ssn[1]), stake1_1)
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", nextBuffer, ssn[2]), stake1_2)
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", ssn[1], nextBuffer), stake1_1)
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", ssn[2], nextBuffer), stake1_2)
	AssertEqual(Field(p.StZIL, "totalstakeamount"), StrAdd(totalstakeamount, stake1_stzil, stake1_1, stake1_2))
	AssertEqual(Field(p.StZIL, "total_supply"), StrAdd(total_supply, Field(p.StZIL, "balances", addr1)))

	//offchain-tool calls ChownStakeConfirmSwap(addr2), expecting success
	tx, _ = AssertSuccess(p.StZIL.WithUser(sdk.Cfg.VerifierKey).ChownStakeConfirmSwap(addr2))
	AssertEvent(tx, Event{p.Zimpl.Addr, "ConfirmDelegatorSwap", ParamsMap{"initial_deleg": addr2, "new_deleg": nextBuffer}})
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", addr2), "")
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", nextBuffer, ssn[1]), StrAdd(stake1_1, stake2_1))
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", ssn[1], nextBuffer), StrAdd(stake1_1, stake2_1))
	AssertEqual(Field(p.Zimpl, "ssn_deleg_amt", ssn[2], nextBuffer), stake1_2)
	AssertEqual(Field(p.StZIL, "totalstakeamount"), StrAdd(totalstakeamount, stake1_stzil, stake1_1, stake1_2, stake2_1))
	AssertEqual(Field(p.StZIL, "total_supply"), StrAdd(total_supply, Field(p.StZIL, "balances", addr1), Field(p.StZIL, "balances", addr2)))

	tr.NextCycle(p)
	//key2 delegates through StZIL
	//this isn't a part of transfer process, but delegate can happen before offchain-tool calls
	AssertSuccess(p.StZIL.WithUser(key2).DelegateStake(stake2_stzil))
	tools := tr.NextCycleOffchain(p)

	//nextBuffer becomes active
	activeBuffer := p.GetActiveBuffer().Addr
	AssertEqual(nextBuffer, activeBuffer)

	//offchain tool calls ChownStakeReDelegate for each SSN when new cycle starts
	tx = tools.TxLogMap["ChownStakeReDelegate_"+ssn[1]].Tx
	AssertTransition(tx, Transition{
		p.Zimpl.Addr, //sender
		"ReDelegateStakeSuccessCallBack",
		activeBuffer,
		"0",
		ParamsMap{"ssnaddr": ssn[1], "amount": StrAdd(stake1_1, stake2_1)},
	})
	tx = tools.TxLogMap["ChownStakeReDelegate_"+ssn[2]].Tx
	AssertTransition(tx, Transition{
		p.Zimpl.Addr, //sender
		"ReDelegateStakeSuccessCallBack",
		activeBuffer,
		"0",
		ParamsMap{"ssnaddr": ssn[2], "amount": stake1_2},
	})

	total := "0"
	for _, ssn := range sdk.Cfg.SsnAddrs {
		if tmp := Field(p.Zimpl, "deposit_amt_deleg", activeBuffer, ssn); tmp != "" {
			total = StrAdd(total, tmp)
		}
	}
	AssertEqual(total, StrAdd(stake1_1, stake1_2, stake2_1, stake2_stzil))
	total = "0"
	for _, ssn := range sdk.Cfg.SsnAddrs {
		if tmp := Field(p.Zimpl, "ssn_deleg_amt", ssn, activeBuffer); tmp != "" {
			total = StrAdd(total, tmp)
		}
	}
	AssertEqual(total, StrAdd(stake1_1, stake1_2, stake2_1, stake2_stzil))
	AssertEqual(Field(p.StZIL, "totalstakeamount"), StrAdd(totalstakeamount, stake1_stzil, stake1_1, stake1_2, stake2_1, stake2_stzil))
	AssertEqual(Field(p.StZIL, "total_supply"), StrAdd(total_supply, Field(p.StZIL, "balances", addr1), Field(p.StZIL, "balances", addr2)))
}

func (tr *Transitions) ChownStakeManySsnSuccess() {
	Start("Chown Stake Success")

	p := tr.DeployAndUpgrade()

	chownStakeSetup(tr, p)

	key1, addr1, _, _, ssn, _, userStake := chownStakeDefineParams(p)
	total_supply := Field(p.StZIL, "total_supply")
	totalstakeamount := Field(p.StZIL, "totalstakeamount")
	ssnlist := []string{sdk.Cfg.StZilSsnAddress, ssn[1], ssn[2], ssn[3], ssn[4], ssn[5]}

	//key1 delegates to main contract
	for _, ssnaddr := range ssnlist {
		AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssnaddr, userStake))
	}

	//key1 waits 2 reward cycles (they should have no buffered depo in current/prev cycles, else swap request will fail)
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)
	nextBuffer := p.GetBufferToSwapWith().Addr

	//key1 claims rewards
	for _, ssnaddr := range ssnlist {
		AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssnaddr))
	}

	//key1 requests swap
	tx, _ := AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})

	//offchain-tool calls ChownStakeConfirmSwap(addr1), expecting success
	tx, _ = AssertSuccess(p.StZIL.WithUser(sdk.Cfg.VerifierKey).ChownStakeConfirmSwap(addr1))
	AssertEvent(tx, Event{p.Zimpl.Addr, "ConfirmDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})

	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	//nextBuffer becomes active
	activeBuffer := p.GetActiveBuffer().Addr
	AssertEqual(nextBuffer, activeBuffer)

	//check balances
	total := "0"
	for _, ssn := range sdk.Cfg.SsnAddrs {
		if tmp := Field(p.Zimpl, "deposit_amt_deleg", activeBuffer, ssn); tmp != "" {
			total = StrAdd(total, tmp)
		}
	}
	AssertEqual(total, StrMul(userStake, "6"))
	total = "0"
	for _, ssn := range sdk.Cfg.SsnAddrs {
		if tmp := Field(p.Zimpl, "ssn_deleg_amt", ssn, activeBuffer); tmp != "" {
			total = StrAdd(total, tmp)
		}
	}
	AssertEqual(total, StrMul(userStake, "6"))
	AssertEqual(Field(p.StZIL, "totalstakeamount"), StrAdd(totalstakeamount, StrMul(userStake, "6")))
	AssertEqual(Field(p.StZIL, "total_supply"), StrAdd(total_supply, Field(p.StZIL, "balances", addr1)))
}

func (tr *Transitions) ChownStakeStZilErrors() {
	Start("Chown Stake StZIL errors")

	p := tr.DeployAndUpgrade()

	chownStakeSetup(tr, p)

	key1, addr1, key2, addr2, ssn, _, userStake := chownStakeDefineParams(p)

	//key1 delegates to main contract
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[1], userStake))

	//key1 waits 2 reward cycles
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)
	nextBuffer := p.GetBufferToSwapWith()

	//key1 claims rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[1]))

	//offchain-tool calls ChownStakeConfirmSwap(addr1), but addr1 didn't called RequestDelegatorSwap before, expecting error
	tx, _ := p.StZIL.WithUser(sdk.Cfg.VerifierKey).ChownStakeConfirmSwap(addr1)
	AssertError(tx, "ChownStakeSwapRequestNotFound")

	//key1 requests swap with NOT buffer address
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(ssn[2]))

	//call ChownStake for addr1, expecting error
	tx, _ = p.StZIL.WithUser(sdk.Cfg.VerifierKey).ChownStakeConfirmSwap(addr1)
	AssertError(tx, "BufferAddrUnknown")

	//key1 requests swap with NOT next buffer address
	activeBuffer := p.GetActiveBuffer()
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(activeBuffer.Addr))

	//call ChownStake for addr1, expecting error
	tx, _ = p.StZIL.WithUser(sdk.Cfg.VerifierKey).ChownStakeConfirmSwap(addr1)
	AssertTransition(tx, Transition{
		activeBuffer.Addr, //sender
		"RejectDelegatorSwap",
		p.Zproxy.Addr,
		"0",
		ParamsMap{"requestor": addr1},
	})

	//key1 withdraws some amount, then requests swap
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeAmt(ssn[1], userStake))
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer.Addr))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer.Addr}})

	//call ChownStake for addr1, expecting error
	tx, _ = p.StZIL.WithUser(sdk.Cfg.VerifierKey).ChownStakeConfirmSwap(addr1)
	AssertTransition(tx, Transition{
		nextBuffer.Addr, //sender
		"RejectDelegatorSwap",
		p.Zproxy.Addr,
		"0",
		ParamsMap{"requestor": addr1},
	})

	//key2 has no deposits, but made swap request
	tx, _ = AssertSuccess(p.Zproxy.WithUser(key2).RequestDelegatorSwap(nextBuffer.Addr))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr2, "new_deleg": nextBuffer.Addr}})
	AssertEqual(Field(p.Zimpl, "deleg_swap_request", addr2), nextBuffer.Addr)

	//call ChownStake for addr2, expecting swap reject
	tx, _ = AssertSuccess(p.StZIL.WithUser(sdk.Cfg.AdminKey).ChownStakeConfirmSwap(addr2))
	AssertTransition(tx, Transition{
		nextBuffer.Addr, //sender
		"RejectDelegatorSwap",
		p.Zproxy.Addr,
		"0",
		ParamsMap{"requestor": addr2},
	})
	AssertEvent(tx, Event{p.Zimpl.Addr, "RejectDelegatorSwap", ParamsMap{"requestor": addr2, "new_deleg": nextBuffer.Addr}})
}

func (tr *Transitions) ChownStakeZimplErrors() {
	Start("Chown Stake Zimpl Errors")

	p := tr.DeployAndUpgrade()

	chownStakeSetup(tr, p)
	key1, addr1, _, _, ssn, _, userStake := chownStakeDefineParams(p)

	//key1 delegates to main contract, expecting success
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[1], userStake))

	//key1 requests delegator swap, but he has buffered deposit, expecting DelegHasBufferedDeposit
	nextBuffer := p.GetBufferToSwapWith().Addr
	tx, _ := p.Zproxy.RequestDelegatorSwap(nextBuffer)
	AssertZimplError(tx, -8)

	tr.NextCycle(p)
	tr.NextCycleOffchain(p)
	nextBuffer = p.GetBufferToSwapWith().Addr

	//key1 requests delegator swap, but he has buffered deposit in previous cycle, expecting DelegHasBufferedDeposit
	tx, _ = p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer)
	AssertZimplError(tx, -8)

	tr.NextCycle(p)
	nextBuffer = p.GetBufferToSwapWith().Addr

	//key1 requests delegator swap, but he has unclaimed rewards, expecting DelegHasUnwithdrawRewards
	tx, _ = p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer)
	AssertZimplError(tx, -12)

	//key1 claims rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[1]))

	//next buffer has no deposit/rewards in this test, so key1 can RequestDelegatorSwap
	AssertEqual(Field(p.Zimpl, "deposit_amt_deleg", nextBuffer), "")

	tr.NextCycleOffchain(p)

	tx, _ = AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})
}

func chownStakeDefineParams(p *contracts.Protocol) (string, string, string, string, []string, string, string) {
	key1 := sdk.Cfg.Key1
	addr1 := sdk.Cfg.Addr1
	key2 := sdk.Cfg.Key2
	addr2 := sdk.Cfg.Addr2
	ssn := []string{"0x1000000000000000000000000000000000000000", "0x1000000000000000000000000000000000000001",
		"0x1000000000000000000000000000000000000002", "0x1000000000000000000000000000000000000003",
		"0x1000000000000000000000000000000000000004", "0x1000000000000000000000000000000000000005"}
	minStake := Field(p.Zimpl, "minstake")
	userStake := ToZil(10)
	return key1, addr1, key2, addr2, ssn, minStake, userStake
}

func chownStakeSetup(tr *Transitions, p *contracts.Protocol) {
	_, _, _, _, ssn, minStake, _ := chownStakeDefineParams(p)

	prevWallet := p.Zproxy.Contract.Wallet

	//add test SSNs to main staking contract
	p.Zproxy.UpdateWallet(sdk.Cfg.AdminKey)
	AssertSuccess(p.Zproxy.AddSSN(ssn[0], "SSN 0"))
	AssertSuccess(p.Zproxy.AddSSN(ssn[1], "SSN 1"))
	AssertSuccess(p.Zproxy.AddSSN(ssn[2], "SSN 2"))
	AssertSuccess(p.Zproxy.AddSSN(ssn[3], "SSN 3"))
	AssertSuccess(p.Zproxy.AddSSN(ssn[4], "SSN 4"))
	AssertSuccess(p.Zproxy.AddSSN(ssn[5], "SSN 5"))
	AssertSuccess(p.Zproxy.DelegateStake(ssn[0], minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn[1], minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn[2], minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn[3], minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn[4], minStake))
	AssertSuccess(p.Zproxy.DelegateStake(ssn[5], minStake))

	p.Zproxy.Contract.Wallet = prevWallet

	//ssns will become active on the next cycle
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)
}

func (tr *Transitions) ChownStakeRequireDrainBuffer() {
	Start("Chown Stake Drain Buffer")

	p := tr.DeployAndUpgrade()

	chownStakeSetup(tr, p)

	key1, addr1, _, _, ssn, _, userStake := chownStakeDefineParams(p)

	//key1 delegates to main contract
	AssertSuccess(p.Zproxy.WithUser(key1).DelegateStake(ssn[1], userStake))

	//after 3 cycles all buffers are empty
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	//next cycle
	tr.NextCycle(p)
	nextBuffer := p.GetBufferToSwapWith().Addr

	//quick swap sequence!

	//key1 claims rewards
	AssertSuccess(p.Zproxy.WithUser(key1).WithdrawStakeRewards(ssn[1]))

	//key1 requests swap
	tx, _ := AssertSuccess(p.Zproxy.WithUser(key1).RequestDelegatorSwap(nextBuffer))
	AssertEvent(tx, Event{p.Zimpl.Addr, "RequestDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})

	//offchain-tool calls ChownStakeConfirmSwap(addr1) before DrainBuffer(), expecting error
	tx, _ = p.StZIL.WithUser(sdk.Cfg.VerifierKey).ChownStakeConfirmSwap(addr1)
	AssertError(tx, "BufferNotDrained")

	//drain buffer
	tr.NextCycleOffchain(p)

	//offchain-tool re-calls ChownStakeConfirmSwap(addr1) after DrainBuffer(), expecting success
	tx, _ = AssertSuccess(p.StZIL.WithUser(sdk.Cfg.VerifierKey).ChownStakeConfirmSwap(addr1))
	AssertEvent(tx, Event{p.Zimpl.Addr, "ConfirmDelegatorSwap", ParamsMap{"initial_deleg": addr1, "new_deleg": nextBuffer}})
}
