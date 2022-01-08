package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) DrainBuffer() {
	Start("DrainBuffer")

	p := tr.DeployAndUpgrade()

	//try to drain not buffer address, expecting error
	txn, _ := p.Aimpl.DrainBuffer(p.Aimpl.Addr)
	AssertError(txn, "BufferAddrUnknown")

	//active buffer is p.Buffer[0] now, delegate
	AssertSuccess(p.Aimpl.DelegateStake(ToZil(100)))

	//we need wait 2 reward cycles, in order to pass AssertNoBufferedDepositLessOneCycle, AssertNoBufferedDeposit checks
	p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)
	sdk.IncreaseBlocknum(10)
	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))
	//active buffer is p.Buffers[1] now
	sdk.IncreaseBlocknum(10)
	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	//we don't repeat complicated reward calculation logic from ssnlist#UpdateStakeReward
	//instead we took already calculated values for the current setup and stakes
	expectedBufferRewards := "4"
	expectedHolderRewards := "97"

	//active buffer is p.Buffers[0] now
	//we are at the very beginning of next reward cycle, stakes weren't delegated yet
	BufferToDrain := p.Buffers[0]
	txn, _ = p.Aimpl.DrainBuffer(BufferToDrain.Addr)

	AssertTransition(txn, Transition{
		p.Aimpl.Addr,       //sender
		"ClaimRewards",     //tag
		BufferToDrain.Addr, //recipient
		"0",                //amount
		ParamsMap{},
	})

	AssertTransition(txn, Transition{
		p.Zimpl.Addr, //sender
		"AddFunds",
		BufferToDrain.Addr,
		expectedBufferRewards,
		ParamsMap{},
	})

	AssertTransition(txn, Transition{
		p.Zimpl.Addr, //sender
		"WithdrawStakeRewardsSuccessCallBack",
		BufferToDrain.Addr,
		"0",
		ParamsMap{"rewards": expectedBufferRewards},
	})

	// Holder rewards for initial funds
	AssertTransition(txn, Transition{
		p.Zimpl.Addr, //sender
		"AddFunds",
		p.Holder.Addr,
		expectedHolderRewards,
		ParamsMap{},
	})

	AssertTransition(txn, Transition{
		p.Zimpl.Addr, //sender
		"WithdrawStakeRewardsSuccessCallBack",
		p.Holder.Addr,
		"0",
		ParamsMap{"rewards": expectedHolderRewards},
	})

	// Check aZIL balance
	totalRewards := StrAdd(expectedHolderRewards, expectedBufferRewards)
	AssertEqual(p.Aimpl.Field("_balance"), totalRewards)
	AssertEqual(p.Aimpl.Field("autorestakeamount"), totalRewards)

	// Send Swap transactions
	AssertTransition(txn, Transition{
		p.GetBuffer().Addr, //sender
		"RequestDelegatorSwap",
		p.Zproxy.Addr,
		"0",
		ParamsMap{"new_deleg_addr": p.Holder.Addr},
	})

	AssertTransition(txn, Transition{
		p.Holder.Addr, //sender
		"ConfirmDelegatorSwap",
		p.Zproxy.Addr,
		"0",
		ParamsMap{"requestor": p.GetBuffer().Addr},
	})

	//after buffer has been drained, it disappears from Zimpl state
	AssertEqual("", p.Zimpl.Field("ssn_deleg_amt", sdk.Cfg.AzilSsnAddress, BufferToDrain.Addr))

	//try to drain buffer, not existent at Zimpl
	//error should not be thrown
	txn, _ = p.Aimpl.DrainBuffer(BufferToDrain.Addr)
	AssertTransition(txn, Transition{
		p.Aimpl.Addr, //sender
		"ClaimRewards",
		p.Holder.Addr,
		"0",
		ParamsMap{},
	})
}
