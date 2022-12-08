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

func DeployZilliqaStaking(sdk *AvelySDK, celestials *Celestials, log *Log) *ZilliqaStaking {
	log.Debug("start to deploy zilliqa staking contracts")

	//deploy gzil
	gzil, err := NewGzil(sdk, celestials.Admin)
	if err != nil {
		log.Fatal("deploy Gzil error = " + err.Error())
	}
	log.Debug("deploy Gzil succeed, address = " + gzil.Addr)

	//deploy Zproxy
	zproxy, err := NewZproxy(sdk, celestials.Admin)
	if err != nil {
		log.Fatal("deploy Zproxy error = " + err.Error())
	}
	log.Debug("deploy Zproxy succeed, address = " + zproxy.Addr)

	//deploy Zimpl
	zimpl, err := NewZimpl(sdk, zproxy.Addr, gzil.Addr, celestials.Admin)
	if err != nil {
		log.Fatal("deploy Zimpl error = " + err.Error())
	}
	log.Debug("deploy Zimpl succeed, address = " + zimpl.Addr)

	return NewZilliqaStaking(zproxy, zimpl, gzil)
}

func SetupZilliqaStaking(sdk *AvelySDK, celestials *Celestials, log *Log) {

	//Restore Zproxy
	Zproxy, err := RestoreZproxy(sdk, sdk.Cfg.ZproxyAddr)
	if err != nil {
		log.Fatal("Restore Zproxy error = " + err.Error())
	}
	log.Debug("Restore Zproxy succeed, address = " + Zproxy.Addr)
	Zproxy.WithWallet(celestials.Admin)

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
