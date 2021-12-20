package transitions

import (
	"Azil/test/contracts"
	"Azil/test/helpers"
	"github.com/Zilliqa/gozilliqa-sdk/core"
)

func (tr *Transitions) DeployAndUpgrade() (*contracts.Zproxy, *contracts.Zimpl, *contracts.AZil, *contracts.BufferContract, *contracts.HolderContract) {
	log.Info("start to deploy")

	//deploy gzil
	gzil, err := contracts.NewGzil(tr.cfg.AdminKey)
	if err != nil {
		log.Fatal("deploy Gzil error = " + err.Error())
	}
	log.Success("deploy Gzil succeed, address = " + gzil.Addr)

	//deploy Zproxy
	Zproxy, err := contracts.NewZproxy(tr.cfg.AdminKey)
	if err != nil {
		log.Fatal("deploy Zproxy error = " + err.Error())
	}
	log.Success("deploy Zproxy succeed, address = " + Zproxy.Addr)

	//deploy Zimpl
	Zimpl, err := contracts.NewZimpl(tr.cfg.AdminKey, Zproxy.Addr, gzil.Addr)
	if err != nil {
		log.Fatal("deploy Zimpl error = " + err.Error())
	}
	log.Success("deploy Zimpl succeed, address = " + Zimpl.Addr)

	//deploy azil
	Aimpl, err := contracts.NewAZilContract(tr.cfg.AdminKey, tr.cfg.AzilSsnAddress, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy aZil error = " + err.Error())
	}
	log.Success("deploy aZil succeed, address = " + Aimpl.Addr)

	//deploy buffer
	Buffer, err := contracts.NewBufferContract(tr.cfg.AdminKey, Aimpl.Addr, tr.cfg.AzilSsnAddress, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy buffer error = " + err.Error())
	}
	log.Success("deploy buffer succeed, address = " + Buffer.Addr)

	//deploy holder
	Holder, err := contracts.NewHolderContract(tr.cfg.AdminKey, Aimpl.Addr, tr.cfg.AzilSsnAddress, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy holder error = " + err.Error())
	}
	log.Success("deploy holder succeed, address = " + Holder.Addr)

	log.Info("start to upgrade")
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
	t.AssertSuccess(Zproxy.AddSSN(tr.cfg.AzilSsnAddress, "aZil SSN"))
	t.AssertSuccess(Zproxy.UpdateVerifierRewardAddr("0x" + tr.cfg.Verifier))
	t.AssertSuccess(Zproxy.UpdateVerifier("0x" + tr.cfg.Verifier))
	t.AssertSuccess(Zproxy.UpdateStakingParameters(zil(1000), zil(10))) //minstake (ssn not active if less), mindelegstake
	t.AssertSuccess(Zproxy.Unpause())

	//we need our SSN to be active, so delegating some stake
	t.AssertSuccess(Aimpl.DelegateStake(zil(1000)))
	t.AssertEqual(Zimpl.Field("direct_deposit_deleg", "0x"+Buffer.Addr, tr.cfg.AzilSsnAddress, "1"), zil(1000))

	//we need to delegate something from Holder, in order to make Zimpl know holder's address
	t.AssertSuccess(Holder.DelegateStake(zil(tr.cfg.HolderInitialDelegateZil)))

	//SSN will become active on next cycle
	Zproxy.UpdateWallet(tr.cfg.VerifierKey)
	//we need to increase blocknum, in order to Gzil won't mint anything. Really minting is over.
	helpers.IncreaseBlocknum(10)
	t.AssertSuccess(Zproxy.AssignStakeReward(tr.cfg.AzilSsnAddress, tr.cfg.AzilSsnRewardSharePercent))

	log.AddShortcut("Zproxy", "0x"+Zproxy.Addr)
	log.AddShortcut("Zimpl", "0x"+Zimpl.Addr)
	log.AddShortcut("Gzil", "0x"+gzil.Addr)
	log.AddShortcut("Aimpl", "0x"+Aimpl.Addr)
	log.AddShortcut("Buffer", "0x"+Buffer.Addr)
	log.AddShortcut("Holder", "0x"+Holder.Addr)
	log.Success("upgrade succeed")

	return Zproxy, Zimpl, Aimpl, Buffer, Holder
}
