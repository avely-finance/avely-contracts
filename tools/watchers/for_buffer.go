package main

import (
	"github.com/tidwall/gjson"
	"net/url"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/Zilliqa/gozilliqa-sdk/subscription"
)

func tryDrainBuffer(p *Protocol) {
	// p.Aimpl.DrainBuffer(addr)
}

func main() {
	log := NewLog()
	config := NewConfig("testnet")
	sdk := NewAvelySDK(*config)
	protocol := RestoreFromState(sdk, log)

	currentLrc := -1 // default valus should be unreal

	u := url.URL{Scheme: "wss", Host: "dev-ws.zilliqa.com", Path: ""}
	subscriber := subscription.BuildNewBlockSubscriber(u)
	_, ec, msg := subscriber.Start()
	cancel := false

	log.Success("Start subscribsion")

	for {
		if cancel {
			_ = subscriber.Ws.Close()
			break
		}
		select {
		case message := <-msg:
			block := gjson.Get(string(message), "values.0.value.TxBlock.header")

			if block.Get("BlockNum").Exists() {
				blockNo := block.Get("BlockNum").String()
				lrc := protocol.GetLastRewardCycle()

				log.Success("New block #" + blockNo + " is mined")
				log.Success("Current Last Reward Cycle is", lrc)

				if (lrc == currentLrc) {
					log.Success("Last Reward Cycle is not changed. Skip")
				} else {
					log.Success("New Last Reward Cycle!")
					tryDrainBuffer(protocol)

					currentLrc = lrc
				}
			}

		case err := <-ec:
			log.Success("Get error: ", err.Error())
			cancel = true
		}

	}
}
