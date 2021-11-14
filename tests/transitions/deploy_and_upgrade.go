package transitions

import (
	"Azil/test/deploy"
	// "github.com/Zilliqa/gozilliqa-sdk/core"
	"log"
)

// this is a help function : )
func (t *Testing) DeployAndUpgrade() (*deploy.StubStakingContract, *deploy.BufferContract) {
	// log.Println("start to deploy proxy contract")
	// proxy, err := deploy.NewProxy(key1)
	// if err != nil {
	// 	t.LogError("deploy proxy error = ", err)
	// }
	// log.Println("deploy proxy succeed, address = ", proxy.Addr)

	log.Println("start to deploy stubStaking contract")
	stubStaking, err1 := deploy.NewStubStakingContract(key1)
	if err1 != nil {
		t.LogError("deploy stubStaking error = ", err1)
	}
	log.Println("deploy stubStaking succeed, address = ", stubStaking.Addr)

	bufferContract, err1 := deploy.NewBufferContract(key1)
	if err1 != nil {
		t.LogError("deploy buffer error = ", err1)
	}
	log.Println("deploy buffer succeed, address = ", bufferContract.Addr)

	if _, err := bufferContract.ChangeProxyStakingContractAddress(stubStaking.Addr); err != nil {
		t.LogError("failed to change buffer's staking contract address; error = ", err)
	}

	// log.Println("start to upgrade")
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

	return stubStaking, bufferContract
}
