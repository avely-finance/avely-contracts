package transitions

import (
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsAdmin() {

	Start("IsAdmin")

	p := tr.DeployAndUpgrade()

	// Use non-admin user for StZIL, expecting errors
	p.StZIL.UpdateWallet(sdk.Cfg.Key2)

	tx, _ := p.StZIL.IncreaseAutoRestakeAmount(ToZil(1))
	AssertError(tx, p.StZIL.ErrorCode("AdminValidationFailed"))

	tx, _ = p.StZIL.PerformAutoRestake()
	AssertError(tx, p.StZIL.ErrorCode("AdminValidationFailed"))
	tx, _ = p.StZIL.ClaimRewards(p.GetBuffer().Addr, sdk.Cfg.SsnAddrs[0])
	AssertError(tx, p.StZIL.ErrorCode("AdminValidationFailed"))
	readyBlocks := []string{}
	tx, _ = p.StZIL.ClaimWithdrawal(readyBlocks)
	AssertError(tx, p.StZIL.ErrorCode("AdminValidationFailed"))
	tx, _ = p.StZIL.ChownStakeReDelegate(sdk.Cfg.Addr3, ToZil(1))
	AssertError(tx, p.StZIL.ErrorCode("AdminValidationFailed"))
}
