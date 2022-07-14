package contracts

import (
	"github.com/Zilliqa/gozilliqa-sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
)

type ZilliqaStaking struct {
	Zproxy *Zproxy
	Zimpl  *Zimpl
	Gzil   *Gzil
}

func NewZilliqaStaking(zproxy *Zproxy, zimpl *Zimpl, gzil *Gzil) *ZilliqaStaking {
	return &ZilliqaStaking{
		Zproxy: zproxy,
		Zimpl:  zimpl,
		Gzil:   gzil,
	}
}

func DeployZilliqaStaking(sdk *AvelySDK, log *Log) *ZilliqaStaking {
	log.Debug("start to deploy zilliqa staking contracts")

	//deploy gzil
	gzil, err := NewGzil(sdk)
	if err != nil {
		log.Fatal("deploy Gzil error = " + err.Error())
	}
	log.Debug("deploy Gzil succeed, address = " + gzil.Addr)

	//deploy Zproxy
	zproxy, err := NewZproxy(sdk)
	if err != nil {
		log.Fatal("deploy Zproxy error = " + err.Error())
	}
	log.Debug("deploy Zproxy succeed, address = " + zproxy.Addr)

	//deploy Zimpl
	zimpl, err := NewZimpl(sdk, zproxy.Addr, gzil.Addr)
	if err != nil {
		log.Fatal("deploy Zimpl error = " + err.Error())
	}
	log.Debug("deploy Zimpl succeed, address = " + zimpl.Addr)

	return NewZilliqaStaking(zproxy, zimpl, gzil)
}

func SetupZilliqaStaking(sdk *AvelySDK, log *Log) {

	//Restore Zproxy
	Zproxy, err := RestoreZproxy(sdk, sdk.Cfg.ZproxyAddr)
	if err != nil {
		log.Fatal("Restore Zproxy error = " + err.Error())
	}
	log.Debug("Restore Zproxy succeed, address = " + Zproxy.Addr)

	args := []core.ContractValue{
		{
			"newImplementation",
			"ByStr20",
			sdk.Cfg.ZimplAddr,
		},
	}
	CheckTx(Zproxy.Call("UpgradeTo", args, "0"))
	for _, ssnaddr := range sdk.Cfg.SsnAddrs {
		CheckTx(Zproxy.AddSSN(ssnaddr, ssnaddr))
	}
	CheckTx(Zproxy.UpdateVerifierRewardAddr(sdk.Cfg.Verifier))
	CheckTx(Zproxy.UpdateVerifier(sdk.Cfg.Verifier))
	CheckTx(Zproxy.UpdateStakingParameters(ToZil(sdk.Cfg.SsnInitialDelegateZil), ToZil(10))) //minstake (ssn not active if less), mindelegstake
	CheckTx(Zproxy.Unpause())

	//we need our SSN to be active, so delegating some stake to each
	for _, ssnaddr := range sdk.Cfg.SsnAddrs {
		CheckTx(Zproxy.DelegateStake(ssnaddr, ToZil(sdk.Cfg.SsnInitialDelegateZil)))
	}

	// SSN will become active on next cycle
	//we need to increase blocknum, in order to Gzil won't mint anything. Really minting is over.
	sdk.IncreaseBlocknum(2)
	Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)
	CheckTx(Zproxy.AssignStakeReward(sdk.Cfg.StZilSsnAddress, sdk.Cfg.StZilSsnRewardShare))
}

func Deploy(sdk *AvelySDK, log *Log) *Protocol {
	log.Debug("start to deploy")

	zilliqa := DeployZilliqaStaking(sdk, log)

	// deploy stzil
	StZIL, err := NewStZILContract(sdk, sdk.Cfg.Owner, zilliqa.Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy StZIL error = " + err.Error())
	}
	log.Debug("deploy StZIL succeed, address = " + StZIL.Addr)

	// deploy buffer
	Buffer, err := NewBufferContract(sdk, StZIL.Addr, zilliqa.Zproxy.Addr)
	if err != nil {
		log.Fatal("deploy buffer error = " + err.Error())
	}
	log.Debug("deploy buffer succeed, address = " + Buffer.Addr)
	buffers := []*BufferContract{Buffer}

	// deploy holder
	Holder, err := NewHolderContract(sdk, StZIL.Addr, zilliqa.Zproxy.Addr)
	if err != nil {
		log.Fatal("deploy holder error = " + err.Error())
	}
	log.Debug("deploy holder succeed, address = " + Holder.Addr)

	return NewProtocol(zilliqa.Zproxy, zilliqa.Zimpl, StZIL, buffers, Holder)
}

// Restore ZProxy + Zimpl and deploy new versions of StZIL, Buffer and Holder
func DeployOnlyAvely(sdk *AvelySDK, log *Log) *Protocol {
	log.Debug("start to DeployOnlyAvely")

	//Restore Zproxy
	Zproxy, err := RestoreZproxy(sdk, sdk.Cfg.ZproxyAddr)
	if err != nil {
		log.Fatal("Restore Zproxy error = " + err.Error())
	}
	log.Debug("Restore Zproxy succeed, address = " + Zproxy.Addr)

	//Restore Zimpl
	Zimpl, err := RestoreZimpl(sdk, sdk.Cfg.ZimplAddr, sdk.Cfg.ZproxyAddr, sdk.Cfg.GzilAddr)
	if err != nil {
		log.Fatal("Restore Zimpl error = " + err.Error())
	}
	log.Debug("Restore Zimpl succeed, address = " + Zimpl.Addr)

	// deploy stzil
	StZIL, err := NewStZILContract(sdk, sdk.Cfg.Owner, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy StZIL error = " + err.Error())
	}
	log.Debug("deploy StZIL succeed, address = " + StZIL.Addr)

	// deploy buffer
	Buffer, err := NewBufferContract(sdk, StZIL.Addr, Zproxy.Addr)
	if err != nil {
		log.Fatal("deploy buffer error = " + err.Error())
	}
	log.Debug("deploy buffer succeed, address = " + Buffer.Addr)
	buffers := []*BufferContract{Buffer}

	// deploy holder
	Holder, err := NewHolderContract(sdk, StZIL.Addr, Zproxy.Addr)
	if err != nil {
		log.Fatal("deploy holder error = " + err.Error())
	}
	log.Debug("deploy holder succeed, address = " + Holder.Addr)

	return NewProtocol(Zproxy, Zimpl, StZIL, buffers, Holder)
}

func RestoreFromState(sdk *AvelySDK, log *Log) *Protocol {
	log.Debug("start to Restoreialize from state")

	//Restore Zproxy
	Zproxy, err := RestoreZproxy(sdk, sdk.Cfg.ZproxyAddr)
	if err != nil {
		log.Fatal("Restore Zproxy error = " + err.Error())
	}
	log.Debug("Restore Zproxy succeed, address = " + Zproxy.Addr)

	//Restore Zimpl
	Zimpl, err := RestoreZimpl(sdk, sdk.Cfg.ZimplAddr, sdk.Cfg.ZproxyAddr, sdk.Cfg.GzilAddr)
	if err != nil {
		log.Fatal("Restore Zimpl error = " + err.Error())
	}
	log.Debug("Restore Zimpl succeed, address = " + Zimpl.Addr)

	// Restore stzil
	StZIL, err := RestoreStZILContract(sdk, sdk.Cfg.StZilAddr, sdk.GetAddressFromPrivateKey(sdk.Cfg.OwnerKey), Zimpl.Addr)
	if err != nil {
		log.Fatal("Restore StZIL error = " + err.Error())
	}
	log.Debug("Restore StZIL succeed, address = " + StZIL.Addr)

	// Restore buffers
	buffers := []*BufferContract{}
	for _, addr := range sdk.Cfg.BufferAddrs {
		Buffer, err := RestoreBufferContract(sdk, addr, StZIL.Addr, Zproxy.Addr)
		if err != nil {
			log.Fatal("Restore buffer error = " + err.Error())
		}
		log.Debug("Restore buffer succeed, address = " + Buffer.Addr)

		buffers = append(buffers, Buffer)
	}

	// Restore holder
	Holder, err := RestoreHolderContract(sdk, sdk.Cfg.HolderAddr, StZIL.Addr, Zproxy.Addr)
	if err != nil {
		log.Fatal("Restore holder error = " + err.Error())
	}
	log.Debug("Restore holder succeed, address = " + Holder.Addr)

	return NewProtocol(Zproxy, Zimpl, StZIL, buffers, Holder)
}
