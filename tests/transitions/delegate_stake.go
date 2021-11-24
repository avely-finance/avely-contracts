package transitions

import (
//"log"
)

func (t *Testing) DelegateStakeSuccess() {
	t.LogStart("DelegateStake")

	// deploy smart contract
	stubStakingContract, aZilContract, bufferContract, _ := t.DeployAndUpgrade()

	_, err := aZilContract.DelegateStake(zil10)
	if err != nil {
		t.LogError("DelegateStake", err)
	}

	stubStakingState := stubStakingContract.LogContractStateJson()
	t.AssertContain(stubStakingState, "_balance\":\""+zil10)
	t.AssertContain(stubStakingState, "buff_deposit_deleg\":{\""+"0x"+bufferContract.Addr+"\":{\""+aZilSSNAddress+"\":{\"1\":\""+zil10)

	aZilState := aZilContract.LogContractStateJson()
	t.AssertContain(aZilState, "_balance\":\"0")
	t.AssertContain(aZilState, "\"totalstakeamount\":\""+zil10+"\",\"totaltokenamount\":\""+azil10+"\"")
	t.AssertContain(aZilState, "balances\":{\""+"0x"+admin+"\":\""+azil10)

	t.LogEnd("DelegateStake")
}
