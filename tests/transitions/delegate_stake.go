package transitions

import (
	//"log"
	"Azil/test/deploy"
)

func (t *Testing) DelegateStakeSuccess() {
	t.LogStart("DelegateStake: Stake 10 ZIL")

	stubStakingContract, aZilContract, bufferContract, _ := t.DeployAndUpgrade()

	_, err := aZilContract.DelegateStake(zil(10))
	if err != nil {
		t.LogError("DelegateStake", err)
	}

	stubStakingState := stubStakingContract.LogContractStateJson()
	t.AssertContain(stubStakingState, "_balance\":\""+zil(10))
	t.AssertContain(stubStakingState, "buff_deposit_deleg\":{\""+"0x"+bufferContract.Addr+"\":{\""+aZilSSNAddress+"\":{\"1\":\""+zil(10))

	aZilState := aZilContract.LogContractStateJson()
	t.AssertContain(aZilState, "_balance\":\"0")
	t.AssertContain(aZilState, "\"totalstakeamount\":\""+zil(10)+"\",\"totaltokenamount\":\""+azil(10)+"\"")
	t.AssertContain(aZilState, "balances\":{\""+"0x"+admin+"\":\""+azil(10))
}

func (t *Testing) DelegateStakeBuffersRotation() {
	t.LogStart("DelegateStake: Buffers rotation")

	stubStakingContract, aZilContract, bufferContract, _ := t.DeployAndUpgrade()

	anotherBufferContract, err1 := deploy.NewBufferContract(adminKey, aZilContract.Addr, aZilSSNAddress, stubStakingContract.Addr)
	if err1 != nil {
		t.LogError("Deploy buffer error = ", err1)
	}

	new_buffers := []string{"0x" + bufferContract.Addr, "0x" + bufferContract.Addr, "0x" + anotherBufferContract.Addr}

	aZilContract.ChangeBuffers(new_buffers)
	stubStakingContract.AssignStakeReward() // move to the next cycle

	_, err := aZilContract.DelegateStake(zil(10))
	if err != nil {
		t.LogError("DelegateStake", err)
	}

	stubStakingState := stubStakingContract.LogContractStateJson()
	t.AssertContain(stubStakingState, "_balance\":\""+zil(10))

	// lastrewardcycly = 2; buffers has 3 elements
	// => active buffer = buffers[ 2 % 3 ] = buffers[2] = anotherBufferConract
	t.AssertContain(stubStakingState, "buff_deposit_deleg\":{\""+"0x"+anotherBufferContract.Addr+"\":{\""+aZilSSNAddress+"\":{\"2\":\""+zil(10))
}
