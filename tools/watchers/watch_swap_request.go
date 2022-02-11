package main

import (
	"flag"

	"github.com/avely-finance/avely-contracts/sdk/actions"
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	"github.com/sirupsen/logrus"
)

type SwapRequestWatcher struct {
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

	watcher := &SwapRequestWatcher{
		gap:        *gapPtr,
		runAtBlock: -1,
	}

	log.Debug("Start swap request watcher")
	blockWatcher := utils.CreateBlockWatcher(url)
	blockWatcher.AddObserver(watcher)
	blockWatcher.Start()
}

func (w *SwapRequestWatcher) Notify(blockNum int) {
	if (blockNum - w.runAtBlock) > w.gap {
		log.WithFields(logrus.Fields{"block_number": blockNum}).Debug("Mined block")
		actions.ConfirmSwapRequests(protocol)
		w.runAtBlock = blockNum
	} else {
		log.WithFields(logrus.Fields{
			"block_number": blockNum,
			"current_gap":  (blockNum - w.runAtBlock),
			"expected_gap": w.gap,
		}).Debug("Block mined, skip")
	}
}
