package main

import (
	"flag"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

var log *Log
var sdk *AvelySDK

func main() {
	log = NewLog()
	config := NewConfig("local")
	sdk = NewAvelySDK(*config)

	shortcuts := map[string]string{
		"azilssn":  config.AzilSsnAddress,
		"addr1":    "0x" + config.Addr1,
		"addr2":    "0x" + config.Addr2,
		"addr3":    "0x" + config.Addr3,
		"admin":    "0x" + config.Admin,
		"verifier": "0x" + config.Verifier,
	}
	log.AddShortcuts(shortcuts)

	cmd := flag.String("cmd", "default", "specific command")
	addrPtr := flag.String("addr", "default", "an entity address")

	flag.Parse()

	if *cmd == "deploy" {
		deployAvely()
	} else {
		// for non-deploy commands we need initialize protocol from config
		p := RestoreFromState(sdk, log)
		addr := *addrPtr

		switch *cmd {
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

func deployBuffer(p *Protocol) {
	buffer, err := p.DeployBuffer()

	if err != nil {
		log.Fatalf("Buffer deploy failed with error: ", err)
	} else {
		log.Success("Buffer deploy is successfully compelted. Address: " + buffer.Addr)
	}
}

func syncBuffers(p *Protocol) {
	p.SyncBufferAndHolder()
}

func drainBuffer(p *Protocol, buffer_addr string) {
	tx, err := p.Aimpl.DrainBuffer(buffer_addr)

	if err != nil {
		log.Fatalf("Drain failed with error: ", err)
	} else {
		log.Success("Drain is successfully compelted. Tx: " + tx.ID)
	}
}
