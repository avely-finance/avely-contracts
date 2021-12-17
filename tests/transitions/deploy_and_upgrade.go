package transitions

import (
	"Azil/test/contracts"
	"Azil/test/helpers"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"log"
)

func (tr *Transitions) DeployAndUpgrade() (*contracts.Zproxy, *contracts.Zimpl, *contracts.AZil, *contracts.BufferContract, *contracts.HolderContract) {
	log.Println("start to deploy")

	//deploy gzil
	gzil, err := contracts.NewGzil(adminKey)
	if err != nil {
		t.LogError("deploy Gzil error = ", err)
	}
	log.Println("deploy Gzil succeed, address = ", gzil.Addr)

	//deploy Zproxy
	Zproxy, err := contracts.NewZproxy(adminKey)
	if err != nil {
		t.LogError("deploy Zproxy error = ", err)
	}
	log.Println("deploy Zproxy succeed, address = ", Zproxy.Addr)

	//deploy Zimpl
	Zimpl, err := contracts.NewZimpl(adminKey, Zproxy.Addr, gzil.Addr)
	if err != nil {
		t.LogError("deploy Zimpl error = ", err)
	}
	log.Println("deploy Zimpl succeed, address = ", Zimpl.Addr)

	//deploy azil
	Aimpl, err := contracts.NewAZilContract(adminKey, AZIL_SSN_ADDRESS, Zimpl.Addr)
	if err != nil {
		t.LogError("deploy aZil error = ", err)
	}
	log.Println("deploy aZil succeed, address = ", Aimpl.Addr)

	//deploy buffer
	Buffer, err := contracts.NewBufferContract(adminKey, Aimpl.Addr /*aimpl_address*/, AZIL_SSN_ADDRESS, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		t.LogError("deploy buffer error = ", err)
	}
	log.Println("deploy buffer succeed, address = ", Buffer.Addr)

	//deploy holder
	Holder, err := contracts.NewHolderContract(adminKey, Aimpl.Addr /*aimpl_address*/, AZIL_SSN_ADDRESS, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		t.LogError("deploy holder error = ", err)
	}
	log.Println("deploy holder succeed, address = ", Holder.Addr)

	log.Println("start to upgrade")
	/********************************************************************
	* Upgrade buffer/holder
	********************************************************************/
	new_buffers := []string{"0x" + Buffer.Addr}
	t.AssertSuccess(Aimpl.ChangeBuffers(new_buffers))
	t.AssertSuccess(Aimpl.ChangeHolderAddress(Holder.Addr))

	/********************************************************************
	* Upgrade Zproxy, make some initial actions
	********************************************************************/
	args := []core.ContractValue{
		{
			"newImplementation",
			"ByStr20",
			"0x" + Zimpl.Addr,
		},
	}
	t.AssertSuccess(Zproxy.Call("UpgradeTo", args, "0"))
	t.AssertSuccess(Zproxy.AddSSN(AZIL_SSN_ADDRESS, "aZil SSN"))
	t.AssertSuccess(Zproxy.UpdateVerifierRewardAddr("0x" + verifier))
	t.AssertSuccess(Zproxy.UpdateVerifier("0x" + verifier))
	t.AssertSuccess(Zproxy.UpdateStakingParameters(zil(1000), zil(10))) //minstake (ssn not active if less), mindelegstake
	t.AssertSuccess(Zproxy.Unpause())

	//we need our SSN to be active, so delegating some stake
	t.AssertSuccess(Aimpl.DelegateStake(zil(1000)))
	t.AssertEqual(Zimpl.Field("direct_deposit_deleg", "0x"+Buffer.Addr, AZIL_SSN_ADDRESS, "1"), zil(1000))

	//we need to delegate something from Holder, in order to make Zimpl know holder's address
	t.AssertSuccess(Holder.DelegateStake(zil(HOLDER_INITIAL_DELEGATE_ZIL)))

	//SSN will become active on next cycle
	Zproxy.UpdateWallet(verifierKey)
	//we need to increase blocknum, in order to Gzil won't mint anything. Really minting is over.
	helpers.IncreaseBlocknum(10)
	t.AssertSuccess(Zproxy.AssignStakeReward(AZIL_SSN_ADDRESS, AZIL_SSN_REWARD_SHARE_PERCENT))

	log.Println("upgrade succeed")
	t.AddShortcut("Zproxy", "0x"+Zproxy.Addr)
	t.AddShortcut("Zimpl", "0x"+Zimpl.Addr)
	t.AddShortcut("Gzil", "0x"+gzil.Addr)
	t.AddShortcut("Aimpl", "0x"+Aimpl.Addr)
	t.AddShortcut("Buffer", "0x"+Buffer.Addr)
	t.AddShortcut("Holder", "0x"+Holder.Addr)

	return Zproxy, Zimpl, Aimpl, Buffer, Holder
}
