package main

import (
	"flag"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	"reflect"
)

var log *Log
var sdk *AvelySDK

func main() {
	chainPtr := flag.String("chain", "local", "chain")
	cmd := flag.String("cmd", "default", "specific command")
	usrPtr := flag.String("usr", "default", "an user ID")
	amountPtr := flag.Int("amount", 0, "an amount of action")

	flag.Parse()

	chain := *chainPtr
	amount := *amountPtr
	usr := *usrPtr

	log = NewLog()
	config := NewConfig(chain)
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

	p := RestoreFromState(sdk, log)

	setupUsr(p, usr)

	switch *cmd {
	case "delegate":
		delegate(p, amount)
	default:
		log.Fatal("Unknown command")
	}

	log.Success("Done")
}

func setupUsr(p *Protocol, usr string) {
	if usr == "default" {
		log.Fatal("Undefined user")
	}

	pr := reflect.ValueOf(sdk.Cfg)
	keyValue := reflect.Indirect(pr).FieldByName("Key" + usr)

	key := keyValue.Interface().(string)

	p.Aproxy.UpdateWallet(key)

	log.Success("Wallet has been updates to Key" + usr)
}

func delegate(p *Protocol, amount int) {
	if amount <= 0 {
		log.Fatal("Amount should be greater than 0")
	}

	tx, err := p.Aproxy.DelegateStake(utils.ToZil(amount))

	if err != nil {
		log.Fatalf("Delegate failed with error:", tx)
	} else {
		log.Success("Delegate is successfully compelted. Tx: " + tx.ID)
	}
}
