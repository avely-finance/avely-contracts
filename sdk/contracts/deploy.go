package contracts

import (
	"log"
	"runtime"
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

func Deploy(sdk *AvelySDK, log *Log) *Protocol {
	log.Info("start to deploy")

	//deploy gzil
	gzil, err := NewGzil(sdk)
	if err != nil {
		log.Fatal("deploy Gzil error = " + err.Error())
	}
	log.Info("deploy Gzil succeed, address = " + gzil.Addr)

	//deploy Zproxy
	Zproxy, err := NewZproxy(sdk)
	if err != nil {
		log.Fatal("deploy Zproxy error = " + err.Error())
	}
	log.Info("deploy Zproxy succeed, address = " + Zproxy.Addr)

	//deploy Zimpl
	Zimpl, err := NewZimpl(sdk, Zproxy.Addr, gzil.Addr)
	if err != nil {
		log.Fatal("deploy Zimpl error = " + err.Error())
	}
	log.Info("deploy Zimpl succeed, address = " + Zimpl.Addr)

	// deploy azil
	Aimpl, err := NewAZilContract(sdk, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy aZil error = " + err.Error())
	}
	log.Info("deploy aZil succeed, address = " + Aimpl.Addr)

	// deploy buffer
	Buffer, err := NewBufferContract(sdk, Aimpl.Addr, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy buffer error = " + err.Error())
	}
	log.Info("deploy buffer succeed, address = " + Buffer.Addr)
	buffers := []*BufferContract{Buffer}

	// deploy holder
	Holder, err := NewHolderContract(sdk, Aimpl.Addr, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy holder error = " + err.Error())
	}
	log.Info("deploy holder succeed, address = " + Holder.Addr)

	return NewProtocol(Zproxy, Zimpl, Aimpl, buffers, Holder)
}

// Restore ZProxy + Zimpl and deploy new versions of Azil, Buffer and Holder
func DeployOnlyAvely(sdk *AvelySDK, log *Log) *Protocol {
	log.Info("start to DeployOnlyAvely")

	//Restore Zproxy
	Zproxy, err := RestoreZproxy(sdk, sdk.Cfg.ZproxyAddr)
	if err != nil {
		log.Fatal("Restore Zproxy error = " + err.Error())
	}
	log.Info("Restore Zproxy succeed, address = " + Zproxy.Addr)

	//Restore Zimpl
	Zimpl, err := RestoreZimpl(sdk, sdk.Cfg.ZimplAddr, sdk.Cfg.ZproxyAddr, sdk.Cfg.GzilAddr)
	if err != nil {
		log.Fatal("Restore Zimpl error = " + err.Error())
	}
	log.Info("Restore Zimpl succeed, address = " + Zimpl.Addr)

	// deploy azil
	Aimpl, err := NewAZilContract(sdk, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy aZil error = " + err.Error())
	}
	log.Info("deploy aZil succeed, address = " + Aimpl.Addr)

	// deploy buffer
	Buffer, err := NewBufferContract(sdk, Aimpl.Addr, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy buffer error = " + err.Error())
	}
	log.Info("deploy buffer succeed, address = " + Buffer.Addr)
	buffers := []*BufferContract{Buffer}

	// deploy holder
	Holder, err := NewHolderContract(sdk, Aimpl.Addr, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy holder error = " + err.Error())
	}
	log.Info("deploy holder succeed, address = " + Holder.Addr)

	return NewProtocol(Zproxy, Zimpl, Aimpl, buffers, Holder)
}

func RestoreFromState(sdk *AvelySDK, log *Log) *Protocol {
	log.Info("start to Restoreialize from state")

	//Restore Zproxy
	Zproxy, err := RestoreZproxy(sdk, sdk.Cfg.ZproxyAddr)
	if err != nil {
		log.Fatal("Restore Zproxy error = " + err.Error())
	}
	log.Info("Restore Zproxy succeed, address = " + Zproxy.Addr)

	//Restore Zimpl
	Zimpl, err := RestoreZimpl(sdk, sdk.Cfg.ZimplAddr, sdk.Cfg.ZproxyAddr, sdk.Cfg.GzilAddr)
	if err != nil {
		log.Fatal("Restore Zimpl error = " + err.Error())
	}
	log.Info("Restore Zimpl succeed, address = " + Zimpl.Addr)

	// Restore azil
	Aimpl, err := RestoreAZilContract(sdk, sdk.Cfg.AzilAddr, Zimpl.Addr)
	if err != nil {
		log.Fatal("Restore aZil error = " + err.Error())
	}
	log.Info("Restore aZil succeed, address = " + Aimpl.Addr)

	// Restore buffers
	buffers := []*BufferContract{}
	for _, addr := range sdk.Cfg.BufferAddrs {
		Buffer, err := RestoreBufferContract(sdk, addr, Aimpl.Addr, Zproxy.Addr, Zimpl.Addr)
		if err != nil {
			log.Fatal("Restore buffer error = " + err.Error())
		}
		log.Info("Restore buffer succeed, address = " + Buffer.Addr)

		buffers = append(buffers, Buffer)
	}

	// Restore holder
	Holder, err := RestoreHolderContract(sdk, sdk.Cfg.HolderAddr, Aimpl.Addr, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		log.Fatal("Restore holder error = " + err.Error())
	}
	log.Info("Restore holder succeed, address = " + Holder.Addr)

	return NewProtocol(Zproxy, Zimpl, Aimpl, buffers, Holder)
}

//TODO: move this function to core/sdk.go, rename to CheckTx
func check(tx *transaction.Transaction, err error) (*transaction.Transaction, error) {
	if err != nil {
		_, file, no, _ := runtime.Caller(1)
		log.Fatal("TRANSACTION FAILED, " + file + ":" + strconv.Itoa(no))
	}
	return tx, err
}
