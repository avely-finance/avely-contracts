package transitions

import (
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"

	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) DrainBuffer() {
	Start("DrainBuffer")

	p := tr.DeployAndUpgrade()

	rewardsFee := "1000" //10% of feeDenom=10000
	AssertSuccess(p.StZIL.WithUser(sdk.Cfg.OwnerKey).ChangeRewardsFee(rewardsFee))
	p.StZIL.UpdateWallet(sdk.Cfg.AdminKey) //back to admin

	ssnForInput := p.GetSsnAddressForInput()
	activeBuffer := p.GetActiveBuffer().Addr
	AssertSuccess(p.StZIL.DelegateStake(ToZil(10)))

	//we need wait 2 reward cycles, in order to pass AssertNoBufferedDepositLessOneCycle, AssertNoBufferedDeposit checks
	tr.NextCycle(p)
	tr.NextCycleOffchain(p)

	tr.NextCycle(p)

	bufferToDrain := p.GetBufferToDrain().Addr
	AssertEqual(bufferToDrain, activeBuffer)
	//try to consolidate in holder without rewards claiming, excepting Zimpl.DelegHasUnwithdrawRewards -12 error
	txn, _ := p.StZIL.ConsolidateInHolder(bufferToDrain)
	AssertZimplError(txn, -12)

	tools := tr.NextCycleOffchain(p)

	AssertEqual(strconv.Itoa(p.Zimpl.GetLastRewardCycle()), Field(p.StZIL, "buffer_drained_cycle", bufferToDrain))

	txn = tools.TxLogMap["ClaimRewardsBuffer_"+ssnForInput].Tx
	AssertTransition(txn, Transition{
		p.StZIL.Addr,   //sender
		"ClaimRewards", //tag
		bufferToDrain,  //recipient
		"0",            //amount
		ParamsMap{},
	})

	txn = tools.TxLogMap["ClaimRewardsHolder_"+sdk.Cfg.StZilSsnAddress].Tx
	//holder has HolderInitialDelegateZil on sdk.Cfg.StZilSsnAddress
	AssertTransition(txn, Transition{
		p.StZIL.Addr,   //sender
		"ClaimRewards", //tag
		p.Holder.Addr,  //recipient
		"0",            //amount
		ParamsMap{},
	})

	// check Swap transactions
	txn = tools.TxLogMap["ConsolidateInHolder"].Tx
	AssertTransition(txn, Transition{
		bufferToDrain, //sender
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
		ParamsMap{"requestor": bufferToDrain},
	})

	//rewards claiming from holder/empty ssn address will not return errors
	AssertSuccess(p.StZIL.ClaimRewards(p.Holder.Addr, sdk.Cfg.SsnAddrs[0]))

	//rewards claiming from holder/non-ssn address will return DelegDoesNotExistAtSSN error
	txn, _ = p.StZIL.ClaimRewards(p.Holder.Addr, core.ZeroAddr)
	AssertError(txn, p.StZIL.ErrorCode("DelegDoesNotExistAtSSN"))

	//rewards claiming from buffer/empty ssn address will return DelegDoesNotExistAtSSN error
	txn, _ = p.StZIL.ClaimRewards(bufferToDrain, sdk.Cfg.SsnAddrs[0])
	AssertError(txn, p.StZIL.ErrorCode("DelegDoesNotExistAtSSN"))

	//rewards claiming from buffer/non-ssn address will return DelegDoesNotExistAtSSN error
	txn, _ = p.StZIL.ClaimRewards(bufferToDrain, core.ZeroAddr)
	AssertError(txn, p.StZIL.ErrorCode("DelegDoesNotExistAtSSN"))

	//rewards claiming from non-buffer address, expecting BufferAddrUnknown error
	txn, _ = p.StZIL.ClaimRewards(core.ZeroAddr, core.ZeroAddr)
	AssertError(txn, p.StZIL.ErrorCode("BufferOrHolderValidationFailed"))

	//repeat consolidate, excepting BufferAlreadyDrained error
	//we don't call complete tools.DrainBufferAuto(p), because it will not run twice for same buffer/lrc
	txn, _ = p.StZIL.ConsolidateInHolder(bufferToDrain)
	AssertError(txn, p.StZIL.ErrorCode("BufferAlreadyDrained"))

	//repeat consolidate with non-buffer address, excepting BufferAddrUnknown error
	txn, _ = p.StZIL.ConsolidateInHolder(core.ZeroAddr)
	AssertError(txn, p.StZIL.ErrorCode("BufferAddrUnknown"))
}

func checkRewards(p *Protocol, txn *transaction.Transaction, bufferToDrain string) {
	totalFee := "0"
	treasuryAddr := p.StZIL.GetTreasuryAddress()
	treasuryBalance := sdk.GetBalance(treasuryAddr[2:])

	// ssnlist#UpdateStakeReward has complex logic based on a fee and comission calculations
	// since we use extra small numbers (not QA 10 ^ 12) all calculations are rounded
	// and all assigned rewards are credited to one SSN node
	bufferRewards := StrAdd(sdk.Cfg.StZilSsnRewardShare, sdk.Cfg.StZilSsnRewardShare)
	AssertEqual(bufferRewards, "100")

	AssertTransition(txn, Transition{
		p.Zimpl.Addr, //sender
		"AddFunds",
		bufferToDrain,
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
		p.StZIL.Addr,
		bufferRewards,
		ParamsMap{},
	})
	AssertTransition(txn, Transition{
		p.StZIL.Addr, //sender
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
		p.StZIL.Addr,
		bufferRewards,
		ParamsMap{},
	})
	AssertTransition(txn, Transition{
		p.StZIL.Addr, //sender
		"AddFunds",
		treasuryAddr,
		rewardsFeeValue,
		ParamsMap{},
	})

	// Check StZIL balance
	totalRewards := "149" // "100" from Buffer + "49" from Holder[]
	totalRewards = StrSub(totalRewards, totalFee)
	AssertEqual(Field(p.StZIL, "_balance"), totalRewards)
	AssertEqual(Field(p.StZIL, "autorestakeamount"), totalRewards)
	//check if treasury balance increased properly
	AssertEqual(StrAdd(treasuryBalance, totalFee), sdk.GetBalance(treasuryAddr[2:]))
}
