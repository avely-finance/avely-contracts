package main

import (
	"flag"
	"reflect"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
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
	usrPtr := flag.String("usr", "default", "an user ID")
	amountPtr := flag.Int("amount", 0, "an amount of action")

	flag.Parse()
	p := RestoreFromState(sdk, log)

	amount := *amountPtr
	usr := *usrPtr

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

	p.Aimpl.UpdateWallet(key)

	log.Success("Wallet has been updates to Key" + usr)
}

func delegate(p *Protocol, amount int) {
	if amount <= 0 {
		log.Fatal("Amount should be greater than 0")
	}

	tx, err := p.Aimpl.DelegateStake(utils.ToZil(amount))

	if err != nil {
		log.Fatalf("Delegate failed with error:", tx)
	} else {
		log.Success("Delegate is successfully compelted. Tx: " + tx.ID)
	}
}
