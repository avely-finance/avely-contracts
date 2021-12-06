package transitions

import (
	"Azil/test/deploy"
	"log"
)

var FetcherContract deploy.FetcherContract

func (t *Testing) DeployAndUpgrade() (*deploy.StubStakingContract, *deploy.AZil, *deploy.BufferContract, *deploy.HolderContract) {
	log.Println("start to deploy")

	//deploy AzilLib
	azilUtils, err1 := deploy.NewAzilUtilsContract(adminKey)
	if err1 != nil {
		t.LogError("deploy AzilUtils error = ", err1)
	}
	log.Println("deploy AzilUtils succeed, address = ", azilUtils.Addr)

	//deploy stubStakingContract
	stubStakingContract, err1 := deploy.NewStubStakingContract(adminKey)
	if err1 != nil {
		t.LogError("deploy stubStaking error = ", err1)
	}
	log.Println("deploy stubStaking succeed, address = ", stubStakingContract.Addr)

	//deploy azil
	aZilContract, err1 := deploy.NewAZilContract(adminKey, azilUtils.Addr, aZilSSNAddress, stubStakingContract.Addr)
	if err1 != nil {
		t.LogError("deploy aZil error = ", err1)
	}
	log.Println("deploy aZil succeed, address = ", aZilContract.Addr)

	//deploy buffer
	bufferContract, err1 := deploy.NewBufferContract(adminKey, aZilContract.Addr /*aimpl_address*/, aZilSSNAddress, stubStakingContract.Addr)
	if err1 != nil {
		t.LogError("deploy buffer error = ", err1)
	}
	log.Println("deploy buffer succeed, address = ", bufferContract.Addr)

	//deploy holder
	holderContract, err1 := deploy.NewHolderContract(adminKey, aZilContract.Addr /*aimpl_address*/, aZilSSNAddress, stubStakingContract.Addr)
	if err1 != nil {
		t.LogError("deploy holder error = ", err1)
	}
	log.Println("deploy holder succeed, address = ", holderContract.Addr)

	//deploy fetcher
	fetcherContract, err1 := deploy.NewFetcherContract(adminKey, azilUtils.Addr, aZilContract.Addr, stubStakingContract.Addr)
	FetcherContract = *fetcherContract
	if err1 != nil {
		t.LogError("deploy fetcher error = ", err1)
	}
	log.Println("deploy fetcher succeed, address = ", fetcherContract.Addr)

	//upgrade contracts with correct field values
	log.Println("start to upgrade")

	if _, err := stubStakingContract.AddSSN(aZilSSNAddress); err != nil {
		t.LogError("failed to stubStakingContract.AddSSN(aZilSSNAddress); error = ", err)
	}

	new_buffers := []string{"0x" + bufferContract.Addr}
	if _, err := aZilContract.ChangeBuffers(new_buffers); err != nil {
		t.LogError("failed to change aZil's buffer contract address; error = ", err)
	}
	if _, err := aZilContract.ChangeHolderAddress(holderContract.Addr); err != nil {
		t.LogError("failed to change aZil's holder contract address; error = ", err)
	}

	log.Println("upgrade succeed")

	t.AddDebug("stubStakingContract", "0x"+stubStakingContract.Addr)
	t.AddDebug("aZilContract", "0x"+aZilContract.Addr)
	t.AddDebug("bufferContract", "0x"+bufferContract.Addr)
	t.AddDebug("holderContract", "0x"+holderContract.Addr)
	t.AddDebug("fetcherContract", "0x"+FetcherContract.Addr)
	t.AddDebug("azilUtilsContract", "0x"+azilUtils.Addr)

	return stubStakingContract, aZilContract, bufferContract, holderContract
}
