package main

import (
	"flag"

	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/tests/transitions"
	"github.com/sirupsen/logrus"
)

var log *Log
var sdk *AvelySDK

func main() {
	chainPtr := flag.String("chain", "local", "chain")
	cmdPtr := flag.String("cmd", "default", "specific command")
	addrPtr := flag.String("addr", "", "address")
	amountPtr := flag.Int("amount", 0, "an amount of action")

	flag.Parse()

	cmd := *cmdPtr
	addr := *addrPtr
	amount := *amountPtr

	log = NewLog()
	config := NewConfig(*chainPtr)
	sdk = NewAvelySDK(*config)

	switch cmd {
	case "deploy":
		zilliqa := DeployZilliqaStaking(sdk, log)
		log.WithFields(logrus.Fields{
			"gzil":    zilliqa.Gzil.Addr,
			"proxy":   zilliqa.Zproxy.Addr,
			"ssnlist": zilliqa.Zimpl.Addr,
		}).Info("zilliqa staking deployed")
	case "setup":
		SetupZilliqaStaking(sdk, log)
	case "next_cycle":
		p := RestoreFromState(sdk, log)
		tr := transitions.InitTransitions(sdk)
		tr.NextCycle(p)
	case "add_ssn":
		Zproxy, err := RestoreZproxy(sdk, sdk.Cfg.ZproxyAddr)
		if err != nil {
			log.Fatal("Restore Zproxy error = " + err.Error())
		}
		log.Debug("Restore Zproxy succeed, address = " + Zproxy.Addr)
		Zproxy.AddSSN(addr, addr)
	default:
		log.Fatal("Unknown command")
	}

	log.Info("Done")
}
