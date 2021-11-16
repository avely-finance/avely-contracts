package transitions

import (
	"Azil/test/deploy"
	"log"
)

func (t *Testing) DeployAndUpgrade() (*deploy.StubStakingContract, *deploy.AZil, *deploy.BufferContract,  *deploy.HolderContract) {
	log.Println("start to deploy")

	//deploy stubStakingContract
	stubStakingContract, err1 := deploy.NewStubStakingContract(adminKey)
	if err1 != nil {
		t.LogError("deploy stubStaking error = ", err1)
	}
	log.Println("deploy stubStaking succeed, address = ", stubStakingContract.Addr)

	//deploy azil
	aZilContract, err1 := deploy.NewAZilContract(adminKey, aZilSSNAddress, stubStakingContract.Addr)
	if err1 != nil {
		t.LogError("deploy aZil error = ", err1)
	}
	log.Println("deploy aZil succeed, address = ", aZilContract.Addr)

	//deploy buffer
	bufferContract, err1 := deploy.NewBufferContract(adminKey, aZilSSNAddress, stubStakingContract.Addr)
	if err1 != nil {
		t.LogError("deploy buffer error = ", err1)
	}
	log.Println("deploy buffer succeed, address = ", bufferContract.Addr)

	//deploye holder
	holderContract, err1 := deploy.NewHolderContract(adminKey, aZilSSNAddress, stubStakingContract.Addr)
	if err1 != nil {
		t.LogError("deploy holder error = ", err1)
	}
	log.Println("deploy holder succeed, address = ", holderContract.Addr)

	log.Println("start to upgrade")

	if _, err := stubStakingContract.AddSSN(aZilSSNAddress); err != nil {
		t.LogError("failed to stubStakingContract.AddSSN(aZilSSNAddress); error = ", err)
	}
	if _, err := aZilContract.ChangeBufferAddress(bufferContract.Addr); err != nil {
		t.LogError("failed to change aZil's buffer contract address; error = ", err)
	}
	if _, err := aZilContract.ChangeHolderAddress(holderContract.Addr); err != nil {
		t.LogError("failed to change aZil's holder contract address; error = ", err)
	}

	log.Println("upgrade succeed")



	return stubStakingContract, aZilContract, bufferContract, holderContract
}
