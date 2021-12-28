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

	cmd := flag.String("cmd", "help", "specific command")

	flag.Parse()

	switch *cmd {
	case "deploy":
		deployAvely()
	default:
		log.Fatal("Unknown command")
	}

	log.Success("Done")
}

func deployAvely() {
	p := DeployOnlyAvely(sdk, log)
	p.SyncBufferAndHolder()
}
