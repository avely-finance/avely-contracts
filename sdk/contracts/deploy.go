package contracts

import (
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

func Deploy(sdk *AvelySDK, celestials *Celestials, log *Log) *Protocol {
	log.Debug("start to deploy")

	zilliqa := DeployZilliqaStaking(sdk, celestials, log)

	// deploy stzil
	StZIL, err := NewStZILContract(sdk, "0x"+celestials.Owner.DefaultAccount.Address, zilliqa.Zimpl.Addr, celestials.Admin)
	if err != nil {
		log.Fatal("deploy StZIL error = " + err.Error())
	}
	log.Debug("deploy StZIL succeed, address = " + StZIL.Addr)

	// deploy buffer
	Buffer, err := NewBufferContract(sdk, StZIL.Addr, zilliqa.Zproxy.Addr, celestials.Admin)
	if err != nil {
		log.Fatal("deploy buffer error = " + err.Error())
	}
	log.Debug("deploy buffer succeed, address = " + Buffer.Addr)
	buffers := []*BufferContract{Buffer}

	// deploy holder
	Holder, err := NewHolderContract(sdk, StZIL.Addr, zilliqa.Zproxy.Addr, celestials.Admin)
	if err != nil {
		log.Fatal("deploy holder error = " + err.Error())
	}
	log.Debug("deploy holder succeed, address = " + Holder.Addr)

	// deploy treasury
	Treasury, err := NewTreasuryContract(sdk, sdk.Cfg.Owner, celestials.Admin)
	if err != nil {
		log.Fatal("deploy Treasury error = " + err.Error())
	}
	log.Debug("deploy Treasury succeed, address = " + Treasury.Addr)

	return NewProtocol(zilliqa.Zproxy, zilliqa.Zimpl, StZIL, buffers, Holder, Treasury)
}

// Restore ZProxy + Zimpl and deploy new versions of StZIL, Buffer and Holder
func DeployOnlyAvely(sdk *AvelySDK, celestials *Celestials, log *Log) *Protocol {
	log.Debug("start to DeployOnlyAvely")

	//Restore Zproxy
	Zproxy, err := RestoreZproxy(sdk, sdk.Cfg.ZproxyAddr)
	if err != nil {
		log.Fatal("Restore Zproxy error = " + err.Error())
	}
	log.Debug("Restore Zproxy succeed, address = " + Zproxy.Addr)

	//Restore Zimpl
	Zimpl, err := RestoreZimpl(sdk, sdk.Cfg.ZimplAddr)
	if err != nil {
		log.Fatal("Restore Zimpl error = " + err.Error())
	}
	log.Debug("Restore Zimpl succeed, address = " + Zimpl.Addr)

	// deploy stzil
	StZIL, err := NewStZILContract(sdk, "0x"+celestials.Owner.DefaultAccount.Address, Zimpl.Addr, celestials.Admin)

	if err != nil {
		log.Fatal("deploy StZIL error = " + err.Error())
	}
	log.Debug("deploy StZIL succeed, address = " + StZIL.Addr)

	// deploy buffer
	Buffer, err := NewBufferContract(sdk, StZIL.Addr, Zproxy.Addr, celestials.Admin)
	if err != nil {
		log.Fatal("deploy buffer error = " + err.Error())
	}
	log.Debug("deploy buffer succeed, address = " + Buffer.Addr)
	buffers := []*BufferContract{Buffer}

	// deploy holder
	Holder, err := NewHolderContract(sdk, StZIL.Addr, Zproxy.Addr, celestials.Admin)
	if err != nil {
		log.Fatal("deploy holder error = " + err.Error())
	}
	log.Debug("deploy holder succeed, address = " + Holder.Addr)

	// deploy treasury
	Treasury, err := NewTreasuryContract(sdk, sdk.Cfg.Owner, celestials.Admin)
	if err != nil {
		log.Fatal("deploy Treasury error = " + err.Error())
	}
	log.Debug("deploy Treasury succeed, address = " + StZIL.Addr)

	return NewProtocol(Zproxy, Zimpl, StZIL, buffers, Holder, Treasury)
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
	Zimpl, err := RestoreZimpl(sdk, sdk.Cfg.ZimplAddr)
	if err != nil {
		log.Fatal("Restore Zimpl error = " + err.Error())
	}
	log.Debug("Restore Zimpl succeed, address = " + Zimpl.Addr)

	// Restore stzil
	StZIL, err := RestoreStZILContract(sdk, sdk.Cfg.StZilAddr)
	if err != nil {
		log.Fatal("Restore StZIL error = " + err.Error())
	}
	log.Debug("Restore StZIL succeed, address = " + StZIL.Addr)

	// Restore buffers
	buffers := []*BufferContract{}
	for _, addr := range sdk.Cfg.BufferAddrs {
		Buffer, err := RestoreBufferContract(sdk, addr)
		if err != nil {
			log.Fatal("Restore buffer error = " + err.Error())
		}
		log.Debug("Restore buffer succeed, address = " + Buffer.Addr)

		buffers = append(buffers, Buffer)
	}

	// Restore holder
	Holder, err := RestoreHolderContract(sdk, sdk.Cfg.HolderAddr)
	if err != nil {
		log.Fatal("Restore holder error = " + err.Error())
	}
	log.Debug("Restore holder succeed, address = " + Holder.Addr)

	// Restore treasury
	Treasury, err := RestoreTreasuryContract(sdk, sdk.Cfg.TreasuryAddr)
	if err != nil {
		log.Fatal("Restore Treasury error = " + err.Error())
	}
	log.Debug("Restore Treasury succeed, address = " + Treasury.Addr)

	return NewProtocol(Zproxy, Zimpl, StZIL, buffers, Holder, Treasury)
}
