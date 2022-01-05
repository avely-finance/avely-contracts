package main

import (
	"flag"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

var log *Log
var sdk *AvelySDK

func main() {
	chainPtr := flag.String("chain", "local", "chain")
	cmdPtr := flag.String("cmd", "default", "specific command")
	addrPtr := flag.String("addr", "default", "an entity address")

	flag.Parse()

	cmd := *cmdPtr

	log = NewLog()
	config := NewConfig(*chainPtr)
	sdk = NewAvelySDK(*config)

	shortcuts := map[string]string{
		"azilssn":  config.AzilSsnAddress,
		"addr1":    config.Addr1,
		"addr2":    config.Addr2,
		"addr3":    config.Addr3,
		"admin":    config.Admin,
		"verifier": config.Verifier,
	}
	log.AddShortcuts(shortcuts)

	if cmd == "deploy" {
		deployAvely()
	} else {
		// for non-deploy commands we need initialize protocol from config
		p := RestoreFromState(sdk, log)
		addr := *addrPtr

		switch cmd {
		case "from_bech32":
			convertFromBech32Addr(addr)
		case "to_bech32":
			convertToBech32Addr(addr)
		case "show_tx":
			showTx(p, addr)
		case "deploy_buffer":
			deployBuffer(p)
		case "sync_buffers":
			syncBuffers(p)
		case "drain_buffer":
			drainBuffer(p, addr)
		default:
			log.Fatal("Unknown command")
		}
	}

	log.Success("Done")
}

func deployAvely() {
	p := DeployOnlyAvely(sdk, log)
	p.SyncBufferAndHolder()
}

func showTx(p *Protocol, tx_addr string) {
	provider := p.Aimpl.Contract.Provider
	tx, err := provider.GetTransaction(tx_addr)

	log.Successf("Tx: ", tx)
	log.Successf("Err: ", err)
}

func convertFromBech32Addr(addr32 string) {
	addr, err := bech32.FromBech32Addr(addr32)

	if err != nil {
		log.Fatalf("Convert failed with err: ", err)
	}

	log.Success("Converted address: " + addr)
}

func convertToBech32Addr(addr32 string) {
	addr, err := bech32.ToBech32Address(addr32)

	if err != nil {
		log.Fatalf("Convert failed with err: ", err)
	}

	log.Success("Converted address: " + addr)
}

func deployBuffer(p *Protocol) {
	buffer, err := p.DeployBuffer()

	if err != nil {
		log.Fatalf("Buffer deploy failed with error: ", err)
	}
	log.Success("Buffer deploy is successfully compelted. Address: " + buffer.Addr)
}

func syncBuffers(p *Protocol) {
	p.SyncBufferAndHolder()
}

func drainBuffer(p *Protocol, buffer_addr string) {
	tx, err := p.Aimpl.DrainBuffer(buffer_addr)

	if err != nil {
		log.Fatalf("Drain failed with error: ", err)
	}
	log.Success("Drain is successfully compelted. Tx: " + tx.ID)
}
