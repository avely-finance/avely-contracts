package main

import (
	"flag"
	"github.com/tidwall/gjson"

	"github.com/Zilliqa/gozilliqa-sdk/subscription"
	"github.com/avely-finance/avely-contracts/sdk/actions"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

var log *Log
var sdk *AvelySDK

// If Last reward cycly has been changed, then:
//   1. Drain Buffer
//   2. ReDelegate stakes from other SSNs
//   3. Autorestake funds
func main() {
	chainPtr := flag.String("chain", "local", "chain")

	flag.Parse()

	log = NewLog()
	config := NewConfig(*chainPtr)
	sdk = NewAvelySDK(*config)
	protocol := RestoreFromState(sdk, log)

	currentLrc := -1 // init value should be unreal

	url := sdk.GetWsURL()
	subscriber := subscription.BuildNewBlockSubscriber(url)
	_, ec, msg := subscriber.Start()
	cancel := false

	log.Success("Start subscription")

	for {
		if cancel {
			_ = subscriber.Ws.Close()
			break
		}
		select {
		case message := <-msg:
			blockNum, found := getBlockNum(message)
			if !found {
				continue
			}
			lrc := protocol.GetLastRewardCycle()
			log.Successf("New block #%d is mined. Current Last Reward Cycle is %d", blockNum, lrc)

			if lrc == currentLrc {
				log.Success("Last Reward Cycle is not changed. Skip")
			} else {
				log.Success("New Last Reward Cycle!")
				actions.DrainBuffer(protocol, lrc)
				actions.ChownStakeReDelegate(protocol)
				actions.AutoRestake(protocol)

				currentLrc = lrc
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
