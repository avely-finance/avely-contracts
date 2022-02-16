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
	protocol   *contracts.Protocol
	log        *core.Log
}

func main() {
	chainPtr := flag.String("chain", "local", "chain")
	gapPtr := flag.Int("gap", 5, "gap between blocks")

	flag.Parse()

	config := core.NewConfig(*chainPtr)
	sdk := core.NewAvelySDK(*config)
	log := core.NewLog()
	log.SetOutputStdout()
	log.AddSlackHook(sdk.Cfg.Slack.HookUrl, sdk.Cfg.Slack.LogLevel)
	protocol := contracts.RestoreFromState(sdk, log)
	url := sdk.GetWsURL()

	claimWatcher := &ClaimWithdrawalWatcher{
		gap:        *gapPtr,
		runAtBlock: -1,
		protocol:   protocol,
		log:        log,
	}

	log.Debug("Start claim withdrawal watcher")
	blockWatcher := utils.CreateBlockWatcher(url, log)
	blockWatcher.AddObserver(claimWatcher)
	blockWatcher.Start()
}

func (cww *ClaimWithdrawalWatcher) Notify(blockNum int) {
	protocol := cww.protocol
	action := actions.NewAdminActions(cww.log)

	if (blockNum - cww.runAtBlock) > cww.gap {
		cww.log.WithFields(logrus.Fields{"block_number": blockNum}).Debug("Mined block")
		action.ClaimWithdrawal(protocol)
		cww.runAtBlock = blockNum
	} else {
		cww.log.WithFields(logrus.Fields{
			"block_number": blockNum,
			"current_gap":  (blockNum - cww.runAtBlock),
			"expected_gap": cww.gap,
		}).Debug("Block mined, skip")
	}
}
