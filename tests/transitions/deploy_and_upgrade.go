package transitions

import (
	"Azil/test/deploy"
	"log"
)

func (t *Testing) DeployAndUpgrade() (*deploy.StubStakingContract, *deploy.AZil, *deploy.BufferContract) {
	log.Println("start to deploy stubStaking contract")
	stubStakingContract, err1 := deploy.NewStubStakingContract(adminKey)
	if err1 != nil {
		t.LogError("deploy stubStaking error = ", err1)
	}
	log.Println("deploy stubStaking succeed, address = ", stubStakingContract.Addr)

	aZilContract, err1 := deploy.NewAZilContract(adminKey, aZilSSNAddress, stubStakingContract.Addr)
	if err1 != nil {
		t.LogError("deploy aZil error = ", err1)
	}
	log.Println("deploy aZil succeed, address = ", aZilContract.Addr)

	bufferContract, err1 := deploy.NewBufferContract(adminKey, aZilSSNAddress, stubStakingContract.Addr)
	if err1 != nil {
		t.LogError("deploy buffer error = ", err1)
	}
	log.Println("deploy buffer succeed, address = ", bufferContract.Addr)

	log.Println("start to upgrade")

	if _, err := aZilContract.ChangeBufferAddress(bufferContract.Addr); err != nil {
		t.LogError("failed to change aZil's buffer contract address; error = ", err)
	}

	log.Println("upgrade succeed")

	return stubStakingContract, aZilContract, bufferContract
}
