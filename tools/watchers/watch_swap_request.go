package main

import (
	"flag"
	"github.com/avely-finance/avely-contracts/sdk/actions"
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
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
	protocol = contracts.RestoreFromState(sdk, log)
	url := sdk.GetWsURL()

	claimWatcher := &SwapRequestWatcher{
		gap:        *gapPtr,
		runAtBlock: -1,
	}

	log.Success("Start swap request watcher")
	blockWatcher := utils.CreateBlockWatcher(url)
	blockWatcher.AddObserver(claimWatcher)
	blockWatcher.Start()
}

func (w *SwapRequestWatcher) Notify(blockNum int) {
	if (blockNum - w.runAtBlock) > w.gap {
		log.Successf("Mined block #%d.", blockNum)
		actions.ConfirmSwapRequests(protocol)
		w.runAtBlock = blockNum
	} else {
		log.Successf("Mined block #%d, but gap=%d <= %d, skip.", blockNum, (blockNum - w.runAtBlock), w.gap)
	}
}
