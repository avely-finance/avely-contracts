package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) SlashSSN() {

	Start("SlashSSN")

	// deploy smart contract
	p := tr.DeployAndUpgrade()
	stakeAmt := ToZil(15)

	/*******************************************************************************
	 * 0. delegator (sdk.Cfg.Addr2) delegate 15 zil and try to withdrawal
	 *******************************************************************************/
	p.StZIL.SetSigner(bob)
	AssertSuccess(p.StZIL.DelegateStake(stakeAmt))
	ssnNode1 := sdk.Cfg.SsnAddrs[0]
	txn, _ := p.StZIL.SlashSSN(ToStZil(10), ssnNode1)

	AssertError(txn, p.StZIL.ErrorCode("AdminValidationFailed"))

	/*******************************************************************************
	 * 1. non delegator(sdk.Cfg.Addr4) try to withdraw stake, should fail
	 *******************************************************************************/
	Start("WithdwarStakeAmount, Do it under admin")
	adminAddr := utils.GetAddressByWallet(celestials.Admin)
	AssertSuccess(p.StZIL.Transfer(adminAddr, stakeAmt))

	p.StZIL.SetSigner(celestials.Admin)

	AssertSuccess(p.StZIL.SlashSSN(stakeAmt, ssnNode1))
	bnum1 := txn.Receipt.EpochNum

	withdrawal := Dig(p.StZIL, "withdrawal_pending", bnum1, adminAddr).Withdrawal()

	AssertEqual(withdrawal.TokenAmount.String(), stakeAmt)
	AssertEqual(withdrawal.StakeAmount.String(), stakeAmt)
}
