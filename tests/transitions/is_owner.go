package transitions

import (
	"github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	. "github.com/avely-finance/avely-contracts/tests/helpers"
)

func (tr *Transitions) IsOwner() {

	Start("IsOwner")

	p := tr.DeployAndUpgrade()

	// Use non-owner user for StZIL, expecting errors
	p.StZIL.SetSigner(bob)
	randomAddr := utils.GetAddressByWallet(eve)

	tx, _ := p.StZIL.ChangeAdmin(randomAddr)
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	tx, _ = p.StZIL.ChangeOwner(randomAddr)
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	new_buffers := []string{p.GetBuffer().Addr, p.GetBuffer().Addr}
	tx, _ = p.StZIL.ChangeBuffers(new_buffers)
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	tx, _ = p.StZIL.AddSSN(randomAddr)
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	tx, _ = p.StZIL.RemoveSSN(randomAddr)
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	tx, _ = p.StZIL.SetHolderAddress(randomAddr)
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	tx, _ = p.StZIL.ChangeRewardsFee("100")
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	tx, _ = p.StZIL.ChangeTreasuryAddress(randomAddr)
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	tx, _ = p.StZIL.ChangeZimplAddress(randomAddr)
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	tx, _ = p.StZIL.UpdateStakingParameters(ToZil(100))
	AssertError(tx, p.StZIL.ErrorCode("CodeNotOwner"))

	// Holder
	tx, _ = p.Holder.ChangeOwner(randomAddr)
	AssertError(tx, p.Holder.ErrorCode("CodeNotOwner"))

	p.Holder.SetSigner(bob)
	tx, _ = p.Holder.ChangeZimplAddress(randomAddr)
	AssertError(tx, p.Holder.ErrorCode("CodeNotOwner"))

	tx, _ = p.Holder.ChangeZproxyAddress(randomAddr)
	AssertError(tx, p.Holder.ErrorCode("CodeNotOwner"))
}
