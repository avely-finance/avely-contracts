package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) WithdrawTokenAmountWithFee() {
	Start("WithdrawTokenAmountWithFee")

	// deploy smart contract
	p := tr.DeployAndUpgrade()
	bobAddr := utils.GetAddressByWallet(bob)
	aliceAddr := utils.GetAddressByWallet(alice)

	p.StZIL.SetSigner(celestials.Owner)
	fees := ToStZil(2)
	AssertSuccess(p.StZIL.UpdateStakingParameters("0", "0", fees))
	AssertSuccess(p.StZIL.ChangeWithdrawalFeeAddress(aliceAddr))

	p.StZIL.SetSigner(bob)

	tokens := ToZil(12)
	AssertSuccess(p.StZIL.DelegateStake(tokens))
	txn, _ := p.StZIL.WithdrawTokensAmt(tokens)
	bnum1 := txn.Receipt.EpochNum

	AssertEvent(txn, Event{
		Sender:    p.StZIL.Addr,
		EventName: "TransferSuccess",
		Params:    ParamsMap{"sender": bobAddr, "recipient": aliceAddr, "amount": fees},
	})

	withdrawal := Dig(p.StZIL, "withdrawal_pending", bnum1, bobAddr).Withdrawal()

	AssertEqual(withdrawal.TokenAmount.String(), ToStZil(10))
	AssertEqual(withdrawal.StakeAmount.String(), ToStZil(10))

	// alice should have withdrawal fees
	AssertEqual(Field(p.StZIL, "balances", aliceAddr), fees)

	tx, _ := p.StZIL.WithdrawTokensAmt(ToStZil(1))
	AssertError(tx, p.StZIL.ErrorCode("WithdrawTokensAmtLessThanFee"))
}
