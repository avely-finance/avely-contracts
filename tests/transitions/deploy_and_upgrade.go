package transitions

import (
	"Azil/test/deploy"
	// "github.com/Zilliqa/gozilliqa-sdk/core"
	"log"
)

// this is a help function : )
func (t *Testing) DeployAndUpgrade() (*deploy.StubStakingContract, *deploy.AZil, *deploy.BufferContract) {
	// log.Println("start to deploy proxy contract")
	// proxy, err := deploy.NewProxy(key1)
	// if err != nil {
	// 	t.LogError("deploy proxy error = ", err)
	// }
	// log.Println("deploy proxy succeed, address = ", proxy.Addr)

	log.Println("start to deploy stubStaking contract")
	stubStakingContract, err1 := deploy.NewStubStakingContract(key1)
	if err1 != nil {
		t.LogError("deploy stubStaking error = ", err1)
	}
	log.Println("deploy stubStaking succeed, address = ", stubStakingContract.Addr)

	aZilContract, err1 := deploy.NewAZilContract(key1, aZilSSNAddress, stubStakingContract.Addr)
	if err1 != nil {
		t.LogError("deploy aZil error = ", err1)
	}
	log.Println("deploy aZil succeed, address = ", aZilContract.Addr)

	bufferContract, err1 := deploy.NewBufferContract(key1, aZilSSNAddress, stubStakingContract.Addr)
	if err1 != nil {
		t.LogError("deploy buffer error = ", err1)
	}
	log.Println("deploy buffer succeed, address = ", bufferContract.Addr)

	log.Println("start to upgrade")

	if _, err := aZilContract.ChangeBufferAddress(bufferContract.Addr); err != nil {
		t.LogError("failed to change aZil's buffer contract address; error = ", err)
	}

	// args := []core.ContractValue{
	// 	{
	// 		"newImplementation",
	// 		"ByStr20",
	// 		"0x" + ssnlist.Addr,
	// 	},
	// }
	// _, err2 := proxy.Call("UpgradeTo", args,"0")
	// if err2 != nil {
	// 	t.LogError("UpgradeTo failed", err2)
	// }
	// log.Println("upgrade succeed")

	return stubStakingContract, aZilContract, bufferContract
}
