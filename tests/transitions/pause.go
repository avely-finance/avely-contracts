package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) Pause() {
	Start("IsAimpl")

	p := tr.DeployAndUpgrade()

	callPaused(p)
	callPauseNonAdmin(p)
}

func callPauseNonAdmin(p *contracts.Protocol) {
	//call pause/unpause admin transitions with non-admin user; expecting errors
	p.Aimpl.UpdateWallet(sdk.Cfg.Key1)

	tx, _ := p.Aimpl.Pause()
	AssertError(tx, "AdminValidationFailed")

	tx, _ = p.Aimpl.Unpause()
	AssertError(tx, "AdminValidationFailed")
}

func callPaused(p *contracts.Protocol) {
	//call user's transitions, when contract is paused; expecting errors
	AssertSuccess(p.Aimpl.Pause())

	tx, _ := p.Aproxy.DelegateStake(ToZil(10))
	AssertError(tx, "Paused")

	p.Aproxy.ZilBalanceOf(sdk.Cfg.Addr1)
	tx = sdk.TxLast
	AssertError(tx, "Paused")

	tx, _ = p.Aproxy.WithdrawStakeAmt(ToZil(10))
	AssertError(tx, "Paused")

	tx, _ = p.Aproxy.CompleteWithdrawal()
	AssertError(tx, "Paused")

	AssertSuccess(p.Aimpl.Unpause())
}
