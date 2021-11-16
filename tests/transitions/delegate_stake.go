package transitions

import (
// "log"
)

func (t *Testing) DelegateStakeSuccess() {
	t.LogStart("DelegateStake")

	// deploy smart contract
	stubStakingContract, aZilContract, bufferContract, _ := t.DeployAndUpgrade()

	_, err := aZilContract.DelegateStake(tenzil)
	if err != nil {
		t.LogError("DelegateStake", err)
	}

	stubStakingState := stubStakingContract.LogContractStateJson()
	t.AssertContain(stubStakingState, "_balance\":\""+tenzil)
	t.AssertContain(stubStakingState, "buff_deposit_deleg\":{\""+"0x"+bufferContract.Addr+"\":{\""+aZilSSNAddress+"\":{\"1\":\""+tenzil)

	aZilState := aZilContract.LogContractStateJson()
	t.AssertContain(aZilState, "_balance\":\"0")
	t.AssertContain(aZilState, "\"totalstakeamount\":\"" + tenzil+ "\",\"totaltokenamount\":\"" + tenzil+ "\"")
	t.AssertContain(aZilState, "balances\":{\""+"0x"+admin+"\":\""+tenzil)

	t.LogEnd("DelegateStake")
}
