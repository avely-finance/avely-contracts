package transitions

import (
//"log"
)

func (t *Testing) DelegateStakeSuccess() {
	t.LogStart("DelegateStake")

	// deploy smart contract
	stubStakingContract, aZilContract, bufferContract, _ := t.DeployAndUpgrade()

	_, err := aZilContract.DelegateStake(unit10)
	if err != nil {
		t.LogError("DelegateStake", err)
	}

	stubStakingState := stubStakingContract.LogContractStateJson()
	t.AssertContain(stubStakingState, "_balance\":\"" + unit10)
	t.AssertContain(stubStakingState, "buff_deposit_deleg\":{\"" + "0x" + bufferContract.Addr + "\":{\"" + aZilSSNAddress + "\":{\"1\":\"" + unit10)

	aZilState := aZilContract.LogContractStateJson()
	t.AssertContain(aZilState, "_balance\":\"0")
	t.AssertContain(aZilState, "\"totalstakeamount\":\"" + unit10 + "\",\"totaltokenamount\":\"" + unit10 + "\"")
	t.AssertContain(aZilState, "balances\":{\""+"0x" + admin + "\":\"" + unit10)

	tx, err := aZilContract.ZilBalanceOf(admin)
	if err != nil {
		t.LogError("ZilBalanceOf", err)
	}
	t.LogPrettyReceipt(tx)


	t.LogEnd("DelegateStake")
}
