package contracts

import (
	"log"
	"runtime"
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

func Deploy(sdk *AvelySDK, log *Log) *Protocol {
	log.Debug("start to deploy")

	//deploy gzil
	gzil, err := NewGzil(sdk)
	if err != nil {
		log.Fatal("deploy Gzil error = " + err.Error())
	}
	log.Debug("deploy Gzil succeed, address = " + gzil.Addr)

	//deploy Zproxy
	Zproxy, err := NewZproxy(sdk)
	if err != nil {
		log.Fatal("deploy Zproxy error = " + err.Error())
	}
	log.Debug("deploy Zproxy succeed, address = " + Zproxy.Addr)

	//deploy Zimpl
	Zimpl, err := NewZimpl(sdk, Zproxy.Addr, gzil.Addr)
	if err != nil {
		log.Fatal("deploy Zimpl error = " + err.Error())
	}
	log.Debug("deploy Zimpl succeed, address = " + Zimpl.Addr)

	// deploy azil
	Azil, err := NewAZilContract(sdk, sdk.Cfg.Owner, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy aZil error = " + err.Error())
	}
	log.Debug("deploy aZil succeed, address = " + Azil.Addr)

	// deploy buffer
	Buffer, err := NewBufferContract(sdk, Azil.Addr, Zproxy.Addr)
	if err != nil {
		log.Fatal("deploy buffer error = " + err.Error())
	}
	log.Debug("deploy buffer succeed, address = " + Buffer.Addr)
	buffers := []*BufferContract{Buffer}

	// deploy holder
	Holder, err := NewHolderContract(sdk, Azil.Addr, Zproxy.Addr)
	if err != nil {
		log.Fatal("deploy holder error = " + err.Error())
	}
	log.Debug("deploy holder succeed, address = " + Holder.Addr)

	return NewProtocol(Zproxy, Zimpl, Azil, buffers, Holder)
}

// Restore ZProxy + Zimpl and deploy new versions of Azil, Buffer and Holder
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

	// deploy azil
	Azil, err := NewAZilContract(sdk, sdk.Cfg.Owner, Zimpl.Addr)
	if err != nil {
		log.Fatal("deploy aZil error = " + err.Error())
	}
	log.Debug("deploy aZil succeed, address = " + Azil.Addr)

	// deploy buffer
	Buffer, err := NewBufferContract(sdk, Azil.Addr, Zproxy.Addr)
	if err != nil {
		log.Fatal("deploy buffer error = " + err.Error())
	}
	log.Debug("deploy buffer succeed, address = " + Buffer.Addr)
	buffers := []*BufferContract{Buffer}

	// deploy holder
	Holder, err := NewHolderContract(sdk, Azil.Addr, Zproxy.Addr)
	if err != nil {
		log.Fatal("deploy holder error = " + err.Error())
	}
	log.Debug("deploy holder succeed, address = " + Holder.Addr)

	return NewProtocol(Zproxy, Zimpl, Azil, buffers, Holder)
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

	// Restore azil
	Azil, err := RestoreAZilContract(sdk, sdk.Cfg.AzilAddr, sdk.GetAddressFromPrivateKey(sdk.Cfg.OwnerKey), Zimpl.Addr)
	if err != nil {
		log.Fatal("Restore aZil error = " + err.Error())
	}
	log.Debug("Restore aZil succeed, address = " + Azil.Addr)

	// Restore buffers
	buffers := []*BufferContract{}
	for _, addr := range sdk.Cfg.BufferAddrs {
		Buffer, err := RestoreBufferContract(sdk, addr, Azil.Addr, Zproxy.Addr)
		if err != nil {
			log.Fatal("Restore buffer error = " + err.Error())
		}
		log.Debug("Restore buffer succeed, address = " + Buffer.Addr)

		buffers = append(buffers, Buffer)
	}

	// Restore holder
	Holder, err := RestoreHolderContract(sdk, sdk.Cfg.HolderAddr, Azil.Addr, Zproxy.Addr)
	if err != nil {
		log.Fatal("Restore holder error = " + err.Error())
	}
	log.Debug("Restore holder succeed, address = " + Holder.Addr)

	return NewProtocol(Zproxy, Zimpl, Azil, buffers, Holder)
}

//TODO: move this function to core/sdk.go, rename to CheckTx
func check(tx *transaction.Transaction, err error) (*transaction.Transaction, error) {
	if err != nil {
		_, file, no, _ := runtime.Caller(1)
		log.Fatal("TRANSACTION FAILED, " + file + ":" + strconv.Itoa(no))
	}
	return tx, err
}
