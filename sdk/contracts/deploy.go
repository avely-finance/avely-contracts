package contracts

import (
	"runtime"
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

func Deploy(sdk *AvelySDK, lgr *Log) *Protocol {
	lgr.Info("start to deploy")

	//deploy gzil
	gzil, err := NewGzil(sdk)
	if err != nil {
		lgr.Fatal("deploy Gzil error = " + err.Error())
	}
	lgr.Info("deploy Gzil succeed, address = " + gzil.Addr)

	//deploy Zproxy
	Zproxy, err := NewZproxy(sdk)
	if err != nil {
		lgr.Fatal("deploy Zproxy error = " + err.Error())
	}
	lgr.Info("deploy Zproxy succeed, address = " + Zproxy.Addr)

	//deploy Zimpl
	Zimpl, err := NewZimpl(sdk, Zproxy.Addr, gzil.Addr)
	if err != nil {
		lgr.Fatal("deploy Zimpl error = " + err.Error())
	}
	lgr.Info("deploy Zimpl succeed, address = " + Zimpl.Addr)

	// deploy azil
	Aimpl, err := NewAZilContract(sdk, Zimpl.Addr)
	if err != nil {
		lgr.Fatal("deploy aZil error = " + err.Error())
	}
	lgr.Info("deploy aZil succeed, address = " + Aimpl.Addr)

	// deploy buffer
	Buffer, err := NewBufferContract(sdk, Aimpl.Addr, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		lgr.Fatal("deploy buffer error = " + err.Error())
	}
	lgr.Info("deploy buffer succeed, address = " + Buffer.Addr)
	buffers := []*BufferContract{Buffer}

	// deploy holder
	Holder, err := NewHolderContract(sdk, Aimpl.Addr, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		lgr.Fatal("deploy holder error = " + err.Error())
	}
	lgr.Info("deploy holder succeed, address = " + Holder.Addr)

	return NewProtocol(Zproxy, Zimpl, Aimpl, buffers, Holder)
}

// Restore ZProxy + Zimpl and deploy new versions of Azil, Buffer and Holder
func DeployOnlyAvely(sdk *AvelySDK, lgr *Log) *Protocol {
	lgr.Info("start to DeployOnlyAvely")

	//Restore Zproxy
	Zproxy, err := RestoreZproxy(sdk, sdk.Cfg.ZproxyAddr)
	if err != nil {
		lgr.Fatal("Restore Zproxy error = " + err.Error())
	}
	lgr.Info("Restore Zproxy succeed, address = " + Zproxy.Addr)

	//Restore Zimpl
	Zimpl, err := RestoreZimpl(sdk, sdk.Cfg.ZimplAddr, sdk.Cfg.ZproxyAddr, sdk.Cfg.GzilAddr)
	if err != nil {
		lgr.Fatal("Restore Zimpl error = " + err.Error())
	}
	lgr.Info("Restore Zimpl succeed, address = " + Zimpl.Addr)

	// deploy azil
	Aimpl, err := NewAZilContract(sdk, Zimpl.Addr)
	if err != nil {
		lgr.Fatal("deploy aZil error = " + err.Error())
	}
	lgr.Info("deploy aZil succeed, address = " + Aimpl.Addr)

	// deploy buffer
	Buffer, err := NewBufferContract(sdk, Aimpl.Addr, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		lgr.Fatal("deploy buffer error = " + err.Error())
	}
	lgr.Info("deploy buffer succeed, address = " + Buffer.Addr)
	buffers := []*BufferContract{Buffer}

	// deploy holder
	Holder, err := NewHolderContract(sdk, Aimpl.Addr, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		lgr.Fatal("deploy holder error = " + err.Error())
	}
	lgr.Info("deploy holder succeed, address = " + Holder.Addr)

	return NewProtocol(Zproxy, Zimpl, Aimpl, buffers, Holder)
}

func RestoreFromState(sdk *AvelySDK, lgr *Log) *Protocol {
	lgr.Info("start to Restoreialize from state")

	//Restore Zproxy
	Zproxy, err := RestoreZproxy(sdk, sdk.Cfg.ZproxyAddr)
	if err != nil {
		lgr.Fatal("Restore Zproxy error = " + err.Error())
	}
	lgr.Info("Restore Zproxy succeed, address = " + Zproxy.Addr)

	//Restore Zimpl
	Zimpl, err := RestoreZimpl(sdk, sdk.Cfg.ZimplAddr, sdk.Cfg.ZproxyAddr, sdk.Cfg.GzilAddr)
	if err != nil {
		lgr.Fatal("Restore Zimpl error = " + err.Error())
	}
	lgr.Info("Restore Zimpl succeed, address = " + Zimpl.Addr)

	// Restore azil
	Aimpl, err := RestoreAZilContract(sdk, sdk.Cfg.AzilAddr, Zimpl.Addr)
	if err != nil {
		lgr.Fatal("Restore aZil error = " + err.Error())
	}
	lgr.Info("Restore aZil succeed, address = " + Aimpl.Addr)

	// Restore buffers
	buffers := []*BufferContract{}
	for _, addr := range sdk.Cfg.BufferAddrs {
		Buffer, err := RestoreBufferContract(sdk, addr, Aimpl.Addr, Zproxy.Addr, Zimpl.Addr)
		if err != nil {
			lgr.Fatal("Restore buffer error = " + err.Error())
		}
		lgr.Info("Restore buffer succeed, address = " + Buffer.Addr)

		buffers = append(buffers, Buffer)
	}

	// Restore holder
	Holder, err := RestoreHolderContract(sdk, sdk.Cfg.HolderAddr, Aimpl.Addr, Zproxy.Addr, Zimpl.Addr)
	if err != nil {
		lgr.Fatal("Restore holder error = " + err.Error())
	}
	lgr.Info("Restore holder succeed, address = " + Holder.Addr)

	return NewProtocol(Zproxy, Zimpl, Aimpl, buffers, Holder)
}

//TODO: move this function to core/sdk.go, rename to CheckTx
func check(tx *transaction.Transaction, err error) (*transaction.Transaction, error) {
	if err != nil {
		_, file, no, _ := runtime.Caller(1)
		lgr.Fatal("TRANSACTION FAILED, " + file + ":" + strconv.Itoa(no))
	}
	return tx, err
}
