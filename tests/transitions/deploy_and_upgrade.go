package transitions

import (
	"Azil/test/deploy"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"log"
)

func (t *Testing) DeployAndUpgrade() (*deploy.Zproxy, *deploy.Zimpl, *deploy.AZil, *deploy.BufferContract, *deploy.HolderContract) {
	log.Println("start to deploy")

	//deploy gzil
	gzil, err := deploy.NewGzil(adminKey)
	if err != nil {
		t.LogError("deploy Gzil error = ", err)
	}
	log.Println("deploy Gzil succeed, address = ", gzil.Addr)

	//deploy Zproxy
	Zproxy, err := deploy.NewZproxy(adminKey)
	if err != nil {
		t.LogError("deploy Zproxy error = ", err)
	}
	log.Println("deploy Zproxy succeed, address = ", Zproxy.Addr)

	//deploy Zimpl
	Zimpl, err := deploy.NewZimpl(adminKey, Zproxy.Addr, gzil.Addr)
	if err != nil {
		t.LogError("deploy Zimpl error = ", err)
	}
	log.Println("deploy Zimpl succeed, address = ", Zimpl.Addr)

	//deploy azil
	Aimpl, err := deploy.NewAZilContract(adminKey, AZIL_SSN_ADDRESS, Zproxy.Addr)
	if err != nil {
		t.LogError("deploy aZil error = ", err)
	}
	log.Println("deploy aZil succeed, address = ", Aimpl.Addr)

	//deploy buffer
	Buffer, err := deploy.NewBufferContract(adminKey, Aimpl.Addr /*aimpl_address*/, AZIL_SSN_ADDRESS, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		t.LogError("deploy buffer error = ", err)
	}
	log.Println("deploy buffer succeed, address = ", Buffer.Addr)

	//deploy holder
	Holder, err := deploy.NewHolderContract(adminKey, Aimpl.Addr /*aimpl_address*/, AZIL_SSN_ADDRESS, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		t.LogError("deploy holder error = ", err)
	}
	log.Println("deploy holder succeed, address = ", Holder.Addr)

	log.Println("start to upgrade")
	/********************************************************************
	* Upgrade buffer/holder
	********************************************************************/
	new_buffers := []string{"0x" + Buffer.Addr}
	if _, err := Aimpl.ChangeBuffers(new_buffers); err != nil {
		t.LogError("failed to change aZil's buffer contract address; error = ", err)
	}
	if _, err := Aimpl.ChangeHolderAddress(Holder.Addr); err != nil {
		t.LogError("failed to change aZil's holder contract address; error = ", err)
	}

	/********************************************************************
	* Upgrade Zproxy
	********************************************************************/
	args := []core.ContractValue{
		{
			"newImplementation",
			"ByStr20",
			"0x" + Zimpl.Addr,
		},
	}
	_, err = Zproxy.Call("UpgradeTo", args, "0")
	if err != nil {
		t.LogError("Zproxy UpgradeTo failed", err)
	}
	Zproxy.AddSSN(AZIL_SSN_ADDRESS, "aZil SSN")
	Zproxy.UpdateVerifierRewardAddr("0x" + verifier)
	Zproxy.UpdateVerifier("0x" + verifier)
	Zproxy.UpdateStakingParameters(zil(1000), zil(10)) //minstake (ssn not active if less), mindelegstake
	Zproxy.Unpause()

	//we need our SSN to be active, so delegating some stake
	_, err = Aimpl.DelegateStake(zil(1000))
	if err != nil {
		t.LogError("DelegateStake", err)
	}
	t.AssertEqual(Zimpl.Field("direct_deposit_deleg", "0x"+Buffer.Addr, AZIL_SSN_ADDRESS, "1"), zil(1000))

	//SSN will become active on next cycle
	Zproxy.UpdateWallet(verifierKey)
	Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT)

	log.Println("upgrade succeed")
	t.AddDebug("Zproxy", "0x"+Zproxy.Addr)
	t.AddDebug("Zimpl", "0x"+Zimpl.Addr)
	t.AddDebug("Gzil", "0x"+gzil.Addr)
	t.AddDebug("Aimpl", "0x"+Aimpl.Addr)
	t.AddDebug("Buffer", "0x"+Buffer.Addr)
	t.AddDebug("Holder", "0x"+Holder.Addr)

	return Zproxy, Zimpl, Aimpl, Buffer, Holder
}
