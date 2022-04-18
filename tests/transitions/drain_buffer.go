package transitions

import (
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/avely-finance/avely-contracts/sdk/actions"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"

	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) DrainBuffer() {
	Start("DrainBuffer")

	p := tr.DeployAndUpgrade()
	buffer2, _ := p.DeployBuffer()
	buffer3, _ := p.DeployBuffer()
	p.Buffers = append(p.Buffers, buffer2, buffer3)
	p.SyncBuffers()
	p.SetupShortcuts(GetLog())

	rewardsFee := "1000" //10% of feeDenom=10000
	treasuryAddr := sdk.Cfg.Addr3
	AssertSuccess(p.Azil.WithUser(sdk.Cfg.OwnerKey).ChangeRewardsFee(rewardsFee))
	AssertSuccess(p.Azil.WithUser(sdk.Cfg.OwnerKey).ChangeTreasuryAddress(treasuryAddr))
	p.Azil.UpdateWallet(sdk.Cfg.AdminKey) //back to admin
	tools := actions.NewAdminActions(GetLog())
	tools.SetTestMode(true)

	AssertSuccess(p.Azil.DelegateStake(ToZil(10)))

	//we need wait 2 reward cycles, in order to pass AssertNoBufferedDepositLessOneCycle, AssertNoBufferedDeposit checks
	p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)
	sdk.IncreaseBlocknum(10)
	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	sdk.IncreaseBlocknum(10)
	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	bufferToDrain := p.GetBufferToDrain().Addr

	//try to consolidate in holder without rewards claiming, excepting Zimpl.DelegHasUnwithdrawRewards -12 error
	txn, _ := p.Azil.ConsolidateInHolder(bufferToDrain)
	AssertZimplError(txn, -12)

	tools.DrainBufferAuto(p)
	AssertEqual(strconv.Itoa(p.Zimpl.GetLastRewardCycle()), Field(p.Azil, "buffer_drained_cycle", bufferToDrain))

	txn = tools.TxLogMap["ClaimRewardsBuffer_"+sdk.Cfg.AzilSsnAddress].Tx
	AssertTransition(txn, Transition{
		p.Azil.Addr,    //sender
		"ClaimRewards", //tag
		bufferToDrain,  //recipient
		"0",            //amount
		ParamsMap{},
	})

	txn = tools.TxLogMap["ClaimRewardsHolder_"+sdk.Cfg.AzilSsnAddress].Tx
	AssertTransition(txn, Transition{
		p.Azil.Addr,    //sender
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
	AssertSuccess(p.Azil.ClaimRewards(p.Holder.Addr, sdk.Cfg.SsnAddrs[0]))

	//rewards claiming from holder/non-ssn address will return DelegDoesNotExistAtSSN error
	txn, _ = p.Azil.ClaimRewards(p.Holder.Addr, core.ZeroAddr)
	AssertError(txn, "DelegDoesNotExistAtSSN")

	//rewards claiming from buffer/empty ssn address will return DelegDoesNotExistAtSSN error
	txn, _ = p.Azil.ClaimRewards(bufferToDrain, sdk.Cfg.SsnAddrs[0])
	AssertError(txn, "DelegDoesNotExistAtSSN")

	//rewards claiming from buffer/non-ssn address will return DelegDoesNotExistAtSSN error
	txn, _ = p.Azil.ClaimRewards(bufferToDrain, core.ZeroAddr)
	AssertError(txn, "DelegDoesNotExistAtSSN")

	//rewards claiming from non-buffer address, expecting BufferAddrUnknown error
	txn, _ = p.Azil.ClaimRewards(core.ZeroAddr, core.ZeroAddr)
	AssertError(txn, "BufferOrHolderValidationFailed")

	//repeat consolidate, excepting BufferAlreadyDrained error
	//we don't call complete tools.DrainBufferAuto(p), because it will not run twice for same buffer/lrc
	txn, _ = p.Azil.ConsolidateInHolder(bufferToDrain)
	AssertError(txn, "BufferAlreadyDrained")

	//repeat consolidate with non-buffer address, excepting BufferAddrUnknown error
	txn, _ = p.Azil.ConsolidateInHolder(core.ZeroAddr)
	AssertError(txn, "BufferAddrUnknown")
}

func checkRewards(p *Protocol, txn *transaction.Transaction, bufferToDrain string) {
	totalFee := "0"
	treasuryAddr := p.Azil.GetTreasuryAddress()
	treasuryBalance := sdk.GetBalance(treasuryAddr[2:])

	// ssnlist#UpdateStakeReward has complex logic based on a fee and comission calculations
	// since we use extra small numbers (not QA 10 ^ 12) all calculations are rounded
	// and all assigned rewards are credited to one SSN node
	bufferRewards := StrAdd(sdk.Cfg.AzilSsnRewardShare, sdk.Cfg.AzilSsnRewardShare)
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
	//check if treasury balance increased properly
	AssertEqual(StrAdd(treasuryBalance, totalFee), sdk.GetBalance(treasuryAddr[2:]))
}
