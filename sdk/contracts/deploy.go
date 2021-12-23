package contracts

import (
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

func Deploy(sdk *AvelySDK, log *Log) (*Protocol) {
	log.Info("start to deploy")

	//deploy gzil
	gzil, err := NewGzil(sdk)
	if err != nil {
		log.Fatal("deploy Gzil error = " + err.Error())
	}
	log.Success("deploy Gzil succeed, address = " + gzil.Addr)

	//deploy Zproxy
	Zproxy, err := NewZproxy(sdk)
	if err != nil {
		log.Fatal("deploy Zproxy error = " + err.Error())
	}
	log.Success("deploy Zproxy succeed, address = " + Zproxy.Addr)

	//deploy Zimpl
	Zimpl, err := NewZimpl(sdk, Zproxy.Addr, gzil.Addr)
	if err != nil {
		log.Fatal("deploy Zimpl error = " + err.Error())
	}
	log.Success("deploy Zimpl succeed, address = " + Zimpl.Addr)

	// deploy azil
	Aimpl, err := NewAZilContract(sdk, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy aZil error = " + err.Error())
	}
	log.Success("deploy aZil succeed, address = " + Aimpl.Addr)

	// deploy buffer
	Buffer, err := NewBufferContract(sdk, Aimpl.Addr, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy buffer error = " + err.Error())
	}
	log.Success("deploy buffer succeed, address = " + Buffer.Addr)

	// deploy holder
	Holder, err := NewHolderContract(sdk, Aimpl.Addr, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy holder error = " + err.Error())
	}
	log.Success("deploy holder succeed, address = " + Holder.Addr)

	return NewProtocol(Zproxy, Zimpl, Aimpl, Buffer, Holder)
}
