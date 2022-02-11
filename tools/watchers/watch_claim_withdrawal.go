package main

import (
	"flag"

	"github.com/avely-finance/avely-contracts/sdk/actions"
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	"github.com/sirupsen/logrus"
)

type ClaimWithdrawalWatcher struct {
	gap        int
	runAtBlock int
}

var log *core.Log
var sdk *core.AvelySDK
var protocol *contracts.Protocol

func main() {

	chainPtr := flag.String("chain", "local", "chain")
	gapPtr := flag.Int("gap", 5, "gap between blocks")

	flag.Parse()

	config := core.NewConfig(*chainPtr)
	sdk = core.NewAvelySDK(*config)
	log = core.NewLog()
	log.SetOutputStdout()
	log.AddSlackHook(sdk.Cfg.Slack.HookUrl, sdk.Cfg.Slack.LogLevel)
	protocol = contracts.RestoreFromState(sdk, log)
	url := sdk.GetWsURL()

	claimWatcher := &ClaimWithdrawalWatcher{
		gap:        *gapPtr,
		runAtBlock: -1,
	}

	log.Debug("Start claim withdrawal watcher")
	blockWatcher := utils.CreateBlockWatcher(url)
	blockWatcher.AddObserver(claimWatcher)
	blockWatcher.Start()
}

func (cww *ClaimWithdrawalWatcher) Notify(blockNum int) {
	if (blockNum - cww.runAtBlock) > cww.gap {
		log.WithFields(logrus.Fields{"block_number": blockNum}).Debug("Mined block")
		actions.ClaimWithdrawal(protocol)
		cww.runAtBlock = blockNum
	} else {
		log.WithFields(logrus.Fields{
			"block_number": blockNum,
			"current_gap":  (blockNum - cww.runAtBlock),
			"expected_gap": cww.gap,
		}).Debug("Block mined, skip")
	}
}
