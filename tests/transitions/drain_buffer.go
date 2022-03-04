package transitions

import (
	"strconv"

	"github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) DrainBuffer() {
	Start("CompleteWithdrawal - success")

	p := tr.DeployAndUpgrade()
	rewardsFee := "1000" //10% of feeDenom=10000
	treasuryAddr := sdk.Cfg.Addr3
	totalFee := "0"
	AssertSuccess(p.Azil.WithUser(sdk.Cfg.OwnerKey).ChangeRewardsFee(rewardsFee))
	AssertSuccess(p.Azil.WithUser(sdk.Cfg.OwnerKey).ChangeTreasuryAddress(treasuryAddr))
	p.Azil.UpdateWallet(sdk.Cfg.AdminKey) //back to admin
	treasuryBalance := sdk.GetBalance(treasuryAddr[2:])

	AssertSuccess(p.Azil.DelegateStake(ToZil(10)))

	txn, _ := p.Azil.DrainBuffer(p.Azil.Addr)
	AssertError(txn, "BufferAddrUnknown")

	//we need wait 2 reward cycles, in order to pass AssertNoBufferedDepositLessOneCycle, AssertNoBufferedDeposit checks
	p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)
	sdk.IncreaseBlocknum(10)

	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))
	sdk.IncreaseBlocknum(10)
	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	bufferAddr := p.GetBuffer().Addr
	txn, _ = p.Azil.DrainBuffer(bufferAddr)
	AssertEqual(strconv.Itoa(p.Zimpl.GetLastRewardCycle()), Field(p.Azil, "buffer_drained_cycle", bufferAddr))

	AssertTransition(txn, Transition{
		p.Azil.Addr,        //sender
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

	//transfer rewards fee to treasury
	//rewardsFee * bufferRewards / feeDenom = 1000 * 100 / 10000 = 0.1 * 100 = 10
	rewardsFeeValue := "10"
	totalFee = StrAdd(totalFee, rewardsFeeValue)
	AssertTransition(txn, Transition{
		p.GetBuffer().Addr, //sender
		"ClaimRewardsSuccessCallBack",
		p.Azil.Addr,
		bufferRewards,
		ParamsMap{},
	})
	AssertTransition(txn, Transition{
		p.Azil.Addr, //sender
		"AddFunds",
		treasuryAddr,
		rewardsFeeValue,
		ParamsMap{},
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

	//transfer rewards fee to treasury
	//rewardsFee * holderRewards / feeDenom = 1000 * 49 / 10000 = 0.1 * 49 = 4.9 = 4
	rewardsFeeValue = "4"
	totalFee = StrAdd(totalFee, rewardsFeeValue)
	AssertTransition(txn, Transition{
		p.GetBuffer().Addr, //sender
		"ClaimRewardsSuccessCallBack",
		p.Azil.Addr,
		bufferRewards,
		ParamsMap{},
	})
	AssertTransition(txn, Transition{
		p.Azil.Addr, //sender
		"AddFunds",
		treasuryAddr,
		rewardsFeeValue,
		ParamsMap{},
	})

	// Check aZIL balance
	totalRewards := "149" // "100" from Buffer + "49" from Holder[]
	totalRewards = StrSub(totalRewards, totalFee)
	AssertEqual(Field(p.Azil, "_balance"), totalRewards)
	AssertEqual(Field(p.Azil, "autorestakeamount"), totalRewards)

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

	//repeat drain buffer, excepting error
	txn, _ = p.Azil.DrainBuffer(bufferAddr)
	AssertError(txn, "BufferAlreadyDrained")

	//try to drain buffer, not existent at main staking contract
	//error should not be thrown
	new_buffers := []string{core.ZeroAddr}
	AssertSuccess(p.Azil.WithUser(sdk.Cfg.OwnerKey).ChangeBuffers(new_buffers))
	txn, _ = p.Azil.WithUser(sdk.Cfg.AdminKey).DrainBuffer(core.ZeroAddr)
	AssertTransition(txn, Transition{
		p.Azil.Addr, //sender
		"ClaimRewards",
		p.Holder.Addr,
		"0",
		ParamsMap{},
	})

	//check if treasury balance increased properly
	AssertEqual(StrAdd(treasuryBalance, totalFee), sdk.GetBalance(treasuryAddr[2:]))
}
