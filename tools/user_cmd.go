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
		"stzilssn": config.StZilSsnAddress,
		"addr1":    config.Addr1,
		"addr2":    config.Addr2,
		"addr3":    config.Addr3,
		"admin":    config.Admin,
		"verifier": config.Verifier,
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

	log.Info("Done")
}

func setupUsr(p *Protocol, usr string) {
	if usr == "default" {
		log.Fatal("Undefined user")
	}

	pr := reflect.ValueOf(sdk.Cfg)
	keyValue := reflect.Indirect(pr).FieldByName("Key" + usr)

	key := keyValue.Interface().(string)

	p.StZIL.UpdateWallet(key)

	log.Info("Wallet has been updates to Key" + usr)
}

func delegate(p *Protocol, amount int) {
	if amount <= 0 {
		log.Fatal("Amount should be greater than 0")
	}

	tx, err := p.StZIL.DelegateStake(utils.ToZil(amount))

	if err != nil {
		log.Fatal("Delegate failed with error:" + err.Error())
	} else {
		log.Info("Delegate is successfully compelted. Tx: " + tx.ID)
	}
}