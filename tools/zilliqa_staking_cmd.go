package main

import (
	"flag"

	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/sirupsen/logrus"
)

var log *Log
var sdk *AvelySDK

func main() {
	chainPtr := flag.String("chain", "local", "chain")
	cmdPtr := flag.String("cmd", "default", "specific command")

	flag.Parse()

	cmd := *cmdPtr

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
		p := RestoreFromState(sdk, log)
		SetupZilliqaStaking(p)
	default:
		log.Fatal("Unknown command")
	}

	log.Info("Done")
}
