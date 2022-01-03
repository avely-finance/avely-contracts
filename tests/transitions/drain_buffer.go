package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) DrainBuffer() {
	Start("CompleteWithdrawal - success")

	p := tr.DeployAndUpgrade()

	AssertSuccess(p.Aproxy.DelegateStake(ToZil(10)))

	txn, err := p.Aimpl.DrainBuffer(p.Aimpl.Addr)
	AssertError(txn, err, "BufferAddrUnknown")

	//we need wait 2 reward cycles, in order to pass AssertNoBufferedDepositLessOneCycle, AssertNoBufferedDeposit checks
	p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)
	sdk.IncreaseBlocknum(10)
	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))
	sdk.IncreaseBlocknum(10)
	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	txn, _ = p.Aimpl.DrainBuffer(p.GetBuffer().Addr)

	AssertTransition(txn, Transition{
		p.Aimpl.Addr,       //sender
		"ClaimRewards",     //tag
		p.GetBuffer().Addr, //recipient
		"0",                //amount
		ParamsMap{},
	})

	// ssnlist#UpdateStakeReward has complex logic based on a fee and comission calculations
	// since we use extra small numbers (not QA 10 ^ 12) all calculations are rounded
	// and all assigned rewards are credited to one SSN node
	bufferRewards := StrAdd(sdk.Cfg.AzilSsnRewardShare, sdk.Cfg.AzilSsnRewardShare)
	AssertEqual(bufferRewards, "100")

	AssertTransition(txn, Transition{
		p.Zimpl.Addr, //sender
		"AddFunds",
		p.GetBuffer().Addr,
		bufferRewards,
		ParamsMap{},
	})

	AssertTransition(txn, Transition{
		p.Zimpl.Addr, //sender
		"WithdrawStakeRewardsSuccessCallBack",
		p.GetBuffer().Addr,
		"0",
		ParamsMap{"rewards": bufferRewards},
	})

	// Holder rewards for initial funds
	holderRewards := "49"
	AssertTransition(txn, Transition{
		p.Zimpl.Addr, //sender
		"AddFunds",
		p.Holder.Addr,
		holderRewards,
		ParamsMap{},
	})

	AssertTransition(txn, Transition{
		p.Zimpl.Addr, //sender
		"WithdrawStakeRewardsSuccessCallBack",
		p.Holder.Addr,
		"0",
		ParamsMap{"rewards": holderRewards},
	})

	// Check aZIL balance
	totalRewards := "149" // "100" from Buffer + "49" from Holder[]
	AssertEqual(p.Aimpl.Field("_balance"), totalRewards)
	AssertEqual(p.Aimpl.Field("autorestakeamount"), totalRewards)

	// Send Swap transactions
	AssertTransition(txn, Transition{
		p.GetBuffer().Addr, //sender
		"RequestDelegatorSwap",
		p.Zproxy.Addr,
		"0",
		ParamsMap{"new_deleg_addr": "0x" + p.Holder.Addr},
	})

	AssertTransition(txn, Transition{
		p.Holder.Addr, //sender
		"ConfirmDelegatorSwap",
		p.Zproxy.Addr,
		"0",
		ParamsMap{"requestor": "0x" + p.GetBuffer().Addr},
	})

	//try to drain buffer, not existent at main staking contract
	//error should not be thrown
	new_buffers := []string{"0x0000000000000000000000000000000000000000"}
	AssertSuccess(p.Aimpl.ChangeBuffers(new_buffers))
	txn, _ = p.Aimpl.DrainBuffer("0000000000000000000000000000000000000000")
	AssertTransition(txn, Transition{
		p.Aimpl.Addr, //sender
		"ClaimRewards",
		p.Holder.Addr,
		"0",
		ParamsMap{},
	})
}
