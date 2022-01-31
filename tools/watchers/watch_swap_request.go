package main

import (
	"flag"

	"github.com/Zilliqa/gozilliqa-sdk/subscription"
	"github.com/avely-finance/avely-contracts/sdk/actions"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/tidwall/gjson"
)

var log *Log
var sdk *AvelySDK

// Every N blocks watch delegator's swap requests, isseed by SSNList->RequestDelegatorSwap transition
func main() {
	chainPtr := flag.String("chain", "local", "chain")
	gapPtr := flag.Int("gap", 5, "gap between blocks")

	flag.Parse()

	log = NewLog()
	log.SetOutputStdout()
	config := NewConfig(*chainPtr)
	sdk = NewAvelySDK(*config)
	protocol := RestoreFromState(sdk, log)

	url := sdk.GetWsURL()
	subscriber := subscription.BuildNewBlockSubscriber(url)
	_, ec, msg := subscriber.Start()
	cancel := false

	runAtBlock := -1 // init value should be unreal

	log.Successf("Start swap request watcher, chain=%s, gap=%d", *chainPtr, *gapPtr)
	log.Successf("Start subscription, url=%s", url.String())

	for {
		if cancel {
			_ = subscriber.Ws.Close()
			break
		}
		select {
		case message := <-msg:
			blockNum, found := getBlockNum(message)
			if found && (blockNum-runAtBlock) > *gapPtr {
				log.Successf("Mined block #%d", blockNum)
				actions.ConfirmSwapRequests(protocol)
				runAtBlock = blockNum
			} else if found {
				log.Successf("Mined block #%d, but gap=%d <= %d, skip", blockNum, (blockNum - runAtBlock), *gapPtr)
			}
		case err := <-ec:
			log.Success("Get error: ", err.Error())
			cancel = true
		}

	}
}

func getBlockNum(message []byte) (int, bool) {
	block := gjson.Get(string(message), "values.0.value.TxBlock.header")
	if block.Get("BlockNum").Exists() {
		blockNo := int(block.Get("BlockNum").Int())
		return blockNo, true
	}
	return -1, false
}
