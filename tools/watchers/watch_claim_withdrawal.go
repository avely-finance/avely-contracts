package main

import (
	"flag"

	"github.com/avely-finance/avely-contracts/sdk/actions"
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
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

	log.Info("Start claim withdrawal watcher")
	blockWatcher := utils.CreateBlockWatcher(url)
	blockWatcher.AddObserver(claimWatcher)
	blockWatcher.Start()
}

func (cww *ClaimWithdrawalWatcher) Notify(blockNum int) {
	if (blockNum - cww.runAtBlock) > cww.gap {
		log.Debugf("Mined block #%d.", blockNum)
		actions.ClaimWithdrawal(protocol)
		cww.runAtBlock = blockNum
	} else {
		log.Debugf("Mined block #%d, but gap=%d <= %d, skip.", blockNum, (blockNum - cww.runAtBlock), cww.gap)
	}
}
