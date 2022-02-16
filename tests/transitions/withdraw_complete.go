package transitions

import (
	"strconv"

	"github.com/avely-finance/avely-contracts/sdk/actions"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) CompleteWithdrawalSuccess() {

	Start("CompleteWithdrawal - success")
	readyBlocks := []string{}

	p := tr.DeployAndUpgrade()

	p.Azil.UpdateWallet(sdk.Cfg.Key1)
	AssertSuccess(p.Azil.DelegateStake(ToZil(10)))

	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	sdk.IncreaseBlocknum(10)

	AssertSuccess(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))

	p.Azil.UpdateWallet(sdk.Cfg.AdminKey)
	AssertSuccess(p.Azil.DrainBuffer(p.GetBuffer().Addr))

	p.Azil.UpdateWallet(sdk.Cfg.Key1)
	tx, _ := AssertSuccess(p.Azil.WithdrawStakeAmt(ToAzil(10)))

	block1 := tx.Receipt.EpochNum
	tx, _ = p.Azil.CompleteWithdrawal()
	AssertEvent(tx, Event{p.Azil.Addr, "NoUnbondedStake", ParamsMap{}})

	p.Azil.UpdateWallet(sdk.Cfg.Key2)
	tx, _ = p.Azil.CompleteWithdrawal()
	AssertEvent(tx, Event{p.Azil.Addr, "NoUnbondedStake", ParamsMap{}})

	p.Azil.UpdateWallet(sdk.Cfg.AdminKey)
	readyBlocks = append(readyBlocks, block1)
	tx, _ = p.Azil.ClaimWithdrawal(readyBlocks)
	AssertError(tx, "ClaimWithdrawalNoUnbonded")

	delta, _ := strconv.ParseInt(StrAdd(Field(p.Zimpl, "bnum_req"), "1"), 10, 32)
	sdk.IncreaseBlocknum(int32(delta))

	tx, err := p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare)
	AssertSuccess(tx, err)

	unbondedWithdrawalsBlocks := p.GetClaimWithdrawalBlocks()
	AssertEqual(readyBlocks[0], strconv.Itoa(unbondedWithdrawalsBlocks[0]))
	actions.NewAdminActions(GetLog()).ShowClaimWithdrawal(p)

	withdrawal := Dig(p.Azil, "withdrawal_pending_of_delegator", sdk.Cfg.Addr1, block1).Withdrawal()
	AssertEqual(withdrawal.TokenAmount.String(), ToAzil(10))
	AssertEqual(withdrawal.StakeAmount.String(), ToAzil(10))

	p.Azil.UpdateWallet(sdk.Cfg.AdminKey)
	tx, _ = p.Azil.ClaimWithdrawal(readyBlocks)
	AssertTransition(tx, Transition{
		p.Azil.Addr,          //sender
		"CompleteWithdrawal", //tag
		p.Holder.Addr,        //recipient
		"0",                  //amount
		ParamsMap{},
	})
	AssertEvent(tx, Event{p.Holder.Addr, "AddFunds", ParamsMap{"funder": p.Zimpl.Addr, "amount": ToZil(10)}})

	AssertTransition(tx, Transition{
		p.Holder.Addr,                       //sender
		"CompleteWithdrawalSuccessCallBack", //tag
		p.Azil.Addr,                         //recipient
		ToZil(10),                           //amount
		ParamsMap{},
	})

	p.Azil.UpdateWallet(sdk.Cfg.Key1)
	tx, _ = p.Azil.CompleteWithdrawal()
	AssertEvent(tx, Event{p.Azil.Addr, "CompleteWithdrawal", ParamsMap{"amount": ToZil(10), "delegator": sdk.Cfg.Addr1}})
	AssertTransition(tx, Transition{
		p.Azil.Addr,
		"CompleteWithdrawalSuccessCallBack",
		sdk.Cfg.Addr1,
		"0",
		ParamsMap{"amount": ToZil(10)},
	})
	AssertTransition(tx, Transition{
		p.Azil.Addr,
		"AddFunds",
		sdk.Cfg.Addr1,
		ToZil(10),
		ParamsMap{},
	})

	withdrawal = Dig(p.Azil, "withdrawal_pending_of_delegator", sdk.Cfg.Addr1, block1).Withdrawal()
	AssertEqual(withdrawal.TokenAmount.String(), "0")
	AssertEqual(withdrawal.StakeAmount.String(), "0")

	AssertEqual(Field(p.Azil, "totalstakeamount"), ToZil(1000))
	AssertEqual(Field(p.Azil, "totaltokenamount"), ToAzil(1000))
	AssertEqual(Field(p.Azil, "tmp_complete_withdrawal_available"), "0")
	AssertEqual(Field(p.Azil, "balances", sdk.Cfg.Admin), ToAzil(1000))

	AssertEqual(Field(p.Azil, "withdrawal_unbonded"), "{}")
	AssertEqual(Field(p.Azil, "withdrawal_pending"), "{}")
}
