package main

import (
	"flag"

	"github.com/avely-finance/avely-contracts/sdk/actions"
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	"github.com/sirupsen/logrus"
)

type LastRewardCycleWatcher struct {
	currentLrc int
	protocol   *contracts.Protocol
	log        *core.Log
}

func main() {
	chainPtr := flag.String("chain", "local", "chain")

	flag.Parse()

	config := core.NewConfig(*chainPtr)
	sdk := core.NewAvelySDK(*config)
	log := core.NewLog()
	log.SetOutputStdout()
	log.AddSlackHook(sdk.Cfg.Slack.HookUrl, sdk.Cfg.Slack.LogLevel)
	protocol := contracts.RestoreFromState(sdk, log)
	url := sdk.GetWsURL()

	watcher := &LastRewardCycleWatcher{
		currentLrc: -1,
		protocol:   protocol,
		log:        log,
	}

	log.Debug("Start last reward cycle watcher")
	blockWatcher := utils.CreateBlockWatcher(url, log)
	blockWatcher.AddObserver(watcher)
	blockWatcher.Start()
}

// If Last reward cycly has been changed, then:
//   1. Drain Buffer
//   2. ReDelegate stakes from other SSNs
//   3. Autorestake funds
func (w *LastRewardCycleWatcher) Notify(blockNum int) {
	protocol := w.protocol
	action := actions.NewAdminActions(w.log)
	lrc := protocol.Zimpl.GetLastRewardCycle()

	if lrc == w.currentLrc {
		w.log.WithFields(logrus.Fields{
			"block_number": blockNum,
			"lrc":          lrc,
		}).Debug("Block mined, last reward cycle not changed.")
	} else {
		w.log.WithFields(logrus.Fields{
			"block_number": blockNum,
			"lrc":          lrc,
		}).Info("Block mined, New Last Reward Cycle.")

		action.DrainBuffer(protocol, lrc)
		showOnly := false
		action.ChownStakeReDelegate(protocol, showOnly)
		action.AutoRestake(protocol)

		w.currentLrc = lrc
	}
}
