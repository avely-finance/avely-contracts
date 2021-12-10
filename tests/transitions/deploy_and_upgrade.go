package transitions

import (
	"Azil/test/deploy"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"log"
)

func (t *Testing) DeployAndUpgrade() (*deploy.StubStakingContract, *deploy.AZil, *deploy.BufferContract, *deploy.HolderContract) {
	log.Println("start to deploy")

	//deploy stubStakingContract
	stubStakingContract, err := deploy.NewStubStakingContract(adminKey)
	if err != nil {
		t.LogError("deploy stubStaking error = ", err)
	}
	log.Println("deploy stubStaking succeed, address = ", stubStakingContract.Addr)

	//deploy azil
	aZilContract, err := deploy.NewAZilContract(adminKey, aZilSSNAddress, stubStakingContract.Addr)
	if err != nil {
		t.LogError("deploy aZil error = ", err)
	}
	log.Println("deploy aZil succeed, address = ", aZilContract.Addr)

	//deploy buffer
	bufferContract, err := deploy.NewBufferContract(adminKey, aZilContract.Addr /*aimpl_address*/, aZilSSNAddress, stubStakingContract.Addr)
	if err != nil {
		t.LogError("deploy buffer error = ", err)
	}
	log.Println("deploy buffer succeed, address = ", bufferContract.Addr)

	//deploy holder
	holderContract, err := deploy.NewHolderContract(adminKey, aZilContract.Addr /*aimpl_address*/, aZilSSNAddress, stubStakingContract.Addr)
	if err != nil {
		t.LogError("deploy holder error = ", err)
	}
	log.Println("deploy holder succeed, address = ", holderContract.Addr)

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

	return stubStakingContract, aZilContract, bufferContract, holderContract
}

func (t *Testing) DeployAndUpgradeOriginal() (*deploy.Zproxy, *deploy.Zimpl) {
	log.Println("start to deploy")

	//deploy gzil
	gzil, err := deploy.NewGzil(adminKey)
	if err != nil {
		t.LogError("deploy Gzil error = ", err)
	}
	log.Println("deploy Gzil succeed, address = ", gzil.Addr)

	//deploy zproxy
	zproxy, err := deploy.NewZproxy(adminKey)
	if err != nil {
		t.LogError("deploy Zproxy error = ", err)
	}
	log.Println("deploy Zproxy succeed, address = ", zproxy.Addr)

	//deploy zimpl
	zimpl, err := deploy.NewZimpl(adminKey, zproxy.Addr, gzil.Addr)
	if err != nil {
		t.LogError("deploy Zimpl error = ", err)
	}
	log.Println("deploy Zimpl succeed, address = ", zimpl.Addr)

	//upgrade contracts with correct field values
	log.Println("start to upgrade")

	args := []core.ContractValue{
		{
			"newImplementation",
			"ByStr20",
			"0x" + zimpl.Addr,
		},
	}
	_, err = zproxy.Call("UpgradeTo", args, "0")
	if err != nil {
		t.LogError("Zproxy UpgradeTo failed", err)
	}

	log.Println("upgrade succeed")

	t.AddDebug("Zproxy", "0x"+zproxy.Addr)
	t.AddDebug("Zimpl", "0x"+zimpl.Addr)
	t.AddDebug("Gzil", "0x"+gzil.Addr)

	return zproxy, zimpl

}
